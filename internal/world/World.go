package world

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/trmigor/distr-model/internal/errors"
	"github.com/trmigor/distr-model/internal/messages"
	"github.com/trmigor/distr-model/internal/network"
	"github.com/trmigor/distr-model/internal/process"
)

// World represents the whole distributed system.
type World struct {
	Network       *network.Network
	ProcessesList []*process.Process
	Associates    map[string]process.WorkFunction
}

// New creates a new instance of a world.
func New() *World {
	return &World{
		Network:       network.New(),
		ProcessesList: make([]*process.Process, 0),
		Associates:    make(map[string]process.WorkFunction),
	}
}

// CreateProcess creates a new process for acquired node.
func (w *World) CreateProcess(node int32) int32 {
	p := process.New(node)
	if node >= int32(len(w.ProcessesList)) {
		w.ProcessesList = append(w.ProcessesList, make([]*process.Process, int(node)-len(w.ProcessesList)+1)...)
	}
	w.ProcessesList[node] = p
	w.Network.RegisterProcess(node, p)
	return node
}

// RegisterWorkFunction registers a process working function.
func (w *World) RegisterWorkFunction(function []byte, wf process.WorkFunction) {
	w.Associates[string(function)] = wf
}

// AssignWorkFunction assigns a new function for the process with given node.
func (w *World) AssignWorkFunction(node int32, function []byte) errors.ErrorCode {
	if node < 0 || node >= int32(len(w.ProcessesList)) {
		return errors.ItemNotFound
	}
	dp := w.ProcessesList[node]
	if dp == nil {
		return errors.ItemNotFound
	}

	it, ok := w.Associates[string(function)]
	if !ok {
		return errors.ItemNotFound
	}

	dp.RegisterWorkFunction(function, it)
	return errors.OK
}

// ParseConfig parses the configuration file and launches the model.
func (w *World) ParseConfig(name []byte) bool {
	data, err := ioutil.ReadFile(string(name))
	if err != nil {
		return false
	}

	bidirected := 1
	timeout := 0

	dataLines := strings.Split(string(data), "\n")
	for i := 0; i < len(dataLines); i++ {
		if len(dataLines[i]) == 0 || (len(dataLines[i]) > 0 && dataLines[i][0] == ';') {
			continue
		}

		var startprocess, endprocess, from, to, arg int32
		var latency int32 = 1
		timer := 0
		var errorRate float64
		var id, msg []byte

		if read, err := fmt.Sscanf(dataLines[i], "processes %d %d", &startprocess, &endprocess); err == nil && read == 2 {
			for i := startprocess; i <= endprocess; i++ {
				w.CreateProcess(i)
			}
			continue
		}

		if read, err := fmt.Sscanf(dataLines[i], "bidirected %d", &bidirected); err == nil && read == 1 {
			continue
		}

		if read, err := fmt.Sscanf(dataLines[i], "errorRate %f", &errorRate); err == nil && read == 1 {
			w.Network.SetErrorRate(errorRate)
			continue
		}

		if read, err := fmt.Sscanf(dataLines[i], "link from all to all latency %d", &latency); (err == nil && read == 1) || dataLines[i] == "link from all to all" {
			w.Network.AddLinksAllToAll(bidirected != 0, latency)
			continue
		}

		if read, err := fmt.Sscanf(dataLines[i], "setprocesses %d %d %s", &startprocess, &endprocess, &id); err == nil && read == 3 {
			for i := startprocess; i <= endprocess; i++ {
				res := w.AssignWorkFunction(i, id)
				if res != errors.OK {
					return false
				}
			}
			continue
		}

		if read, err := fmt.Sscanf(dataLines[i], "send from %d to %d %s %d", &from, &to, &msg, &arg); read == 4 && err == nil {
			w.Network.SendMessage(from, to, messages.NewMessageByArgs(messages.NewMessageArg(msg), messages.NewMessageArg(arg)))
			continue
		}

		if read, err := fmt.Sscanf(dataLines[i], "send from %d to %d %s", &from, &to, &msg); read == 3 && err == nil {
			w.Network.SendMessage(from, to, messages.NewMessageByArgs(messages.NewMessageArg(msg)))
			continue
		}

		if read, err := fmt.Sscanf(dataLines[i], "wait %d", &timeout); read == 1 && err == nil {
			time.Sleep(time.Duration(timeout) * time.Second)
			continue
		}

		if read, err := fmt.Sscanf(dataLines[i], "launch timer %d", &timer); read == 1 && err == nil {
			go network.TimerSender(w.Network, timer)
			continue
		}

		if read, err := fmt.Sscanf(dataLines[i], "link from %d to %d latency %d", &from, &to, &latency); err == nil && read == 3 {
			w.Network.CreateLink(from, to, bidirected != 0, latency)
			continue
		}

		if read, err := fmt.Sscanf(dataLines[i], "link from %d to %d", &from, &to); err == nil && read == 2 {
			w.Network.CreateLink(from, to, bidirected != 0, latency)
			continue
		}

		if read, err := fmt.Sscanf(dataLines[i], "link from %d to all latency %d", &from, &latency); err == nil && read == 2 {
			w.Network.AddLinksToAll(from, bidirected != 0, latency)
			continue
		}

		if read, err := fmt.Sscanf(dataLines[i], "link from %d to all", &from); err == nil && read == 1 {
			w.Network.AddLinksToAll(from, bidirected != 0, latency)
			continue
		}

		if read, err := fmt.Sscanf(dataLines[i], "link from all to %d latency %d", &to, &latency); err == nil && read == 2 {
			w.Network.AddLinksFromAll(to, bidirected != 0, latency)
			continue
		}

		if read, err := fmt.Sscanf(dataLines[i], "link from all to %d", &to); err == nil && read == 1 {
			w.Network.AddLinksFromAll(to, bidirected != 0, latency)
			continue
		}

		res := "Unknown directive in input file: " + dataLines[i]
		panic(res)
	}
	return true
}

// Stop terminates the model work.
// It should be called at the end of model usage.
func (w *World) Stop() {
	w.Network.Stop()
	for _, p := range w.ProcessesList {
		if p != nil {
			p.Stop()
		}
	}
}
