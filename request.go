package biliapi

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type requestOptions struct {
	WBI          bool
	RequireLogin bool
}

func (c *Client) getJSON(ctx context.Context, endpoint string, query url.Values, opts requestOptions, dest any) error {
	return c.getJSONAllow(ctx, endpoint, query, opts, dest, nil)
}

func (c *Client) getJSONAllow(ctx context.Context, endpoint string, query url.Values, opts requestOptions, dest any, allowCodes map[int]bool) error {
	if opts.WBI {
		var err error
		query, err = c.signWBI(ctx, query)
		if err != nil {
			return err
		}
	}
	req, err := c.newRequest(ctx, http.MethodGet, endpoint, query, nil)
	if err != nil {
		return err
	}
	if opts.RequireLogin && !c.creds.IsLoggedInCandidate() {
		return ErrMissingCredentials
	}
	return c.doJSONAllow(req, dest, allowCodes)
}

func (c *Client) postFormJSON(ctx context.Context, endpoint string, form url.Values, opts requestOptions, dest any) error {
	if opts.RequireLogin && !c.creds.IsLoggedInCandidate() {
		return ErrMissingCredentials
	}
	req, err := c.newRequest(ctx, http.MethodPost, endpoint, nil, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.doJSON(req, dest)
}

func (c *Client) getRaw(ctx context.Context, endpoint string, query url.Values, opts requestOptions) ([]byte, error) {
	if opts.WBI {
		var err error
		query, err = c.signWBI(ctx, query)
		if err != nil {
			return nil, err
		}
	}
	req, err := c.newRequest(ctx, http.MethodGet, endpoint, query, nil)
	if err != nil {
		return nil, err
	}
	if opts.RequireLogin && !c.creds.IsLoggedInCandidate() {
		return nil, ErrMissingCredentials
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := readResponseBody(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &BiliError{HTTPStatus: resp.StatusCode, Message: string(body)}
	}
	return body, nil
}

func (c *Client) getPlainJSON(ctx context.Context, endpoint string, query url.Values, opts requestOptions, dest any) error {
	body, err := c.getRaw(ctx, endpoint, query, opts)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, dest)
}

func (c *Client) newRequest(ctx context.Context, method, endpoint string, query url.Values, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	if query != nil {
		q := u.Query()
		for key, values := range query {
			for _, value := range values {
				q.Add(key, value)
			}
		}
		u.RawQuery = q.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", c.userAgent)
	if c.referer != "" {
		req.Header.Set("Referer", c.referer)
	}
	for _, cookie := range c.creds.Cookies() {
		req.AddCookie(cookie)
	}
	return req, nil
}

func (c *Client) doJSON(req *http.Request, dest any) error {
	return c.doJSONAllow(req, dest, nil)
}

func (c *Client) doJSONAllow(req *http.Request, dest any, allowCodes map[int]bool) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := readResponseBody(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &BiliError{HTTPStatus: resp.StatusCode, Message: string(body)}
	}
	var envelope Response[json.RawMessage]
	if err := json.Unmarshal(body, &envelope); err != nil {
		return err
	}
	if envelope.Code != 0 && !allowCodes[envelope.Code] {
		return &BiliError{
			Code:    envelope.Code,
			Message: envelope.Message,
			TTL:     envelope.TTL,
			Data:    envelope.Data,
		}
	}
	if dest == nil || bytes.Equal(envelope.Data, []byte("null")) || len(envelope.Data) == 0 {
		return nil
	}
	return json.Unmarshal(envelope.Data, dest)
}

func readResponseBody(resp *http.Response) ([]byte, error) {
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	encoding := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Encoding")))
	switch encoding {
	case "deflate":
		zr, err := zlib.NewReader(bytes.NewReader(raw))
		if err == nil {
			defer zr.Close()
			return io.ReadAll(zr)
		}
		fr := flate.NewReader(bytes.NewReader(raw))
		defer fr.Close()
		return io.ReadAll(fr)
	default:
		return raw, nil
	}
}

func cloneValues(values url.Values) url.Values {
	out := make(url.Values, len(values))
	for key, vals := range values {
		out[key] = append([]string(nil), vals...)
	}
	return out
}

func setString(values url.Values, key, value string) {
	if value != "" {
		values.Set(key, value)
	}
}

func setInt(values url.Values, key string, value int) {
	if value != 0 {
		values.Set(key, strconv.Itoa(value))
	}
}

func setInt64(values url.Values, key string, value int64) {
	if value != 0 {
		values.Set(key, strconv.FormatInt(value, 10))
	}
}

func setBool01(values url.Values, key string, value bool) {
	if value {
		values.Set(key, "1")
	}
}

func setBool(values url.Values, key string, value *bool) {
	if value == nil {
		return
	}
	values.Set(key, strconv.FormatBool(*value))
}

func setCSVInt64(values url.Values, key string, ids []int64) {
	if len(ids) == 0 {
		return
	}
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, strconv.FormatInt(id, 10))
	}
	values.Set(key, strings.Join(parts, ","))
}

func requireAIDOrBVID(aid int64, bvid string) error {
	if aid == 0 && bvid == "" {
		return fmt.Errorf("%w: aid or bvid is required", ErrInvalidParams)
	}
	if aid != 0 && bvid != "" {
		return fmt.Errorf("%w: aid and bvid are mutually exclusive", ErrInvalidParams)
	}
	return nil
}

func requirePositive(name string, value int64) error {
	if value <= 0 {
		return fmt.Errorf("%w: %s must be positive", ErrInvalidParams, name)
	}
	return nil
}
