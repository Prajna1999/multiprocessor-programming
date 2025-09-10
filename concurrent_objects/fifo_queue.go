package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Lock-free single-enqueuer / single-dequeuer queue

type SingleProducerQueue struct {
	items []interface{}
	head  int64 //Index where next deq() will read
	tail  int64 // Index where next enq() will write
}

func NewSingleProducerQueue(capacity int) *SingleProducerQueue {
	return &SingleProducerQueue{
		items: make([]interface{}, capacity),
		head:  0,
		tail:  0,
	}
}

func (q *SingleProducerQueue) Enq(item interface{}) bool {
	currentTail := atomic.LoadInt64(&q.tail)
	currentHead := atomic.LoadInt64(&q.head)

	if currentTail-currentHead >= int64(len(q.items)) {
		return false // Queue full
	}
	q.items[currentTail%int64(len(q.items))] = item
	fmt.Printf("📥 Enqueuer: Added %v at position %d\n", item, currentTail%int64(len(q.items)))

	// Update the tail atomically (this is the only synchronization point)
	atomic.StoreInt64(&q.tail, currentTail+1)
	return true
}

func (q *SingleProducerQueue) Deq() interface{} {
	currentHead := atomic.LoadInt64(&q.head)
	currentTail := atomic.LoadInt64(&q.tail)

	if currentHead >= currentTail {
		return nil //Queue empty
	}
	item := q.items[currentHead%int64(len(q.items))]
	fmt.Printf("📤 Dequeuer: Removed %v from position %d\n", item, currentHead%int64(len(q.items)))

	// Update head atomically
	atomic.StoreInt64(&q.head, currentHead+1)
	return item

}

func main() {
	q := NewSingleProducerQueue(5)
	var wg sync.WaitGroup

	fmt.Println("🎬 Testing Single-Producer/Single-Consumer Queue\n")
	wg.Add(2)

	// single enqueuer thread
	go func() {
		defer wg.Done()
		for i := 'A'; i <= 'E'; i++ {
			q.Enq(string(i))
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Println("✅ Enqueuer finished")
	}()

	// Single dequeuer thread
	go func() {
		defer wg.Done()
		time.Sleep(250 * time.Millisecond)

		for i := 0; i < 5; i++ {
			item := q.Deq()
			if item != nil {
				fmt.Printf("🎯 Got: %v\n", item)
			} else {
				fmt.Printf("❌ Queue was empty")
			}
			time.Sleep(150 * time.Millisecond)
		}
		fmt.Printf("✅ Dequeuer finished: head %d, taik %d", q.head, q.tail)
	}()
	wg.Wait()
}
