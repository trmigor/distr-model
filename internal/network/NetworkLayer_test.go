package network

import (
	"github.com/trmigor/distr-model/internal/errors"
	"github.com/trmigor/distr-model/internal/messages"
	"math/rand"
	"testing"
	"time"
)

func TestSetErrorRate(t *testing.T) {
	rand.Seed((time.Now().UnixNano()))
	nl := New()
	defer nl.Stop()
	value := rand.Float64()
	nl.SetErrorRate(value)
	if value != nl.ErrorRate {
		t.Errorf("Results mismatch: expected %v, got %v", value, nl.ErrorRate)
	}
}

func TestTick(t *testing.T) {
	nl := New()
	defer nl.Stop()
	time.Sleep(1*time.Second + 500*time.Millisecond)
	if nl.Tick != 1 {
		t.Errorf("Results mismatch: expected %v, got %v", 1, nl.Tick)
	}
}

func TestCreateLinkSame(t *testing.T) {
	nl := New()
	defer nl.Stop()

	nl.CreateLink(0, 0, false, 0)
	if len(nl.networkMap) > 0 {
		t.Errorf("Created a loop")
	}
}

func TestCreateLinkNonBidir(t *testing.T) {
	rand.Seed((time.Now().UnixNano()))
	nl := New()
	defer nl.Stop()

	cost := int32(rand.Intn(100) + 100)
	nl.CreateLink(0, 1, false, cost)

	if m, ok := nl.networkMap[0]; ok {
		if v, ok := m[1]; ok {
			if v != cost {
				t.Errorf("Results mismatch: expected %v, got %v", cost, v)
			}
		} else {
			t.Errorf("Value is not inserted")
		}
	} else {
		t.Errorf("Map is not inserted")
	}
	if _, ok := nl.networkMap[1]; ok {
		t.Errorf("Created inverted link")
	}
}

func TestCreateLinkBidir(t *testing.T) {
	rand.Seed((time.Now().UnixNano()))
	nl := New()
	defer nl.Stop()

	cost := int32(rand.Intn(100) + 100)
	nl.CreateLink(0, 1, true, cost)

	if m, ok := nl.networkMap[0]; ok {
		if v, ok := m[1]; ok {
			if v != cost {
				t.Errorf("Results mismatch: expected %v, got %v", cost, v)
			}
		} else {
			t.Errorf("Value is not inserted")
		}
	} else {
		t.Errorf("Map is not inserted")
	}

	if m, ok := nl.networkMap[1]; ok {
		if v, ok := m[0]; ok {
			if v != cost {
				t.Errorf("Results of inverted link mismatch: expected %v, got %v", cost, v)
			}
		} else {
			t.Errorf("Value for inverted link is not inserted")
		}
	} else {
		t.Errorf("Map for inverted link is not inserted")
	}
}

func TestGetLinkSpecial(t *testing.T) {
	nl := New()
	defer nl.Stop()

	if res := nl.GetLink(-1, 0); res != 0 {
		t.Errorf("Results mismatch: expected %v, got %v", 0, res)
	}

	if res := nl.GetLink(0, 0); res != 0 {
		t.Errorf("Results mismatch: expected %v, got %v", 0, res)
	}
}

func TestGetLinkNotFound(t *testing.T) {
	nl := New()
	defer nl.Stop()

	if res := nl.GetLink(0, 1); res != -1 {
		t.Errorf("Results mismatch: expected %v, got %v", -1, res)
	}
}

func TestGetLink(t *testing.T) {
	rand.Seed((time.Now().UnixNano()))
	nl := New()
	defer nl.Stop()

	cost := int32(rand.Intn(100) + 100)
	nl.CreateLink(0, 1, false, cost)

	if res := nl.GetLink(0, 1); res != cost {
		t.Errorf("Results mismatch: expected %v, got %v", cost, res)
	}
}

func TestSendBytesTooBig(t *testing.T) {
	nl := New()
	defer nl.Stop()

	if res := nl.SendBytes(0, 1, []byte("Lorem")); res != errors.SizeTooBig {
		t.Errorf("Results mismatch: expected %v, got %v", errors.SizeTooBig, res)
	}
}

func TestSendBytesTimeOut(t *testing.T) {
	nl := New()
	defer nl.Stop()

	nl.networkSize = 2
	nl.ErrorRate = 1

	if res := nl.SendBytes(0, 1, []byte("Lorem")); res != errors.TimeOut {
		t.Errorf("Results mismatch: expected %v, got %v", errors.TimeOut, res)
	}
}

func TestSendBytesNil(t *testing.T) {
	nl := New()
	defer nl.Stop()

	nl.networkSize = 2
	nl.QueueMap = make([]*messages.MessageQueue, 2)

	if res := nl.SendBytes(0, 1, []byte("Lorem")); res != errors.ItemNotFound {
		t.Errorf("Results mismatch: expected %v, got %v", errors.ItemNotFound, res)
	}
}

func TestSendBytesNoConnection(t *testing.T) {
	nl := New()
	defer nl.Stop()

	nl.networkSize = 2
	nl.QueueMap = make([]*messages.MessageQueue, 2)
	nl.QueueMap[0] = messages.NewMessageQueue()
	nl.QueueMap[1] = messages.NewMessageQueue()

	if res := nl.SendBytes(0, 1, []byte("Lorem")); res != errors.ItemNotFound {
		t.Errorf("Results mismatch: expected %v, got %v", errors.ItemNotFound, res)
	}
}

func TestSendBytes(t *testing.T) {
	rand.Seed((time.Now().UnixNano()))
	nl := New()
	defer nl.Stop()

	nl.networkSize = 2
	nl.QueueMap = make([]*messages.MessageQueue, 2)
	nl.QueueMap[0] = messages.NewMessageQueue()
	nl.QueueMap[1] = messages.NewMessageQueue()

	cost := int32(rand.Intn(100) + 100)
	nl.CreateLink(0, 1, true, cost)

	if res := nl.SendBytes(0, 1, []byte("Lorem")); res != errors.OK {
		t.Errorf("Results mismatch: expected %v, got %v", errors.OK, res)
	}

	if nl.QueueMap[1].Size() != 1 {
		t.Errorf("Message is not sent")
	}
}

func TestSendMessage(t *testing.T) {
	rand.Seed((time.Now().UnixNano()))
	nl := New()
	defer nl.Stop()

	nl.networkSize = 2
	nl.QueueMap = make([]*messages.MessageQueue, 2)
	nl.QueueMap[0] = messages.NewMessageQueue()
	nl.QueueMap[1] = messages.NewMessageQueue()

	cost := int32(rand.Intn(100) + 100)
	nl.CreateLink(0, 1, true, cost)

	msg := messages.NewMessage(0, 1, []byte("Lorem"))

	if res := nl.SendMessage(0, 1, msg); res != errors.OK {
		t.Errorf("Results mismatch: expected %v, got %v", errors.OK, res)
	}

	if nl.QueueMap[1].Size() != 1 {
		t.Errorf("Message is not sent")
	}
}

func TestSendMessageBroadCast(t *testing.T) {
	rand.Seed((time.Now().UnixNano()))
	nl := New()
	defer nl.Stop()

	nl.networkSize = 2
	nl.QueueMap = make([]*messages.MessageQueue, 2)
	nl.QueueMap[0] = messages.NewMessageQueue()
	nl.QueueMap[1] = messages.NewMessageQueue()

	cost := int32(rand.Intn(100) + 100)
	nl.CreateLink(0, 1, true, cost)

	msg := messages.NewMessage(0, 1, []byte("Lorem"))

	if res := nl.SendMessage(0, -1, msg); res != errors.OK {
		t.Errorf("Results mismatch: expected %v, got %v", errors.OK, res)
	}

	if nl.QueueMap[1].Size() != 1 {
		t.Errorf("Message is not sent")
	}
}

func TestAddLinksToAllNonBidir(t *testing.T) {
	rand.Seed((time.Now().UnixNano()))
	nl := New()
	defer nl.Stop()

	nl.networkSize = 5

	cost := int32(rand.Intn(100) + 100)
	nl.AddLinksToAll(0, false, cost)

	if m, ok := nl.networkMap[0]; ok {
		if _, ok := m[0]; ok {
			t.Errorf("Created a loop")
		}
		for i := 1; i < 5; i++ {
			if v, ok := m[int32(i)]; ok {
				if v != cost {
					t.Errorf("Results mismatch: expected %v, got %v", cost, v)
				}
			} else {
				t.Errorf("Value is not inserted")
			}
		}
	} else {
		t.Errorf("Map is not inserted")
	}

	for i := 1; i < 5; i++ {
		if _, ok := nl.networkMap[int32(i)]; ok {
			t.Errorf("Created inverted link")
		}
	}
}

func TestAddLinksToAllBidir(t *testing.T) {
	rand.Seed((time.Now().UnixNano()))
	nl := New()
	defer nl.Stop()

	nl.networkSize = 5

	cost := int32(rand.Intn(100) + 100)
	nl.AddLinksToAll(0, true, cost)

	if m, ok := nl.networkMap[0]; ok {
		if _, ok := m[0]; ok {
			t.Errorf("Created a loop")
		}
		for i := 1; i < 5; i++ {
			if v, ok := m[int32(i)]; ok {
				if v != cost {
					t.Errorf("Results mismatch: expected %v, got %v", cost, v)
				}
			} else {
				t.Errorf("Value is not inserted")
			}
		}
	} else {
		t.Errorf("Map is not inserted")
	}

	for i := 1; i < 5; i++ {
		if m, ok := nl.networkMap[int32(i)]; ok {
			if v, ok := m[0]; ok {
				if v != cost {
					t.Errorf("Results mismatch: expected %v, got %v", cost, v)
				}
			} else {
				t.Errorf("Value for inverted link is not inserted")
			}
		} else {
			t.Errorf("Map for inverted link is not inserted")
		}
	}
}

func TestAddLinksFromAllNonBidir(t *testing.T) {
	rand.Seed((time.Now().UnixNano()))
	nl := New()
	defer nl.Stop()

	nl.networkSize = 5

	cost := int32(rand.Intn(100) + 100)
	nl.AddLinksFromAll(0, false, cost)

	for i := 1; i < 5; i++ {
		if m, ok := nl.networkMap[int32(i)]; ok {
			if v, ok := m[0]; ok {
				if v != cost {
					t.Errorf("Results mismatch: expected %v, got %v", cost, v)
				}
			} else {
				t.Errorf("Value for inverted link is not inserted")
			}
		} else {
			t.Errorf("Map for inverted link is not inserted")
		}
	}

	if _, ok := nl.networkMap[0]; ok {
		t.Errorf("Created link")
	}
}

func TestAddLinksFromAllBidir(t *testing.T) {
	rand.Seed((time.Now().UnixNano()))
	nl := New()
	defer nl.Stop()

	nl.networkSize = 5

	cost := int32(rand.Intn(100) + 100)
	nl.AddLinksFromAll(0, true, cost)

	if m, ok := nl.networkMap[0]; ok {
		if _, ok := m[0]; ok {
			t.Errorf("Created a loop")
		}
		for i := 1; i < 5; i++ {
			if v, ok := m[int32(i)]; ok {
				if v != cost {
					t.Errorf("Results mismatch: expected %v, got %v", cost, v)
				}
			} else {
				t.Errorf("Value is not inserted")
			}
		}
	} else {
		t.Errorf("Map is not inserted")
	}

	for i := 1; i < 5; i++ {
		if m, ok := nl.networkMap[int32(i)]; ok {
			if v, ok := m[0]; ok {
				if v != cost {
					t.Errorf("Results mismatch: expected %v, got %v", cost, v)
				}
			} else {
				t.Errorf("Value for inverted link is not inserted")
			}
		} else {
			t.Errorf("Map for inverted link is not inserted")
		}
	}
}

func TestAddLinksAllToAll(t *testing.T) {
	rand.Seed((time.Now().UnixNano()))
	nl := New()
	defer nl.Stop()

	nl.networkSize = 5

	cost := int32(rand.Intn(100) + 100)
	nl.AddLinksAllToAll(true, cost)

	for i := int32(0); i < 5; i++ {
		if m, ok := nl.networkMap[i]; ok {
			for j := int32(0); j < 5; j++ {
				if v, ok := m[j]; ok {
					if i == j {
						t.Errorf("Created a loop")
					}
					if v != cost {
						t.Errorf("Results mismatch: expected %v, got %v", cost, v)
					}
				} else {
					if i != j {
						t.Errorf("Value is not inserted")
					}
				}
			}
		} else {
			t.Errorf("Map is not inserted")
		}
	}
}

func TestNeibs(t *testing.T) {
	nl := New()
	defer nl.Stop()

	nl.networkSize = 5

	nl.AddLinksToAll(0, true, 0)

	neibs := nl.Neibs(0)
	if neibs.Contains(0) {
		t.Errorf("Created a loop")
	}
	for i := int32(1); i < 5; i++ {
		if !neibs.Contains(i) {
			t.Errorf("Neighbour not found: %v", i)
		}
	}

	for i := int32(1); i < 5; i++ {
		neibs := nl.Neibs(i)
		if !neibs.Contains(int32(0)) {
			t.Errorf("Neighbour not found: %v", i)
		}
	}
}

type process struct {
	nl *Network
	mq *messages.MessageQueue
}

func (p *process) NetworkLayer() **Network {
	return &p.nl
}

func (p *process) WorkerMessagesQueue() *messages.MessageQueue {
	return p.mq
}

func TestRegisterProcess(t *testing.T) {
	nl := New()
	defer nl.Stop()

	p := &process{
		nl: nil,
		mq: messages.NewMessageQueue(),
	}

	if res := nl.RegisterProcess(0, p); res != errors.OK {
		t.Errorf("Results mismatch: expected %v, got %v", errors.OK, res)
	}

	if len(nl.QueueMap) != 1 {
		t.Errorf("Not inserted to queue map")
	}

	if nl.networkSize != 1 {
		t.Errorf("Network size not changed")
	}
}

func TestRegisterProcessDuplicate(t *testing.T) {
	nl := New()
	defer nl.Stop()

	p := &process{
		nl: nil,
		mq: messages.NewMessageQueue(),
	}

	if res := nl.RegisterProcess(0, p); res != errors.OK {
		t.Errorf("Results mismatch: expected %v, got %v", errors.OK, res)
	}

	if res := nl.RegisterProcess(0, p); res != errors.DuplicateItems {
		t.Errorf("Results mismatch: expected %v, got %v", errors.DuplicateItems, res)
	}

	if len(nl.QueueMap) != 1 {
		t.Errorf("Queue map size mismatch: expected %v, got %v", 1, len(nl.QueueMap))
	}

	if nl.networkSize != 1 {
		t.Errorf("Network size mismatch: expected %v, got %v", 1, len(nl.QueueMap))
	}
}

func TestTimerSender(t *testing.T) {
	rand.Seed((time.Now().UnixNano()))
	nl := New()

	nl.networkSize = 2
	nl.QueueMap = make([]*messages.MessageQueue, 2)
	nl.QueueMap[0] = messages.NewMessageQueue()
	nl.QueueMap[1] = messages.NewMessageQueue()

	cost := int32(rand.Intn(100) + 100)
	nl.CreateLink(0, 1, true, cost)

	go TimerSender(nl, 1)

	time.Sleep(1500 * time.Millisecond)
	nl.Stop()

	if nl.QueueMap[0].Size() != 2 {
		t.Errorf("Message is not sent")
	}

	if nl.QueueMap[1].Size() != 2 {
		t.Errorf("Message is not sent")
	}
}
