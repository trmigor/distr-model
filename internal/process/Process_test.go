package process

import (
	"reflect"
	"testing"
	"time"

	"github.com/trmigor/distr-model/internal/messages"
	"github.com/trmigor/distr-model/internal/network"
	"github.com/trmigor/distr-model/pkg/set"
)

func TestProcess_Neibs(t *testing.T) {
	tests := []struct {
		name string
		want *set.Set
	}{
		{"Valid", set.New()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(0)
			defer p.Stop()
			p1 := New(1)
			defer p1.Stop()
			p2 := New(2)
			defer p2.Stop()

			nl := network.New()
			nl.RegisterProcess(p.Node, p)
			nl.RegisterProcess(p1.Node, p)
			nl.RegisterProcess(p2.Node, p)
			nl.CreateLink(0, 1, false, 0)

			tt.want.Insert(int32(1))

			if got := p.Neibs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Process.Neibs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcess_RegisterWorkFunction(t *testing.T) {
	type args struct {
		prefix []byte
		wf     WorkFunction
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid", args{[]byte("SETX"), func(context *Process, m *messages.Message) bool { return true }}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(0)
			defer p.Stop()
			p.RegisterWorkFunction(tt.args.prefix, tt.args.wf)
			if len(p.workers) != 1 {
				t.Errorf("Process.RegisterWorkFunction(): worker not inserted")
			}
		})
	}
}

func TestProcess_IsMyMessage(t *testing.T) {
	type args struct {
		prefix  []byte
		message []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Star", args{[]byte(""), []byte("*")}, true},
		{"LongPrefix", args{[]byte("Lorem ipsum"), []byte("")}, false},
		{"NotMatching", args{[]byte("Lorem"), []byte("Larem_ipsum")}, false},
		{"Matching", args{[]byte("Lorem"), []byte("Lorem_ipsum")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(0)
			defer p.Stop()

			if got := p.IsMyMessage(tt.args.prefix, tt.args.message); got != tt.want {
				t.Errorf("Process.IsMyMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_workerThreadExecutor(t *testing.T) {
	t.Run("Tick", func(t *testing.T) {
		nl := network.New()
		defer nl.Stop()
		p := New(0)
		defer p.Stop()
		nl.RegisterProcess(0, p)

		f := func(context *Process, m *messages.Message) bool {
			return true
		}

		msg := &messages.Message{
			DeliveryTime: -1,
		}

		p.MessagesQueue.Enqueue(msg)

		p.RegisterWorkFunction([]byte("*TIME"), f)
		time.Sleep(500 * time.Millisecond)
	})
}
