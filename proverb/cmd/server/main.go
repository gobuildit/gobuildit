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
//
// The server binary starts an HTTP server on port 80 and servers Go proverbs in
// JSON.
package main

import (
	"log"

	"github.com/gobuildit/gobuildit/proverb"
)

func main() {
	addr := ":80"
	log.Printf("listening on %s", addr)

	s := proverb.NewServer(addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
}
