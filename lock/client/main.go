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

// The client sends a GET request and reads nothing from the response, causing a
// poorly implemented server to block.
package main

import (
	"log"
	"net"
)

func main() {
	println("dialing")
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	println("sending GET request")
	_, err = conn.Write([]byte("GET / HTTP/1.1\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write([]byte("Host: localhost\r\n\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	println("blocking and never reading")
	select {}
}
