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

// The web binary runs a number guessing game in a server.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const addr = "localhost:8080"

func main() {
	// Ensure the random numbers from the rand package are different each time.
	rand.Seed(time.Now().Unix())
	log.Printf("starting server on %s\n", addr)
	http.Handle("/guesses", newGame())
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("http.ListenAndServe error: %s", err)
	}
}

// newGame creates a guessingGame with a random number.
func newGame() *guessingGame {
	return &guessingGame{
		num: rand.Intn(10) + 1,
	}
}

// guessingGame holds all the state of the game.
type guessingGame struct {
	num int

	// mu synchronizes read and write access to `over`.
	mu   sync.Mutex
	over bool
}

// ServeHTTP ensures guessingGame implements the http.Handler interface and
// responds to POST methods with a JSON body holding a client's guess. If the
// guess is correct, the handler returns a 200 OK status. Otherwise, it returns
// a 418 I'm a teapot status, indicating the game is not yet over.
func (gg *guessingGame) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Ensure all responses are of content type application JSON.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Ensure the HTTP method is a POST.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"message":"only HTTP POST methods are supported"}`))
		return
	}

	// Read guess from the request.
	var g guess
	err := json.NewDecoder(r.Body).Decode(&g)
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Printf("failed to close response body: %s", closeErr)
		}
	}()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"failed to read JSON body"}`))
		return
	}

	status, body := gg.guess(g.Number)

	// Write the result of the guess to the client connection.
	w.WriteHeader(status)
	w.Write([]byte(body))
}

func (gg *guessingGame) guess(number int) (int, string) {
	// Uninitialized variables default to their zero value.
	var (
		status int
		body   string
	)
	// Enter the critical section. Reads and writes of `gg.over` are
	// synchronized across all client requests. Note that the locking does not
	// occur around network I/O.
	gg.mu.Lock()
	// The switch statement below will always first check if the game is over.
	//
	// Note: "case expressions... are evaluated left-to-right and top-to-bottom;
	// the first one that equals the switch expression triggers execution of the
	// statements of the associated case; the other cases are skipped."
	//
	// See https://golang.org/ref/spec#Switch_statements
	switch {
	case gg.over:
		status = http.StatusBadRequest
		body = `{"message":"game is already over"}`
	case gg.num == number:
		gg.over = true
		status = http.StatusCreated
		body = fmt.Sprintf(`{"message":"Correct! %d is the secret number! You win!"}`, number)
	default:
		status = http.StatusTeapot
		body = `{"message":"guess again"}`
	}
	gg.mu.Unlock()

	return status, body
}

// guess represents a client's submitted guess.
type guess struct {
	Number int `json:"number"`
}

// String implements the Stringer interface.
func (g guess) String() string {
	return fmt.Sprintf("guess{ number = %d }", g.Number)
}
