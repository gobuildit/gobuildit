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

// The server consumes JSON and Protobuf payloads and reports on their
// respective sizes.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gobuildit/gobuildit/payload"
	"github.com/golang/protobuf/proto"
)

const addr = "localhost:8080"

func main() {
	log.Printf("starting server on %s", addr)

	http.HandleFunc("/", payloadHandler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("http.ListenAndServe error: %s", err)
	}
}

func payloadHandler(rw http.ResponseWriter, req *http.Request) {
	cType := req.Header.Get("Content-Type")
	cLength := req.Header.Get("Content-Length")

	log.Printf("%s = %s bytes\n", cType, cLength)
	defer req.Body.Close()
	if cType == payload.ContentTypeProtobuf {
		printProtobuf(req.Body)
	} else {
		printJSON(req.Body)
	}

	rw.WriteHeader(http.StatusCreated)
}

func printProtobuf(reqBody io.Reader) {
	m, err := unmarshalProtobuf(reqBody)
	if err != nil {
		log.Println("failed to unmarshal protobuf body: %s", err)
		return
	}
	if err := proto.MarshalText(os.Stdout, m); err != nil {
		log.Println("failed to marshal protobuf to text: %s", err)
		return
	}
}

func printJSON(reqBody io.Reader) {
	m, err := unmarshalJSON(reqBody)
	if err != nil {
		log.Println("failed to unmarshal request body: %s", err)
		return
	}
	body, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Println("failed to marshal movie with indentation: %s", err)
		return
	}
	fmt.Println(string(body))
}

func unmarshalJSON(rc io.Reader) (payload.JSONMovie, error) {
	m := payload.JSONMovie{}
	if err := json.NewDecoder(rc).Decode(&m); err != nil {
		return payload.JSONMovie{}, err
	}
	return m, nil
}

func unmarshalProtobuf(rc io.Reader) (*payload.Movie, error) {
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	m := &payload.Movie{}
	if err := proto.Unmarshal(b, m); err != nil {
		return nil, err
	}
	return m, nil
}
