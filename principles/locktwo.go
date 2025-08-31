package main

import (
	"fmt"
	"sync"
	"time"
)

type LockTwo struct {
	flag [2]bool // flag[0] for thread 0 and flag[1] for thread 1

}

func NewLockTwo() *LockTwo {
	return &LockTwo{
		flag: [2]bool{false, false},
	}
}

func (l *LockTwo) Lock(threadID int) {
	other := 1 - threadID
	// first check the other's flag
	for l.flag[other] {
		// wait
	}
	// then set my flag
	l.flag[threadID] = true
	fmt.Printf("Thread %d: Set my flag \n", threadID)
	fmt.Printf("Thread %d: Entered critical section!\n", threadID)
}

func (l *LockTwo) Unlock(threadID int) {
	l.flag[threadID] = false // the thread unlocks the resource
	fmt.Printf("Thread %d: Released lock\n", threadID)
}

func main() {
	lock := NewLockTwo()
	var wg sync.WaitGroup
	wg.Add(2)

	// start both threads roughly at the same time
	go func() {
		defer wg.Done()
		lock.Lock(0)
		fmt.Println("Thread 0: Doing work...")
		time.Sleep(100 * time.Millisecond)
		lock.Unlock(0)
	}()

	go func() {
		defer wg.Done()
		lock.Lock(1)
		fmt.Println("Thread 1: Doing work...")
		time.Sleep(100 * time.Millisecond)
		lock.Unlock(1)
	}()

	// add a timeout to detect deadlock
	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true

	}()
	select {
	case <-done:
		fmt.Println("SUCCESS: Both threads completed")
	case <-time.After(2 * time.Second):
		fmt.Println("DEADLOCK DETECTED: Threads are stuck!")
	}

}
