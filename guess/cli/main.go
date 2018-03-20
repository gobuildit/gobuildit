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

// The cli binary is a command line interface number guessing game.
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Ensure the random numbers from the rand package are different each time.
	rand.Seed(time.Now().Unix())

	// Initialize the secret number to a value between 1 and 10.
	secretNumber := rand.Intn(10) + 1

	// Wrap Stdin (which conforms to the io.Reader interface) in a bufio.Scanner,
	// a type that provides a number of niceties for handling textual I/O.
	// The bufio.Scanner allows callers to specify a split function. By default,
	// the bufio.Scanner uses bufio.ScanLines, which fits the use case here of
	// reading only one line at a time.
	r := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to the Guessing Game!")
	fmt.Println("Rules: Guess the secret number.")

	// Loop until the game is over.
	for {
		fmt.Print("Guess: ")
		// Scan one line of input.
		r.Scan()

		// rawGuess is still a string at this point.
		rawGuess, err := r.Text(), r.Err()
		if err != nil {
			fmt.Printf("failed to read from Stdin: %s", err)
			continue
		}
		// Remove the newline after the guess.
		trimmed := strings.TrimSpace(rawGuess)
		fmt.Printf("You entered: %q\n", trimmed)

		// Attempt to convert the guess into a number.
		num, err := strconv.Atoi(trimmed)
		if err != nil {
			fmt.Printf("%s is not a number. Try again!\n", trimmed)
			continue
		}

		// Keep looping if the guess is not correct.
		if num != secretNumber {
			fmt.Printf("Nope! %d is not correct. Keep guesing!\n", num)
			continue
		}

		// User has guessed the number. Exit the loop.
		fmt.Printf("Correct! %d is the secret number! You win!\n", num)
		break
	}

	fmt.Println("GAME OVER")
}
