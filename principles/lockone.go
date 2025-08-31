package main

import (
	"fmt"
	"sync"
	"time"
)

type LockOne struct {
	flag [2]bool // flag[0] for thread 0 and flag[1] for thread 1

}

func NewLockOne() *LockOne {
	return &LockOne{
		flag: [2]bool{false, false},
	}
}

func (l *LockOne) Lock(threadID int) {
	other := 1 - threadID
	// toggle the thread that wants to enter
	l.flag[threadID] = true
	fmt.Printf("Thread %d: Set my flag, now checking other's flag\n", threadID)
	for l.flag[other] {
		// wait while other thread wants to enter
		// this is where the deadlock happend
	}

	fmt.Printf("Thread %d: Entered critical section!\n", threadID)
}

func (l *LockOne) Unlock(threadID int) {
	l.flag[threadID] = false // the thread unlocks the resource
	fmt.Printf("Thread %d: Released lock\n", threadID)
}

func main() {
	lock := NewLockOne()
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
