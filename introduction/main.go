// simulating race condition
package main

import (
	"fmt"
	"sync"
	"time"
)

var sharedValue int

func increment(name string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Read the value
	temp := sharedValue
	fmt.Printf("%s read: %d\n", name, temp)

	// Simulate some processing time
	time.Sleep(10 * time.Millisecond)

	// Write back the incremented value
	sharedValue = temp + 1
	fmt.Printf("%s wrote: %d\n", name, sharedValue)
}

func main() {

	// run this multiple times to see different results
	for run := 1; run <= 5; run++ {
		sharedValue = 10

		var wg sync.WaitGroup
		wg.Add(2)
		go increment("Goroutine A", &wg)
		// uncomment this line to hotfix race condition (mostly)
		// time.Sleep(100 * time.Millisecond)
		go increment("Goroutine B", &wg)

		wg.Wait()
		fmt.Printf("Run %d - Final value :%d (expected:12)\n\n", run, sharedValue)

		time.Sleep(1 * time.Millisecond)
	}

}
