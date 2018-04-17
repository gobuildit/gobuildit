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

import "testing"

func TestBuildOAuthHeader(t *testing.T) {
	oauthData := map[string]string{
		"oauth_consumer_key":     "xvz1evFS4wEEPTGEFPHBog",
		"oauth_nonce":            "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		"oauth_signature":        "tnnArxj06cWHq44gCs1OSKk/jLY=",
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        "1318622958",
		"oauth_token":            "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
		"oauth_version":          "1.0",
	}
	sig := "tnnArxj06cWHq44gCs1OSKk%2FjLY%3D"

	got := buildOAuthHeader(oauthData, sig)
	want := `OAuth oauth_consumer_key="xvz1evFS4wEEPTGEFPHBog", oauth_nonce="kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg", oauth_signature="tnnArxj06cWHq44gCs1OSKk%2FjLY%3D", oauth_signature_method="HMAC-SHA1", oauth_timestamp="1318622958", oauth_token="370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb", oauth_version="1.0"`

	if want != got {
		t.Fatalf("want %v, got %s", want, got)
	}
}

func TestEncodeString(t *testing.T) {
	tt := map[string]string{
		"Ladies + Gentlemen": "Ladies%20%2B%20Gentlemen",
		"An encoded string!": "An%20encoded%20string%21",
		"Dogs, Cats & Mice":  "Dogs%2C%20Cats%20%26%20Mice",
		"☃":                  "%E2%98%83",
	}

	for in, want := range tt {
		if got := percentEncode(in); got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	}
}

func TestGenerateSig(t *testing.T) {
	oauth := map[string]string{
		"oauth_consumer_key":     "xvz1evFS4wEEPTGEFPHBog",
		"oauth_nonce":            "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        "1318622958",
		"oauth_token":            "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
		"oauth_version":          "1.0",
	}
	data := map[string]string{
		"status":           "Hello Ladies + Gentlemen, a signed OAuth request!",
		"include_entities": "true",
	}

	consumerSecret := "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw"
	tokenSecret := "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE"
	got := generateSignature(consumerSecret, tokenSecret, "POST", "https://api.twitter.com/1.1/statuses/update.json",
		oauth, data)
	want := "hCtSmYh+iHYCEqBWrE7C7hYmtUk="

	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}
