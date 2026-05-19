package biliapi

import (
	"context"
	"encoding/json"
	"net/url"
)

type MiscService struct {
	client *Client
}

type Timestamp struct {
	Now int64 `json:"now"`
}

type ReportTimestamp struct {
	Timestamp int64 `json:"timestamp"`
	Now       int64 `json:"now"`
}

type BuvidInfo struct {
	B3   string `json:"b_3"`
	B4   string `json:"b_4"`
	BNut string `json:"b_nut"`
}

type Buvid3Info struct {
	Buvid string `json:"buvid"`
}

type RTCTimestamp struct {
	Timestamp int64 `json:"timestamp"`
}

type MathJaxParams struct {
	Tex string
}

func (p MathJaxParams) values() url.Values {
	v := url.Values{}
	setString(v, "tex", p.Tex)
	return v
}

type MathJaxResult struct {
	SVG string `json:"svg"`
}

type ShortLinkInfo struct {
	URL         string          `json:"url"`
	Content     string          `json:"content"`
	Title       string          `json:"title"`
	RedirectURL string          `json:"redirect_url"`
	Extra       json.RawMessage `json:"extra,omitempty"`
}

func (s *MiscService) Now(ctx context.Context) (*Timestamp, error) {
	var out Timestamp
	err := s.client.getJSON(ctx, endpointNow, nil, requestOptions{}, &out)
	return &out, err
}

func (s *MiscService) ReportNow(ctx context.Context) (*ReportTimestamp, error) {
	var out ReportTimestamp
	err := s.client.getJSON(ctx, endpointReportNow, nil, requestOptions{}, &out)
	return &out, err
}

func (s *MiscService) Buvid(ctx context.Context) (*BuvidInfo, error) {
	var out BuvidInfo
	err := s.client.getJSON(ctx, endpointBuvid, nil, requestOptions{}, &out)
	return &out, err
}

func (s *MiscService) Buvid3(ctx context.Context) (*Buvid3Info, error) {
	var out Buvid3Info
	err := s.client.getJSON(ctx, endpointGetBuvid, nil, requestOptions{}, &out)
	return &out, err
}

func (s *MiscService) RTCTimestamp(ctx context.Context) (*RTCTimestamp, error) {
	var out RTCTimestamp
	err := s.client.getJSON(ctx, endpointRTCTime, nil, requestOptions{}, &out)
	return &out, err
}

func (s *MiscService) ServerDateJS(ctx context.Context) ([]byte, error) {
	return s.client.getRaw(ctx, endpointServerDate, nil, requestOptions{})
}

func (s *MiscService) MathJax(ctx context.Context, params MathJaxParams) (*MathJaxResult, error) {
	var out MathJaxResult
	err := s.client.getJSON(ctx, endpointMathJax, params.values(), requestOptions{}, &out)
	return &out, err
}

type ShortLinkParams struct {
	Build int
	BVID  string
	AID   int64
	URL   string
}

func (p ShortLinkParams) values() url.Values {
	v := url.Values{}
	setInt(v, "build", p.Build)
	setString(v, "bvid", p.BVID)
	setInt64(v, "aid", p.AID)
	setString(v, "url", p.URL)
	return v
}

func (s *MiscService) ShortLink(ctx context.Context, params ShortLinkParams) (*ShortLinkInfo, error) {
	var out ShortLinkInfo
	err := s.client.getJSON(ctx, endpointShortLink, params.values(), requestOptions{}, &out)
	return &out, err
}
