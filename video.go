package biliapi

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

type VideoService struct {
	client *Client
}

type VideoViewParams struct {
	AID  int64
	BVID string
}

func (p VideoViewParams) values() url.Values {
	v := url.Values{}
	setInt64(v, "aid", p.AID)
	setString(v, "bvid", p.BVID)
	return v
}

type VideoView struct {
	BVID               string          `json:"bvid"`
	AID                int64           `json:"aid"`
	Videos             int             `json:"videos"`
	TID                int64           `json:"tid"`
	TIDV2              int64           `json:"tid_v2"`
	TName              string          `json:"tname"`
	TNameV2            string          `json:"tname_v2"`
	Copyright          int             `json:"copyright"`
	Pic                string          `json:"pic"`
	Title              string          `json:"title"`
	PubDate            int64           `json:"pubdate"`
	CTime              int64           `json:"ctime"`
	Desc               string          `json:"desc"`
	DescV2             []DescV2Item    `json:"desc_v2"`
	State              int             `json:"state"`
	Duration           int             `json:"duration"`
	Forward            int64           `json:"forward"`
	MissionID          int64           `json:"mission_id"`
	RedirectURL        string          `json:"redirect_url"`
	Rights             ArchiveRights   `json:"rights"`
	Owner              Owner           `json:"owner"`
	Stat               ArchiveStat     `json:"stat"`
	ArgueInfo          ArgueInfo       `json:"argue_info"`
	Dynamic            string          `json:"dynamic"`
	CID                int64           `json:"cid"`
	Dimension          Dimension       `json:"dimension"`
	Premiere           json.RawMessage `json:"premiere"`
	TeenageMode        int             `json:"teenage_mode"`
	IsChargeableSeason bool            `json:"is_chargeable_season"`
	IsStory            bool            `json:"is_story"`
	IsUpowerExclusive  bool            `json:"is_upower_exclusive"`
	IsUpowerPlay       bool            `json:"is_upower_play"`
	IsUpowerPreview    bool            `json:"is_upower_preview"`
	NoCache            bool            `json:"no_cache"`
	Pages              []Page          `json:"pages"`
	Subtitle           Subtitle        `json:"subtitle"`
	UGCSeason          json.RawMessage `json:"ugc_season"`
	Staff              json.RawMessage `json:"staff"`
	IsSeasonDisplay    bool            `json:"is_season_display"`
	UserGarb           json.RawMessage `json:"user_garb"`
	HonorReply         json.RawMessage `json:"honor_reply"`
	LikeIcon           string          `json:"like_icon"`
	NeedJumpBV         bool            `json:"need_jump_bv"`
	DisableShowUpInfo  bool            `json:"disable_show_up_info"`
	IsStoryPlay        int             `json:"is_story_play"`
	IsViewSelf         bool            `json:"is_view_self"`
}

func (s *VideoService) View(ctx context.Context, params VideoViewParams) (*VideoView, error) {
	if err := requireAIDOrBVID(params.AID, params.BVID); err != nil {
		return nil, err
	}
	var out VideoView
	err := s.client.getJSON(ctx, endpointVideoView, params.values(), requestOptions{WBI: true}, &out)
	return &out, err
}

type VideoDetail struct {
	View      VideoView       `json:"View"`
	Card      json.RawMessage `json:"Card"`
	Tags      json.RawMessage `json:"Tags"`
	Reply     json.RawMessage `json:"Reply"`
	Related   []VideoView     `json:"Related"`
	Spec      json.RawMessage `json:"Spec"`
	HotShare  json.RawMessage `json:"hot_share"`
	Emergency json.RawMessage `json:"emergency"`
}

func (s *VideoService) ViewDetail(ctx context.Context, params VideoViewParams) (*VideoDetail, error) {
	if err := requireAIDOrBVID(params.AID, params.BVID); err != nil {
		return nil, err
	}
	var out VideoDetail
	err := s.client.getJSON(ctx, endpointVideoViewDetail, params.values(), requestOptions{WBI: true}, &out)
	return &out, err
}

func (s *VideoService) Desc(ctx context.Context, params VideoViewParams) (string, error) {
	if err := requireAIDOrBVID(params.AID, params.BVID); err != nil {
		return "", err
	}
	var out string
	err := s.client.getJSON(ctx, endpointVideoDesc, params.values(), requestOptions{}, &out)
	return out, err
}

func (s *VideoService) PageList(ctx context.Context, params VideoViewParams) ([]Page, error) {
	if err := requireAIDOrBVID(params.AID, params.BVID); err != nil {
		return nil, err
	}
	var out []Page
	err := s.client.getJSON(ctx, endpointVideoPageList, params.values(), requestOptions{}, &out)
	return out, err
}

type PlayURLParams struct {
	AID         int64
	BVID        string
	CID         int64
	QN          int
	FnVer       int
	FnVal       int
	FourK       bool
	Platform    string
	HighQuality bool
}

func (p PlayURLParams) values() url.Values {
	v := url.Values{}
	setInt64(v, "avid", p.AID)
	setString(v, "bvid", p.BVID)
	setInt64(v, "cid", p.CID)
	setInt(v, "qn", p.QN)
	setInt(v, "fnver", p.FnVer)
	setInt(v, "fnval", p.FnVal)
	setBool01(v, "fourk", p.FourK)
	setString(v, "platform", p.Platform)
	setBool01(v, "high_quality", p.HighQuality)
	return v
}

type PlayURL struct {
	From              string          `json:"from"`
	Result            string          `json:"result"`
	Message           string          `json:"message"`
	Quality           int             `json:"quality"`
	Format            string          `json:"format"`
	Timelength        int64           `json:"timelength"`
	AcceptFormat      string          `json:"accept_format"`
	AcceptDescription []string        `json:"accept_description"`
	AcceptQuality     []int           `json:"accept_quality"`
	VideoCodecid      int             `json:"video_codecid"`
	SeekParam         string          `json:"seek_param"`
	SeekType          string          `json:"seek_type"`
	DURL              []PlayDURL      `json:"durl"`
	Dash              *Dash           `json:"dash"`
	SupportFormats    []SupportFormat `json:"support_formats"`
	LastPlayTime      int64           `json:"last_play_time"`
	LastPlayCID       int64           `json:"last_play_cid"`
}

type PlayDURL struct {
	Order     int      `json:"order"`
	Length    int64    `json:"length"`
	Size      int64    `json:"size"`
	AHead     string   `json:"ahead"`
	VHead     string   `json:"vhead"`
	URL       string   `json:"url"`
	BackupURL []string `json:"backup_url"`
}

type Dash struct {
	Duration         int64           `json:"duration"`
	MinBufferTime    float64         `json:"minBufferTime"`
	MinBufferTimeAlt float64         `json:"min_buffer_time"`
	Video            []DashMedia     `json:"video"`
	Audio            []DashMedia     `json:"audio"`
	Dolby            json.RawMessage `json:"dolby"`
	Flac             json.RawMessage `json:"flac"`
}

type DashMedia struct {
	ID              int         `json:"id"`
	BaseURL         string      `json:"baseUrl"`
	BaseURLAlt      string      `json:"base_url"`
	BackupURL       []string    `json:"backupUrl"`
	BackupURLAlt    []string    `json:"backup_url"`
	Bandwidth       int64       `json:"bandwidth"`
	MimeType        string      `json:"mimeType"`
	MimeTypeAlt     string      `json:"mime_type"`
	Codecs          string      `json:"codecs"`
	Width           int         `json:"width"`
	Height          int         `json:"height"`
	FrameRate       string      `json:"frameRate"`
	FrameRateAlt    string      `json:"frame_rate"`
	Sar             string      `json:"sar"`
	StartWithSAP    int         `json:"startWithSap"`
	StartWithSAPAlt int         `json:"start_with_sap"`
	SegmentBase     SegmentBase `json:"SegmentBase"`
	SegmentBaseAlt  SegmentBase `json:"segment_base"`
	Codecid         int         `json:"codecid"`
}

type SegmentBase struct {
	Initialization string `json:"Initialization"`
	IndexRange     string `json:"indexRange"`
}

type SupportFormat struct {
	Quality        int      `json:"quality"`
	Format         string   `json:"format"`
	NewDescription string   `json:"new_description"`
	DisplayDesc    string   `json:"display_desc"`
	Superscript    string   `json:"superscript"`
	Codecs         []string `json:"codecs"`
}

func (s *VideoService) PlayURL(ctx context.Context, params PlayURLParams) (*PlayURL, error) {
	if err := requireAIDOrBVID(params.AID, params.BVID); err != nil {
		return nil, err
	}
	if params.CID <= 0 {
		return nil, requirePositive("cid", params.CID)
	}
	var out PlayURL
	err := s.client.getJSON(ctx, endpointVideoPlayURL, params.values(), requestOptions{WBI: true}, &out)
	return &out, err
}

type PlayerInfoParams struct {
	AID  int64
	BVID string
	CID  int64
}

func (p PlayerInfoParams) values() url.Values {
	v := url.Values{}
	setInt64(v, "aid", p.AID)
	setString(v, "bvid", p.BVID)
	setInt64(v, "cid", p.CID)
	return v
}

type PlayerInfo struct {
	AID          int64           `json:"aid"`
	BVID         string          `json:"bvid"`
	AllowBP      bool            `json:"allow_bp"`
	NoShare      bool            `json:"no_share"`
	CID          int64           `json:"cid"`
	MaxLimit     int             `json:"max_limit"`
	PageNo       int             `json:"page_no"`
	HasNext      bool            `json:"has_next"`
	IPInfo       json.RawMessage `json:"ip_info"`
	LoginMid     int64           `json:"login_mid"`
	LoginMidHash string          `json:"login_mid_hash"`
	IsOwner      bool            `json:"is_owner"`
	Name         string          `json:"name"`
	Permission   string          `json:"permission"`
	LevelInfo    LevelInfo       `json:"level_info"`
	Vip          VipInfo         `json:"vip"`
	AnswerStatus int             `json:"answer_status"`
	Subtitle     Subtitle        `json:"subtitle"`
	ViewPoints   json.RawMessage `json:"view_points"`
}

func (s *VideoService) PlayerInfo(ctx context.Context, params PlayerInfoParams) (*PlayerInfo, error) {
	if err := requireAIDOrBVID(params.AID, params.BVID); err != nil {
		return nil, err
	}
	if params.CID <= 0 {
		return nil, requirePositive("cid", params.CID)
	}
	var out PlayerInfo
	err := s.client.getJSON(ctx, endpointVideoPlayerV2, params.values(), requestOptions{WBI: true}, &out)
	return &out, err
}

func (s *VideoService) Stat(ctx context.Context, params VideoViewParams) (*ArchiveStat, error) {
	if err := requireAIDOrBVID(params.AID, params.BVID); err != nil {
		return nil, err
	}
	var out ArchiveStat
	err := s.client.getJSON(ctx, endpointVideoStat, params.values(), requestOptions{}, &out)
	return &out, err
}

type VideoHasLike struct {
	Like bool `json:"like"`
}

func (v *VideoHasLike) UnmarshalJSON(data []byte) error {
	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		v.Like = n != 0
		return nil
	}
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		v.Like = b
		return nil
	}
	var obj struct {
		Like json.RawMessage `json:"like"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	if len(obj.Like) == 0 {
		return nil
	}
	if err := json.Unmarshal(obj.Like, &n); err == nil {
		v.Like = n != 0
		return nil
	}
	if err := json.Unmarshal(obj.Like, &b); err != nil {
		return err
	}
	v.Like = b
	return nil
}

func (s *VideoService) HasLike(ctx context.Context, params VideoViewParams) (*VideoHasLike, error) {
	if err := requireAIDOrBVID(params.AID, params.BVID); err != nil {
		return nil, err
	}
	var out VideoHasLike
	err := s.client.getJSON(ctx, endpointVideoHasLike, params.values(), requestOptions{RequireLogin: true}, &out)
	return &out, err
}

type VideoCoins struct {
	Multiply int `json:"multiply"`
}

func (s *VideoService) Coins(ctx context.Context, params VideoViewParams) (*VideoCoins, error) {
	if err := requireAIDOrBVID(params.AID, params.BVID); err != nil {
		return nil, err
	}
	var out VideoCoins
	err := s.client.getJSON(ctx, endpointVideoCoins, params.values(), requestOptions{RequireLogin: true}, &out)
	return &out, err
}

type VideoFavoured struct {
	Count    int64 `json:"count"`
	Favoured bool  `json:"favoured"`
}

func (s *VideoService) Favoured(ctx context.Context, aid int64) (*VideoFavoured, error) {
	if aid <= 0 {
		return nil, requirePositive("aid", aid)
	}
	v := url.Values{}
	setInt64(v, "aid", aid)
	var out VideoFavoured
	err := s.client.getJSON(ctx, endpointVideoFavoured, v, requestOptions{RequireLogin: true}, &out)
	return &out, err
}

func (s *VideoService) Related(ctx context.Context, params VideoViewParams) ([]VideoView, error) {
	if err := requireAIDOrBVID(params.AID, params.BVID); err != nil {
		return nil, err
	}
	var out []VideoView
	err := s.client.getJSON(ctx, endpointVideoRelated, params.values(), requestOptions{}, &out)
	return out, err
}

type OnlineTotal struct {
	Total      string          `json:"total"`
	Count      string          `json:"count"`
	ShowSwitch json.RawMessage `json:"show_switch"`
}

func (s *VideoService) OnlineTotal(ctx context.Context, aid, cid int64, bvid string) (*OnlineTotal, error) {
	if err := requireAIDOrBVID(aid, bvid); err != nil {
		return nil, err
	}
	if cid <= 0 {
		return nil, requirePositive("cid", cid)
	}
	v := url.Values{}
	setInt64(v, "aid", aid)
	setString(v, "bvid", bvid)
	setInt64(v, "cid", cid)
	var out OnlineTotal
	err := s.client.getJSON(ctx, endpointVideoOnline, v, requestOptions{}, &out)
	return &out, err
}

type VideoShot struct {
	PVData   string   `json:"pvdata"`
	ImgXLen  int      `json:"img_x_len"`
	ImgYLen  int      `json:"img_y_len"`
	ImgXSize int      `json:"img_x_size"`
	ImgYSize int      `json:"img_y_size"`
	Image    []string `json:"image"`
	Index    []int64  `json:"index"`
}

func (s *VideoService) Shot(ctx context.Context, aid, cid int64, bvid string) (*VideoShot, error) {
	if err := requireAIDOrBVID(aid, bvid); err != nil {
		return nil, err
	}
	if cid <= 0 {
		return nil, requirePositive("cid", cid)
	}
	v := url.Values{}
	setInt64(v, "aid", aid)
	setString(v, "bvid", bvid)
	setInt64(v, "cid", cid)
	var out VideoShot
	err := s.client.getJSON(ctx, endpointVideoShot, v, requestOptions{}, &out)
	return &out, err
}

type AIConclusion struct {
	ModelResult json.RawMessage `json:"model_result"`
	Stid        string          `json:"stid"`
	Status      int             `json:"status"`
	LikeNum     int64           `json:"like_num"`
	DislikeNum  int64           `json:"dislike_num"`
}

func (s *VideoService) AIConclusion(ctx context.Context, bvid string, cid int64, upMid int64) (*AIConclusion, error) {
	v := url.Values{}
	setString(v, "bvid", bvid)
	setInt64(v, "cid", cid)
	setInt64(v, "up_mid", upMid)
	var out AIConclusion
	err := s.client.getJSON(ctx, endpointVideoAIGet, v, requestOptions{}, &out)
	return &out, err
}

type VideoLikeParams struct {
	AID  int64
	BVID string
	Like bool
}

func (s *VideoService) Like(ctx context.Context, params VideoLikeParams) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	if err := requireAIDOrBVID(params.AID, params.BVID); err != nil {
		return err
	}
	form := params.videoIDForm()
	if params.Like {
		form.Set("like", "1")
	} else {
		form.Set("like", "2")
	}
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointVideoLike, form, requestOptions{RequireLogin: true}, nil)
}

type VideoCoinParams struct {
	AID        int64
	BVID       string
	Multiply   int
	SelectLike bool
}

func (s *VideoService) Coin(ctx context.Context, params VideoCoinParams) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	if err := requireAIDOrBVID(params.AID, params.BVID); err != nil {
		return err
	}
	form := params.videoIDForm()
	setInt(form, "multiply", params.Multiply)
	setBool01(form, "select_like", params.SelectLike)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointVideoCoin, form, requestOptions{RequireLogin: true}, nil)
}

type VideoFavoriteParams struct {
	AID         int64
	BVID        string
	AddMediaIDs []int64
	DelMediaIDs []int64
}

func (s *VideoService) Favorite(ctx context.Context, params VideoFavoriteParams) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := params.videoIDForm()
	setCSVInt64(form, "add_media_ids", params.AddMediaIDs)
	setCSVInt64(form, "del_media_ids", params.DelMediaIDs)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointVideoFav, form, requestOptions{RequireLogin: true}, nil)
}

func (s *VideoService) Triple(ctx context.Context, params VideoViewParams) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	if err := requireAIDOrBVID(params.AID, params.BVID); err != nil {
		return err
	}
	form := params.videoIDForm()
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointVideoTriple, form, requestOptions{RequireLogin: true}, nil)
}

func (p VideoViewParams) videoIDForm() url.Values {
	v := url.Values{}
	setInt64(v, "aid", p.AID)
	setString(v, "bvid", p.BVID)
	return v
}

func (p VideoLikeParams) videoIDForm() url.Values {
	return VideoViewParams{AID: p.AID, BVID: p.BVID}.videoIDForm()
}

func (p VideoCoinParams) videoIDForm() url.Values {
	return VideoViewParams{AID: p.AID, BVID: p.BVID}.videoIDForm()
}

func (p VideoFavoriteParams) videoIDForm() url.Values {
	v := url.Values{}
	if p.AID != 0 {
		v.Set("rid", strconv.FormatInt(p.AID, 10))
	} else {
		setString(v, "bvid", p.BVID)
	}
	v.Set("type", "2")
	return v
}
