package main

import (
	"fmt"
	"sync"
	"time"
)

type Peterson struct {
	flag [2]bool
	turn int
}

func (p *Peterson) Lock(threadID int) {
	other := 1 - threadID
	p.flag[threadID] = true
	fmt.Printf("🚪 Thread %d: I want to enter\n", threadID)
	p.turn = other
	fmt.Printf("🤝 Thread %d: Setting turn = %d\n", threadID, other)

	for p.flag[other] && p.turn == other {
		// spin wait
	}
	fmt.Printf("🎯 Thread %d: ACTUALLY ENTERING critical section now!\n", threadID)
}

func (p *Peterson) Unlock(threadID int) {
	p.flag[threadID] = false
	fmt.Printf("🔓 Thread %d: LEAVING critical section\n", threadID)
}

func main() {
	peterson := &Peterson{}
	var wg sync.WaitGroup
	wg.Add(2)

	// Both threads start simultaneously
	go func() {
		defer wg.Done()
		peterson.Lock(0)
		fmt.Println("   Thread 0: Working in critical section...")
		time.Sleep(200 * time.Millisecond) // Do some work
		peterson.Unlock(0)
	}()

	go func() {
		defer wg.Done()
		peterson.Lock(1)
		fmt.Println("   Thread 1: Working in critical section...")
		time.Sleep(200 * time.Millisecond) // Do some work
		peterson.Unlock(1)
	}()

	wg.Wait()
	fmt.Println("Both threads completed")
}
