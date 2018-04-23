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

package proverb_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gobuildit/gobiuldit/proverb"
)

func TestRandomHandler(t *testing.T) {
	s := newStubStore(proverb.Proverb{
		Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=3m42s",
		Content: "Concurrency is not parallelism.",
	})
	h := proverb.NewRandomHandler(s)
	recorder := httptest.NewRecorder()

	h.ServeHTTP(recorder, &http.Request{})

	if recorder.Code != http.StatusOK {
		t.Fatalf("want %v, got %v", http.StatusOK, recorder.Code)
	}

	got := recorder.Header().Get("Content-type")
	want := "application/json; charset=utf-8"
	if got != want {
		t.Fatalf("want %v, got %v", want, got)
	}

	actual := proverb.Proverb{}
	err := json.Unmarshal(recorder.Body.Bytes(), &actual)
	if err != nil {
		t.Fatalf("expected unmarshal JSON to succeed, got %v", err)
	}

	expected := proverb.Proverb{
		Link:    "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=3m42s",
		Content: "Concurrency is not parallelism.",
	}
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}
}

func newStubStore(p proverb.Proverb) stubStore {
	return stubStore{
		randomReturns: p,
	}
}

type stubStore struct {
	randomReturns proverb.Proverb
}

func (s stubStore) Random() proverb.Proverb {
	return s.randomReturns
}
