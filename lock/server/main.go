package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

const (
	// payloadBytes is the number of bytes written back to the client. The value
	// ensures all writes to a client socket fill the TCP buffer. The TCP buffer
	// size is controlled by kernel configuration. To learn more about the TCP
	// configuration, see: http://fasterdata.es.net/host-tuning/
	payloadBytes = 1024 * 1024
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

	msg := []byte(strings.Repeat(fmt.Sprintf("%d", count), payloadBytes))
	w.Write(msg)
}

// GOOD: This version of root holds the lock as shortly as possible,
// incrementing count and storing a copy of the current count value, before
// releasing the lock. If a client is slow to read, the handler may still be
// invoked a second time without creating any lock contention.
/*
func root(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	current := count
	mu.Unlock()

	msg := []byte(strings.Repeat(fmt.Sprintf("%d", current), payloadBytes))
	w.Write([]byte(msg))
}
*/
