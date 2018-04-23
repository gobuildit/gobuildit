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
package proverb

import (
	"math/rand"
	"time"
)

// Proverb represents a particular proverb with a corresponding link to learn
// more.
type Proverb struct {
	Link    string `json:"link"`
	Content string `json:"content"`
}

// NewProverbStore initializes a ProverbStore.
func NewInMemProverbStore() *InMemProverbStore {
	rand.Seed(time.Now().Unix())
	return &InMemProverbStore{
		proverbs: []Proverb{
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=2m48s",
				Content: "Don't communicate by sharing memory, share memory by communicating.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=3m42s",
				Content: "Concurrency is not parallelism.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=4m20s",
				Content: "Channels orchestrate; mutexes serialize.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=5m17s",
				Content: "The bigger the interface, the weaker the abstraction.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=6m25s",
				Content: "Make the zero value useful.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=7m36s",
				Content: "interface{} says nothing.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=8m43s",
				Content: "Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=9m28s",
				Content: "A little copying is better than a little dependency.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=11m10s",
				Content: "Syscall must always be guarded with build tags.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=11m53s",
				Content: "Cgo must always be guarded with build tags.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=12m37s",
				Content: "Cgo is not Go.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=13m49s",
				Content: "With the unsafe package there are no guarantees.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=14m35s",
				Content: "Clear is better than clever.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=15m22s",
				Content: "Reflection is never clear.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=16m13s",
				Content: "Errors are values.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=17m25s",
				Content: "Don't just check errors, handle them gracefully.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=18m09s",
				Content: "Design the architecture, name the components, document the details.",
			},
			{
				Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=19m07s",
				Content: "Documentation is for users.",
			},
			{
				Link:    "https://github.com/golang/go/wiki/CodeReviewComments#dont-panic",
				Content: "Don't panic.",
			},
		},
	}
}

// ProverbStore provides access to the Go Proverbs.
type InMemProverbStore struct {
	proverbs []Proverb
}

// Random returns a random proverb.
func (p *InMemProverbStore) Random() Proverb {
	return p.proverbs[rand.Intn(len(p.proverbs))]
}
