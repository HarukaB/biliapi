package biliapi

import (
	"encoding/json"
	"net/http"
)

type Credentials struct {
	SESSDATA        string `json:"SESSDATA,omitempty"`
	BiliJCT         string `json:"bili_jct,omitempty"`
	DedeUserID      string `json:"DedeUserID,omitempty"`
	DedeUserIDCKMd5 string `json:"DedeUserID__ckMd5,omitempty"`
	SID             string `json:"sid,omitempty"`
	RefreshToken    string `json:"refresh_token,omitempty"`
}

func ParseCredentialsJSON(data []byte) (Credentials, error) {
	var creds Credentials
	err := json.Unmarshal(data, &creds)
	return creds, err
}

func (c Credentials) IsZero() bool {
	return c.SESSDATA == "" &&
		c.BiliJCT == "" &&
		c.DedeUserID == "" &&
		c.DedeUserIDCKMd5 == "" &&
		c.SID == "" &&
		c.RefreshToken == ""
}

func (c Credentials) IsLoggedInCandidate() bool {
	return c.SESSDATA != ""
}

func (c Credentials) CSRF() string {
	return c.BiliJCT
}

func (c Credentials) RequireCSRF() error {
	if c.BiliJCT == "" {
		return ErrMissingCSRF
	}
	return nil
}

func (c Credentials) Cookies() []*http.Cookie {
	var cookies []*http.Cookie
	add := func(name, value string, httpOnly bool) {
		if value == "" {
			return
		}
		cookies = append(cookies, &http.Cookie{
			Name:     name,
			Value:    value,
			Path:     "/",
			Domain:   ".bilibili.com",
			Secure:   true,
			HttpOnly: httpOnly,
		})
	}
	add("SESSDATA", c.SESSDATA, true)
	add("bili_jct", c.BiliJCT, false)
	add("DedeUserID", c.DedeUserID, false)
	add("DedeUserID__ckMd5", c.DedeUserIDCKMd5, false)
	add("sid", c.SID, false)
	return cookies
}
