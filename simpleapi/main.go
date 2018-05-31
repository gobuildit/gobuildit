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
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("missing DATABASE_URL")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open connection to database: %s", err)
	}
	http.Handle("/developers", withLogging(withAuthentication(&developerHandler{db: db})))

	http.Handle("/healthz", withLogging(http.HandlerFunc(index)))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to start HTTP server: %s", err)
	}
}

func withLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		log.Printf("%s %s", r.Method, r.URL.Path)
	})
}

func validAuth(auth string) bool {
	token := strings.TrimPrefix(auth, "Bearer ")
	if token != "secret" {
		return false
	}
	return true
}

func withAuthentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !validAuth(auth) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json;charset=utf-8")
	w.Write([]byte(`{"healthy":true}`))
}

type developerHandler struct {
	db *sql.DB
}

func (d *developerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json;charset=utf-8")
	rows, err := d.db.Query("SELECT id, name, age FROM developers")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var ds []developer
	for rows.Next() {
		var d developer
		if err := rows.Scan(&d.ID, &d.Name, &d.Age); err != nil {
			log.Printf("failed to scan developer: %s", err)
			continue
		}
		ds = append(ds, d)
	}
	if err := rows.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(developersList{Developers: ds})
}

type developersList struct {
	Developers []developer `json:"developers"`
}

type developer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}
