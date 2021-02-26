package messages

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

func TestNewMessage(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	from := int32(rand.Int())
	to := int32(rand.Int())
	body := []byte("Lorem ipsum dolor sit amet")
	msg := NewMessage(from, to, body)

	if !bytes.Equal(msg.Body, body) {
		t.Errorf("Body mismatch: expected %v, got %v", body, msg.Body)
	}

	if msg.From != from {
		t.Errorf("Source mismatch: expected %v, got %v", from, msg.From)
	}

	if msg.To != to {
		t.Errorf("Destination mismatch: expected %v, got %v", to, msg.To)
	}
}

func TestAppend(t *testing.T) {
	body := []byte("Lorem ")
	mArg := NewMessageArg([]byte("ipsum"))
	msg := NewMessage(-1, -1, body)
	msg.append(mArg)
	if !bytes.Equal(msg.Body, append(body, mArg.Body...)) {
		t.Errorf("Result mismatch: expected %v, got %v", append(body, mArg.Body...), msg.Body)
	}
}

func TestNewMessageByArgs(t *testing.T) {
	mArg1 := NewMessageArg([]byte("Lorem "))
	mArg2 := NewMessageArg([]byte("ipsum"))
	msg := NewMessageByArgs(mArg1, mArg2)
	if !bytes.Equal(msg.Body, append(mArg1.Body, mArg2.Body...)) {
		t.Errorf("Result mismatch: expected %v, got %v", append(mArg1.Body, mArg2.Body...), msg.Body)
	}
}

func TestGetInt32(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	first := int32(rand.Int())
	second := int32(rand.Int())
	mArg1 := NewMessageArg(first)
	mArg2 := NewMessageArg(second)

	msg := NewMessageByArgs(mArg1, mArg2)
	got := msg.GetInt32()
	if got != first {
		t.Errorf("Result mismatch: expected %v, got %v", first, got)
	}
}

func TestGetInt32Panic(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	first := int64(rand.Int())
	mArg1 := NewMessageArg(first)

	msg := NewMessageByArgs(mArg1)

	defer func() {
		if recover() == nil {
			t.Errorf("No panic occured")
		}
	}()
	_ = msg.GetInt32()
}

func TestGetInt64(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	first := int64(rand.Uint64())
	second := int64(rand.Uint64())
	mArg1 := NewMessageArg(first)
	mArg2 := NewMessageArg(second)

	msg := NewMessageByArgs(mArg1, mArg2)
	got := msg.GetInt64()
	if got != first {
		t.Errorf("Result mismatch: expected %v, got %v", first, got)
	}
}

func TestGetInt64Panic(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	first := int32(rand.Int())
	mArg1 := NewMessageArg(first)

	msg := NewMessageByArgs(mArg1)

	defer func() {
		if recover() == nil {
			t.Errorf("No panic occured")
		}
	}()
	_ = msg.GetInt64()
}

func TestGetString(t *testing.T) {
	first := []byte("Lorem ipsum")
	second := []byte("dolor sit amet")
	mArg1 := NewMessageArg(first)
	mArg2 := NewMessageArg(second)

	msg := NewMessageByArgs(mArg1, mArg2)
	got := msg.GetString()
	if !bytes.Equal(got, first) {
		t.Errorf("Result mismatch: expected %v, got %v", first, got)
	}
}

func TestGetStringPanic(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	first := int32(rand.Int())
	mArg1 := NewMessageArg(first)

	msg := NewMessageByArgs(mArg1)

	defer func() {
		if recover() == nil {
			t.Errorf("No panic occured")
		}
	}()
	_ = msg.GetString()
}

func TestGetData(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	first := int32(rand.Int())
	second := int64(rand.Uint64())
	third := []byte("Lorem ipsum")

	mArg1 := NewMessageArg(first)
	mArg2 := NewMessageArg(second)
	mArg3 := NewMessageArg(third)

	msg := NewMessageByArgs(mArg1, mArg2, mArg3)
	res1 := msg.GetData()
	res2 := msg.GetData()
	res3 := msg.GetData()

	if res1.(int32) != first {
		t.Errorf("Result mismatch: expected %v, got %v", first, res1)
	}
	if res2.(int64) != second {
		t.Errorf("Result mismatch: expected %v, got %v", second, res2)
	}
	if !bytes.Equal(res3.([]byte), third) {
		t.Errorf("Result mismatch: expected %v, got %v", third, res3)
	}
}

func TestGetDataPanic(t *testing.T) {
	msg := NewMessage(-1, -1, []byte("Lorem"))

	defer func() {
		if recover() == nil {
			t.Errorf("No panic occured")
		}
	}()
	_ = msg.GetData()
}

func TestGetDataOutOfBound(t *testing.T) {
	msg := NewMessage(-1, -1, []byte(""))

	res := msg.GetData()
	if res != nil {
		t.Errorf("Result mismatch: expected %v, got %v", nil, res)
	}
}

func TestGreater(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var msg1, msg2 Message
	msg1.DeliveryTime = int64(rand.Intn(5))
	msg2.DeliveryTime = int64(rand.Intn(5) + 5)

	if Greater(&msg1, &msg2) {
		t.Errorf("Result mismatch: expected %v > %v == false, got true", msg1.DeliveryTime, msg2.DeliveryTime)
	}
}
