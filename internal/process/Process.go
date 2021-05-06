package process

import (
	"github.com/trmigor/distr-model/internal/messages"
	"github.com/trmigor/distr-model/internal/network"
	"github.com/trmigor/distr-model/pkg/set"
	"github.com/trmigor/distr-model/user/context"
	"time"
)

// WorkFunction represents a process working function.
type WorkFunction func(context *Process, m *messages.Message) int32

// Process models a real distributed process.
type Process struct {
	MessagesQueue *messages.MessageQueue
	Network       *network.Network
	Node          int32
	Context       map[string]context.Context
	workerThread  chan bool
	stopFlag      bool
	workers       []WorkFunction
}

// New returns a valid Process instance.
func New(node int32) *Process {
	res := &Process{
		MessagesQueue: messages.NewMessageQueue(),
		Node:          node,
		Context:       context.Contexts,
		workerThread:  make(chan bool),
		workers:       make([]WorkFunction, 0),
	}
	go workerThreadExecutor(res)
	return res
}

// Stop terminates the worker goroutine.
func (p *Process) Stop() {
	p.stopFlag = true
	<-p.workerThread
}

// Neibs returns all the neighbours of the process in its network.
func (p *Process) Neibs() *set.Set {
	return p.Network.Neibs(p.Node)
}

// RegisterWorkFunction registers a working function for the process.
func (p *Process) RegisterWorkFunction(prefix []byte, wf WorkFunction) {
	p.workers = append(p.workers, wf)
}

// IsMyMessage checks whether a message is for the process.
func (p *Process) IsMyMessage(prefix []byte, message []byte) bool {
	if len(message) > 0 && message[0] == '*' {
		return true
	}
	if len(prefix)+1 >= len(message) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if prefix[i] != message[i] {
			return false
		}
	}
	return message[len(prefix)] == '_'
}

func workerThreadExecutor(dp *Process) {
	for !dp.stopFlag {
		if dp.MessagesQueue.Size() > 0 && dp.Network.Tick >= dp.MessagesQueue.Peek().DeliveryTime {
			m := dp.MessagesQueue.Dequeue()
			for _, worker := range dp.workers {
				if worker(dp, m) != 0 {
					break
				}
			}
		}
		time.Sleep(time.Millisecond)
	}
	dp.workerThread <- true
}

// NetworkLayer returns the process network pointer for implementing network.Process interface.
func (p *Process) NetworkLayer() **network.Network {
	return &p.Network
}

// WorkerMessagesQueue returns the process message queue for implementing network.Process interface.
func (p *Process) WorkerMessagesQueue() *messages.MessageQueue {
	return p.MessagesQueue
}
