package biliapi

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestNewClientOptionsOrder(t *testing.T) {
	custom := &http.Client{}
	client := NewClient(WithHTTPClient(custom), WithTimeout(3*time.Second))
	if client.httpClient != custom {
		t.Fatal("NewClient() did not use custom HTTP client")
	}
	if client.httpClient.Timeout != 3*time.Second {
		t.Fatalf("timeout = %s, want %s", client.httpClient.Timeout, 3*time.Second)
	}

	custom = &http.Client{}
	client = NewClient(WithTimeout(4*time.Second), WithHTTPClient(custom))
	if client.httpClient != custom {
		t.Fatal("NewClient() did not use custom HTTP client")
	}
	if client.httpClient.Timeout != 4*time.Second {
		t.Fatalf("timeout = %s, want %s", client.httpClient.Timeout, 4*time.Second)
	}
}

func TestNewClientPreservesCustomHTTPClientTimeout(t *testing.T) {
	custom := &http.Client{Timeout: time.Second}
	client := NewClient(WithHTTPClient(custom))
	if client.httpClient.Timeout != time.Second {
		t.Fatalf("timeout = %s, want %s", client.httpClient.Timeout, time.Second)
	}
}

func TestNewClientAttachesJarAndCredentialsToCustomHTTPClient(t *testing.T) {
	custom := &http.Client{}
	client := NewClient(WithHTTPClient(custom), WithCredentials(Credentials{
		SESSDATA: "sess",
		BiliJCT:  "csrf",
	}))
	if custom.Jar == nil {
		t.Fatal("custom HTTP client jar was not initialized")
	}
	if client.jar != custom.Jar {
		t.Fatal("client jar and custom HTTP client jar differ")
	}

	u, err := url.Parse("https://api.bilibili.com/")
	if err != nil {
		t.Fatal(err)
	}
	cookies := custom.Jar.Cookies(u)
	if len(cookies) == 0 {
		t.Fatal("credentials were not written to the custom HTTP client jar")
	}
}
