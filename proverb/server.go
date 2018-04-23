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
// Package proverb provides a store of Go proverbs and an HTTP server for those
// proverbs.
package proverb

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Server contains the http Server and registers all proverbs handlers.
type Server struct {
	http *http.Server
}

// NewServer initializes a new server without starting it.
func NewServer(addr string) *Server {
	mux := http.NewServeMux()
	s := NewInMemProverbStore()
	mux.Handle("/proverbs/random", NewRandomHandler(s))

	return &Server{
		http: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

// ListenAndServe starts the Server.
func (s *Server) ListenAndServe() error {
	return s.http.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.http.Shutdown(ctx)
}

// RandomHandler returns a random proverbs.
type RandomHandler struct {
	store ProverbStore
}

// ProverbStore defines the interface that provides random proverbs.
type ProverbStore interface {
	Random() Proverb
}

// NewRandomHandler initializes a RandomHandler.
func NewRandomHandler(s ProverbStore) http.Handler {
	return &RandomHandler{
		store: s,
	}
}

// ServeHTTP satisfies the http.Handler interface.
func (r *RandomHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	p := r.store.Random()

	rw.Header().Add("Content-type", "application/json; charset=utf-8")
	if err := json.NewEncoder(rw).Encode(p); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
