package priorityq

import (
	"container/heap"
	"testing"
)

func TestPriorityQ(t *testing.T) {
	// Some items and their priorities.
	items := map[interface{}]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			Value:    value,
			Priority: priority,
			Index:    i,
		}
		i++
	}
	heap.Init(&pq)

	// Insert a new item and then modify its priority.
	item := &Item{
		Value:    "orange",
		Priority: 1,
	}
	heap.Push(&pq, item)
	pq.update(item, item.Value, 5)

	// Check top element
	if pq.Len() > 0 {
		if pq.Top().(*Item).Value != "apple" {
			t.Error("Minimum priority element of the priority queue is wrong!")
		}
	}

	// Take the items out; they arrive in decreasing priority order.
	item = heap.Pop(&pq).(*Item)
	if item.Priority != 5 || item.Value != "orange" {
		t.Error("Wrong element was poped!")
	}
	item = heap.Pop(&pq).(*Item)
	if item.Priority != 4 || item.Value != "pear" {
		t.Error("Wrong element was poped!")
	}
	item = heap.Pop(&pq).(*Item)
	if item.Priority != 3 || item.Value != "banana" {
		t.Error("Wrong element was poped!")
	}
	item = heap.Pop(&pq).(*Item)
	if item.Priority != 2 || item.Value != "apple" {
		t.Error("Wrong element was poped!")
	}

}
