package biliapi

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mixinKeyEncTab = []int{
	46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35,
	27, 43, 5, 49, 33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13,
	37, 48, 7, 16, 24, 55, 40, 61, 26, 17, 0, 1, 60, 51, 30, 4,
	22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11, 36, 20, 34, 44, 52,
}

type wbiState struct {
	mu        sync.Mutex
	imgKey    string
	subKey    string
	mixinKey  string
	updatedAt time.Time
	now       func() time.Time
}

func newWBIState(now func() time.Time) *wbiState {
	return &wbiState{now: now}
}

func (c *Client) signWBI(ctx context.Context, values url.Values) (url.Values, error) {
	keys, err := c.wbiKeys(ctx)
	if err != nil {
		return nil, err
	}
	wts := c.wbi.now().Unix()
	return encodeWBI(values, keys.mixinKey, wts), nil
}

type wbiKeys struct {
	imgKey   string
	subKey   string
	mixinKey string
}

func (c *Client) wbiKeys(ctx context.Context) (wbiKeys, error) {
	c.wbi.mu.Lock()
	if c.wbi.mixinKey != "" && c.wbi.now().Sub(c.wbi.updatedAt) < 12*time.Hour {
		keys := wbiKeys{imgKey: c.wbi.imgKey, subKey: c.wbi.subKey, mixinKey: c.wbi.mixinKey}
		c.wbi.mu.Unlock()
		return keys, nil
	}
	c.wbi.mu.Unlock()

	var nav navWBIResponse
	if err := c.getJSONAllow(ctx, endpointNav, nil, requestOptions{}, &nav, map[int]bool{-101: true}); err != nil {
		return wbiKeys{}, err
	}
	imgKey, err := wbiKeyFromURL(nav.WBIImg.ImgURL)
	if err != nil {
		return wbiKeys{}, err
	}
	subKey, err := wbiKeyFromURL(nav.WBIImg.SubURL)
	if err != nil {
		return wbiKeys{}, err
	}
	mixin := genMixinKey(imgKey + subKey)

	c.wbi.mu.Lock()
	c.wbi.imgKey = imgKey
	c.wbi.subKey = subKey
	c.wbi.mixinKey = mixin
	c.wbi.updatedAt = c.wbi.now()
	c.wbi.mu.Unlock()
	return wbiKeys{imgKey: imgKey, subKey: subKey, mixinKey: mixin}, nil
}

func genMixinKey(raw string) string {
	var b strings.Builder
	for _, index := range mixinKeyEncTab {
		if index >= len(raw) {
			continue
		}
		b.WriteByte(raw[index])
		if b.Len() == 32 {
			break
		}
	}
	return b.String()
}

func wbiKeyFromURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	base := path.Base(u.Path)
	key := strings.TrimSuffix(base, path.Ext(base))
	if key == "" {
		return "", fmt.Errorf("%w: empty wbi key in %q", ErrInvalidParams, raw)
	}
	return key, nil
}

func encodeWBI(values url.Values, mixinKey string, wts int64) url.Values {
	signed := cloneValues(values)
	signed.Set("wts", strconv.FormatInt(wts, 10))

	keys := make([]string, 0, len(signed))
	for key := range signed {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		vals := append([]string(nil), signed[key]...)
		sort.Strings(vals)
		for _, value := range vals {
			parts = append(parts, encodeURIComponent(key)+"="+encodeURIComponent(filterWBISignValue(value)))
		}
	}
	query := strings.Join(parts, "&")
	sum := md5.Sum([]byte(query + mixinKey))
	signed.Set("w_rid", hex.EncodeToString(sum[:]))
	return signed
}

func filterWBISignValue(value string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '!', '\'', '(', ')', '*':
			return -1
		default:
			return r
		}
	}, value)
}

func encodeURIComponent(s string) string {
	escaped := url.PathEscape(s)
	escaped = strings.ReplaceAll(escaped, "+", "%20")
	return escaped
}

type navWBIResponse struct {
	WBIImg WBIImage `json:"wbi_img"`
}

type WBIImage struct {
	ImgURL string `json:"img_url"`
	SubURL string `json:"sub_url"`
}

func (c *Client) setWBIKeysForTest(imgKey, subKey string, now time.Time) {
	c.wbi.mu.Lock()
	defer c.wbi.mu.Unlock()
	c.wbi.imgKey = imgKey
	c.wbi.subKey = subKey
	c.wbi.mixinKey = genMixinKey(imgKey + subKey)
	c.wbi.updatedAt = now
	c.wbi.now = func() time.Time { return now }
}

func rawMessage(data any) json.RawMessage {
	b, _ := json.Marshal(data)
	return b
}
