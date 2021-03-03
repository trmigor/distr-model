package priorityq

import (
	"container/heap"
)

// An Item is something we manage in a priority queue.
type Item struct {
	value    interface{} // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

// Len returns the PriorityQueue length
func (pq PriorityQueue) Len() int {
	return len(pq)
}

// Less returns is the j-th element priority of the pq PriorityQueue less than i-th element priority
func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use "greater than" here.
	return pq[i].priority > pq[j].priority
}

// Swap reverse j-th and i-th elements location in the pq PriorityQueue
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push adds items to the PriorityQueue
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes the highest-priority item from the PriorityQueue
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// Top gets the highest-priority item from the PriorityQueue
func (pq *PriorityQueue) Top() interface{} {
	queue := *pq
	return (queue)[len(queue)-1]
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value interface{}, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}
