package world

import (
	"testing"

	"github.com/trmigor/distr-model/internal/errors"
	"github.com/trmigor/distr-model/internal/messages"
	"github.com/trmigor/distr-model/internal/process"
)

func TestCreateProcess(t *testing.T) {
	w := New()
	defer w.Stop()

	if n := w.CreateProcess(0); n != 0 {
		t.Errorf("Created process for the wrong node: expected 0, got %v", n)
	}
	if w.ProcessesList[0] == nil {
		t.Errorf("Invalid process reference")
	}
}

func TestAssignWorkFunction(t *testing.T) {
	w := New()
	defer w.Stop()

	if n := w.CreateProcess(0); n != 0 {
		t.Errorf("Created process for the wrong node: expected 0, got %v", n)
	}
	w.RegisterWorkFunction([]byte("Lorem"), func(context *process.Process, m *messages.Message) int32 { return 0 })
	if res := w.AssignWorkFunction(0, []byte("Lorem")); res != errors.OK {
		t.Errorf("Results mismatch: expected %v, got %v", errors.OK, res)
	}
}

func TestAssignWorkFunctionTooBigNode(t *testing.T) {
	w := New()
	defer w.Stop()

	if n := w.CreateProcess(0); n != 0 {
		t.Errorf("Created process for the wrong node: expected 0, got %v", n)
	}
	w.RegisterWorkFunction([]byte("Lorem"), func(context *process.Process, m *messages.Message) int32 { return 0 })
	if res := w.AssignWorkFunction(1, []byte("Lorem")); res != errors.ItemNotFound {
		t.Errorf("Results mismatch: expected %v, got %v", errors.ItemNotFound, res)
	}
}

func TestAssignWorkFunctionTooSmallNode(t *testing.T) {
	w := New()
	defer w.Stop()

	if n := w.CreateProcess(0); n != 0 {
		t.Errorf("Created process for the wrong node: expected 0, got %v", n)
	}
	w.RegisterWorkFunction([]byte("Lorem"), func(context *process.Process, m *messages.Message) int32 { return 0 })
	if res := w.AssignWorkFunction(-1, []byte("Lorem")); res != errors.ItemNotFound {
		t.Errorf("Results mismatch: expected %v, got %v", errors.ItemNotFound, res)
	}
}

func TestAssignWorkFunctionNilProcessReference(t *testing.T) {
	w := New()
	defer w.Stop()

	if n := w.CreateProcess(2); n != 2 {
		t.Errorf("Created process for the wrong node: expected 0, got %v", n)
	}
	w.RegisterWorkFunction([]byte("Lorem"), func(context *process.Process, m *messages.Message) int32 { return 0 })
	if res := w.AssignWorkFunction(0, []byte("Lorem")); res != errors.ItemNotFound {
		t.Errorf("Results mismatch: expected %v, got %v", errors.ItemNotFound, res)
	}
}

func TestAssignWorkFunctionWithoutWorkFunction(t *testing.T) {
	w := New()
	defer w.Stop()

	if n := w.CreateProcess(0); n != 0 {
		t.Errorf("Created process for the wrong node: expected 0, got %v", n)
	}
	if res := w.AssignWorkFunction(0, []byte("Lorem")); res != errors.ItemNotFound {
		t.Errorf("Results mismatch: expected %v, got %v", errors.ItemNotFound, res)
	}
}

func TestRegisterWorkFunction(t *testing.T) {
	w := New()
	defer w.Stop()

	w.RegisterWorkFunction([]byte("Lorem"), func(context *process.Process, m *messages.Message) int32 { return 0 })
	if res := w.Associates["Lorem"](process.New(0), messages.NewMessage(0, 1, []byte("Lorem"))); res != 0 {
		t.Errorf("Results mismatch: expected 0, got %v", res)
	}
}

func TestParseConfig(t *testing.T) {
	w := New()
	defer w.Stop()

	if w.ParseConfig([]byte("../../configs/config.data")) != true {
		t.Errorf("Results mismatch: expected true, got false")
	}
}

func TestParseConfigWrongName(t *testing.T) {
	w := New()
	defer w.Stop()

	if w.ParseConfig([]byte("../../configs/wrongConfig.data")) != false {
		t.Errorf("Results mismatch: expected false, got true")
	}
}
