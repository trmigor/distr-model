package network

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/trmigor/distr-model/internal/errors"
	"github.com/trmigor/distr-model/internal/messages"
	"github.com/trmigor/distr-model/pkg/set"
)

func TestNetwork_SetErrorRate(t *testing.T) {
	type args struct {
		rate float64
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid", args{0.5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nl := New()
			defer nl.Stop()
			nl.SetErrorRate(tt.args.rate)
			if nl.ErrorRate != tt.args.rate {
				t.Errorf("Network.SetErrorRate(): error rate not set")
			}
		})
	}
}

func TestNetwork_CreateLink(t *testing.T) {
	t.Run("Same", func(t *testing.T) {
		nl := New()
		defer nl.Stop()

		nl.CreateLink(0, 0, false, 0)
		if len(nl.networkMap) > 0 {
			t.Errorf("Created a loop")
		}
	})

	t.Run("NonBidir", func(t *testing.T) {
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
	})

	t.Run("Bidir", func(t *testing.T) {
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
	})
}

func TestNetwork_GetLink(t *testing.T) {
	type args struct {
		p1 int32
		p2 int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{"Special", args{-1, 1}, 0},
		{"Same", args{0, 0}, 0},
		{"Common", args{0, 1}, 1},
		{"NotFound", args{1, 2}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nl := New()
			defer nl.Stop()
			nl.CreateLink(0, 1, false, 1)
			if got := nl.GetLink(tt.args.p1, tt.args.p2); got != tt.want {
				t.Errorf("Network.GetLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetwork_SendBytes(t *testing.T) {
	type args struct {
		fromProcess int32
		toProcess   int32
		msg         []byte
	}
	tests := []struct {
		name string
		args args
		want errors.ErrorCode
	}{
		{"TooBig", args{0, 4, []byte{65, 1, 0, 0, 0}}, errors.SizeTooBig},
		{"TimeOut", args{0, 1, []byte{65, 1, 0, 0, 0}}, errors.TimeOut},
		{"NoProcess", args{0, 3, []byte{65, 1, 0, 0, 0}}, errors.ItemNotFound},
		{"NoLink", args{0, 2, []byte{65, 1, 0, 0, 0}}, errors.ItemNotFound},
		{"Valid", args{0, 1, []byte{65, 1, 0, 0, 0}}, errors.OK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nl := New()
			defer nl.Stop()

			nl.QueueMap = make([]*messages.MessageQueue, 4)
			nl.QueueMap[0] = messages.NewMessageQueue()
			nl.QueueMap[1] = messages.NewMessageQueue()
			nl.QueueMap[2] = messages.NewMessageQueue()
			nl.networkSize = 4
			nl.CreateLink(0, 1, true, 1)

			if tt.name == "TimeOut" {
				nl.SetErrorRate(1)
			}

			if got := nl.SendBytes(tt.args.fromProcess, tt.args.toProcess, tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Network.SendBytes() = %v, want %v", got, tt.want)
			}
			if tt.name == "Valid" {
				if nl.QueueMap[1].Size() != 1 {
					t.Errorf("Message is not sent")
				}
			}
		})
	}
}

func TestNetwork_SendMessage(t *testing.T) {
	type args struct {
		fromProcess int32
		toProcess   int32
		msg         *messages.Message
	}
	tests := []struct {
		name string
		args args
		want errors.ErrorCode
	}{
		{"Direct", args{0, 1, messages.NewMessageByArgs(messages.NewMessageArg(int32(1)))}, errors.OK},
		{"Broadcast", args{0, -1, messages.NewMessageByArgs(messages.NewMessageArg(int32(1)))}, errors.OK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nl := New()
			defer nl.Stop()

			nl.QueueMap = make([]*messages.MessageQueue, 2)
			nl.QueueMap[0] = messages.NewMessageQueue()
			nl.QueueMap[1] = messages.NewMessageQueue()
			nl.networkSize = 2
			nl.CreateLink(0, 1, true, 1)

			if got := nl.SendMessage(tt.args.fromProcess, tt.args.toProcess, tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Network.SendMessage() = %v, want %v", got, tt.want)
			}
			if nl.QueueMap[1].Size() != 1 {
				t.Errorf("Message is not sent")
			}
		})
	}
}

func TestNetwork_AddLinksToAll(t *testing.T) {
	type args struct {
		from          int32
		bidirectional bool
		latency       int32
	}
	tests := []struct {
		name string
		args args
	}{
		{"NonBidir", args{0, false, 1}},
		{"Bidir", args{0, true, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nl := New()
			defer nl.Stop()

			nl.QueueMap = make([]*messages.MessageQueue, 3)
			nl.QueueMap[0] = messages.NewMessageQueue()
			nl.QueueMap[1] = messages.NewMessageQueue()
			nl.QueueMap[2] = messages.NewMessageQueue()
			nl.networkSize = 3

			nl.AddLinksToAll(tt.args.from, tt.args.bidirectional, tt.args.latency)

			for i := int32(0); i < nl.networkSize; i++ {
				if i == tt.args.from {
					continue
				}
				if nl.networkMap[tt.args.from][i] != tt.args.latency {
					t.Errorf("Network.AddLinksToAll(): link %v %v not created", tt.args.from, i)
				}
			}

			if tt.args.bidirectional {
				for i := int32(0); i < nl.networkSize; i++ {
					if i == tt.args.from {
						continue
					}
					if nl.networkMap[i][tt.args.from] != tt.args.latency {
						t.Errorf("Network.AddLinksToAll(): link %v %v not created", i, tt.args.from)
					}
				}
			} else {
				for i := int32(0); i < nl.networkSize; i++ {
					if i == tt.args.from {
						continue
					}
					if _, ok := nl.networkMap[i]; ok {
						t.Errorf("Network.AddLinksToAll(): link %v %v created", i, tt.args.from)
					}
				}
			}
		})
	}
}

func TestNetwork_AddLinksFromAll(t *testing.T) {
	type args struct {
		to            int32
		bidirectional bool
		latency       int32
	}
	tests := []struct {
		name string
		args args
	}{
		{"NonBidir", args{0, false, 1}},
		{"Bidir", args{0, true, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nl := New()
			defer nl.Stop()

			nl.QueueMap = make([]*messages.MessageQueue, 3)
			nl.QueueMap[0] = messages.NewMessageQueue()
			nl.QueueMap[1] = messages.NewMessageQueue()
			nl.QueueMap[2] = messages.NewMessageQueue()
			nl.networkSize = 3

			nl.AddLinksFromAll(tt.args.to, tt.args.bidirectional, tt.args.latency)

			for i := int32(0); i < nl.networkSize; i++ {
				if i == tt.args.to {
					continue
				}
				if nl.networkMap[i][tt.args.to] != tt.args.latency {
					t.Errorf("Network.AddLinksToAll(): link %v %v not created", i, tt.args.to)
				}
			}

			if tt.args.bidirectional {
				for i := int32(0); i < nl.networkSize; i++ {
					if i == tt.args.to {
						continue
					}
					if nl.networkMap[tt.args.to][i] != tt.args.latency {
						t.Errorf("Network.AddLinksToAll(): link %v %v not created", i, tt.args.to)
					}
				}
			} else {
				if _, ok := nl.networkMap[tt.args.to]; ok {
					t.Errorf("Network.AddLinksToAll(): inverted links created")
				}
			}
		})
	}
}

func TestNetwork_AddLinksAllToAll(t *testing.T) {
	type args struct {
		bidirectional bool
		latency       int32
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid", args{true, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nl := New()
			defer nl.Stop()

			nl.QueueMap = make([]*messages.MessageQueue, 3)
			nl.QueueMap[0] = messages.NewMessageQueue()
			nl.QueueMap[1] = messages.NewMessageQueue()
			nl.QueueMap[2] = messages.NewMessageQueue()
			nl.networkSize = 3

			nl.AddLinksAllToAll(tt.args.bidirectional, tt.args.latency)
			for i := int32(0); i < nl.networkSize; i++ {
				for j := int32(0); j < nl.networkSize; j++ {
					if nl.networkMap[i][j] != 0 {
						t.Errorf("Network.AddLinksToAll(): link %v %v not created", i, j)
					}
				}
			}
		})
	}
}

func TestNetwork_Neibs(t *testing.T) {
	type args struct {
		from int32
	}
	tests := []struct {
		name string
		args args
		want *set.Set
	}{
		{"NoNeibs", args{2}, set.New()},
		{"OneNeib", args{1}, set.New()},
		{"TwoNeibs", args{0}, set.New()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nl := New()
			defer nl.Stop()

			nl.QueueMap = make([]*messages.MessageQueue, 3)
			nl.QueueMap[0] = messages.NewMessageQueue()
			nl.QueueMap[1] = messages.NewMessageQueue()
			nl.QueueMap[2] = messages.NewMessageQueue()
			nl.networkSize = 3

			nl.CreateLink(0, 1, false, 0)
			nl.CreateLink(0, 2, false, 0)
			nl.CreateLink(1, 2, false, 0)

			if tt.name == "OneNeib" {
				tt.want.Insert(int32(2))
			}

			if tt.name == "TwoNeibs" {
				tt.want.Insert(int32(1))
				tt.want.Insert(int32(2))
			}

			if got := nl.Neibs(tt.args.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Network.Neibs() = %v, want %v", got, tt.want)
			}
		})
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

func TestNetwork_RegisterProcess(t *testing.T) {
	type args struct {
		node int32
		dp   Process
	}
	tests := []struct {
		name string
		args args
		want errors.ErrorCode
	}{
		{"Valid", args{1, &process{nl: nil, mq: messages.NewMessageQueue()}}, errors.OK},
		{"Duplicate", args{0, &process{nl: nil, mq: messages.NewMessageQueue()}}, errors.DuplicateItems},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nl := New()
			defer nl.Stop()

			if got := nl.RegisterProcess(0, &process{nl: nil, mq: messages.NewMessageQueue()}); !reflect.DeepEqual(got, errors.OK) {
				t.Errorf("Network.RegisterProcess() = %v, want %v", got, errors.OK)
			}

			if got := nl.RegisterProcess(tt.args.node, tt.args.dp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Network.RegisterProcess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimerSender(t *testing.T) {
	type args struct {
		nl  *Network
		nap int
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid", args{New(), 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.nl.networkSize = 2
			tt.args.nl.QueueMap = make([]*messages.MessageQueue, 2)
			tt.args.nl.QueueMap[0] = messages.NewMessageQueue()
			tt.args.nl.QueueMap[1] = messages.NewMessageQueue()
			tt.args.nl.CreateLink(0, 1, true, 0)

			go TimerSender(tt.args.nl, tt.args.nap)

			time.Sleep(1500 * time.Duration(tt.args.nap) * time.Millisecond)
			tt.args.nl.Stop()

			if tt.args.nl.QueueMap[0].Size() != 2 {
				t.Errorf("Message is not sent")
			}

			if tt.args.nl.QueueMap[1].Size() != 2 {
				t.Errorf("Message is not sent")
			}
		})
	}
}
