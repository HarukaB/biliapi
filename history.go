package biliapi

import (
	"context"
	"encoding/json"
	"net/url"
)

type HistoryService struct {
	client *Client
}

type HistoryCursorParams struct {
	Max      int64
	Business string
	ViewAt   int64
	Type     string
	PS       int
}

func (p HistoryCursorParams) values() url.Values {
	v := url.Values{}
	setInt64(v, "max", p.Max)
	setString(v, "business", p.Business)
	setInt64(v, "view_at", p.ViewAt)
	setString(v, "type", p.Type)
	setInt(v, "ps", p.PS)
	return v
}

type HistoryCursor struct {
	Cursor HistoryCursorPage `json:"cursor"`
	List   []HistoryItem     `json:"list"`
	Tab    []HistoryTab      `json:"tab"`
}

type HistoryCursorPage struct {
	Max      int64  `json:"max"`
	ViewAt   int64  `json:"view_at"`
	Business string `json:"business"`
	PS       int    `json:"ps"`
}

type HistoryTab struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type HistoryItem struct {
	Title      string          `json:"title"`
	LongTitle  string          `json:"long_title"`
	Cover      string          `json:"cover"`
	Covers     []string        `json:"covers"`
	URI        string          `json:"uri"`
	History    HistoryRef      `json:"history"`
	Videos     int             `json:"videos"`
	AuthorName string          `json:"author_name"`
	AuthorFace string          `json:"author_face"`
	AuthorMID  int64           `json:"author_mid"`
	ViewAt     int64           `json:"view_at"`
	Progress   int64           `json:"progress"`
	Badge      string          `json:"badge"`
	ShowTitle  string          `json:"show_title"`
	Duration   int64           `json:"duration"`
	Current    string          `json:"current"`
	Total      int             `json:"total"`
	NewDesc    string          `json:"new_desc"`
	IsFinish   int             `json:"is_finish"`
	IsFav      int             `json:"is_fav"`
	KID        int64           `json:"kid"`
	TagName    string          `json:"tag_name"`
	LiveStatus int             `json:"live_status"`
	Extra      json.RawMessage `json:"extra,omitempty"`
}

type HistoryRef struct {
	OID      int64  `json:"oid"`
	Epid     int64  `json:"epid"`
	BVID     string `json:"bvid"`
	Page     int    `json:"page"`
	CID      int64  `json:"cid"`
	Part     string `json:"part"`
	Business string `json:"business"`
	DT       int    `json:"dt"`
}

func (s *HistoryService) Cursor(ctx context.Context, params HistoryCursorParams) (*HistoryCursor, error) {
	var out HistoryCursor
	err := s.client.getJSON(ctx, endpointHistoryCursor, params.values(), requestOptions{RequireLogin: true}, &out)
	return &out, err
}

type HistoryLegacyParams struct {
	PN int
	PS int
}

func (p HistoryLegacyParams) values() url.Values {
	v := url.Values{}
	setInt(v, "pn", p.PN)
	setInt(v, "ps", p.PS)
	return v
}

type HistoryLegacy struct {
	List []HistoryItem `json:"list"`
}

func (h *HistoryLegacy) UnmarshalJSON(data []byte) error {
	var list []HistoryItem
	if err := json.Unmarshal(data, &list); err == nil {
		h.List = list
		return nil
	}
	type historyLegacy HistoryLegacy
	var obj historyLegacy
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*h = HistoryLegacy(obj)
	return nil
}

func (s *HistoryService) Legacy(ctx context.Context, params HistoryLegacyParams) (*HistoryLegacy, error) {
	var out HistoryLegacy
	err := s.client.getJSON(ctx, endpointHistoryLegacy, params.values(), requestOptions{RequireLogin: true}, &out)
	return &out, err
}

type ShadowStatus struct {
	Shadow bool `json:"shadow"`
}

func (s *ShadowStatus) UnmarshalJSON(data []byte) error {
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		s.Shadow = b
		return nil
	}
	type shadowStatus ShadowStatus
	var obj shadowStatus
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*s = ShadowStatus(obj)
	return nil
}

func (s *HistoryService) Shadow(ctx context.Context) (*ShadowStatus, error) {
	var out ShadowStatus
	err := s.client.getJSON(ctx, endpointHistoryShadow, nil, requestOptions{RequireLogin: true}, &out)
	return &out, err
}

type ToViewList struct {
	Count int64        `json:"count"`
	List  []ToViewItem `json:"list"`
}

type ToViewItem struct {
	AID       int64         `json:"aid"`
	Videos    int           `json:"videos"`
	TID       int64         `json:"tid"`
	TName     string        `json:"tname"`
	Copyright int           `json:"copyright"`
	Pic       string        `json:"pic"`
	Title     string        `json:"title"`
	PubDate   int64         `json:"pubdate"`
	CTime     int64         `json:"ctime"`
	Desc      string        `json:"desc"`
	State     int           `json:"state"`
	Duration  int           `json:"duration"`
	Rights    ArchiveRights `json:"rights"`
	Owner     Owner         `json:"owner"`
	Stat      ArchiveStat   `json:"stat"`
	Dynamic   string        `json:"dynamic"`
	CID       int64         `json:"cid"`
	Dimension Dimension     `json:"dimension"`
	BVID      string        `json:"bvid"`
	Viewed    bool          `json:"viewed"`
	Pages     []Page        `json:"pages"`
}

func (s *HistoryService) ToView(ctx context.Context) (*ToViewList, error) {
	var out ToViewList
	err := s.client.getJSON(ctx, endpointHistoryToView, nil, requestOptions{RequireLogin: true}, &out)
	return &out, err
}

type HistoryDeleteItem struct {
	Business string `json:"business"`
	ID       int64  `json:"id"`
}

func (s *HistoryService) Delete(ctx context.Context, items []HistoryDeleteItem) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	body, _ := json.Marshal(items)
	form := url.Values{}
	form.Set("kid", string(body))
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointHistoryDelete, form, requestOptions{RequireLogin: true}, nil)
}

func (s *HistoryService) Clear(ctx context.Context) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{"csrf": {s.client.creds.CSRF()}}
	return s.client.postFormJSON(ctx, endpointHistoryClear, form, requestOptions{RequireLogin: true}, nil)
}

func (s *HistoryService) SetShadow(ctx context.Context, shadow bool) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{}
	if shadow {
		form.Set("switch", "1")
	} else {
		form.Set("switch", "0")
	}
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointHistoryShadowSet, form, requestOptions{RequireLogin: true}, nil)
}

func (s *HistoryService) AddToView(ctx context.Context, aid int64) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{}
	setInt64(form, "aid", aid)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointHistoryToViewAdd, form, requestOptions{RequireLogin: true}, nil)
}

func (s *HistoryService) DeleteToView(ctx context.Context, aid int64) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{}
	setInt64(form, "aid", aid)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointHistoryToViewDel, form, requestOptions{RequireLogin: true}, nil)
}

func (s *HistoryService) ClearToView(ctx context.Context) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{"csrf": {s.client.creds.CSRF()}}
	return s.client.postFormJSON(ctx, endpointHistoryToViewClear, form, requestOptions{RequireLogin: true}, nil)
}
