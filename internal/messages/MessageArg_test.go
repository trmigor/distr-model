package messages

import (
	"reflect"
	"testing"
)

func TestNewMessageArg(t *testing.T) {
	type args struct {
		q interface{}
	}
	tests := []struct {
		name string
		args args
		want *MessageArg
	}{
		{"Int32", args{int32(1)}, &MessageArg{[]byte{65, 1, 0, 0, 0}}},
		{"Int64", args{int64(1)}, &MessageArg{[]byte{66, 1, 0, 0, 0, 0, 0, 0, 0}}},
		{"String", args{[]byte("Lorem")}, &MessageArg{append([]byte("CLorem"), 0)}},
		{"Panic", args{complex(0, 1)}, &MessageArg{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Panic" {
				defer func() {
					if recover() == nil {
						t.Errorf("No panic occured")
					}
				}()
			}
			if got := NewMessageArg(tt.args.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessageArg() = %v, want %v", got, tt.want)
			}
		})
	}
}
