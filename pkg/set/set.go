package set

type void struct{}

var member void

// Set is a container that contains a set of unique objects
type Set map[interface{}]void

// New creates an empty set
func New() *Set {
	s := make(Set)
	return &s
}

// Insert adds a new value to set if it is not yet added
func (s *Set) Insert(value interface{}) {
	(*s)[value] = member
}

// Contains check if set contains the requested value
func (s *Set) Contains(value interface{}) bool {
	_, contains := (*s)[value]
	return contains
}

// Erase deletes value from set if it is contained
func (s *Set) Erase(value interface{}) {
	if s.Contains(value) {
		delete(*s, value)
	}
}

// Size returns number of set elements
func (s *Set) Size() int {
	return len(*s)
}

// Empty checks whether the set is empty
func (s *Set) Empty() bool {
	return s.Size() == 0
}

// Clear removes all elements from a set
func (s *Set) Clear() {
	*s = make(Set)
}

// Union returns a mathematical union of two sets
func Union(s1 *Set, s2 *Set) *Set {
	res := New()
	for v := range *s1 {
		res.Insert(v)
	}
	for v := range *s2 {
		res.Insert(v)
	}
	return res
}

// Intersection returns a mathematical intersection of two sets
func Intersection(s1 *Set, s2 *Set) *Set {
	res := New()
	for v := range *s1 {
		if s2.Contains(v) {
			res.Insert(v)
		}
	}
	return res
}

// Difference returns a mathematical difference of two sets
func Difference(s1 *Set, s2 *Set) *Set {
	res := New()
	for v := range *s1 {
		if !s2.Contains(v) {
			res.Insert(v)
		}
	}
	return res
}
