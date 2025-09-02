package main

import (
	"fmt"
	"sync"
	"time"
)

// Lamport's Bakery Algorithm
type Bakery struct {
	choosing []bool // choosing[i] = true if thread i is choosing a number
	number   []int  // number[i] = ticket number for thread i
	n        int    // number of threads
}

func NewBakery(n int) *Bakery {
	return &Bakery{
		choosing: make([]bool, n),
		number:   make([]int, n),
		n:        n,
	}
}

func (b *Bakery) Lock(threadID int) {
	fmt.Printf("🎫 Thread %d: Going to take a ticket\n", threadID)

	// Step 1: I'm choosing a number
	b.choosing[threadID] = true

	// Step 2: Take the maximum number + 1
	maxNumber := 0
	for j := 0; j < b.n; j++ {
		if b.number[j] > maxNumber {
			maxNumber = b.number[j]
		}
	}
	b.number[threadID] = maxNumber + 1
	fmt.Printf("🔢 Thread %d: Got ticket number %d\n", threadID, b.number[threadID])

	// Step 3: Done choosing
	b.choosing[threadID] = false

	// Step 4: Wait for all threads with smaller numbers (or same number but smaller ID)
	for j := 0; j < b.n; j++ {
		if j == threadID {
			continue // Skip myself
		}

		// Wait for thread j to finish choosing
		fmt.Printf("⏳ Thread %d: Waiting for Thread %d to finish choosing...\n", threadID, j)
		for b.choosing[j] {
			// spin wait
		}

		// Wait while thread j has a smaller ticket (or same ticket but smaller ID)
		for b.number[j] != 0 && b.hasHigherPriority(j, threadID) {
			fmt.Printf("⏳ Thread %d: Waiting for Thread %d (has priority)\n", threadID, j)
			time.Sleep(10 * time.Millisecond) // Small delay to see the waiting
		}
	}

	fmt.Printf("🏆 Thread %d: ENTERED CRITICAL SECTION with ticket %d!\n", threadID, b.number[threadID])
}

// Returns true if thread j has higher priority than thread i
func (b *Bakery) hasHigherPriority(j, i int) bool {
	return b.number[j] < b.number[i] || (b.number[j] == b.number[i] && j < i)
}

func (b *Bakery) Unlock(threadID int) {
	b.number[threadID] = 0 // Give up my ticket
	fmt.Printf("🔓 Thread %d: RELEASED LOCK (gave up ticket)\n", threadID)
}

func main() {
	const numThreads = 4
	bakery := NewBakery(numThreads)
	var wg sync.WaitGroup

	fmt.Printf("🥖 Starting Bakery Algorithm test with %d threads\n\n", numThreads)

	// Start all threads at different times to see ticket assignment
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			bakery.Lock(id)

			// Critical section
			fmt.Printf("💼 Thread %d: Working in critical section...\n", id)
			time.Sleep(200 * time.Millisecond)

			bakery.Unlock(id)
		}(i)

		// Stagger thread starts slightly
		time.Sleep(150 * time.Millisecond)
	}

	wg.Wait()
	fmt.Println("\n🎉 All threads completed!")
}
