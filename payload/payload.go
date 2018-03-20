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

// Package payload provides JSON and Protobuf representations of a Movie.
//
// To regenerate the movie.pb.go file, run `go generate`, which will invoke the
// directive below:
//
//go:generate protoc --go_out=. movie.proto

package payload

import (
	"fmt"
	"time"
)

const (
	// ContentTypeProtobuf is MIME type for protobuf content types.
	ContentTypeProtobuf = "application/vnd.google.protobuf"
	// ContentTypeJSON is MIME type for JSON.
	ContentTypeJSON = "application/json; charset=utf-8"
)

// JSONMovie holds identifying information about a released film.
type JSONMovie struct {
	Title    string       `json:"title"`
	Director JSONPerson   `json:"director"`
	Cast     []JSONPerson `json:"cast"`
	Release  time.Time    `json:"release"`
}

// JSONPerson represents an individual involved with a film production.
type JSONPerson struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// String implements the Stringer interface.
func (p JSONPerson) String() string {
	return fmt.Sprintf("%s %s", p.FirstName, p.LastName)
}
