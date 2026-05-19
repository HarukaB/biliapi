package biliapi

import (
	"context"
	"encoding/json"
	"net/url"
)

type SearchService struct {
	client *Client
}

type SearchDefault struct {
	TrackID   string `json:"trackid"`
	ShowName  string `json:"show_name"`
	Name      string `json:"name"`
	GotoType  int    `json:"goto_type"`
	GotoValue string `json:"goto_value"`
	URL       string `json:"url"`
	ExpStr    string `json:"exp_str"`
}

func (s *SearchService) Default(ctx context.Context) (*SearchDefault, error) {
	var out SearchDefault
	err := s.client.getJSON(ctx, endpointSearchDefault, nil, requestOptions{WBI: true}, &out)
	return &out, err
}

type SearchSquare struct {
	Type     string          `json:"type"`
	Title    string          `json:"title"`
	TrackID  string          `json:"trackid"`
	List     []SearchHotItem `json:"list"`
	TopList  []SearchHotItem `json:"top_list"`
	Trending SearchTrending  `json:"trending"`
	Extra    json.RawMessage `json:"extra,omitempty"`
}

type SearchTrending struct {
	Title   string          `json:"title"`
	TrackID string          `json:"trackid"`
	List    []SearchHotItem `json:"list"`
	TopList []SearchHotItem `json:"top_list"`
}

type SearchSquareParams struct {
	Limit    int
	Platform string
}

func (p SearchSquareParams) values() url.Values {
	v := url.Values{}
	setInt(v, "limit", p.Limit)
	setString(v, "platform", p.Platform)
	return v
}

type SearchHotItem struct {
	Keyword  string          `json:"keyword"`
	ShowName string          `json:"show_name"`
	Word     string          `json:"word"`
	Position int             `json:"position"`
	Pos      int             `json:"pos"`
	HotID    int64           `json:"hot_id"`
	Icon     string          `json:"icon"`
	URL      string          `json:"url"`
	Heat     int64           `json:"heat"`
	Extra    json.RawMessage `json:"extra,omitempty"`
}

func (s *SearchService) Square(ctx context.Context, params ...SearchSquareParams) (*SearchSquare, error) {
	p := SearchSquareParams{Limit: 10, Platform: "web"}
	if len(params) > 0 {
		p = params[0]
		if p.Limit == 0 {
			p.Limit = 10
		}
		if p.Platform == "" {
			p.Platform = "web"
		}
	}
	var out SearchSquare
	err := s.client.getJSON(ctx, endpointSearchSquare, p.values(), requestOptions{WBI: true}, &out)
	if out.Title == "" {
		out.Title = out.Trending.Title
	}
	if out.TrackID == "" {
		out.TrackID = out.Trending.TrackID
	}
	if len(out.List) == 0 {
		out.List = out.Trending.List
	}
	if len(out.TopList) == 0 {
		out.TopList = out.Trending.TopList
	}
	return &out, err
}

type SearchHotword struct {
	Code   int             `json:"code"`
	Result SearchHotResult `json:"result"`
	List   []SearchHotItem `json:"list"`
}

type SearchHotResult struct {
	TopList []SearchHotItem `json:"top_list"`
}

func (s *SearchService) Hotword(ctx context.Context) (*SearchHotword, error) {
	var out SearchHotword
	err := s.client.getPlainJSON(ctx, endpointSearchHotword, nil, requestOptions{}, &out)
	if len(out.Result.TopList) == 0 {
		out.Result.TopList = out.List
	}
	return &out, err
}

type SearchSuggestParams struct {
	Term      string
	MainVer   string
	Highlight bool
}

func (p SearchSuggestParams) values() url.Values {
	v := url.Values{}
	setString(v, "term", p.Term)
	setString(v, "main_ver", p.MainVer)
	setBool01(v, "highlight", p.Highlight)
	return v
}

type SearchSuggest struct {
	Result SearchSuggestResult `json:"result"`
}

type SearchSuggestResult struct {
	Tag []SearchSuggestItem `json:"tag"`
}

type SearchSuggestItem struct {
	Value    string `json:"value"`
	Term     string `json:"term"`
	Ref      int    `json:"ref"`
	Name     string `json:"name"`
	SPID     int64  `json:"spid"`
	TermType int    `json:"term_type"`
	SubType  string `json:"sub_type"`
}

func (s *SearchService) Suggest(ctx context.Context, params SearchSuggestParams) (*SearchSuggest, error) {
	var out SearchSuggest
	err := s.client.getPlainJSON(ctx, endpointSearchSuggest, params.values(), requestOptions{}, &out)
	return &out, err
}

type SearchAllParams struct {
	Keyword  string
	Page     int
	PageSize int
}

func (p SearchAllParams) values() url.Values {
	v := url.Values{}
	setString(v, "keyword", p.Keyword)
	setInt(v, "page", p.Page)
	setInt(v, "page_size", p.PageSize)
	return v
}

type SearchAll struct {
	SeID           string             `json:"seid"`
	Page           int                `json:"page"`
	PageSize       int                `json:"pagesize"`
	NumResults     int64              `json:"numResults"`
	NumPages       int                `json:"numPages"`
	SuggestKeyword string             `json:"suggest_keyword"`
	RqtType        string             `json:"rqt_type"`
	CostTime       json.RawMessage    `json:"cost_time"`
	Result         []SearchResultBand `json:"result"`
	ShowColumn     int                `json:"show_column"`
	InBlackKey     int                `json:"in_black_key"`
	InWhiteKey     int                `json:"in_white_key"`
}

type SearchResultBand struct {
	ResultType string            `json:"result_type"`
	Data       []json.RawMessage `json:"data"`
}

func (s *SearchService) All(ctx context.Context, params SearchAllParams) (*SearchAll, error) {
	var out SearchAll
	err := s.client.getJSON(ctx, endpointSearchAll, params.values(), requestOptions{WBI: true}, &out)
	return &out, err
}

type SearchType string

const (
	SearchTypeVideo   SearchType = "video"
	SearchTypeBangumi SearchType = "media_bangumi"
	SearchTypeMovie   SearchType = "media_ft"
	SearchTypeLive    SearchType = "live"
	SearchTypeArticle SearchType = "article"
	SearchTypeTopic   SearchType = "topic"
	SearchTypeUser    SearchType = "bili_user"
)

type SearchTypeParams struct {
	Keyword    string
	SearchType SearchType
	Page       int
	PageSize   int
	Order      string
	Duration   int
	TID        int
}

func (p SearchTypeParams) values() url.Values {
	v := url.Values{}
	setString(v, "keyword", p.Keyword)
	setString(v, "search_type", string(p.SearchType))
	setInt(v, "page", p.Page)
	setInt(v, "page_size", p.PageSize)
	setString(v, "order", p.Order)
	setInt(v, "duration", p.Duration)
	setInt(v, "tids", p.TID)
	return v
}

type SearchTypeResult struct {
	SeID           string            `json:"seid"`
	Page           int               `json:"page"`
	PageSize       int               `json:"pagesize"`
	NumResults     int64             `json:"numResults"`
	NumPages       int               `json:"numPages"`
	SuggestKeyword string            `json:"suggest_keyword"`
	RqtType        string            `json:"rqt_type"`
	CostTime       json.RawMessage   `json:"cost_time"`
	Result         []json.RawMessage `json:"result"`
}

func (s *SearchService) Type(ctx context.Context, params SearchTypeParams) (*SearchTypeResult, error) {
	var out SearchTypeResult
	err := s.client.getJSON(ctx, endpointSearchType, params.values(), requestOptions{WBI: true}, &out)
	return &out, err
}
