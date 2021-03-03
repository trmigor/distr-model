package messages

import (
	"bytes"
	"testing"
)

func TestNewMessageQueue(t *testing.T) {
	mq := NewMessageQueue()
	if len(mq.queue) != 0 {
		t.Errorf("Expected empty message queue")
	}
}

func TestPeek(t *testing.T) {
	mq := NewMessageQueue()

	msg1 := NewMessage(-1, -1, []byte("Lorem ipsum"))
	msg1.DeliveryTime = 1
	mq.Enqueue(msg1)

	msg2 := NewMessage(-1, -1, []byte("dolor sit amet"))
	msg2.DeliveryTime = 2
	mq.Enqueue(msg2)

	if !bytes.Equal(mq.Peek().Body, msg1.Body) {
		t.Errorf("Result mismatch: expected %v, got %v", msg1.Body, mq.Peek().Body)
	}
}

func TestEnqueue(t *testing.T) {
	mq := NewMessageQueue()
	msg := NewMessage(-1, -1, []byte("Lorem ipsum"))
	mq.Enqueue(msg)
	if !bytes.Equal(mq.Peek().Body, msg.Body) {
		t.Errorf("Result mismatch: expected %v, got %v", msg.Body, mq.Peek().Body)
	}
}

func TestDequeue(t *testing.T) {
	mq := NewMessageQueue()
	msg := NewMessage(-1, -1, []byte("Lorem ipsum"))
	mq.Enqueue(msg)
	res := mq.Dequeue()
	if !bytes.Equal(res.Body, msg.Body) {
		t.Errorf("Result mismatch: expected %v, got %v", msg.Body, res.Body)
	}
}

func TestSize(t *testing.T) {
	mq := NewMessageQueue()
	msg := NewMessage(-1, -1, []byte("Lorem ipsum"))
	mq.Enqueue(msg)
	size := mq.Size()
	if size != 1 {
		t.Errorf("Result mismatch: expected %v, got %v", 1, size)
	}
}
