package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	mu    sync.Mutex
	count int
)

func main() {
	http.HandleFunc("/", root)

	http.ListenAndServe(":8080", nil)
}

// BAD: This version of root holds the lock around the write to the client,
// which means if the client is slow to read, then the handler cannot quickly
// release the lock and respond to the next request.
func root(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	count++

	// Simulate a slow client based on the presence of a query parameter.
	if _, ok := r.URL.Query()["slow"]; ok {
		time.Sleep(10 * time.Second)
	}

	msg := fmt.Sprintf("Count = %d", current)
	w.Write([]byte(msg))
}

// GOOD: This version of root holds the lock for as short as possible,
// incrementing count and storing a copy of the current count value, before
// releasing the lock. If a client is slow to read, the handler may still be
// invoked a second time without creating any lock contention.
/*
func root(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	current := count
	mu.Unlock()

	// Simulate a slow client based on the presence of a query parameter.
	if _, ok := r.URL.Query()["slow"]; ok {
		time.Sleep(10 * time.Second)
	}

	msg := fmt.Sprintf("Count = %d", current)
	w.Write([]byte(msg))
}
*/
