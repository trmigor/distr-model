package model

import (
	"container/heap"
	"sync"

	"github.com/trmigor/distr-model/priorityq"

// A MessageQueue is a PriorityQueue structure with mutex for parallel usage.
type MessageQueue struct {
	queue priorityq.PriorityQueue
	mutex sync.Mutex
}

// NewMessageQueue establishes the heap invariants.
func NewMessageQueue() (mq *MessageQueue){
	heap.Init(&mq.queue)
}

// Dequeue removes and returns the object of the priority queue with the minimum priority.
func (mq *MessageQueue) Dequeue() (ret Message) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	ret = mq.queue.Top().(Message)
	heap.Pop(&mq.queue)

	return ret
}

// Enqueue adds an object to the priority queue.
func (mq *MessageQueue) Enqueue(msg *Message) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	heap.Push(&mq.queue, msg)
}

// Peek returns the object of the priority queue with the minimum priority without removing it.
func (mq *MessageQueue) Peek() Message {
	return mq.queue.Top().(Message)
}

// Size gets the number of elements contained in the priority queue.
func (mq *MessageQueue) Size() int {
	return mq.queue.Len()
}
