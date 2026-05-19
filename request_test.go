package biliapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetJSONDecodesEnvelope(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("User-Agent"); got == "" {
			t.Fatal("missing user-agent")
		}
		_ = json.NewEncoder(w).Encode(Response[NavStat]{
			Code: 0,
			TTL:  1,
			Data: NavStat{Following: 1, Follower: 2, DynamicCount: 3},
		})
	}))
	defer server.Close()

	c := NewClient()
	var out NavStat
	if err := c.getJSON(context.Background(), server.URL, nil, requestOptions{}, &out); err != nil {
		t.Fatal(err)
	}
	if out.Following != 1 || out.Follower != 2 || out.DynamicCount != 3 {
		t.Fatalf("decoded unexpected payload: %#v", out)
	}
}

func TestGetJSONBusinessError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(Response[map[string]string]{
			Code:    -101,
			Message: "账号未登录",
			TTL:     1,
			Data:    map[string]string{"hint": "login"},
		})
	}))
	defer server.Close()

	c := NewClient()
	var out NavStat
	err := c.getJSON(context.Background(), server.URL, nil, requestOptions{}, &out)
	var biliErr *BiliError
	if !errors.As(err, &biliErr) {
		t.Fatalf("expected BiliError, got %T %v", err, err)
	}
	if biliErr.Code != -101 {
		t.Fatalf("BiliError.Code = %d", biliErr.Code)
	}
}

func TestRequireAIDOrBVID(t *testing.T) {
	if err := requireAIDOrBVID(0, ""); err == nil {
		t.Fatal("expected missing id error")
	}
	if err := requireAIDOrBVID(1, "BV1"); err == nil {
		t.Fatal("expected mutually exclusive id error")
	}
	if err := requireAIDOrBVID(1, ""); err != nil {
		t.Fatal(err)
	}
	if err := requireAIDOrBVID(0, "BV1"); err != nil {
		t.Fatal(err)
	}
}

func TestWriteSkeletonRequiresCSRF(t *testing.T) {
	c := NewClient(WithCredentials(Credentials{SESSDATA: "sess"}))
	err := c.Video.Like(context.Background(), VideoLikeParams{AID: 1, Like: true})
	if !errors.Is(err, ErrMissingCSRF) {
		t.Fatalf("Video.Like error = %v", err)
	}
	_, err = c.Comment.Add(context.Background(), CommentAddParams{OID: 1, Type: CommentTypeVideo, Message: "hi"})
	if !errors.Is(err, ErrMissingCSRF) {
		t.Fatalf("Comment.Add error = %v", err)
	}
	err = c.Fav.DeleteFolder(context.Background(), 1)
	if !errors.Is(err, ErrMissingCSRF) {
		t.Fatalf("Fav.DeleteFolder error = %v", err)
	}
	err = c.History.Clear(context.Background())
	if !errors.Is(err, ErrMissingCSRF) {
		t.Fatalf("History.Clear error = %v", err)
	}
}
