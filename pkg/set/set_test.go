package set

import (
	"reflect"
	"testing"
)

func TestSet_Insert(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		s    *Set
		args args
	}{
		{"Valid", New(), args{int32(123)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Insert(tt.args.value)
			if !tt.s.Contains(tt.args.value) {
				t.Errorf("Set.Insert(): not inserted")
			}
		})
	}
}

func TestSet_Contains(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		s    *Set
		in   []interface{}
		args args
		want bool
	}{
		{"Contains", New(), []interface{}{int32(0), int32(1), int32(2)}, args{int32(0)}, true},
		{"Not", New(), []interface{}{int32(0), int32(1), int32(2)}, args{int32(3)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, e := range tt.in {
				tt.s.Insert(e)
			}
			if got := tt.s.Contains(tt.args.value); got != tt.want {
				t.Errorf("Set.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Erase(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		s    *Set
		in   []interface{}
		args args
	}{
		{"Contains", New(), []interface{}{int32(0), int32(1), int32(2)}, args{int32(0)}},
		{"Not", New(), []interface{}{int32(0), int32(1), int32(2)}, args{int32(3)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, e := range tt.in {
				tt.s.Insert(e)
			}
			tt.s.Erase(tt.args.value)
			if tt.s.Contains(tt.args.value) {
				t.Errorf("Set.Erase(): not erased")
			}
		})
	}
}

func TestSet_Size(t *testing.T) {
	tests := []struct {
		name string
		s    *Set
		in   []interface{}
		want int
	}{
		{"Zero", New(), []interface{}{}, 0},
		{"One", New(), []interface{}{int32(0)}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, e := range tt.in {
				tt.s.Insert(e)
			}
			if got := tt.s.Size(); got != tt.want {
				t.Errorf("Set.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Empty(t *testing.T) {
	tests := []struct {
		name string
		s    *Set
		in   []interface{}
		want bool
	}{
		{"Zero", New(), []interface{}{}, true},
		{"One", New(), []interface{}{int32(0)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, e := range tt.in {
				tt.s.Insert(e)
			}
			if got := tt.s.Empty(); got != tt.want {
				t.Errorf("Set.Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Clear(t *testing.T) {
	tests := []struct {
		name string
		s    *Set
		in   []interface{}
	}{
		{"Zero", New(), []interface{}{}},
		{"One", New(), []interface{}{int32(0)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Clear()
			if !tt.s.Empty() {
				t.Errorf("Set.Clear(): not cleared")
			}
		})
	}
}

func TestUnion(t *testing.T) {
	type args struct {
		s1  *Set
		in1 []interface{}
		s2  *Set
		in2 []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Set
		in   []interface{}
	}{
		{
			"Valid",
			args{
				New(),
				[]interface{}{int32(0), int32(1), int32(2)},
				New(),
				[]interface{}{int32(1), int32(2), int32(3)},
			},
			New(),
			[]interface{}{int32(0), int32(1), int32(2), int32(3)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, e := range tt.args.in1 {
				tt.args.s1.Insert(e)
			}
			for _, e := range tt.args.in2 {
				tt.args.s2.Insert(e)
			}
			for _, e := range tt.in {
				tt.want.Insert(e)
			}
			if got := Union(tt.args.s1, tt.args.s2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	type args struct {
		s1  *Set
		in1 []interface{}
		s2  *Set
		in2 []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Set
		in   []interface{}
	}{
		{
			"Valid",
			args{
				New(),
				[]interface{}{int32(0), int32(1), int32(2)},
				New(),
				[]interface{}{int32(1), int32(2), int32(3)},
			},
			New(),
			[]interface{}{int32(1), int32(2)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, e := range tt.args.in1 {
				tt.args.s1.Insert(e)
			}
			for _, e := range tt.args.in2 {
				tt.args.s2.Insert(e)
			}
			for _, e := range tt.in {
				tt.want.Insert(e)
			}
			if got := Intersection(tt.args.s1, tt.args.s2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Intersection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDifference(t *testing.T) {
	type args struct {
		s1  *Set
		in1 []interface{}
		s2  *Set
		in2 []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Set
		in   []interface{}
	}{
		{
			"Valid",
			args{
				New(),
				[]interface{}{int32(0), int32(1), int32(2)},
				New(),
				[]interface{}{int32(1), int32(2), int32(3)},
			},
			New(),
			[]interface{}{int32(0)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, e := range tt.args.in1 {
				tt.args.s1.Insert(e)
			}
			for _, e := range tt.args.in2 {
				tt.args.s2.Insert(e)
			}
			for _, e := range tt.in {
				tt.want.Insert(e)
			}
			if got := Difference(tt.args.s1, tt.args.s2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Difference() = %v, want %v", got, tt.want)
			}
		})
	}
}
