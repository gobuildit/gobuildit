// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The server provides two handlers, one subject to performance problems, the
// other not.
package main

import (
	"fmt"
	"log"
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
	// For synchronization of counters, the atomic package offers many helpful
	// functions. Typically, something like atomic.AddInt64 might be a better
	// choice for synchronized counters, but for the sake of demonstrating some
	// pifalls of locking, a mutex is used for synchronization instead.
	// See: https://golang.org/pkg/sync/atomic/
	mu    sync.Mutex
	count int
)

func main() {
	http.HandleFunc("/", root)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
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
