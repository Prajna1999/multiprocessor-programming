// package main

// import (
// 	"fmt"
// 	"sync"
// 	"sync/atomic"
// 	"time"
// )

// type PetersonAtomic struct {
// 	flag [2]int32 // Use int32 for atomic operations
// 	turn int32
// }

// func (p *PetersonAtomic) Lock(threadID int) {
// 	other := 1 - threadID

// 	atomic.StoreInt32(&p.flag[threadID], 1)
// 	atomic.StoreInt32(&p.turn, int32(other))

// 	// Fixed condition: wait while OTHER wants in AND it's OTHER's turn
// 	for atomic.LoadInt32(&p.flag[other]) == 1 && atomic.LoadInt32(&p.turn) == int32(other) {
// 		// spin wait
// 	}
// 	fmt.Printf("🎯 Thread %d: ACTUALLY ENTERING\n", threadID)
// }

// func (p *PetersonAtomic) Unlock(threadID int) {
// 	atomic.StoreInt32(&p.flag[threadID], 0)
// 	fmt.Printf("🔓 Thread %d: LEAVING critical section\n", threadID)
// }

// func main() {
// 	peterson := &PetersonAtomic{}
// 	var wg sync.WaitGroup
// 	wg.Add(2)

// 	// Both threads start simultaneously
// 	go func() {
// 		defer wg.Done()
// 		peterson.Lock(0)
// 		fmt.Println("   Thread 0: Working in critical section...")
// 		time.Sleep(200 * time.Millisecond) // Do some work
// 		peterson.Unlock(0)
// 	}()

// 	go func() {
// 		defer wg.Done()
// 		peterson.Lock(1)
// 		fmt.Println("   Thread 1: Working in critical section...")
// 		time.Sleep(200 * time.Millisecond) // Do some work
// 		peterson.Unlock(1)
// 	}()

// 	wg.Wait()
// 	fmt.Println("Both threads completed")
// }

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type PetersonAtomic struct {
	flag [2]int32
	turn int32
}

var (
	threadsInCriticalSection int32 = 0
	maxSimultaneous          int32 = 0
)

func (p *PetersonAtomic) Lock(threadID int) {
	other := 1 - threadID

	atomic.StoreInt32(&p.flag[threadID], 1)
	atomic.StoreInt32(&p.turn, int32(other))

	for atomic.LoadInt32(&p.flag[other]) == 1 && atomic.LoadInt32(&p.turn) == int32(other) {
		// spin wait
	}
}

func (p *PetersonAtomic) Unlock(threadID int) {
	atomic.StoreInt32(&p.flag[threadID], 0)
}

func criticalSection(threadID int, peterson *PetersonAtomic) {
	peterson.Lock(threadID)

	// Enter critical section
	current := atomic.AddInt32(&threadsInCriticalSection, 1)

	// Track maximum simultaneous threads
	for {
		max := atomic.LoadInt32(&maxSimultaneous)
		if current <= max || atomic.CompareAndSwapInt32(&maxSimultaneous, max, current) {
			break
		}
	}

	fmt.Printf("Thread %d: ENTERED (total in CS: %d)\n", threadID, current)
	time.Sleep(100 * time.Millisecond)

	// Leave critical section
	atomic.AddInt32(&threadsInCriticalSection, -1)
	fmt.Printf("Thread %d: LEAVING\n", threadID)

	peterson.Unlock(threadID)
}

func main() {
	peterson := &PetersonAtomic{}
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		criticalSection(0, peterson)
	}()

	go func() {
		defer wg.Done()
		criticalSection(1, peterson)
	}()

	wg.Wait()

	fmt.Printf("\nResult: Max simultaneous threads in critical section: %d\n", maxSimultaneous)
	if maxSimultaneous > 1 {
		fmt.Println("❌ MUTUAL EXCLUSION VIOLATED!")
	} else {
		fmt.Println("✅ MUTUAL EXCLUSION PRESERVED!")
	}
}
