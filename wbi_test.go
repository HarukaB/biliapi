package biliapi

import (
	"net/url"
	"testing"
)

func TestGenMixinKey(t *testing.T) {
	got := genMixinKey("7cd084941338484aae1ad9425b84077c" + "4932caff0ff746eab6f01bf08b70ac45")
	want := "ea1db124af3c7062474693fa704f4ff8"
	if got != want {
		t.Fatalf("genMixinKey() = %q, want %q", got, want)
	}
}

func TestEncodeWBI(t *testing.T) {
	params := url.Values{}
	params.Set("foo", "114")
	params.Set("bar", "514")
	params.Set("zab", "1919810")
	got := encodeWBI(params, "ea1db124af3c7062474693fa704f4ff8", 1702204169)
	if got.Get("wts") != "1702204169" {
		t.Fatalf("wts = %q", got.Get("wts"))
	}
	if got.Get("w_rid") != "8f6f2b5b3d485fe1886cec6a0be8c5d4" {
		t.Fatalf("w_rid = %q", got.Get("w_rid"))
	}
}

func TestEncodeURIComponent(t *testing.T) {
	got := encodeURIComponent("五一四 one")
	want := "%E4%BA%94%E4%B8%80%E5%9B%9B%20one"
	if got != want {
		t.Fatalf("encodeURIComponent() = %q, want %q", got, want)
	}
}

func TestWBIKeyFromURL(t *testing.T) {
	got, err := wbiKeyFromURL("https://i0.hdslb.com/bfs/wbi/7cd084941338484aae1ad9425b84077c.png")
	if err != nil {
		t.Fatal(err)
	}
	if got != "7cd084941338484aae1ad9425b84077c" {
		t.Fatalf("wbiKeyFromURL() = %q", got)
	}
}
