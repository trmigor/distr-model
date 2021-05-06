package process

import (
	"github.com/trmigor/distr-model/internal/messages"
	"github.com/trmigor/distr-model/internal/network"
	"math/rand"
	"testing"
	"time"
)

func TestIsMyMessageStar(t *testing.T) {
	p := New(0)
	defer p.Stop()
	if !p.IsMyMessage([]byte(""), []byte("*")) {
		t.Errorf("Messages starting with a star are not processed")
	}
}

func TestIsMyMessageLongPrefix(t *testing.T) {
	p := New(0)
	defer p.Stop()
	if p.IsMyMessage([]byte("Lorem ipsum"), []byte("")) {
		t.Errorf("Messages with too long prefix are not processed")
	}
}

func TestIsMyMessageNotMatching(t *testing.T) {
	p := New(0)
	defer p.Stop()
	if p.IsMyMessage([]byte("Lorem"), []byte("Larem_ipsum")) {
		t.Errorf("Message matches wrong prefix")
	}
}

func TestIsMyMessageMatching(t *testing.T) {
	p := New(0)
	defer p.Stop()
	if !p.IsMyMessage([]byte("Lorem"), []byte("Lorem_ipsum")) {
		t.Errorf("Message does not match right prefix")
	}
}

func TestNeibs(t *testing.T) {
	nl := network.New()
	defer nl.Stop()

	p0 := New(0)
	p1 := New(1)
	p2 := New(2)
	p3 := New(3)
	p4 := New(4)

	defer p0.Stop()
	defer p1.Stop()
	defer p2.Stop()
	defer p3.Stop()
	defer p4.Stop()

	nl.RegisterProcess(0, p0)
	nl.RegisterProcess(1, p1)
	nl.RegisterProcess(2, p2)
	nl.RegisterProcess(3, p3)
	nl.RegisterProcess(4, p4)

	nl.AddLinksToAll(0, true, 0)

	neibs := p0.Neibs()
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

func TestRegisterWorkFunction(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	num := int32(rand.Int())
	f := func(context *Process, m *messages.Message) int32 {
		return num
	}
	p := New(0)
	defer p.Stop()

	p.RegisterWorkFunction([]byte("*TIME"), f)

	res := p.workers[0](p, messages.NewMessage(-1, -1, []byte("Lorem")))
	exp := f(p, messages.NewMessage(-1, -1, []byte("Lorem")))

	if exp != res {
		t.Errorf("Results mismatch: expected %v, got %v", exp, res)
	}
}

func TestTick(t *testing.T) {
	nl := network.New()
	defer nl.Stop()
	p := New(0)
	defer p.Stop()
	nl.RegisterProcess(0, p)

	f := func(context *Process, m *messages.Message) int32 {
		return 1
	}

	msg := &messages.Message{
		DeliveryTime: -1,
	}

	p.MessagesQueue.Enqueue(msg)

	p.RegisterWorkFunction([]byte("*TIME"), f)
	time.Sleep(500 * time.Millisecond)
}
