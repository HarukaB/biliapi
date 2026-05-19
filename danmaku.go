package biliapi

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type DanmakuService struct {
	client *Client
}

type DanmakuXML struct {
	XMLName    xml.Name         `xml:"i"`
	ChatServer string           `xml:"chatserver"`
	ChatID     int64            `xml:"chatid"`
	Mission    int              `xml:"mission"`
	MaxLimit   int              `xml:"maxlimit"`
	State      int              `xml:"state"`
	RealName   int              `xml:"real_name"`
	Source     string           `xml:"source"`
	Items      []DanmakuXMLItem `xml:"d"`
}

type DanmakuXMLItem struct {
	P    string `xml:"p,attr"`
	Text string `xml:",chardata"`
}

type DanmakuXMLParams struct {
	CID int64
}

func (s *DanmakuService) XML(ctx context.Context, params DanmakuXMLParams) (*DanmakuXML, error) {
	if params.CID <= 0 {
		return nil, requirePositive("cid", params.CID)
	}
	v := url.Values{}
	setInt64(v, "oid", params.CID)
	body, err := s.client.getRaw(ctx, endpointDanmakuXMLList, v, requestOptions{})
	if err != nil {
		return nil, err
	}
	var out DanmakuXML
	if err := unmarshalDanmakuXML(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *DanmakuService) XMLByCID(ctx context.Context, cid int64) (*DanmakuXML, error) {
	if cid <= 0 {
		return nil, requirePositive("cid", cid)
	}
	body, err := s.client.getRaw(ctx, fmt.Sprintf("https://comment.bilibili.com/%d.xml", cid), nil, requestOptions{})
	if err != nil {
		return nil, err
	}
	var out DanmakuXML
	if err := unmarshalDanmakuXML(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type DanmakuSegmentParams struct {
	Type    int
	OID     int64
	PID     int64
	Segment int
}

func (p DanmakuSegmentParams) values() url.Values {
	v := url.Values{}
	setInt(v, "type", p.Type)
	setInt64(v, "oid", p.OID)
	setInt64(v, "pid", p.PID)
	setInt(v, "segment_index", p.Segment)
	return v
}

func (s *DanmakuService) Segment(ctx context.Context, params DanmakuSegmentParams) ([]byte, error) {
	return s.client.getRaw(ctx, endpointDanmakuWebSeg, params.values(), requestOptions{WBI: true})
}

type DanmakuViewParams struct {
	Type int
	OID  int64
	PID  int64
}

func (p DanmakuViewParams) values() url.Values {
	v := url.Values{}
	setInt(v, "type", p.Type)
	setInt64(v, "oid", p.OID)
	setInt64(v, "pid", p.PID)
	return v
}

type DanmakuView struct {
	State      int               `json:"state"`
	Text       string            `json:"text"`
	TextSide   string            `json:"text_side"`
	DmSge      json.RawMessage   `json:"dm_sge"`
	Flag       json.RawMessage   `json:"flag"`
	SpecialDms []string          `json:"special_dms"`
	CheckBox   bool              `json:"check_box"`
	Count      int               `json:"count"`
	CommandDms []json.RawMessage `json:"command_dms"`
	DmSetting  json.RawMessage   `json:"dm_setting"`
	ImageDms   []json.RawMessage `json:"image_dms"`
}

func (s *DanmakuService) View(ctx context.Context, params DanmakuViewParams) (*DanmakuView, error) {
	var out DanmakuView
	err := s.client.getJSON(ctx, endpointDanmakuWebView, params.values(), requestOptions{}, &out)
	return &out, err
}

type DanmakuHistoryIndexParams struct {
	Type  int
	OID   int64
	Month string
}

func (p DanmakuHistoryIndexParams) values() url.Values {
	v := url.Values{}
	setInt(v, "type", p.Type)
	setInt64(v, "oid", p.OID)
	setString(v, "month", p.Month)
	return v
}

type DanmakuHistoryIndex struct {
	Months []string `json:"months"`
}

func (s *DanmakuService) HistoryIndex(ctx context.Context, params DanmakuHistoryIndexParams) (*DanmakuHistoryIndex, error) {
	var out DanmakuHistoryIndex
	err := s.client.getJSON(ctx, endpointDanmakuHistoryIdx, params.values(), requestOptions{RequireLogin: true}, &out)
	return &out, err
}

type DanmakuHistoryParams struct {
	Type int
	OID  int64
	Date string
}

func (p DanmakuHistoryParams) values() url.Values {
	v := url.Values{}
	setInt(v, "type", p.Type)
	setInt64(v, "oid", p.OID)
	setString(v, "date", p.Date)
	return v
}

func (s *DanmakuService) HistorySegment(ctx context.Context, params DanmakuHistoryParams) ([]byte, error) {
	return s.client.getRaw(ctx, endpointDanmakuHistorySeg, params.values(), requestOptions{RequireLogin: true})
}

func (s *DanmakuService) HistoryXML(ctx context.Context, params DanmakuHistoryParams) (*DanmakuXML, error) {
	body, err := s.client.getRaw(ctx, endpointDanmakuHistoryXML, params.values(), requestOptions{RequireLogin: true})
	if err != nil {
		return nil, err
	}
	var out DanmakuXML
	if err := unmarshalDanmakuXML(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func unmarshalDanmakuXML(body []byte, out *DanmakuXML) error {
	if err := xml.Unmarshal(body, out); err == nil {
		return nil
	}
	decoded, err := io.ReadAll(transform.NewReader(bytes.NewReader(body), simplifiedchinese.GB18030.NewDecoder()))
	if err != nil {
		return err
	}
	return xml.Unmarshal(decoded, out)
}

type DanmakuThumbStatsParams struct {
	OID  int64
	IDs  []int64
	Type int
}

func (p DanmakuThumbStatsParams) values() url.Values {
	v := url.Values{}
	setInt64(v, "oid", p.OID)
	setCSVInt64(v, "ids", p.IDs)
	setInt(v, "type", p.Type)
	return v
}

type DanmakuThumbStats map[string]int64

func (s *DanmakuService) ThumbStats(ctx context.Context, params DanmakuThumbStatsParams) (*DanmakuThumbStats, error) {
	var out DanmakuThumbStats
	err := s.client.getJSON(ctx, endpointDanmakuThumbStats, params.values(), requestOptions{}, &out)
	return &out, err
}

type DanmakuPostParams struct {
	Type     int
	OID      int64
	Msg      string
	Progress int
	Color    int
	FontSize int
	Mode     int
	Pool     int
}

func (s *DanmakuService) Post(ctx context.Context, params DanmakuPostParams) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{}
	setInt(form, "type", params.Type)
	setInt64(form, "oid", params.OID)
	setString(form, "msg", params.Msg)
	setInt(form, "progress", params.Progress)
	setInt(form, "color", params.Color)
	setInt(form, "fontsize", params.FontSize)
	setInt(form, "mode", params.Mode)
	setInt(form, "pool", params.Pool)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointDanmakuPost, form, requestOptions{RequireLogin: true}, nil)
}

func (s *DanmakuService) Recall(ctx context.Context, cid, dmid int64) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{}
	setInt64(form, "cid", cid)
	setInt64(form, "dmid", dmid)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointDanmakuRecall, form, requestOptions{RequireLogin: true}, nil)
}

func (s *DanmakuService) ThumbUp(ctx context.Context, oid, dmid int64, up bool) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{}
	setInt64(form, "oid", oid)
	form.Set("dmid", strconv.FormatInt(dmid, 10))
	if up {
		form.Set("op", "1")
	} else {
		form.Set("op", "2")
	}
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointDanmakuThumbAdd, form, requestOptions{RequireLogin: true}, nil)
}
