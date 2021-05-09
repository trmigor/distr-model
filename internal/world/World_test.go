package world

import (
	"reflect"
	"testing"

	"github.com/trmigor/distr-model/internal/errors"
	"github.com/trmigor/distr-model/internal/messages"
	"github.com/trmigor/distr-model/internal/process"
)

func TestWorld_CreateProcess(t *testing.T) {
	type args struct {
		node int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{"Valid", args{0}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			defer w.Stop()
			if got := w.CreateProcess(tt.args.node); got != tt.want {
				t.Errorf("World.CreateProcess() = %v, want %v", got, tt.want)
			}
			if len(w.ProcessesList) != 1 || w.ProcessesList[0].Node != 0 {
				t.Errorf("World.CreateProcess(): process created badly")
			}
		})
	}
}

func TestWorld_RegisterWorkFunction(t *testing.T) {
	type args struct {
		function []byte
		wf       process.WorkFunction
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid", args{
			[]byte("SETX"),
			func(context *process.Process, m *messages.Message) bool {
				return true
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			defer w.Stop()
			w.RegisterWorkFunction(tt.args.function, tt.args.wf)
			if _, ok := w.Associates["SETX"]; !ok {
				t.Errorf("World.RegisterWorkFunction(): function is not registered")
			}
		})
	}
}

func TestWorld_AssignWorkFunction(t *testing.T) {
	type fields struct {
		Node int32
	}
	type args struct {
		node     int32
		function []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   errors.ErrorCode
	}{
		{"Valid", fields{0}, args{0, []byte("SETX")}, errors.OK},
		{"InvalidNode", fields{0}, args{-1, []byte("SETX")}, errors.ItemNotFound},
		{"nilProcess", fields{1}, args{0, []byte("SETX")}, errors.ItemNotFound},
		{"InvalidFunc", fields{0}, args{0, []byte("SETY")}, errors.ItemNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			defer w.Stop()
			f := func(context *process.Process, m *messages.Message) bool {
				return true
			}
			w.RegisterWorkFunction([]byte("SETX"), f)
			w.CreateProcess(tt.fields.Node)
			if got := w.AssignWorkFunction(tt.args.node, tt.args.function); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("World.AssignWorkFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorld_ParseConfig(t *testing.T) {
	type args struct {
		name []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"InvalidFile", args{[]byte("../../test/data/config/InvalidFile.data")}, false},
		{"Comment", args{[]byte("../../test/data/config/Comment.data")}, true},
		{"Processes", args{[]byte("../../test/data/config/Processes.data")}, true},
		{"Bidirected", args{[]byte("../../test/data/config/Bidirected.data")}, true},
		{"ErrorRate", args{[]byte("../../test/data/config/ErrorRate.data")}, true},
		{"AllToAll", args{[]byte("../../test/data/config/AllToAll.data")}, true},
		{"AllToAllLatency", args{[]byte("../../test/data/config/AllToAllLatency.data")}, true},
		{"SetProcesses", args{[]byte("../../test/data/config/SetProcesses.data")}, true},
		{"SetProcessesInvalid", args{[]byte("../../test/data/config/SetProcessesInvalid.data")}, false},
		{"SendMsgArg", args{[]byte("../../test/data/config/SendMsgArg.data")}, true},
		{"SendMsg", args{[]byte("../../test/data/config/SendMsg.data")}, true},
		{"Wait", args{[]byte("../../test/data/config/Wait.data")}, true},
		{"LaunchTimer", args{[]byte("../../test/data/config/LaunchTimer.data")}, true},
		{"LinkLatency", args{[]byte("../../test/data/config/LinkLatency.data")}, true},
		{"Link", args{[]byte("../../test/data/config/Link.data")}, true},
		{"LinkToAllLatency", args{[]byte("../../test/data/config/LinkToAllLatency.data")}, true},
		{"LinkToAll", args{[]byte("../../test/data/config/LinkToAll.data")}, true},
		{"LinkFromAllLatency", args{[]byte("../../test/data/config/LinkFromAllLatency.data")}, true},
		{"LinkFromAll", args{[]byte("../../test/data/config/LinkFromAll.data")}, true},
		{"Unknown", args{[]byte("../../test/data/config/Unknown.data")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			defer w.Stop()
			if tt.name == "Unknown" {
				defer func() {
					if recover() == nil {
						t.Errorf("No panic occured")
					}
				}()
			}
			f := func(context *process.Process, m *messages.Message) bool {
				return true
			}
			w.RegisterWorkFunction([]byte("SETX"), f)
			if got := w.ParseConfig(tt.args.name); got != tt.want {
				t.Errorf("World.ParseConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
