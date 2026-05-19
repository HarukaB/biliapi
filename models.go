package biliapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type MixedStringInt = json.RawMessage

type FlexibleString string

func (s *FlexibleString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = FlexibleString(str)
		return nil
	}
	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		*s = FlexibleString(num.String())
		return nil
	}
	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		*s = FlexibleString(strconv.FormatFloat(f, 'f', -1, 64))
		return nil
	}
	return nil
}

func (s FlexibleString) String() string {
	return string(s)
}

type flexibleInt64 int64

func (n *flexibleInt64) UnmarshalJSON(data []byte) error {
	var raw json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if string(raw) == "null" {
		*n = 0
		return nil
	}
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if str == "" {
			*n = 0
			return nil
		}
		value, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		*n = flexibleInt64(value)
		return nil
	}
	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		value, err := num.Int64()
		if err != nil {
			return err
		}
		*n = flexibleInt64(value)
		return nil
	}
	return fmt.Errorf("biliapi: expected integer or integer string, got %s", string(data))
}

type OfficialInfo struct {
	Role  int    `json:"role"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Type  int    `json:"type"`
}

type VipInfo struct {
	Type       int             `json:"type"`
	Status     int             `json:"status"`
	DueDate    int64           `json:"due_date"`
	VipPayType int             `json:"vip_pay_type"`
	ThemeType  int             `json:"theme_type"`
	Label      VipLabel        `json:"label"`
	Extra      json.RawMessage `json:"-"`
}

func (v *VipInfo) UnmarshalJSON(data []byte) error {
	type vipInfo VipInfo
	aux := struct {
		DueDate flexibleInt64 `json:"due_date"`
		*vipInfo
	}{
		vipInfo: (*vipInfo)(v),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	v.DueDate = int64(aux.DueDate)
	return nil
}

type VipLabel struct {
	Path        string `json:"path"`
	Text        string `json:"text"`
	LabelTheme  string `json:"label_theme"`
	TextColor   string `json:"text_color"`
	BgStyle     int    `json:"bg_style"`
	BgColor     string `json:"bg_color"`
	BorderColor string `json:"border_color"`
}

type Owner struct {
	Mid  int64  `json:"mid"`
	Name string `json:"name"`
	Face string `json:"face"`
}

type Dimension struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Rotate int `json:"rotate"`
}

type Page struct {
	CID       int64     `json:"cid"`
	Page      int       `json:"page"`
	From      string    `json:"from"`
	Part      string    `json:"part"`
	Duration  int       `json:"duration"`
	VID       string    `json:"vid"`
	WebLink   string    `json:"weblink"`
	Dimension Dimension `json:"dimension"`
}

type ArchiveStat struct {
	AID        int64  `json:"aid"`
	View       int64  `json:"view"`
	Danmaku    int64  `json:"danmaku"`
	Reply      int64  `json:"reply"`
	Favorite   int64  `json:"favorite"`
	Coin       int64  `json:"coin"`
	Share      int64  `json:"share"`
	NowRank    int64  `json:"now_rank"`
	HisRank    int64  `json:"his_rank"`
	Like       int64  `json:"like"`
	Dislike    int64  `json:"dislike"`
	Evaluation string `json:"evaluation"`
	VT         int64  `json:"vt"`
}

type ArchiveRights struct {
	BP            int `json:"bp"`
	Elec          int `json:"elec"`
	Download      int `json:"download"`
	Movie         int `json:"movie"`
	Pay           int `json:"pay"`
	HD5           int `json:"hd5"`
	NoReprint     int `json:"no_reprint"`
	Autoplay      int `json:"autoplay"`
	UGCPay        int `json:"ugc_pay"`
	IsCooperation int `json:"is_cooperation"`
	UGCPayPreview int `json:"ugc_pay_preview"`
	NoBackground  int `json:"no_background"`
	CleanMode     int `json:"clean_mode"`
	IsSteinGate   int `json:"is_stein_gate"`
	Is360         int `json:"is_360"`
	NoShare       int `json:"no_share"`
	ArcPay        int `json:"arc_pay"`
	FreeWatch     int `json:"free_watch"`
}

type DescV2Item struct {
	RawText string `json:"raw_text"`
	Type    int    `json:"type"`
	BizID   int64  `json:"biz_id"`
}

type Subtitle struct {
	AllowSubmit bool           `json:"allow_submit"`
	List        []SubtitleItem `json:"list"`
}

type SubtitleItem struct {
	ID          int64           `json:"id"`
	Lan         string          `json:"lan"`
	LanDoc      string          `json:"lan_doc"`
	IsLock      bool            `json:"is_lock"`
	AuthorMID   int64           `json:"author_mid"`
	SubtitleURL string          `json:"subtitle_url"`
	Author      json.RawMessage `json:"author"`
}

type ArgueInfo struct {
	ArgueLink string `json:"argue_link"`
	ArgueMsg  string `json:"argue_msg"`
	ArgueType int    `json:"argue_type"`
}

type Cursor struct {
	IsBegin bool  `json:"is_begin"`
	Prev    int64 `json:"prev"`
	Next    int64 `json:"next"`
	IsEnd   bool  `json:"is_end"`
	Mode    int   `json:"mode"`
}

type CursorString struct {
	IsBegin bool   `json:"is_begin"`
	Prev    string `json:"prev"`
	Next    string `json:"next"`
	IsEnd   bool   `json:"is_end"`
	Mode    int    `json:"mode"`
}

type BasicUser struct {
	Mid   int64  `json:"mid"`
	Name  string `json:"name"`
	Sex   string `json:"sex"`
	Face  string `json:"face"`
	Sign  string `json:"sign"`
	Rank  int64  `json:"rank"`
	Level int    `json:"level"`
}

func (u *BasicUser) UnmarshalJSON(data []byte) error {
	type basicUser BasicUser
	aux := struct {
		Mid  flexibleInt64 `json:"mid"`
		Rank flexibleInt64 `json:"rank"`
		*basicUser
	}{
		basicUser: (*basicUser)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	u.Mid = int64(aux.Mid)
	u.Rank = int64(aux.Rank)
	return nil
}
