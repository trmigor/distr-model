package set

import (
	"math/rand"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	s := New()
	value := rand.Int()
	s.Insert(value)
	if !s.Contains(value) {
		t.Errorf("Value not found: %v", value)
	}
}

func TestErase(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	s := New()
	value := rand.Int()
	s.Insert(value)
	s.Erase(value)
	if s.Contains(value) {
		t.Errorf("Found deleted value: %v", value)
	}
}

func TestSize(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	s := New()
	num := rand.Intn(100) + 100
	for i := 0; i < num; i++ {
		s.Insert(i)
	}
	if s.Size() != num {
		t.Errorf("Size do not match: expected %v, got %v", num, s.Size())
	}
	s.Insert(0)
	if s.Size() != num {
		t.Errorf("Size after reinserting do not match: expected %v, got %v", num, s.Size())
	}
}

func TestClear(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	s := New()
	num := rand.Intn(100) + 100
	for i := 0; i < num; i++ {
		s.Insert(i)
	}
	s.Clear()
	if !s.Empty() {
		t.Errorf("Set is not empty")
	}
}

func TestUnion(t *testing.T) {
	s1 := New()
	s2 := New()
	s1.Insert(0)
	s1.Insert(1)
	s2.Insert(1)
	s2.Insert(2)
	s := Union(s1, s2)
	for i := 0; i < 3; i++ {
		if !s.Contains(i) {
			t.Errorf("Value not found: %v", i)
		}
	}
}

func TestIntersection(t *testing.T) {
	s1 := New()
	s2 := New()
	s1.Insert(0)
	s1.Insert(1)
	s2.Insert(1)
	s2.Insert(2)
	s := Intersection(s1, s2)
	if !s.Contains(1) {
		t.Errorf("Value not found: %v", 1)
	}
	if s.Contains(0) || s.Contains(2) {
		t.Errorf("Found extra value")
	}
}

func TestDifference(t *testing.T) {
	s1 := New()
	s2 := New()
	s1.Insert(0)
	s1.Insert(1)
	s2.Insert(1)
	s2.Insert(2)
	s := Difference(s1, s2)
	if !s.Contains(0) {
		t.Errorf("Value not found: %v", 1)
	}
	if s.Contains(1) || s.Contains(2) {
		t.Errorf("Found extra value")
	}
}
