package biliapi

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

type CommentService struct {
	client *Client
}

type CommentType int

const (
	CommentTypeVideo   CommentType = 1
	CommentTypeArticle CommentType = 12
	CommentTypeDynamic CommentType = 17
)

type CommentListParams struct {
	OID   int64
	Type  CommentType
	PN    int
	PS    int
	Sort  int
	NoHot bool
}

func (p CommentListParams) values() url.Values {
	v := url.Values{}
	setInt64(v, "oid", p.OID)
	setInt(v, "type", int(p.Type))
	setInt(v, "pn", p.PN)
	setInt(v, "ps", p.PS)
	setInt(v, "sort", p.Sort)
	setBool01(v, "nohot", p.NoHot)
	return v
}

type CommentMainParams struct {
	OID      int64
	Type     CommentType
	Mode     int
	Next     int64
	PS       int
	SeekRpid int64
}

func (p CommentMainParams) values() url.Values {
	v := url.Values{}
	setInt64(v, "oid", p.OID)
	setInt(v, "type", int(p.Type))
	setInt(v, "mode", p.Mode)
	setInt64(v, "next", p.Next)
	setInt(v, "ps", p.PS)
	setInt64(v, "seek_rpid", p.SeekRpid)
	return v
}

type CommentList struct {
	Page    CommentPage     `json:"page"`
	Config  CommentConfig   `json:"config"`
	Upper   CommentUpper    `json:"upper"`
	Replies []CommentReply  `json:"replies"`
	Hots    []CommentReply  `json:"hots"`
	Notice  json.RawMessage `json:"notice"`
	Folder  json.RawMessage `json:"folder"`
	Lottery json.RawMessage `json:"lottery"`
}

type CommentMain struct {
	Cursor     CommentCursor   `json:"cursor"`
	Replies    []CommentReply  `json:"replies"`
	Top        json.RawMessage `json:"top"`
	TopReplies []CommentReply  `json:"top_replies"`
	Upper      CommentUpper    `json:"upper"`
	Config     CommentConfig   `json:"config"`
	Control    json.RawMessage `json:"control"`
	Note       int             `json:"note"`
}

type CommentPage struct {
	Num    int   `json:"num"`
	Size   int   `json:"size"`
	Count  int64 `json:"count"`
	ACount int64 `json:"acount"`
}

type CommentCursor struct {
	IsBegin bool   `json:"is_begin"`
	Prev    int64  `json:"prev"`
	Next    int64  `json:"next"`
	IsEnd   bool   `json:"is_end"`
	Mode    int    `json:"mode"`
	Name    string `json:"name"`
}

type CommentConfig struct {
	ShowAdmin        int  `json:"showadmin"`
	ShowEntry        int  `json:"showentry"`
	ShowFloor        int  `json:"showfloor"`
	ShowTopic        int  `json:"showtopic"`
	ShowUpFlag       bool `json:"show_up_flag"`
	ReadOnly         bool `json:"read_only"`
	ShowDelLog       bool `json:"show_del_log"`
	ShowBvid         bool `json:"show_bvid"`
	DisableJumpEmote bool `json:"disable_jump_emote"`
}

type CommentUpper struct {
	Mid  int64           `json:"mid"`
	Top  json.RawMessage `json:"top"`
	Vote json.RawMessage `json:"vote"`
}

type CommentReply struct {
	RPID         int64           `json:"rpid"`
	OID          int64           `json:"oid"`
	Type         int             `json:"type"`
	Mid          int64           `json:"mid"`
	Root         int64           `json:"root"`
	Parent       int64           `json:"parent"`
	Dialog       int64           `json:"dialog"`
	Count        int             `json:"count"`
	RCount       int             `json:"rcount"`
	Floor        int             `json:"floor"`
	State        int             `json:"state"`
	FansGrade    int             `json:"fansgrade"`
	Attr         int64           `json:"attr"`
	CTime        int64           `json:"ctime"`
	RPIDStr      string          `json:"rpid_str"`
	RootStr      string          `json:"root_str"`
	ParentStr    string          `json:"parent_str"`
	Like         int64           `json:"like"`
	Action       int             `json:"action"`
	Member       CommentMember   `json:"member"`
	Content      CommentContent  `json:"content"`
	Replies      []CommentReply  `json:"replies"`
	Assist       int             `json:"assist"`
	Folder       json.RawMessage `json:"folder"`
	UpAction     json.RawMessage `json:"up_action"`
	Invisible    bool            `json:"invisible"`
	ReplyControl json.RawMessage `json:"reply_control"`
}

type CommentMember struct {
	Mid            string          `json:"mid"`
	Uname          string          `json:"uname"`
	Sex            string          `json:"sex"`
	Sign           string          `json:"sign"`
	Avatar         string          `json:"avatar"`
	Rank           string          `json:"rank"`
	LevelInfo      LevelInfo       `json:"level_info"`
	Pendant        Pendant         `json:"pendant"`
	Nameplate      json.RawMessage `json:"nameplate"`
	OfficialVerify OfficialInfo    `json:"official_verify"`
	Vip            VipInfo         `json:"vip"`
	FansDetail     json.RawMessage `json:"fans_detail"`
}

type CommentContent struct {
	Message  string                     `json:"message"`
	Plat     int                        `json:"plat"`
	Device   string                     `json:"device"`
	Members  []BasicUser                `json:"members"`
	Emote    map[string]json.RawMessage `json:"emote"`
	JumpURL  map[string]json.RawMessage `json:"jump_url"`
	MaxLine  int                        `json:"max_line"`
	Pictures []json.RawMessage          `json:"pictures"`
}

func (s *CommentService) List(ctx context.Context, params CommentListParams) (*CommentList, error) {
	var out CommentList
	err := s.client.getJSON(ctx, endpointCommentList, params.values(), requestOptions{}, &out)
	return &out, err
}

func (s *CommentService) Main(ctx context.Context, params CommentMainParams) (*CommentMain, error) {
	var out CommentMain
	err := s.client.getJSON(ctx, endpointCommentMain, params.values(), requestOptions{WBI: true}, &out)
	return &out, err
}

type CommentRepliesParams struct {
	OID  int64
	Type CommentType
	Root int64
	PN   int
	PS   int
}

func (p CommentRepliesParams) values() url.Values {
	v := url.Values{}
	setInt64(v, "oid", p.OID)
	setInt(v, "type", int(p.Type))
	setInt64(v, "root", p.Root)
	setInt(v, "pn", p.PN)
	setInt(v, "ps", p.PS)
	return v
}

type CommentReplies struct {
	Page    CommentPage    `json:"page"`
	Replies []CommentReply `json:"replies"`
	Root    CommentReply   `json:"root"`
}

func (s *CommentService) Replies(ctx context.Context, params CommentRepliesParams) (*CommentReplies, error) {
	var out CommentReplies
	err := s.client.getJSON(ctx, endpointCommentReply, params.values(), requestOptions{}, &out)
	return &out, err
}

type CommentDialogParams struct {
	OID    int64
	Type   CommentType
	Root   int64
	Dialog int64
	PN     int
	PS     int
}

func (p CommentDialogParams) values() url.Values {
	v := url.Values{}
	setInt64(v, "oid", p.OID)
	setInt(v, "type", int(p.Type))
	setInt64(v, "root", p.Root)
	setInt64(v, "dialog", p.Dialog)
	setInt(v, "pn", p.PN)
	setInt(v, "ps", p.PS)
	return v
}

func (s *CommentService) Dialog(ctx context.Context, params CommentDialogParams) (*CommentReplies, error) {
	var out CommentReplies
	err := s.client.getJSON(ctx, endpointCommentDialog, params.values(), requestOptions{}, &out)
	return &out, err
}

func (s *CommentService) Hot(ctx context.Context, params CommentListParams) (*CommentList, error) {
	var out CommentList
	err := s.client.getJSON(ctx, endpointCommentHot, params.values(), requestOptions{}, &out)
	return &out, err
}

func (s *CommentService) Info(ctx context.Context, oid int64, typ CommentType, rpid int64) (*CommentReply, error) {
	v := url.Values{}
	setInt64(v, "oid", oid)
	setInt(v, "type", int(typ))
	setInt64(v, "rpid", rpid)
	var out struct {
		Reply *CommentReply `json:"reply"`
		CommentReply
	}
	err := s.client.getJSON(ctx, endpointCommentInfo, v, requestOptions{}, &out)
	if out.Reply != nil {
		return out.Reply, err
	}
	return &out.CommentReply, err
}

type CommentCount struct {
	Count int64 `json:"count"`
}

func (s *CommentService) Count(ctx context.Context, oid int64, typ CommentType) (*CommentCount, error) {
	v := url.Values{}
	setInt64(v, "oid", oid)
	setInt(v, "type", int(typ))
	var out CommentCount
	err := s.client.getJSON(ctx, endpointCommentCount, v, requestOptions{}, &out)
	return &out, err
}

type CommentAddParams struct {
	OID     int64
	Type    CommentType
	Message string
	Root    int64
	Parent  int64
}

func (s *CommentService) Add(ctx context.Context, params CommentAddParams) (*CommentReply, error) {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return nil, err
	}
	form := url.Values{}
	setInt64(form, "oid", params.OID)
	setInt(form, "type", int(params.Type))
	setString(form, "message", params.Message)
	setInt64(form, "root", params.Root)
	setInt64(form, "parent", params.Parent)
	form.Set("csrf", s.client.creds.CSRF())
	var out CommentReply
	err := s.client.postFormJSON(ctx, endpointCommentAdd, form, requestOptions{RequireLogin: true}, &out)
	return &out, err
}

type CommentActionParams struct {
	OID    int64
	Type   CommentType
	RPID   int64
	Action int
}

func (p CommentActionParams) form(csrf string) url.Values {
	form := url.Values{}
	setInt64(form, "oid", p.OID)
	setInt(form, "type", int(p.Type))
	setInt64(form, "rpid", p.RPID)
	setInt(form, "action", p.Action)
	form.Set("csrf", csrf)
	return form
}

func (s *CommentService) Action(ctx context.Context, params CommentActionParams) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	return s.client.postFormJSON(ctx, endpointCommentAction, params.form(s.client.creds.CSRF()), requestOptions{RequireLogin: true}, nil)
}

func (s *CommentService) Hate(ctx context.Context, params CommentActionParams) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	return s.client.postFormJSON(ctx, endpointCommentHate, params.form(s.client.creds.CSRF()), requestOptions{RequireLogin: true}, nil)
}

func (s *CommentService) Delete(ctx context.Context, oid int64, typ CommentType, rpid int64) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{}
	setInt64(form, "oid", oid)
	setInt(form, "type", int(typ))
	setInt64(form, "rpid", rpid)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointCommentDel, form, requestOptions{RequireLogin: true}, nil)
}

func (s *CommentService) Top(ctx context.Context, oid int64, typ CommentType, rpid int64, action int) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{}
	setInt64(form, "oid", oid)
	setInt(form, "type", int(typ))
	setInt64(form, "rpid", rpid)
	setInt(form, "action", action)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointCommentTop, form, requestOptions{RequireLogin: true}, nil)
}

func (s *CommentService) Report(ctx context.Context, oid int64, typ CommentType, rpid int64, reason int, content string) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{}
	setInt64(form, "oid", oid)
	setInt(form, "type", int(typ))
	setInt64(form, "rpid", rpid)
	setInt(form, "reason", reason)
	setString(form, "content", content)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointCommentReport, form, requestOptions{RequireLogin: true}, nil)
}

func rpidString(rpid int64) string {
	return strconv.FormatInt(rpid, 10)
}
