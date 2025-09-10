package main

import (
	"fmt"
	"sync"
	"time"
)

type LockBasedQueue struct {
	items []interface{}
	mutex sync.Mutex
}

func NewLockBasedQueue() *LockBasedQueue {
	return &LockBasedQueue{
		items: make([]interface{}, 0),
	}
}

func (q *LockBasedQueue) Enq(item interface{}, threadID int) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	fmt.Printf("🔒 Thread %d: Acquired lock for enq(%v)\n", threadID, item)
	// append the item after the thread acquires the lock
	q.items = append(q.items, item)
	fmt.Printf("➕ Thread %d: Added %v, queue now: %v\n", threadID, item, q.items)

}

func (q *LockBasedQueue) Deq(threadID int) interface{} {
	// the thread acquires the lock
	q.mutex.Lock()
	defer q.mutex.Unlock()

	fmt.Printf("🔒 Thread %d: Acquired lock for deq()\n", threadID)

	if len(q.items) == 0 {
		fmt.Printf("❌ Thread %d: Queue empty, returning nil\n", threadID)

	}
	item := q.items[0]
	q.items = q.items[1:]
	fmt.Printf("➖ Thread %d: Removed %v, queue now: %v\n", threadID, item, q.items)
	return item
}
func main() {
	queue := NewLockBasedQueue()
	var wg sync.WaitGroup
	fmt.Println("🎬 Testing Lock-based FIFO Queue\n")
	wg.Add(2)

	go func() {
		defer wg.Done()
		queue.Enq("A", 1)
		time.Sleep(50 * time.Millisecond)
		queue.Enq("C", 1)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(25 * time.Millisecond)
		queue.Enq("B", 2)
	}()

	wg.Wait()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("\n--- Now dequeuing ---")

	// test concurrent dequeues
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(id int) {
			defer wg.Done()
			result := queue.Deq(id + 3)
			fmt.Printf("🎯 Thread %d got: %v\n", id+3, result)
		}(i)
		time.Sleep(30 * time.Millisecond)
	}
	wg.Wait()

}
