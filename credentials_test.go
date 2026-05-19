package biliapi

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestCredentials(t *testing.T) {
	creds := Credentials{
		SESSDATA:        "sess",
		BiliJCT:         "csrf",
		DedeUserID:      "42",
		DedeUserIDCKMd5: "hash",
		SID:             "sid",
	}
	if !creds.IsLoggedInCandidate() {
		t.Fatal("expected SESSDATA credentials to be a login candidate")
	}
	if got := creds.CSRF(); got != "csrf" {
		t.Fatalf("CSRF() = %q", got)
	}
	if err := creds.RequireCSRF(); err != nil {
		t.Fatalf("RequireCSRF() returned error: %v", err)
	}
	if len(creds.Cookies()) != 5 {
		t.Fatalf("Cookies() length = %d", len(creds.Cookies()))
	}

	data, err := json.Marshal(creds)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseCredentialsJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.BiliJCT != creds.BiliJCT || parsed.SESSDATA != creds.SESSDATA {
		t.Fatalf("parsed credentials mismatch: %#v", parsed)
	}
}

func TestCredentialsRequireCSRF(t *testing.T) {
	err := (Credentials{SESSDATA: "sess"}).RequireCSRF()
	if !errors.Is(err, ErrMissingCSRF) {
		t.Fatalf("expected ErrMissingCSRF, got %v", err)
	}
}
