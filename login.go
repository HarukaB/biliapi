package biliapi

import (
	"context"
)

type LoginService struct {
	client *Client
}

type NavInfo struct {
	IsLogin        bool         `json:"isLogin"`
	EmailVerified  int          `json:"email_verified"`
	Face           string       `json:"face"`
	LevelInfo      LevelInfo    `json:"level_info"`
	Mid            int64        `json:"mid"`
	MobileVerified int          `json:"mobile_verified"`
	Money          float64      `json:"money"`
	Moral          int          `json:"moral"`
	Official       OfficialInfo `json:"official"`
	OfficialVerify OfficialInfo `json:"officialVerify"`
	Pendant        Pendant      `json:"pendant"`
	Scores         int64        `json:"scores"`
	Uname          string       `json:"uname"`
	VipDueDate     int64        `json:"vipDueDate"`
	VipStatus      int          `json:"vipStatus"`
	VipType        int          `json:"vipType"`
	VipPayType     int          `json:"vip_pay_type"`
	VipThemeType   int          `json:"vip_theme_type"`
	VipLabel       VipLabel     `json:"vip_label"`
	Wallet         Wallet       `json:"wallet"`
	HasShop        bool         `json:"has_shop"`
	ShopURL        string       `json:"shop_url"`
	AllowanceCount int64        `json:"allowance_count"`
	AnswerStatus   int          `json:"answer_status"`
	IsSeniorMember int          `json:"is_senior_member"`
	WBIImg         WBIImage     `json:"wbi_img"`
	IsJury         bool         `json:"is_jury"`
}

type LevelInfo struct {
	CurrentLevel int            `json:"current_level"`
	CurrentMin   int64          `json:"current_min"`
	CurrentExp   int64          `json:"current_exp"`
	NextExp      MixedStringInt `json:"next_exp"`
}

type Pendant struct {
	PID               int64  `json:"pid"`
	Name              string `json:"name"`
	Image             string `json:"image"`
	Expire            int64  `json:"expire"`
	ImageEnhance      string `json:"image_enhance"`
	ImageEnhanceFrame string `json:"image_enhance_frame"`
}

type Wallet struct {
	Mid           int64   `json:"mid"`
	BcoinBalance  float64 `json:"bcoin_balance"`
	CouponBalance float64 `json:"coupon_balance"`
	CouponDueTime int64   `json:"coupon_due_time"`
}

type NavStat struct {
	Following    int64 `json:"following"`
	Follower     int64 `json:"follower"`
	DynamicCount int64 `json:"dynamic_count"`
}

func (s *LoginService) Nav(ctx context.Context) (*NavInfo, error) {
	var out NavInfo
	err := s.client.getJSONAllow(ctx, endpointNav, nil, requestOptions{}, &out, map[int]bool{-101: true})
	return &out, err
}

func (s *LoginService) NavStat(ctx context.Context) (*NavStat, error) {
	var out NavStat
	err := s.client.getJSON(ctx, endpointNavStat, nil, requestOptions{RequireLogin: true}, &out)
	return &out, err
}
