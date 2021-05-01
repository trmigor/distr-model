package network

import (
	mt "github.com/seehuhn/mt19937"
	"github.com/trmigor/distr-model/internal/errors"
	"github.com/trmigor/distr-model/internal/messages"
	"github.com/trmigor/distr-model/pkg/set"
	"math/rand"
	"time"
)

// Network is a network infrastructure. Every process have to register in it.
// It also registers connections between processes and sends messages to them.
type Network struct {
	QueueMap    []*messages.MessageQueue
	ErrorRate   float64
	Rng         *rand.Rand
	Tick        int64
	StopFlag    bool
	networkSize int32
	networkMap  map[int32]map[int32]int32
	globalTimer chan bool
}

func globalTimerExecutor(nl *Network) {
	start := time.Now()
	for !nl.StopFlag {
		cl := time.Since(start)
		nl.Tick = cl.Milliseconds() / 1000
		time.Sleep(100 * time.Millisecond)
	}
	nl.globalTimer <- true
}

// New creates a new instance of a network layer.
func New() *Network {
	nl := &Network{
		Rng:         rand.New(mt.New()),
		networkMap:  make(map[int32]map[int32]int32),
		globalTimer: make(chan bool),
	}
	nl.Rng.Seed(time.Now().UnixNano())
	go globalTimerExecutor(nl)
	return nl
}

// Stop is a destructor, stopping the global timer.
// It should be called at the end of network usage.
func (nl *Network) Stop() {
	nl.StopFlag = true
	<-nl.globalTimer
}

// SetErrorRate sets rate of connection errors.
func (nl *Network) SetErrorRate(rate float64) {
	nl.ErrorRate = rate
}

// CreateLink enables connection between two processes and sets timing cost for message sending.
// Could be bidirectional.
func (nl *Network) CreateLink(from int32, to int32, bidirectional bool, cost int32) {
	if from == to {
		return
	}
	if _, ok := nl.networkMap[from]; !ok {
		nl.networkMap[from] = make(map[int32]int32)
	}
	nl.networkMap[from][to] = cost
	if bidirectional {
		if _, ok := nl.networkMap[to]; !ok {
			nl.networkMap[to] = make(map[int32]int32)
		}
		nl.networkMap[to][from] = cost
	}
}

// GetLink returns the cost of message sending or -1 if there is no connection.
func (nl *Network) GetLink(p1 int32, p2 int32) int32 {
	if p1 < 0 || p1 == p2 {
		return 0
	}
	if m, ok := nl.networkMap[p1]; ok {
		if cost, ok := m[p2]; ok {
			return cost
		}
	}
	return -1
}

// SendMessage models message sending between two processes.
func (nl *Network) SendMessage(fromProcess int32, toProcess int32, msg *messages.Message) errors.ErrorCode {
	if toProcess >= 0 {
		return nl.SendBytes(fromProcess, toProcess, msg.Body)
	}
	for i, mq := range nl.QueueMap {
		if mq != nil {
			nl.SendBytes(fromProcess, int32(i), msg.Body)
		}
	}
	return errors.OK
}

// SendBytes sends a byte vector from one process to another.
func (nl *Network) SendBytes(fromProcess int32, toProcess int32, msg []byte) errors.ErrorCode {
	if toProcess >= nl.networkSize {
		return errors.SizeTooBig
	}
	m := messages.NewMessage(fromProcess, toProcess, msg)
	if nl.ErrorRate > 0 && nl.Rng.Float64() < nl.ErrorRate {
		return errors.TimeOut
	}
	if nl.QueueMap[toProcess] == nil {
		return errors.ItemNotFound
	}
	p := nl.GetLink(fromProcess, toProcess)
	if p < 0 {
		return errors.ItemNotFound
	}
	m.SendTime = nl.Tick
	m.DeliveryTime = nl.Tick + int64(p)
	nl.QueueMap[toProcess].Enqueue(m)
	return errors.OK
}

// AddLinksToAll adds connections from requested process to all of others.
func (nl *Network) AddLinksToAll(from int32, bidirectional bool, latency int32) {
	if _, ok := nl.networkMap[from]; !ok {
		nl.networkMap[from] = make(map[int32]int32)
	}
	for i := int32(0); i < nl.networkSize; i++ {
		if from != i {
			nl.networkMap[from][i] = latency
		}
	}
	if bidirectional {
		for i := int32(0); i < nl.networkSize; i++ {
			if _, ok := nl.networkMap[i]; !ok {
				nl.networkMap[i] = make(map[int32]int32)
			}
			if from != i {
				nl.networkMap[i][from] = latency
			}
		}
	}
}

// AddLinksFromAll adds connections to requested process from all of others.
func (nl *Network) AddLinksFromAll(to int32, bidirectional bool, latency int32) {
	for i := int32(0); i < nl.networkSize; i++ {
		if to != i {
			if _, ok := nl.networkMap[i]; !ok {
				nl.networkMap[i] = make(map[int32]int32)
			}
			nl.networkMap[i][to] = latency
		}
	}
	if bidirectional {
		if _, ok := nl.networkMap[to]; !ok {
			nl.networkMap[to] = make(map[int32]int32)
		}
		for i := int32(0); i < nl.networkSize; i++ {
			if to != i {
				nl.networkMap[to][i] = latency
			}
		}
	}
}

// AddLinksAllToAll adds connections from all processes to all processes, except themselves.
func (nl *Network) AddLinksAllToAll(bidirectional bool, latency int32) {
	for i := int32(0); i < nl.networkSize; i++ {
		nl.AddLinksFromAll(i, bidirectional, latency)
	}
}

// Neibs returns a set of neighbours of requested process.
func (nl *Network) Neibs(from int32) *set.Set {
	res := set.New()
	if m, ok := nl.networkMap[from]; ok {
		for v := range m {
			res.Insert(v)
		}
	}
	return res
}

// Process ensures requirements for process structure.
type Process interface {
	NetworkLayer() **Network
	WorkerMessagesQueue() *messages.MessageQueue
}

// RegisterProcess registers a process in a network.
func (nl *Network) RegisterProcess(node int32, dp Process) errors.ErrorCode {
	*(dp.NetworkLayer()) = nl
	if int(node) >= len(nl.QueueMap) {
		nl.QueueMap = append(nl.QueueMap, make([]*messages.MessageQueue, int(node)-len(nl.QueueMap)+1)...)
	}
	if nl.QueueMap[node] != nil {
		return errors.DuplicateItems
	}
	nl.QueueMap[node] = dp.WorkerMessagesQueue()
	nl.networkSize = int32(len(nl.QueueMap))
	return errors.OK
}

// TimerSender sends timer message every nap seconds.
func TimerSender(nl *Network, nap int) {
	current := int32(0)
	for !nl.StopFlag {
		arg1 := messages.NewMessageArg([]byte("*TIME"))
		arg2 := messages.NewMessageArg(current)
		nl.SendMessage(-1, -1, messages.NewMessageByArgs(arg1, arg2))
		current++
		time.Sleep(time.Duration(nap) * time.Second)
	}
}
