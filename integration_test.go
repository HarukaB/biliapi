package biliapi

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

const fallbackIntegrationBVID = "BV1xx411c7mD"
const fallbackIntegrationSearchKeyword = "猫"

var fallbackIntegrationAIDs = []int64{1700001}

const fallbackIntegrationFavID int64 = 2519150301

func credentialsFromEnv(t *testing.T) Credentials {
	t.Helper()

	if raw := os.Getenv("BILIAPI_CREDENTIALS_JSON"); raw != "" {
		creds, err := ParseCredentialsJSON([]byte(raw))
		if err != nil {
			t.Fatalf("parse BILIAPI_CREDENTIALS_JSON: %v", err)
		}
		return creds
	}

	return Credentials{
		SESSDATA:        os.Getenv("BILIAPI_SESSDATA"),
		BiliJCT:         os.Getenv("BILIAPI_BILI_JCT"),
		DedeUserID:      os.Getenv("BILIAPI_DEDE_USER_ID"),
		DedeUserIDCKMd5: os.Getenv("BILIAPI_DEDE_USER_ID_CKMD5"),
		SID:             os.Getenv("BILIAPI_SID"),
		RefreshToken:    os.Getenv("BILIAPI_REFRESH_TOKEN"),
	}
}

func integrationClient() *Client {
	return NewClient(WithTimeout(20 * time.Second))
}

func integrationContext(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(context.Background(), 45*time.Second)
}

func integrationBVID() string {
	if bvid := os.Getenv("BILIAPI_TEST_BVID"); bvid != "" {
		return bvid
	}
	return fallbackIntegrationBVID
}

func integrationAIDs(t *testing.T) []int64 {
	t.Helper()
	raw := os.Getenv("BILIAPI_TEST_AIDS")
	if raw == "" {
		return append([]int64(nil), fallbackIntegrationAIDs...)
	}

	parts := strings.Split(raw, ",")
	aids := make([]int64, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(strings.TrimPrefix(strings.ToLower(part), "av"))
		if part == "" {
			continue
		}
		aid, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			t.Fatalf("parse BILIAPI_TEST_AIDS item %q: %v", part, err)
		}
		aids = append(aids, aid)
	}
	if len(aids) == 0 {
		t.Fatalf("BILIAPI_TEST_AIDS did not contain any valid aid")
	}
	return aids
}

func integrationFavID(t *testing.T) int64 {
	t.Helper()
	raw := strings.TrimSpace(os.Getenv("BILIAPI_TEST_FAV_ID"))
	if raw == "" {
		return fallbackIntegrationFavID
	}
	raw = strings.TrimPrefix(strings.ToLower(raw), "fav")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		t.Fatalf("parse BILIAPI_TEST_FAV_ID %q: %v", raw, err)
	}
	return id
}

func integrationSearchKeyword() string {
	if keyword := strings.TrimSpace(os.Getenv("BILIAPI_TEST_SEARCH_KEYWORD")); keyword != "" {
		return keyword
	}
	return fallbackIntegrationSearchKeyword
}

func integrationLoggedClient(t *testing.T) *Client {
	t.Helper()
	creds := credentialsFromEnv(t)
	if !creds.IsLoggedInCandidate() {
		t.Skip("set BILIAPI_CREDENTIALS_JSON or BILIAPI_SESSDATA to run credential integration tests")
	}
	return NewClient(WithCredentials(creds), WithTimeout(20*time.Second))
}

func integrationReferenceVideo(t *testing.T, ctx context.Context, client *Client) *VideoView {
	t.Helper()
	bvid := integrationBVID()
	view, err := client.Video.View(ctx, VideoViewParams{BVID: bvid})
	if err != nil {
		t.Fatalf("Video.View(%s) with real API: %v", bvid, err)
	}
	if view.BVID == "" || view.AID <= 0 || view.CID <= 0 || len(view.Pages) == 0 {
		t.Fatalf("Video.View(%s) returned incomplete payload: %#v", bvid, view)
	}
	return view
}

func integrationVideoCases(t *testing.T) []struct {
	name   string
	params VideoViewParams
} {
	t.Helper()
	cases := []struct {
		name   string
		params VideoViewParams
	}{
		{name: "bvid/" + integrationBVID(), params: VideoViewParams{BVID: integrationBVID()}},
	}
	for _, aid := range integrationAIDs(t) {
		cases = append(cases, struct {
			name   string
			params VideoViewParams
		}{
			name:   "aid/av" + strconv.FormatInt(aid, 10),
			params: VideoViewParams{AID: aid},
		})
	}
	return cases
}

func TestIntegrationPublicAPIs(t *testing.T) {
	client := integrationClient()
	ctx, cancel := integrationContext(t)
	defer cancel()

	now, err := client.Misc.Now(ctx)
	if err != nil {
		t.Fatalf("Misc.Now() with real API: %v", err)
	}
	if now.Now <= 0 {
		t.Fatalf("Misc.Now() returned invalid timestamp: %#v", now)
	}

	buvid, err := client.Misc.Buvid(ctx)
	if err != nil {
		t.Fatalf("Misc.Buvid() with real API: %v", err)
	}
	if buvid.B3 == "" || buvid.B4 == "" {
		t.Fatalf("Misc.Buvid() returned incomplete identifiers: %#v", buvid)
	}

	zone, err := client.ClientInfo.Zone(ctx)
	if err != nil {
		t.Fatalf("ClientInfo.Zone() with real API: %v", err)
	}
	if zone.Addr == "" && zone.Country == "" && zone.Province == "" {
		t.Fatalf("ClientInfo.Zone() returned empty location data: %#v", zone)
	}

	suggest, err := client.Search.Suggest(ctx, SearchSuggestParams{Term: "go"})
	if err != nil {
		t.Fatalf("Search.Suggest() with real API: %v", err)
	}
	if len(suggest.Result.Tag) == 0 {
		t.Fatalf("Search.Suggest() returned no suggestions: %#v", suggest)
	}

	defaultSearch, err := client.Search.Default(ctx)
	if err != nil {
		t.Fatalf("Search.Default() with real WBI API: %v", err)
	}
	if defaultSearch.ShowName == "" && defaultSearch.Name == "" {
		t.Fatalf("Search.Default() returned empty payload: %#v", defaultSearch)
	}
}

func TestIntegrationSearchRead(t *testing.T) {
	client := integrationClient()
	ctx, cancel := integrationContext(t)
	defer cancel()
	keyword := integrationSearchKeyword()

	square, err := client.Search.Square(ctx, SearchSquareParams{Limit: 10, Platform: "web"})
	if err != nil {
		t.Fatalf("Search.Square() with real API: %v", err)
	}
	if square.Title == "" && len(square.List) == 0 {
		t.Fatalf("Search.Square() returned empty payload: %#v", square)
	}

	hotword, err := client.Search.Hotword(ctx)
	if err != nil {
		t.Fatalf("Search.Hotword() with real API: %v", err)
	}
	if len(hotword.Result.TopList) == 0 {
		t.Fatalf("Search.Hotword() returned no hotwords: %#v", hotword)
	}

	all, err := client.Search.All(ctx, SearchAllParams{Keyword: keyword, Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("Search.All(%q) with real API: %v", keyword, err)
	}
	if all.NumResults < 0 || len(all.Result) == 0 {
		t.Fatalf("Search.All(%q) returned invalid payload: %#v", keyword, all)
	}

	typed, err := client.Search.Type(ctx, SearchTypeParams{
		Keyword:    keyword,
		SearchType: SearchTypeVideo,
		Page:       1,
		PageSize:   10,
	})
	if err != nil {
		t.Fatalf("Search.Type(%q) with real API: %v", keyword, err)
	}
	if typed.NumResults < 0 || len(typed.Result) == 0 {
		t.Fatalf("Search.Type(%q) returned invalid payload: %#v", keyword, typed)
	}
}

func TestIntegrationVideoReadChain(t *testing.T) {
	client := integrationClient()
	ctx, cancel := integrationContext(t)
	defer cancel()

	for _, tc := range integrationVideoCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			view, err := client.Video.View(ctx, tc.params)
			if err != nil {
				t.Fatalf("Video.View(%s) with real API: %v", tc.name, err)
			}
			if view.BVID == "" || view.AID <= 0 || view.CID <= 0 || len(view.Pages) == 0 {
				t.Fatalf("Video.View(%s) returned incomplete payload: %#v", tc.name, view)
			}

			pages, err := client.Video.PageList(ctx, tc.params)
			if err != nil {
				t.Fatalf("Video.PageList(%s) with real API: %v", tc.name, err)
			}
			if len(pages) == 0 || pages[0].CID <= 0 {
				t.Fatalf("Video.PageList(%s) returned no pages: %#v", tc.name, pages)
			}

			stat, err := client.Video.Stat(ctx, tc.params)
			if isHTTPStatus(err, 404) {
				t.Logf("Video.Stat(%s) returned HTTP 404; archive/stat appears unavailable from current edge", tc.name)
			} else if err != nil {
				t.Fatalf("Video.Stat(%s) with real API: %v", tc.name, err)
			} else if stat.AID != view.AID {
				t.Fatalf("Video.Stat(%s) aid mismatch: got %d want %d", tc.name, stat.AID, view.AID)
			}

			desc, err := client.Video.Desc(ctx, tc.params)
			if err != nil {
				t.Fatalf("Video.Desc(%s) with real API: %v", tc.name, err)
			}
			if view.Desc != "" && desc == "" {
				t.Fatalf("Video.Desc(%s) returned empty desc while Video.View had non-empty desc", tc.name)
			}

			playURL, err := client.Video.PlayURL(ctx, PlayURLParams{
				AID:   tc.params.AID,
				BVID:  tc.params.BVID,
				CID:   pages[0].CID,
				QN:    16,
				FnVer: 0,
				FnVal: 16,
			})
			if err != nil {
				t.Fatalf("Video.PlayURL(%s, cid=%d) with real API: %v", tc.name, pages[0].CID, err)
			}
			if len(playURL.DURL) == 0 && (playURL.Dash == nil || len(playURL.Dash.Video) == 0) {
				t.Fatalf("Video.PlayURL(%s) returned no playable URL: %#v", tc.name, playURL)
			}
		})
	}
}

func TestIntegrationVideoAdditionalRead(t *testing.T) {
	client := integrationClient()
	ctx, cancel := integrationContext(t)
	defer cancel()
	view := integrationReferenceVideo(t, ctx, client)
	params := VideoViewParams{BVID: view.BVID}

	detail, err := client.Video.ViewDetail(ctx, params)
	if isBiliCode(err, -352) {
		t.Logf("Video.ViewDetail(%s) hit Bilibili risk control on this network: %v", view.BVID, err)
	} else if err != nil {
		t.Fatalf("Video.ViewDetail(%s) with real API: %v", view.BVID, err)
	} else if detail.View.BVID == "" || detail.View.AID != view.AID {
		t.Fatalf("Video.ViewDetail(%s) returned mismatched view: %#v", view.BVID, detail.View)
	}

	player, err := client.Video.PlayerInfo(ctx, PlayerInfoParams{BVID: view.BVID, CID: view.CID})
	if err != nil {
		t.Fatalf("Video.PlayerInfo(%s) with real API: %v", view.BVID, err)
	}
	if player.CID != view.CID {
		t.Fatalf("Video.PlayerInfo(%s) cid mismatch: got %d want %d", view.BVID, player.CID, view.CID)
	}

	related, err := client.Video.Related(ctx, params)
	if err != nil {
		t.Fatalf("Video.Related(%s) with real API: %v", view.BVID, err)
	}
	if len(related) == 0 {
		t.Fatalf("Video.Related(%s) returned no related videos", view.BVID)
	}

	online, err := client.Video.OnlineTotal(ctx, view.AID, view.CID, "")
	if err != nil {
		t.Fatalf("Video.OnlineTotal(%s) with real API: %v", view.BVID, err)
	}
	if online.Total == "" && online.Count == "" {
		t.Fatalf("Video.OnlineTotal(%s) returned empty payload: %#v", view.BVID, online)
	}

	shot, err := client.Video.Shot(ctx, view.AID, view.CID, "")
	if err != nil {
		t.Fatalf("Video.Shot(%s) with real API: %v", view.BVID, err)
	}
	if len(shot.Image) == 0 && shot.PVData == "" {
		t.Fatalf("Video.Shot(%s) returned empty payload: %#v", view.BVID, shot)
	}

	conclusion, err := client.Video.AIConclusion(ctx, view.BVID, view.CID, view.Owner.Mid)
	if isBiliCode(err, -403) || isBiliCode(err, -404) || isBiliCode(err, 400) || isBiliCode(err, 30002) {
		t.Logf("Video.AIConclusion(%s) is unavailable for this video: %v", view.BVID, err)
	} else if err != nil {
		t.Fatalf("Video.AIConclusion(%s) with real API: %v", view.BVID, err)
	} else if conclusion.Status < 0 {
		t.Fatalf("Video.AIConclusion(%s) returned invalid status: %#v", view.BVID, conclusion)
	}
}

func TestIntegrationVideoCredentialRead(t *testing.T) {
	client := integrationLoggedClient(t)
	ctx, cancel := integrationContext(t)
	defer cancel()
	view := integrationReferenceVideo(t, ctx, client)
	params := VideoViewParams{BVID: view.BVID}

	hasLike, err := client.Video.HasLike(ctx, params)
	if err != nil {
		t.Fatalf("Video.HasLike(%s) with real API: %v", view.BVID, err)
	}
	_ = hasLike

	coins, err := client.Video.Coins(ctx, params)
	if err != nil {
		t.Fatalf("Video.Coins(%s) with real API: %v", view.BVID, err)
	}
	if coins.Multiply < 0 {
		t.Fatalf("Video.Coins(%s) returned invalid payload: %#v", view.BVID, coins)
	}

	favoured, err := client.Video.Favoured(ctx, view.AID)
	if err != nil {
		t.Fatalf("Video.Favoured(%s) with real API: %v", view.BVID, err)
	}
	if favoured.Count < 0 {
		t.Fatalf("Video.Favoured(%s) returned invalid payload: %#v", view.BVID, favoured)
	}
}

func isHTTPStatus(err error, status int) bool {
	var biliErr *BiliError
	return errors.As(err, &biliErr) && biliErr.HTTPStatus == status
}

func isBiliCode(err error, code int) bool {
	var biliErr *BiliError
	return errors.As(err, &biliErr) && biliErr.Code == code
}

func TestIntegrationCommentAndDanmakuRead(t *testing.T) {
	client := integrationClient()
	ctx, cancel := integrationContext(t)
	defer cancel()

	bvid := integrationBVID()
	view, err := client.Video.View(ctx, VideoViewParams{BVID: bvid})
	if err != nil {
		t.Fatalf("Video.View(%s) with real API: %v", bvid, err)
	}

	count, err := client.Comment.Count(ctx, view.AID, CommentTypeVideo)
	if err != nil {
		t.Fatalf("Comment.Count(%s) with real API: %v", view.BVID, err)
	}
	if count.Count < 0 {
		t.Fatalf("Comment.Count(%s) returned invalid count: %#v", view.BVID, count)
	}

	comments, err := client.Comment.List(ctx, CommentListParams{
		OID:  view.AID,
		Type: CommentTypeVideo,
		PN:   1,
		PS:   5,
	})
	if err != nil {
		t.Fatalf("Comment.List(%s) with real API: %v", view.BVID, err)
	}
	if comments.Page.Count < 0 {
		t.Fatalf("Comment.List(%s) returned invalid page: %#v", view.BVID, comments.Page)
	}

	danmaku, err := client.Danmaku.XML(ctx, DanmakuXMLParams{CID: view.CID})
	if err != nil {
		t.Fatalf("Danmaku.XML(cid=%d) with real API: %v", view.CID, err)
	}
	if danmaku.ChatID <= 0 && len(danmaku.Items) == 0 {
		t.Fatalf("Danmaku.XML(cid=%d) returned empty payload: %#v", view.CID, danmaku)
	}
}

func TestIntegrationCommentAdditionalRead(t *testing.T) {
	client := integrationClient()
	ctx, cancel := integrationContext(t)
	defer cancel()
	view := integrationReferenceVideo(t, ctx, client)

	main, err := client.Comment.Main(ctx, CommentMainParams{
		OID:  view.AID,
		Type: CommentTypeVideo,
		Mode: 3,
		PS:   5,
	})
	if err != nil {
		t.Fatalf("Comment.Main(%s) with real API: %v", view.BVID, err)
	}
	if main.Cursor.Mode < 0 {
		t.Fatalf("Comment.Main(%s) returned invalid cursor: %#v", view.BVID, main.Cursor)
	}

	hot, err := client.Comment.Hot(ctx, CommentListParams{OID: view.AID, Type: CommentTypeVideo, PN: 1, PS: 5})
	if err != nil {
		t.Fatalf("Comment.Hot(%s) with real API: %v", view.BVID, err)
	}
	if hot.Page.Count < 0 {
		t.Fatalf("Comment.Hot(%s) returned invalid page: %#v", view.BVID, hot.Page)
	}

	var first *CommentReply
	if len(main.Replies) > 0 {
		first = &main.Replies[0]
	} else if len(hot.Replies) > 0 {
		first = &hot.Replies[0]
	}
	if first == nil {
		t.Skipf("Comment.Main/Hot(%s) returned no comments to test replies/info/dialog", view.BVID)
	}

	info, err := client.Comment.Info(ctx, view.AID, CommentTypeVideo, first.RPID)
	if err != nil {
		t.Fatalf("Comment.Info(%s, rpid=%d) with real API: %v", view.BVID, first.RPID, err)
	}
	if info.RPID == 0 {
		t.Logf("Comment.Info(%s, rpid=%d) returned an empty payload from deprecated endpoint", view.BVID, first.RPID)
	} else if info.RPID != first.RPID {
		t.Fatalf("Comment.Info(%s) rpid mismatch: got %d want %d", view.BVID, info.RPID, first.RPID)
	}

	replies, err := client.Comment.Replies(ctx, CommentRepliesParams{
		OID:  view.AID,
		Type: CommentTypeVideo,
		Root: first.RPID,
		PN:   1,
		PS:   5,
	})
	if err != nil {
		t.Fatalf("Comment.Replies(%s, root=%d) with real API: %v", view.BVID, first.RPID, err)
	}
	if replies.Page.Count < 0 {
		t.Fatalf("Comment.Replies(%s) returned invalid page: %#v", view.BVID, replies.Page)
	}

	if first.Dialog == 0 {
		t.Logf("Comment.Dialog(%s) skipped because first comment has no dialog id", view.BVID)
		return
	}
	dialog, err := client.Comment.Dialog(ctx, CommentDialogParams{
		OID:    view.AID,
		Type:   CommentTypeVideo,
		Root:   first.RPID,
		Dialog: first.Dialog,
		PN:     1,
		PS:     5,
	})
	if err != nil {
		t.Fatalf("Comment.Dialog(%s, dialog=%d) with real API: %v", view.BVID, first.Dialog, err)
	}
	if dialog.Page.Count < 0 {
		t.Fatalf("Comment.Dialog(%s) returned invalid page: %#v", view.BVID, dialog.Page)
	}
}

func TestIntegrationFavoriteReadChain(t *testing.T) {
	client := integrationClient()
	ctx, cancel := integrationContext(t)
	defer cancel()

	favID := integrationFavID(t)
	info, err := client.Fav.FolderInfo(ctx, FavFolderInfoParams{MediaID: favID})
	if err != nil {
		t.Fatalf("Fav.FolderInfo(%d) with real API: %v", favID, err)
	}
	if info.ID <= 0 && info.FID <= 0 {
		t.Fatalf("Fav.FolderInfo(%d) returned invalid folder info: %#v", favID, info)
	}

	resources, err := client.Fav.ResourceList(ctx, FavResourceListParams{
		MediaID:  favID,
		PN:       1,
		PS:       5,
		Platform: "web",
	})
	if err != nil {
		t.Fatalf("Fav.ResourceList(%d) with real API: %v", favID, err)
	}
	if resources.Info.ID <= 0 && resources.Info.FID <= 0 {
		t.Fatalf("Fav.ResourceList(%d) returned invalid folder info: %#v", favID, resources.Info)
	}
	if resources.Info.MediaCount > 0 && len(resources.Medias) == 0 {
		t.Fatalf("Fav.ResourceList(%d) returned no medias for non-empty folder: %#v", favID, resources)
	}

	ids, err := client.Fav.ResourceIDs(ctx, favID)
	if err != nil {
		t.Fatalf("Fav.ResourceIDs(%d) with real API: %v", favID, err)
	}
	if resources.Info.MediaCount > 0 && len(ids) == 0 {
		t.Fatalf("Fav.ResourceIDs(%d) returned no ids for non-empty folder", favID)
	}
}

func TestIntegrationCredentialsReadOnly(t *testing.T) {
	client := integrationLoggedClient(t)
	ctx, cancel := integrationContext(t)
	defer cancel()

	nav, err := client.Login.Nav(ctx)
	if err != nil {
		t.Fatalf("Login.Nav() with real API: %v", err)
	}
	if !nav.IsLogin {
		t.Fatalf("Login.Nav() returned IsLogin=false; credentials are present but not accepted")
	}
	if nav.Mid <= 0 {
		t.Fatalf("Login.Nav() returned invalid mid: %d", nav.Mid)
	}
	if nav.WBIImg.ImgURL == "" || nav.WBIImg.SubURL == "" {
		t.Fatalf("Login.Nav() did not return WBI image keys: %#v", nav.WBIImg)
	}

	stat, err := client.Login.NavStat(ctx)
	if err != nil {
		t.Fatalf("Login.NavStat() with real API: %v", err)
	}
	if stat.Follower < 0 || stat.Following < 0 {
		t.Fatalf("Login.NavStat() returned impossible counts: %#v", stat)
	}

	myInfo, err := client.User.MyInfo(ctx)
	if err != nil {
		t.Fatalf("User.MyInfo() with real API: %v", err)
	}
	if myInfo.Mid != nav.Mid {
		t.Fatalf("User.MyInfo() mid mismatch: got %d want %d", myInfo.Mid, nav.Mid)
	}

	history, err := client.History.Cursor(ctx, HistoryCursorParams{Type: "archive", PS: 5})
	if err != nil {
		t.Fatalf("History.Cursor() with real API: %v", err)
	}
	if history.Cursor.PS <= 0 {
		t.Fatalf("History.Cursor() returned invalid cursor: %#v", history.Cursor)
	}
	if len(history.Tab) == 0 {
		t.Fatalf("History.Cursor() returned no tabs: %#v", history)
	}
	for i, item := range history.List {
		if item.Title == "" {
			t.Fatalf("History.Cursor() item %d has empty title: %#v", i, item)
		}
		if item.History.Business == "" {
			t.Fatalf("History.Cursor() item %d has empty business: %#v", i, item.History)
		}
	}

	shadow, err := client.History.Shadow(ctx)
	if err != nil {
		t.Fatalf("History.Shadow() with real API: %v", err)
	}
	_ = shadow

	folders, err := client.Fav.CreatedListAll(ctx, FavCreatedListAllParams{UpMid: nav.Mid, Type: 2})
	if err != nil {
		t.Fatalf("Fav.CreatedListAll() with real API: %v", err)
	}
	if folders.Count < 0 {
		t.Fatalf("Fav.CreatedListAll() returned invalid count: %#v", folders)
	}
}

func TestIntegrationHistoryRead(t *testing.T) {
	client := integrationLoggedClient(t)
	ctx, cancel := integrationContext(t)
	defer cancel()

	legacy, err := client.History.Legacy(ctx, HistoryLegacyParams{PN: 1, PS: 5})
	if err != nil {
		t.Fatalf("History.Legacy() with real API: %v", err)
	}
	if len(legacy.List) == 0 {
		t.Logf("History.Legacy() returned no items")
	}

	toView, err := client.History.ToView(ctx)
	if err != nil {
		t.Fatalf("History.ToView() with real API: %v", err)
	}
	if toView.Count < 0 {
		t.Fatalf("History.ToView() returned invalid count: %#v", toView)
	}
}

func TestIntegrationUserRead(t *testing.T) {
	client := integrationLoggedClient(t)
	ctx, cancel := integrationContext(t)
	defer cancel()
	view := integrationReferenceVideo(t, ctx, client)
	mid := view.Owner.Mid
	if mid <= 0 {
		t.Fatalf("reference video owner mid is invalid: %#v", view.Owner)
	}

	space, err := client.User.SpaceInfo(ctx, UserSpaceInfoParams{Mid: mid})
	if err != nil {
		t.Fatalf("User.SpaceInfo(%d) with real API: %v", mid, err)
	}
	if space.Mid != mid {
		t.Fatalf("User.SpaceInfo(%d) mid mismatch: got %d", mid, space.Mid)
	}

	card, err := client.User.Card(ctx, UserCardParams{Mid: mid})
	if err != nil {
		t.Fatalf("User.Card(%d) with real API: %v", mid, err)
	}
	if card.Card.Mid == "" || card.Card.Name == "" {
		t.Fatalf("User.Card(%d) returned invalid card: %#v", mid, card.Card)
	}

	cards, err := client.User.Cards(ctx, UserCardsParams{Mids: []int64{mid}})
	if err != nil {
		t.Fatalf("User.Cards(%d) with real API: %v", mid, err)
	}
	if len(*cards) == 0 {
		t.Fatalf("User.Cards(%d) returned empty map", mid)
	}

	arcs, err := client.User.ArcSearch(ctx, UserArcSearchParams{Mid: mid, PN: 1, PS: 5, Order: "pubdate"})
	if isHTTPStatus(err, 412) || isBiliCode(err, -352) {
		t.Logf("User.ArcSearch(%d) hit Bilibili risk control on this network: %v", mid, err)
	} else if err != nil {
		t.Fatalf("User.ArcSearch(%d) with real API: %v", mid, err)
	} else if arcs.Page.Count < 0 {
		t.Fatalf("User.ArcSearch(%d) returned invalid page: %#v", mid, arcs.Page)
	}

	upStat, err := client.User.UpStat(ctx, mid)
	if err != nil {
		t.Fatalf("User.UpStat(%d) with real API: %v", mid, err)
	}
	if upStat.Likes < 0 {
		t.Fatalf("User.UpStat(%d) returned invalid likes: %#v", mid, upStat)
	}

	relation, err := client.User.RelationStat(ctx, mid)
	if err != nil {
		t.Fatalf("User.RelationStat(%d) with real API: %v", mid, err)
	}
	if relation.Mid != mid {
		t.Fatalf("User.RelationStat(%d) mid mismatch: %#v", mid, relation)
	}

	navNum, err := client.User.NavNum(ctx, mid)
	if err != nil {
		t.Fatalf("User.NavNum(%d) with real API: %v", mid, err)
	}
	if navNum.Video < 0 {
		t.Fatalf("User.NavNum(%d) returned invalid nav num: %#v", mid, navNum)
	}
}

func TestIntegrationHistoryCursorData(t *testing.T) {
	creds := credentialsFromEnv(t)
	if !creds.IsLoggedInCandidate() {
		t.Skip("set BILIAPI_CREDENTIALS_JSON or BILIAPI_SESSDATA to dump history data")
	}

	client := NewClient(WithCredentials(creds), WithTimeout(20*time.Second))
	ctx, cancel := integrationContext(t)
	defer cancel()

	history, err := client.History.Cursor(ctx, HistoryCursorParams{Type: "archive", PS: 5})
	if err != nil {
		t.Fatalf("History.Cursor() with real API: %v", err)
	}

	t.Logf("cursor: max=%d view_at=%d business=%q ps=%d", history.Cursor.Max, history.Cursor.ViewAt, history.Cursor.Business, history.Cursor.PS)
	for i, tab := range history.Tab {
		t.Logf("tab[%d]: type=%q name=%q", i, tab.Type, tab.Name)
	}
	for i, item := range history.List {
		t.Logf(
			"item[%d]: title=%q bvid=%q aid=%d cid=%d page=%d part=%q business=%q author=%q view_at=%d progress=%d duration=%d uri=%q",
			i,
			item.Title,
			item.History.BVID,
			item.History.OID,
			item.History.CID,
			item.History.Page,
			item.History.Part,
			item.History.Business,
			item.AuthorName,
			item.ViewAt,
			item.Progress,
			item.Duration,
			item.URI,
		)
	}
}
