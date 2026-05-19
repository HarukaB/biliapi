package biliapi

import (
	"context"
	"encoding/json"
	"net/url"
)

type UserService struct {
	client *Client
}

type UserSpaceInfoParams struct {
	Mid      int64
	Token    string
	Platform string
}

func (p UserSpaceInfoParams) values() url.Values {
	v := url.Values{}
	setInt64(v, "mid", p.Mid)
	setString(v, "token", p.Token)
	setString(v, "platform", p.Platform)
	return v
}

type UserSpaceInfo struct {
	Mid            int64           `json:"mid"`
	Name           string          `json:"name"`
	Sex            string          `json:"sex"`
	Face           string          `json:"face"`
	FaceNFT        int             `json:"face_nft"`
	FaceNFTType    int             `json:"face_nft_type"`
	Sign           string          `json:"sign"`
	Rank           int64           `json:"rank"`
	Level          int             `json:"level"`
	Jointime       int64           `json:"jointime"`
	Moral          int             `json:"moral"`
	Silence        int             `json:"silence"`
	Coins          float64         `json:"coins"`
	FansBadge      bool            `json:"fans_badge"`
	FansMedal      json.RawMessage `json:"fans_medal"`
	Official       OfficialInfo    `json:"official"`
	Vip            VipInfo         `json:"vip"`
	Pendant        Pendant         `json:"pendant"`
	Nameplate      json.RawMessage `json:"nameplate"`
	UserHonourInfo json.RawMessage `json:"user_honour_info"`
	IsFollowed     bool            `json:"is_followed"`
	TopPhoto       string          `json:"top_photo"`
	Theme          json.RawMessage `json:"theme"`
	SysNotice      json.RawMessage `json:"sys_notice"`
	LiveRoom       json.RawMessage `json:"live_room"`
	Birthday       FlexibleString  `json:"birthday"`
	School         json.RawMessage `json:"school"`
	Profession     json.RawMessage `json:"profession"`
	Tags           json.RawMessage `json:"tags"`
	Series         json.RawMessage `json:"series"`
	IsSeniorMember int             `json:"is_senior_member"`
	McnInfo        json.RawMessage `json:"mcn_info"`
	GaiaResType    int             `json:"gaia_res_type"`
	GaiaData       json.RawMessage `json:"gaia_data"`
	IsRisk         bool            `json:"is_risk"`
	Elec           json.RawMessage `json:"elec"`
}

func (s *UserService) SpaceInfo(ctx context.Context, params UserSpaceInfoParams) (*UserSpaceInfo, error) {
	if params.Mid <= 0 {
		return nil, requirePositive("mid", params.Mid)
	}
	var out UserSpaceInfo
	err := s.client.getJSON(ctx, endpointUserSpaceInfo, params.values(), requestOptions{WBI: true}, &out)
	return &out, err
}

type UserCardParams struct {
	Mid      int64
	Photo    bool
	Relation bool
}

func (p UserCardParams) values() url.Values {
	v := url.Values{}
	setInt64(v, "mid", p.Mid)
	if p.Photo {
		v.Set("photo", "true")
	}
	if p.Relation {
		v.Set("relation", "true")
	}
	return v
}

type UserCard struct {
	Card      UserCardInfo    `json:"card"`
	Space     json.RawMessage `json:"space"`
	Following bool            `json:"following"`
	Archive   bool            `json:"archive"`
	Article   bool            `json:"article"`
	Follower  int64           `json:"follower"`
	LikeNum   int64           `json:"like_num"`
}

type UserCardInfo struct {
	Mid       string          `json:"mid"`
	Name      string          `json:"name"`
	Sex       string          `json:"sex"`
	Face      string          `json:"face"`
	Sign      string          `json:"sign"`
	Rank      string          `json:"rank"`
	LevelInfo LevelInfo       `json:"level_info"`
	Pendant   Pendant         `json:"pendant"`
	Nameplate json.RawMessage `json:"nameplate"`
	Official  OfficialInfo    `json:"official"`
	Vip       VipInfo         `json:"vip"`
	FansBadge bool            `json:"fans_badge"`
	IsDeleted int             `json:"is_deleted"`
}

func (s *UserService) Card(ctx context.Context, params UserCardParams) (*UserCard, error) {
	if params.Mid <= 0 {
		return nil, requirePositive("mid", params.Mid)
	}
	var out UserCard
	err := s.client.getJSON(ctx, endpointUserCard, params.values(), requestOptions{}, &out)
	return &out, err
}

func (s *UserService) MyInfo(ctx context.Context) (*UserSpaceInfo, error) {
	var out UserSpaceInfo
	err := s.client.getJSON(ctx, endpointUserMyInfo, nil, requestOptions{RequireLogin: true}, &out)
	return &out, err
}

type UserCardsParams struct {
	Mids []int64
}

func (p UserCardsParams) values() url.Values {
	v := url.Values{}
	setCSVInt64(v, "uids", p.Mids)
	return v
}

type UserCards map[string]UserCardInfo

func (s *UserService) Cards(ctx context.Context, params UserCardsParams) (*UserCards, error) {
	var out UserCards
	err := s.client.getJSON(ctx, endpointUserCards, params.values(), requestOptions{}, &out)
	return &out, err
}

type UserArcSearchParams struct {
	Mid     int64
	PN      int
	PS      int
	Tid     int
	Keyword string
	Order   string
}

func (p UserArcSearchParams) values() url.Values {
	v := url.Values{}
	setInt64(v, "mid", p.Mid)
	setInt(v, "pn", p.PN)
	setInt(v, "ps", p.PS)
	setInt(v, "tid", p.Tid)
	setString(v, "keyword", p.Keyword)
	setString(v, "order", p.Order)
	return v
}

type UserArcSearch struct {
	List     UserArcSearchList `json:"list"`
	Page     PageInfo          `json:"page"`
	Episodic json.RawMessage   `json:"episodic_button"`
	IsRisk   bool              `json:"is_risk"`
}

type UserArcSearchList struct {
	TList map[string]json.RawMessage `json:"tlist"`
	VList []ArchiveItem              `json:"vlist"`
}

type PageInfo struct {
	PN    int   `json:"pn"`
	PS    int   `json:"ps"`
	Count int64 `json:"count"`
}

type ArchiveItem struct {
	Comment       int64  `json:"comment"`
	TypeID        int64  `json:"typeid"`
	Play          int64  `json:"play"`
	Pic           string `json:"pic"`
	Subtitle      string `json:"subtitle"`
	Description   string `json:"description"`
	Copyright     string `json:"copyright"`
	Title         string `json:"title"`
	Review        int64  `json:"review"`
	Author        string `json:"author"`
	Mid           int64  `json:"mid"`
	Created       int64  `json:"created"`
	Length        string `json:"length"`
	VideoReview   int64  `json:"video_review"`
	AID           int64  `json:"aid"`
	BVID          string `json:"bvid"`
	HideClick     bool   `json:"hide_click"`
	IsPay         int    `json:"is_pay"`
	IsUnionVideo  int    `json:"is_union_video"`
	IsSteinsGate  int    `json:"is_steins_gate"`
	IsLivePlay    int    `json:"is_live_play"`
	IsLessonVideo int    `json:"is_lesson_video"`
	IsAvoided     int    `json:"is_avoided"`
	Attribute     int64  `json:"attribute"`
	Duration      int    `json:"duration"`
}

func (s *UserService) ArcSearch(ctx context.Context, params UserArcSearchParams) (*UserArcSearch, error) {
	if params.Mid <= 0 {
		return nil, requirePositive("mid", params.Mid)
	}
	var out UserArcSearch
	err := s.client.getJSON(ctx, endpointUserArcSearch, params.values(), requestOptions{WBI: true}, &out)
	return &out, err
}

type UserStat struct {
	Archive json.RawMessage `json:"archive"`
	Article json.RawMessage `json:"article"`
	Likes   int64           `json:"likes"`
}

func (s *UserService) UpStat(ctx context.Context, mid int64) (*UserStat, error) {
	if mid <= 0 {
		return nil, requirePositive("mid", mid)
	}
	v := url.Values{}
	setInt64(v, "mid", mid)
	var out UserStat
	err := s.client.getJSON(ctx, endpointUserUpStat, v, requestOptions{}, &out)
	return &out, err
}

type RelationStat struct {
	Mid       int64 `json:"mid"`
	Following int64 `json:"following"`
	Whisper   int64 `json:"whisper"`
	Black     int64 `json:"black"`
	Follower  int64 `json:"follower"`
}

func (s *UserService) RelationStat(ctx context.Context, vmid int64) (*RelationStat, error) {
	if vmid <= 0 {
		return nil, requirePositive("vmid", vmid)
	}
	v := url.Values{}
	setInt64(v, "vmid", vmid)
	var out RelationStat
	err := s.client.getJSON(ctx, endpointUserRelationStat, v, requestOptions{}, &out)
	return &out, err
}

type NavNum struct {
	Video     int64      `json:"video"`
	Channel   NavNumPair `json:"channel"`
	Favourite NavNumPair `json:"favourite"`
	Article   int64      `json:"article"`
	Album     int64      `json:"album"`
	Audio     int64      `json:"audio"`
	Bangumi   int64      `json:"bangumi"`
	Cinema    int64      `json:"cinema"`
	Tag       int64      `json:"tag"`
	Playlist  int64      `json:"playlist"`
	PUGV      int64      `json:"pugv"`
	UPOS      int64      `json:"upos"`
	SeasonNum int64      `json:"season_num"`
	Opus      int64      `json:"opus"`
}

type NavNumPair struct {
	Master int64 `json:"master"`
	Guest  int64 `json:"guest"`
}

func (p *NavNumPair) UnmarshalJSON(data []byte) error {
	var n flexibleInt64
	if err := json.Unmarshal(data, &n); err == nil {
		p.Master = int64(n)
		p.Guest = int64(n)
		return nil
	}
	type navNumPair NavNumPair
	var obj navNumPair
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*p = NavNumPair(obj)
	return nil
}

func (s *UserService) NavNum(ctx context.Context, mid int64) (*NavNum, error) {
	if mid <= 0 {
		return nil, requirePositive("mid", mid)
	}
	v := url.Values{}
	setInt64(v, "mid", mid)
	var out NavNum
	err := s.client.getJSON(ctx, endpointUserNavNum, v, requestOptions{}, &out)
	return &out, err
}

func (s *UserService) SetNotice(ctx context.Context, notice string) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{}
	form.Set("notice", notice)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointUserNoticeSet, form, requestOptions{RequireLogin: true}, nil)
}

func (s *UserService) SetTags(ctx context.Context, tags []string) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{}
	for _, tag := range tags {
		form.Add("tags", tag)
	}
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointUserTagsSet, form, requestOptions{RequireLogin: true}, nil)
}

func (s *UserService) SetTopArchive(ctx context.Context, aid int64) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{}
	setInt64(form, "aid", aid)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointUserTopArcSet, form, requestOptions{RequireLogin: true}, nil)
}

func (s *UserService) CancelTopArchive(ctx context.Context) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{"csrf": {s.client.creds.CSRF()}}
	return s.client.postFormJSON(ctx, endpointUserTopArcCancel, form, requestOptions{RequireLogin: true}, nil)
}
