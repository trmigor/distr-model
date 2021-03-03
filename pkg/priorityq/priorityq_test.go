package priorityq

import (
	"container/heap"
	"reflect"
	"testing"
)

func TestPriorityQueue_Len(t *testing.T) {
	tests := []struct {
		name string
		pq   PriorityQueue
		want int
	}{
		{"Zero", PriorityQueue{}, 0},
		{"One", PriorityQueue{&Item{}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pq.Len(); got != tt.want {
				t.Errorf("PriorityQueue.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriorityQueue_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		pq   PriorityQueue
		args args
		want bool
	}{
		{"Less", PriorityQueue{&Item{Priority: 1}, &Item{Priority: 0}}, args{0, 1}, true},
		{"Equal", PriorityQueue{&Item{Priority: 0}, &Item{Priority: 0}}, args{0, 1}, false},
		{"LeGreaterss", PriorityQueue{&Item{Priority: 0}, &Item{Priority: 1}}, args{0, 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pq.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("PriorityQueue.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriorityQueue_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		pq   PriorityQueue
		args args
	}{
		{"Valid", PriorityQueue{&Item{Priority: 1}, &Item{Priority: 0}}, args{0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prev := tt.pq.Less(tt.args.i, tt.args.j)
			tt.pq.Swap(tt.args.i, tt.args.j)
			curr := tt.pq.Less(tt.args.i, tt.args.j)
			if prev == curr {
				t.Errorf("PriorityQueue.Swap(): swap was not performed")
			}
		})
	}
}

func TestPriorityQueue_Push(t *testing.T) {
	type args struct {
		x interface{}
	}
	tests := []struct {
		name string
		pq   *PriorityQueue
		args args
		want interface{}
	}{
		{"Valid", &PriorityQueue{}, args{&Item{Value: 123}}, &Item{Value: 123}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pq.Push(tt.args.x)
			if !reflect.DeepEqual((*tt.pq)[0], tt.want.(*Item)) {
				t.Errorf("PriorityQueue.Push(): push was not performed")
			}
		})
	}
}

func TestPriorityQueue_Pop(t *testing.T) {
	tests := []struct {
		name string
		pq   *PriorityQueue
		want interface{}
	}{
		{"Valid", &PriorityQueue{}, &Item{Value: 123, Index: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pq.Push(&Item{Value: 123})
			if got := tt.pq.Pop(); !reflect.DeepEqual(got.(*Item), tt.want.(*Item)) {
				t.Errorf("PriorityQueue.Pop() = %v, want %v", got.(*Item), tt.want.(*Item))
			}
		})
	}
}

func TestPriorityQueue_Top(t *testing.T) {
	tests := []struct {
		name string
		pq   *PriorityQueue
		want interface{}
	}{
		{"Valid", &PriorityQueue{}, &Item{Value: 123}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pq.Push(&Item{Value: 123})
			if got := tt.pq.Top(); !reflect.DeepEqual(got.(*Item), tt.want.(*Item)) {
				t.Errorf("PriorityQueue.Top() = %v, want %v", got.(*Item), tt.want.(*Item))
			}
		})
	}
}

func TestPriorityQueue_update(t *testing.T) {
	type args struct {
		value    interface{}
		priority int
	}
	tests := []struct {
		name string
		pq   *PriorityQueue
		args args
		want *Item
	}{
		{"Valid", &PriorityQueue{}, args{int32(123), 456}, &Item{Value: int32(123), Priority: 456}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			heap.Init(tt.pq)
			tt.pq.Push(&Item{Value: 0, Priority: 0})
			tt.pq.update(tt.pq.Top().(*Item), tt.args.value, tt.args.priority)
			if got := tt.pq.Top(); !reflect.DeepEqual(got.(*Item), tt.want) {
				t.Errorf("PriorityQueue.Top() = %v, want %v", got.(*Item), tt.want)
			}
		})
	}
}
