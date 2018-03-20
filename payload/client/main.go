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

// The client sends JSON or Protobuf payloads to a local web server.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gobuildit/gobuildit/payload"
	"github.com/golang/protobuf/proto"
)

func main() {
	protoPayload := flag.Bool("proto", false, "use protobuf for payload")
	flag.Parse()
	err := sendRequest("http://localhost:8080", *protoPayload)
	if err != nil {
		log.Fatalf("failed to send request: %s", err)
	}
}

func sendRequest(url string, useProto bool) error {
	var (
		p           []byte
		err         error
		contentType string
	)
	if useProto {
		contentType = payload.ContentTypeProtobuf
		p, err = protobufPayload()
	} else {
		contentType = payload.ContentTypeJSON
		p, err = jsonPayload()
	}
	if err != nil {
		return err
	}
	return send(url, p, contentType)
}

func send(url string, p []byte, contentType string) error {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(p))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", contentType)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		return fmt.Errorf("Status=%d; want %d", got, want)
	}
	return nil
}

func protobufPayload() ([]byte, error) {
	l, _ := time.LoadLocation("Asia/Tokyo")
	d := time.Date(1954, time.April, 26, 0, 0, 0, 0, l).Format(time.RFC3339)

	m := &payload.Movie{
		Title: "Seven Samurai",
		Director: &payload.Person{
			FirstName: "Akira",
			LastName:  "Kurozawa",
		},
		Cast: []*payload.Person{
			{FirstName: "Toshiro", LastName: "Mifune"},
			{FirstName: "Takashi", LastName: "Shimura"},
			{FirstName: "Tsushima", LastName: "Keiko"},
		},
		Release: d,
	}
	return proto.Marshal(m)
}

func jsonPayload() ([]byte, error) {
	l, _ := time.LoadLocation("Asia/Tokyo")
	m := payload.JSONMovie{
		Title: "Seven Samurai",
		Director: payload.JSONPerson{
			FirstName: "Akira",
			LastName:  "Kurozawa",
		},
		Cast: []payload.JSONPerson{
			{FirstName: "Toshiro", LastName: "Mifune"},
			{FirstName: "Takashi", LastName: "Shimura"},
			{FirstName: "Tsushima", LastName: "Keiko"},
		},
		Release: time.Date(1954, time.April, 26, 0, 0, 0, 0, l),
	}
	return json.Marshal(m)
}
