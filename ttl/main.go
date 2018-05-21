package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	// This goroutine reports on the runtime's number of goroutines and the
	// memory allocator's heap allocations every five seconds.
	go func() {
		var stats runtime.MemStats
		for {
			runtime.ReadMemStats(&stats)
			fmt.Printf("HeapAlloc    = %d\n", stats.HeapAlloc)
			fmt.Printf("NumGoroutine = %d\n", runtime.NumGoroutine())
			time.Sleep(5 * time.Second)
		}
	}()

	// A tight loop to create a large number of ttl.Map types
	for {
		work()
	}
}

// work creates a map, sets a key, and then closes it down.
func work() {
	m := NewMap(5 * time.Minute)
	m.Set("my-key", []byte("my-value"))

	if _, ok := m.Get("my-key"); !ok {
		panic("no value present")
	}
	m.Close()
	// m goes out of scope
}
