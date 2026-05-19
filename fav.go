package biliapi

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

type FavService struct {
	client *Client
}

type FavFolderInfoParams struct {
	MediaID int64
}

func (p FavFolderInfoParams) values() url.Values {
	v := url.Values{}
	setInt64(v, "media_id", p.MediaID)
	return v
}

type FavFolder struct {
	ID         int64           `json:"id"`
	FID        int64           `json:"fid"`
	MID        int64           `json:"mid"`
	Attr       int64           `json:"attr"`
	Title      string          `json:"title"`
	FavState   int             `json:"fav_state"`
	MediaCount int64           `json:"media_count"`
	Cover      string          `json:"cover"`
	Upper      Owner           `json:"upper"`
	CoverType  int             `json:"cover_type"`
	CntInfo    FavCntInfo      `json:"cnt_info"`
	Type       int             `json:"type"`
	Intro      string          `json:"intro"`
	CTime      int64           `json:"ctime"`
	MTime      int64           `json:"mtime"`
	State      int             `json:"state"`
	Extra      json.RawMessage `json:"extra,omitempty"`
}

type FavCntInfo struct {
	Collect int64 `json:"collect"`
	Play    int64 `json:"play"`
	ThumbUp int64 `json:"thumb_up"`
	Share   int64 `json:"share"`
}

func (s *FavService) FolderInfo(ctx context.Context, params FavFolderInfoParams) (*FavFolder, error) {
	var out FavFolder
	err := s.client.getJSON(ctx, endpointFavFolderInfo, params.values(), requestOptions{}, &out)
	return &out, err
}

type FavCreatedListAllParams struct {
	UpMid int64
	Type  int
}

func (p FavCreatedListAllParams) values() url.Values {
	v := url.Values{}
	setInt64(v, "up_mid", p.UpMid)
	setInt(v, "type", p.Type)
	return v
}

type FavCreatedListAll struct {
	Count int64       `json:"count"`
	List  []FavFolder `json:"list"`
}

func (s *FavService) CreatedListAll(ctx context.Context, params FavCreatedListAllParams) (*FavCreatedListAll, error) {
	var out FavCreatedListAll
	err := s.client.getJSON(ctx, endpointFavCreatedListAll, params.values(), requestOptions{}, &out)
	return &out, err
}

type FavCollectedListParams struct {
	UpMid int64
	PN    int
	PS    int
}

func (p FavCollectedListParams) values() url.Values {
	v := url.Values{}
	setInt64(v, "up_mid", p.UpMid)
	setInt(v, "pn", p.PN)
	setInt(v, "ps", p.PS)
	return v
}

type FavCollectedList struct {
	Count int64       `json:"count"`
	List  []FavFolder `json:"list"`
}

func (s *FavService) CollectedList(ctx context.Context, params FavCollectedListParams) (*FavCollectedList, error) {
	var out FavCollectedList
	err := s.client.getJSON(ctx, endpointFavCollectedList, params.values(), requestOptions{}, &out)
	return &out, err
}

type FavResourceListParams struct {
	MediaID  int64
	PN       int
	PS       int
	Keyword  string
	Order    string
	Type     int
	TID      int
	Platform string
}

func (p FavResourceListParams) values() url.Values {
	v := url.Values{}
	setInt64(v, "media_id", p.MediaID)
	setInt(v, "pn", p.PN)
	setInt(v, "ps", p.PS)
	setString(v, "keyword", p.Keyword)
	setString(v, "order", p.Order)
	setInt(v, "type", p.Type)
	setInt(v, "tid", p.TID)
	setString(v, "platform", p.Platform)
	return v
}

type FavResourceList struct {
	Info    FavFolder     `json:"info"`
	Medias  []FavResource `json:"medias"`
	HasMore bool          `json:"has_more"`
	TTL     int           `json:"ttl"`
}

type FavResource struct {
	ID       int64           `json:"id"`
	Type     int             `json:"type"`
	Title    string          `json:"title"`
	Cover    string          `json:"cover"`
	Intro    string          `json:"intro"`
	Page     int             `json:"page"`
	Duration int             `json:"duration"`
	Upper    Owner           `json:"upper"`
	Attr     int64           `json:"attr"`
	CNTInfo  FavCntInfo      `json:"cnt_info"`
	Link     string          `json:"link"`
	CTime    int64           `json:"ctime"`
	PubTime  int64           `json:"pubtime"`
	FavTime  int64           `json:"fav_time"`
	BVID     string          `json:"bvid"`
	Season   json.RawMessage `json:"season"`
	OGV      json.RawMessage `json:"ogv"`
}

func (s *FavService) ResourceList(ctx context.Context, params FavResourceListParams) (*FavResourceList, error) {
	var out FavResourceList
	err := s.client.getJSON(ctx, endpointFavResourceList, params.values(), requestOptions{}, &out)
	return &out, err
}

type FavResourceRef struct {
	ID      int64  `json:"id"`
	Type    int    `json:"type"`
	BVID    string `json:"bvid"`
	BVIDAlt string `json:"bv_id"`
}

func (s *FavService) ResourceIDs(ctx context.Context, mediaID int64) ([]FavResourceRef, error) {
	v := url.Values{}
	setInt64(v, "media_id", mediaID)
	var out []FavResourceRef
	err := s.client.getJSON(ctx, endpointFavResourceIDs, v, requestOptions{}, &out)
	return out, err
}

type FavResourceInfosParams struct {
	Resources []FavResourceID
}

type FavResourceID struct {
	ID   int64
	Type int
}

func (p FavResourceInfosParams) values() url.Values {
	v := url.Values{}
	parts := make([]string, 0, len(p.Resources))
	for _, item := range p.Resources {
		if item.ID == 0 {
			continue
		}
		parts = append(parts, strconv.FormatInt(item.ID, 10)+":"+strconv.Itoa(item.Type))
	}
	if len(parts) > 0 {
		v.Set("resources", strings.Join(parts, ","))
	}
	return v
}

func (s *FavService) ResourceInfos(ctx context.Context, params FavResourceInfosParams) ([]FavResource, error) {
	var out []FavResource
	err := s.client.getJSON(ctx, endpointFavResourceInfos, params.values(), requestOptions{}, &out)
	return out, err
}

type FavFolderAddParams struct {
	Title   string
	Intro   string
	Privacy int
}

func (s *FavService) AddFolder(ctx context.Context, params FavFolderAddParams) (*FavFolder, error) {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return nil, err
	}
	form := url.Values{}
	setString(form, "title", params.Title)
	setString(form, "intro", params.Intro)
	setInt(form, "privacy", params.Privacy)
	form.Set("csrf", s.client.creds.CSRF())
	var out FavFolder
	err := s.client.postFormJSON(ctx, endpointFavFolderAdd, form, requestOptions{RequireLogin: true}, &out)
	return &out, err
}

func (s *FavService) EditFolder(ctx context.Context, mediaID int64, params FavFolderAddParams) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{}
	setInt64(form, "media_id", mediaID)
	setString(form, "title", params.Title)
	setString(form, "intro", params.Intro)
	setInt(form, "privacy", params.Privacy)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointFavFolderEdit, form, requestOptions{RequireLogin: true}, nil)
}

func (s *FavService) DeleteFolder(ctx context.Context, mediaID int64) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{}
	setInt64(form, "media_id", mediaID)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointFavFolderDel, form, requestOptions{RequireLogin: true}, nil)
}

type FavResourceMoveParams struct {
	SourceMediaID int64
	TargetMediaID int64
	Resources     []FavResourceID
}

func (s *FavService) MoveResources(ctx context.Context, params FavResourceMoveParams) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := params.resourceForm()
	setInt64(form, "src_media_id", params.SourceMediaID)
	setInt64(form, "tar_media_id", params.TargetMediaID)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointFavResourceMove, form, requestOptions{RequireLogin: true}, nil)
}

func (s *FavService) CopyResources(ctx context.Context, params FavResourceMoveParams) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := params.resourceForm()
	setInt64(form, "src_media_id", params.SourceMediaID)
	setInt64(form, "tar_media_id", params.TargetMediaID)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointFavResourceCopy, form, requestOptions{RequireLogin: true}, nil)
}

func (s *FavService) DeleteResources(ctx context.Context, mediaID int64, resources []FavResourceID) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{}
	setInt64(form, "media_id", mediaID)
	setFavResources(form, resources)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointFavResourceBatchDel, form, requestOptions{RequireLogin: true}, nil)
}

func (s *FavService) CleanInvalidResources(ctx context.Context, mediaID int64) error {
	if err := s.client.creds.RequireCSRF(); err != nil {
		return err
	}
	form := url.Values{}
	setInt64(form, "media_id", mediaID)
	form.Set("csrf", s.client.creds.CSRF())
	return s.client.postFormJSON(ctx, endpointFavResourceClean, form, requestOptions{RequireLogin: true}, nil)
}

func (p FavResourceMoveParams) resourceForm() url.Values {
	form := url.Values{}
	setFavResources(form, p.Resources)
	return form
}

func setFavResources(values url.Values, resources []FavResourceID) {
	parts := make([]string, 0, len(resources))
	for _, item := range resources {
		if item.ID == 0 {
			continue
		}
		parts = append(parts, strconv.FormatInt(item.ID, 10)+":"+strconv.Itoa(item.Type))
	}
	if len(parts) > 0 {
		values.Set("resources", strings.Join(parts, ","))
	}
}
