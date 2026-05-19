package biliapi

import (
	"context"
	"net/url"
)

type ClientInfoService struct {
	client *Client
}

type ZoneInfo struct {
	Addr       string  `json:"addr"`
	Country    string  `json:"country"`
	Province   string  `json:"province"`
	City       string  `json:"city"`
	ISP        string  `json:"isp"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	ZoneID     int64   `json:"zone_id"`
	CountryID  int64   `json:"country_id"`
	ProvinceID int64   `json:"province_id"`
	CityID     int64   `json:"city_id"`
	ISPCode    int64   `json:"isp_code"`
}

type LiveIPInfo struct {
	Addr     string `json:"addr"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	ISP      string `json:"isp"`
}

func (s *ClientInfoService) Zone(ctx context.Context) (*ZoneInfo, error) {
	var out ZoneInfo
	err := s.client.getJSON(ctx, endpointClientZone, nil, requestOptions{}, &out)
	return &out, err
}

func (s *ClientInfoService) LiveIPInfo(ctx context.Context) (*LiveIPInfo, error) {
	var out LiveIPInfo
	err := s.client.getJSON(ctx, endpointClientLiveIPInfo, nil, requestOptions{}, &out)
	return &out, err
}

func (s *ClientInfoService) AppIP(ctx context.Context) (*LiveIPInfo, error) {
	var out LiveIPInfo
	err := s.client.getJSON(ctx, endpointClientAppIP, nil, requestOptions{}, &out)
	return &out, err
}

type IPInfoParams struct {
	IP string
}

func (s *ClientInfoService) IPInfo(ctx context.Context, params IPInfoParams) (*LiveIPInfo, error) {
	v := url.Values{}
	setString(v, "ip", params.IP)
	var out LiveIPInfo
	err := s.client.getJSON(ctx, endpointClientIPInfo, v, requestOptions{}, &out)
	return &out, err
}
