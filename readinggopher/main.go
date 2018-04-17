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

// The readinggopher tweets links from the Gopher Reading List.
package main

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	readingListJSON = "https://raw.githubusercontent.com/enocom/gopher-reading-list/master/README.json"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	log.Println("Starting Twitter bot")
	conf, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config from environment: %s", err)
	}

	links, err := getLinks(readingListJSON)
	if err != nil {
		log.Fatalf("failed to get link JSON: %s", err)
	}

	for _, l := range links {
		t, err := postTweet(conf, l.String())
		if err != nil {
			log.Fatalf("failed to send tweet: %s", err)
		}

		log.Printf("posted tweet: %s", t)
		time.Sleep(conf.frequency)
	}
}

// config holds all the required information to run the readinggopher binary.
type config struct {
	baseURL        string
	consumerKey    string
	consumerSecret string
	token          string
	tokenSecret    string
	frequency      time.Duration
}

// loadConfig creates a configuration from environments variables to avoid
// hard-coding secret values.
func loadConfig() (config, error) {
	var missing []string
	consumerKey := os.Getenv("CONSUMER_KEY")
	if consumerKey == "" {
		missing = append(missing, "CONSUMER_KEY")
	}
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	if consumerSecret == "" {
		missing = append(missing, "CONSUMER_SECRET")
	}
	token := os.Getenv("TOKEN")
	if token == "" {
		missing = append(missing, "TOKEN")
	}
	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		missing = append(missing, "TOKEN_SECRET")
	}
	frequency, err := time.ParseDuration(os.Getenv("FREQUENCY"))
	if err != nil {
		missing = append(missing, "FREQUENCY")
	}
	if len(missing) > 0 {
		err := fmt.Errorf(
			"missing the following environment variables: %s",
			strings.Join(missing, ", "),
		)
		return config{}, err
	}
	c := config{
		baseURL:        "https://api.twitter.com/1.1/statuses/update.json",
		consumerKey:    consumerKey,
		consumerSecret: consumerSecret,
		token:          token,
		tokenSecret:    tokenSecret,
		frequency:      frequency,
	}
	return c, nil
}

// getLinks assumes a JSON response with the following format:
// {
//    "gopherReadingList": [{
//       "title": "Some blog post title",
//       "url": "https://www.example.com",
//    }]
// }
func getLinks(u string) ([]link, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if want := resp.StatusCode; want != http.StatusOK {
		return nil, fmt.Errorf("want %d, got %d", want, http.StatusOK)
	}

	var r responseJSON
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r.links(), nil
}

// responseJSON represents the JSON format of the Gopher Reading List.
type responseJSON map[string][]link

func (r responseJSON) links() []link {
	return r["gopherReadingList"]
}

// link holds data to identify a notable blog post on Go.
type link struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

func (l link) String() string {
	return fmt.Sprintf("\"%s\" %s", l.Title, l.URL)
}

// tweet represents a published link.
type tweet struct {
	ID        string `json:"id_str"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

// link returns the permalink for a particular Tweet.
func (t tweet) link() string {
	return fmt.Sprintf("https://twitter.com/readinggopher/status/%s", t.ID)
}

func (t tweet) String() string {
	return fmt.Sprintf("{ID = %s, CreatedAt = %s, Text = %s}", t.ID, t.CreatedAt, t.Text)
}

// postTweet publishes a string of text to Twitter with the configured account.
func postTweet(c config, msg string) (tweet, error) {
	data := map[string]string{"status": msg}
	oauthData := map[string]string{
		"oauth_consumer_key":     c.consumerKey,
		"oauth_nonce":            strconv.FormatInt(rand.Int63(), 10),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
		"oauth_token":            c.token,
		"oauth_version":          "1.0",
	}
	sig := generateSignature(c.consumerSecret, c.tokenSecret, http.MethodPost, c.baseURL, oauthData, data)
	oauthData["oauth_signature"] = sig

	oauthHeader := buildOAuthHeader(oauthData, sig)
	body := buildBody(data)
	req, err := http.NewRequest(http.MethodPost, c.baseURL, body)
	if err != nil {
		return tweet{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", oauthHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return tweet{}, err
	}
	defer resp.Body.Close()

	if want := http.StatusOK; resp.StatusCode != want {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return tweet{}, fmt.Errorf("Tweet failed. Status = %d, Body = (failed to read)", resp.StatusCode)

		}
		return tweet{}, fmt.Errorf("Tweet failed. Status = %d, Body = %s", resp.StatusCode, string(body))
	}
	var t tweet
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return tweet{}, fmt.Errorf("failed to decode response body: %s", err)
	}
	return t, nil
}

// generateSignature produces an OAuth 1.0 signature for a request to the
// Twitter API.
func generateSignature(consumerSecret, tokenSecret, method, baseURL string, oauth map[string]string, data map[string]string) string {
	encoded := make(map[string]string)
	for k, v := range oauth {
		encoded[percentEncode(k)] = percentEncode(v)
	}
	for k, v := range data {
		encoded[percentEncode(k)] = percentEncode(v)
	}
	var keys sort.StringSlice
	for k, _ := range encoded {
		keys = append(keys, k)
	}
	keys.Sort()
	var kvs []string
	for _, k := range keys {
		kvs = append(kvs, fmt.Sprintf("%s=%s", k, encoded[k]))
	}
	paramStr := strings.Join(kvs, "&")
	signingKey := fmt.Sprintf("%s&%s", percentEncode(consumerSecret), percentEncode(tokenSecret))
	h := hmac.New(crypto.SHA1.New, []byte(signingKey))
	sigBaseStr := fmt.Sprintf("%s&%s&%s", method, percentEncode(baseURL), percentEncode(paramStr))
	h.Write([]byte(sigBaseStr))
	rawSignature := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(rawSignature)
}

// buildOAuthHeader creates the HTTP header based on OAuth 1.0 requirements for
// the Twitter API.
func buildOAuthHeader(oauthData map[string]string, sig string) string {
	encoded := make(map[string]string)
	for k, v := range oauthData {
		encoded[k] = percentEncode(v)
	}
	data := []string{
		fmt.Sprintf(`OAuth oauth_consumer_key="%s"`, encoded["oauth_consumer_key"]),
		fmt.Sprintf(`oauth_nonce="%s"`, encoded["oauth_nonce"]),
		fmt.Sprintf(`oauth_signature="%s"`, encoded["oauth_signature"]),
		fmt.Sprintf(`oauth_signature_method="%s"`, encoded["oauth_signature_method"]),
		fmt.Sprintf(`oauth_timestamp="%s"`, encoded["oauth_timestamp"]),
		fmt.Sprintf(`oauth_token="%s"`, encoded["oauth_token"]),
		fmt.Sprintf(`oauth_version="%s"`, encoded["oauth_version"]),
	}
	return strings.Join(data, ", ")
}

// buildBody creates a io.Reader with encoded key-value pairs.
func buildBody(data map[string]string) io.Reader {
	values := make(url.Values)
	for k, v := range data {
		values[k] = []string{v}
	}
	return strings.NewReader(values.Encode())
}

// validByte determines if a particular byte needs encoding in Twitter's percent
// encoding scheme.
func validByte(b byte) bool {
	if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') {
		return true
	}
	if b == '-' || b == '.' || b == '~' || b == '_' {
		return true
	}
	return false
}

// percentEncode encodes a string according to Twitter's percent encoding. For details, see:
// https://developer.twitter.com/en/docs/basics/authentication/guides/percent-encoding-parameters
func percentEncode(src string) string {
	bytes := []byte(src)
	var dst string
	for i := 0; i < len(bytes); i++ {
		b := bytes[i]
		if validByte(b) {
			dst += string(b)
			continue
		}
		dst += fmt.Sprintf("%%%X", b)
	}
	return dst
}
