package main

import (
	"fmt"
	"sync"
	"time"
)

// Filter Lock implementation (paste the previous code here)
type Filter struct {
	level  []int
	victim []int
	n      int
}

func NewFilter(n int) *Filter {
	return &Filter{
		level:  make([]int, n),
		victim: make([]int, n),
		n:      n,
	}
}

func (f *Filter) Lock(threadID int) {
	fmt.Printf("🚪 Thread %d: Starting to acquire lock\n", threadID)

	// Thread must pass through levels 1, 2, ..., n-1
	for L := 1; L < f.n; L++ {
		f.level[threadID] = L
		fmt.Printf("📈 Thread %d: Moved to level %d\n", threadID, L)

		f.victim[L] = threadID
		fmt.Printf("🎯 Thread %d: Became victim at level %d\n", threadID, L)

		// Wait while some other thread is at level L or higher AND I'm the victim
		for f.existsThreadAtLevelOrHigher(L, threadID) && f.victim[L] == threadID {
			// spin wait
		}

		fmt.Printf("✅ Thread %d: Passed level %d\n", threadID, L)
	}
	fmt.Printf("🏆 Thread %d: ENTERED CRITICAL SECTION!\n", threadID)
}

func (f *Filter) existsThreadAtLevelOrHigher(L int, threadID int) bool {
	for k := 0; k < f.n; k++ {
		if k != threadID && f.level[k] >= L {
			return true
		}
	}
	return false
}

func (f *Filter) Unlock(threadID int) {
	f.level[threadID] = 0
	fmt.Printf("🔓 Thread %d: RELEASED LOCK (back to level 0)\n", threadID)
}

func main() {
	const numThreads = 3
	filter := NewFilter(numThreads)
	var wg sync.WaitGroup

	fmt.Printf("🎬 Starting Filter Lock test with %d threads\n\n", numThreads)

	// Start all threads roughly simultaneously
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			filter.Lock(id)

			// Critical section
			fmt.Printf("💼 Thread %d: Working in critical section...\n", id)
			time.Sleep(300 * time.Millisecond) // Simulate work

			filter.Unlock(id)
		}(i)

		// Small delay to make output more readable
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
	fmt.Println("\n🎉 All threads completed!")
}
