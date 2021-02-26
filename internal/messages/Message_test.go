package messages

import (
	"reflect"
	"testing"
)

func TestNewMessage(t *testing.T) {
	type args struct {
		from int32
		to   int32
		body []byte
	}
	tests := []struct {
		name string
		args args
		want *Message
	}{
		{"Valid", args{0, 1, []byte("Lorem")}, &Message{From: 0, To: 1, Body: []byte("Lorem")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMessage(tt.args.from, tt.args.to, tt.args.body); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_append(t *testing.T) {
	type args struct {
		a *MessageArg
	}
	tests := []struct {
		name string
		args args
		want *Message
	}{
		{"Valid", args{NewMessageArg(int32(1))}, &Message{From: 0, To: 1, Body: []byte{65, 1, 0, 0, 0}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := NewMessage(0, 1, []byte(""))
			msg.append(tt.args.a)
			if !reflect.DeepEqual(msg, tt.want) {
				t.Errorf("NewMessage() = %v, want %v", msg, tt.want)
			}
		})
	}
}

func TestNewMessageByArgs(t *testing.T) {
	type args struct {
		args []*MessageArg
	}
	tests := []struct {
		name string
		args args
		want *Message
	}{
		{"One", args{[]*MessageArg{NewMessageArg(int32(1))}}, &Message{From: -1, To: -1, Body: []byte{65, 1, 0, 0, 0}}},
		{"Multiple", args{[]*MessageArg{NewMessageArg(int32(1)), NewMessageArg(int32(2))}}, &Message{From: -1, To: -1, Body: []byte{65, 1, 0, 0, 0, 65, 2, 0, 0, 0}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMessageByArgs(tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessageByArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_GetInt32(t *testing.T) {
	type fields struct {
		Body []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   int32
	}{
		{"Valid", fields{[]byte{65, 1, 0, 0, 0}}, 1},
		{"Panic", fields{[]byte{}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &Message{
				Body: tt.fields.Body,
			}
			if tt.name == "Panic" {
				defer func() {
					if recover() == nil {
						t.Errorf("No panic occured")
					}
				}()
			}
			if got := msg.GetInt32(); got != tt.want {
				t.Errorf("Message.GetInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_GetInt64(t *testing.T) {
	type fields struct {
		Body []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{"Valid", fields{[]byte{66, 1, 0, 0, 0, 0, 0, 0, 0}}, 1},
		{"Panic", fields{[]byte{}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &Message{
				Body: tt.fields.Body,
			}
			if tt.name == "Panic" {
				defer func() {
					if recover() == nil {
						t.Errorf("No panic occured")
					}
				}()
			}
			if got := msg.GetInt64(); got != tt.want {
				t.Errorf("Message.GetInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_GetString(t *testing.T) {
	type fields struct {
		Body []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{"Valid", fields{append([]byte("CLorem"), 0)}, []byte("Lorem")},
		{"Panic", fields{[]byte{}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &Message{
				Body: tt.fields.Body,
			}
			if tt.name == "Panic" {
				defer func() {
					if recover() == nil {
						t.Errorf("No panic occured")
					}
				}()
			}
			if got := msg.GetString(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Message.GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_GetData(t *testing.T) {
	type fields struct {
		Body []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{"Int32", fields{[]byte{65, 1, 0, 0, 0}}, int32(1)},
		{"Int64", fields{[]byte{66, 1, 0, 0, 0, 0, 0, 0, 0}}, int64(1)},
		{"String", fields{append([]byte("CLorem"), 0)}, []byte("Lorem")},
		{"Nil", fields{[]byte{}}, nil},
		{"Panic", fields{[]byte{64}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &Message{
				Body: tt.fields.Body,
			}
			if tt.name == "Panic" {
				defer func() {
					if recover() == nil {
						t.Errorf("No panic occured")
					}
				}()
			}
			if got := msg.GetData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Message.GetData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGreater(t *testing.T) {
	type args struct {
		first  *Message
		second *Message
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Greater", args{&Message{DeliveryTime: 1}, &Message{DeliveryTime: 0}}, true},
		{"Equal", args{&Message{DeliveryTime: 1}, &Message{DeliveryTime: 1}}, false},
		{"Less", args{&Message{DeliveryTime: 0}, &Message{DeliveryTime: 1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Greater(tt.args.first, tt.args.second); got != tt.want {
				t.Errorf("Greater() = %v, want %v", got, tt.want)
			}
		})
	}
}
