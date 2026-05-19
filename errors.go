package biliapi

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrMissingCredentials = errors.New("biliapi: missing credentials")
	ErrMissingCSRF        = errors.New("biliapi: missing csrf token")
	ErrInvalidParams      = errors.New("biliapi: invalid parameters")
)

type BiliError struct {
	Code       int
	Message    string
	TTL        int
	Data       json.RawMessage
	HTTPStatus int
}

func (e *BiliError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.HTTPStatus != 0 {
		return fmt.Sprintf("biliapi: http status %d: code=%d message=%s", e.HTTPStatus, e.Code, e.Message)
	}
	return fmt.Sprintf("biliapi: code=%d message=%s", e.Code, e.Message)
}

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    T      `json:"data"`
}
