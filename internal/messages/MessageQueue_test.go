package messages

import (
	"reflect"
	"testing"
)

func TestMessageQueue_Peek(t *testing.T) {
	tests := []struct {
		name  string
		queue []*Message
		want  Message
	}{
		{"One", []*Message{{DeliveryTime: 123}}, Message{DeliveryTime: 123}},
		{"Multiple", []*Message{{DeliveryTime: 456}, {DeliveryTime: 123}}, Message{DeliveryTime: 123}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mq := NewMessageQueue()
			for _, e := range tt.queue {
				mq.Enqueue(e)
			}
			if got := mq.Peek(); !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("MessageQueue.Peek() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func TestMessageQueue_Enqueue(t *testing.T) {
	type args struct {
		msg *Message
	}
	tests := []struct {
		name string
		args args
		want *Message
	}{
		{"Valid", args{NewMessageByArgs(NewMessageArg(int32(123)))}, NewMessageByArgs(NewMessageArg(int32(123)))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mq := NewMessageQueue()
			mq.Enqueue(tt.args.msg)
			if got := mq.Peek(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessageQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageQueue_Dequeue(t *testing.T) {
	tests := []struct {
		name  string
		queue []*Message
		want  Message
	}{
		{"One", []*Message{{DeliveryTime: 123}}, Message{DeliveryTime: 123}},
		{"Multiple", []*Message{{DeliveryTime: 456}, {DeliveryTime: 123}}, Message{DeliveryTime: 123}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mq := NewMessageQueue()
			for _, e := range tt.queue {
				mq.Enqueue(e)
			}
			if got := mq.Dequeue(); !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("MessageQueue.Dequeue() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func TestMessageQueue_Size(t *testing.T) {
	tests := []struct {
		name  string
		queue []*Message
		want  int
	}{
		{"Zero", []*Message{}, 0},
		{"One", []*Message{{DeliveryTime: 123}}, 1},
		{"Two", []*Message{{DeliveryTime: 456}, {DeliveryTime: 123}}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mq := NewMessageQueue()
			for _, e := range tt.queue {
				mq.Enqueue(e)
			}
			if got := mq.Size(); got != tt.want {
				t.Errorf("MessageQueue.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
