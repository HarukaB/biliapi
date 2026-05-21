# Go API 文档

本文面向直接调用 `gopkg.yay.moe/biliapi` 的 Go 程序。它只描述当前包导出的 Go API；具体返回字段以源码中的结构体定义和 B 站线上返回为准。

## 安装与导入

```powershell
go get gopkg.yay.moe/biliapi
```

```go
import "gopkg.yay.moe/biliapi"
```

模块要求 Go 1.18 或更新版本。

## 创建客户端

```go
client := biliapi.NewClient()
```

带登录态：

```go
client := biliapi.NewClient(biliapi.WithCredentials(biliapi.Credentials{
	SESSDATA: "your SESSDATA",
	BiliJCT:  "your bili_jct",
}))
```

带自定义超时：

```go
client := biliapi.NewClient(
	biliapi.WithTimeout(20*time.Second),
	biliapi.WithUserAgent("my-app/1.0"),
)
```

`NewClient` 会初始化默认 `http.Client`、cookie jar、默认 User-Agent、默认 Referer，并绑定所有服务字段：

```go
client.Login
client.User
client.Video
client.Search
client.Comment
client.Danmaku
client.Fav
client.History
client.ClientInfo
client.Misc
```

### Client 选项

| 选项 | 用途 |
| --- | --- |
| `WithHTTPClient(*http.Client)` | 使用调用方提供的 HTTP client。若它没有 `Jar`，包会复用自己的 cookie jar。 |
| `WithTimeout(time.Duration)` | 设置 HTTP client 超时；非正数会被忽略。 |
| `WithUserAgent(string)` | 覆盖默认 User-Agent；空字符串会被忽略。 |
| `WithReferer(string)` | 覆盖默认 Referer；空字符串会被忽略。 |
| `WithCredentials(Credentials)` | 初始化时写入登录凭据，并同步到 cookie jar。 |

### Client 方法

| 方法 | 说明 |
| --- | --- |
| `NewClient(opts ...Option) *Client` | 创建客户端。 |
| `(*Client).SetCredentials(Credentials)` | 更新凭据，并把 Cookie 写入 B 站相关域名。 |
| `(*Client).Credentials() Credentials` | 读取当前凭据副本。 |

## 凭据与登录态

`Credentials` 表示 B 站 Web Cookie 中常用的登录字段：

```go
type Credentials struct {
	SESSDATA        string
	BiliJCT         string
	DedeUserID      string
	DedeUserIDCKMd5 string
	SID             string
	RefreshToken    string
}
```

常用方法：

| 方法 | 说明 |
| --- | --- |
| `ParseCredentialsJSON([]byte) (Credentials, error)` | 从 JSON 反序列化凭据。 |
| `(Credentials).IsZero() bool` | 判断是否没有任何凭据字段。 |
| `(Credentials).IsLoggedInCandidate() bool` | 判断是否至少带有 `SESSDATA`，可尝试登录态接口。 |
| `(Credentials).CSRF() string` | 返回 `bili_jct`。 |
| `(Credentials).RequireCSRF() error` | 缺少 `bili_jct` 时返回 `ErrMissingCSRF`。 |
| `(Credentials).Cookies() []*http.Cookie` | 转成 `.bilibili.com` 域名下的 Cookie。 |

只读登录接口通常需要 `SESSDATA`。会修改账号状态的接口通常同时需要 `SESSDATA` 和 `BiliJCT`，缺失时会先在本地返回错误，不会发出请求。

`RefreshToken` 当前只是随 `Credentials` 解析和保存的字段。这个包尚未实现 Web 端 Cookie 刷新流程，不会自动检查 Cookie 是否需要刷新，不会调用刷新接口，也不会在收到 `-101` 时自动换取新的 `SESSDATA` / `bili_jct`。`RefreshToken` 也不会被写入 Cookie。

上游参考文档中 Cookie 刷新涉及 `cookie/info`、获取 `refresh_csrf`、`cookie/refresh`、`confirm/refresh` 和 SSO 同步等步骤；当前 API 没有封装这些步骤。

## 错误处理

包内有三个哨兵错误：

| 错误 | 触发场景 |
| --- | --- |
| `ErrMissingCredentials` | 方法需要登录态，但当前凭据不足。 |
| `ErrMissingCSRF` | 修改类方法需要 `bili_jct`，但当前凭据没有提供。 |
| `ErrInvalidParams` | 调用参数本身非法，例如同时传入 `AID` 和 `BVID`。 |

B 站业务错误和非 2xx HTTP 响应用 `*BiliError` 表示：

```go
var biliErr *biliapi.BiliError
if errors.As(err, &biliErr) {
	fmt.Println(biliErr.HTTPStatus, biliErr.Code, biliErr.Message)
}
```

`BiliError.Data` 保留原始 `data` 字段，便于排查接口返回的额外信息。

## 调用约定

- 所有网络方法都接收 `context.Context`，取消或超时由调用方控制。
- `VideoViewParams` 这类 ID 参数通常接受 `AID` 或 `BVID` 二选一；同时为空或同时设置会返回 `ErrInvalidParams`。
- 参数结构体里的零值一般不会写入 query/form。需要传 `0` 作为业务值的接口要特别确认源码实现。
- WBI 签名接口由包自动处理；调用方不需要手动传 `wts` 或 `w_rid`。
- 部分接口返回 XML、JavaScript 或二进制段，API 会直接返回结构体或原始 `[]byte`，不会强行包成 JSON。
- B 站 Web API 存在字段漂移；文档中列出的 `json.RawMessage` 字段表示当前包保留原始 JSON，调用方可按需二次解析。

## 快速示例

读取视频信息：

```go
ctx := context.Background()
client := biliapi.NewClient()

view, err := client.Video.View(ctx, biliapi.VideoViewParams{BVID: "BV1xx411c7mD"})
if err != nil {
	return err
}
fmt.Println(view.Title, view.Owner.Name, view.CID)
```

获取搜索建议：

```go
suggest, err := client.Search.Suggest(ctx, biliapi.SearchSuggestParams{Term: "go"})
if err != nil {
	return err
}
for _, item := range suggest.Result.Tag {
	fmt.Println(item.Value)
}
```

调用需要登录态的接口：

```go
client := biliapi.NewClient(biliapi.WithCredentials(biliapi.Credentials{
	SESSDATA: sessdata,
	BiliJCT:  biliJCT,
}))

nav, err := client.Login.Nav(ctx)
if err != nil {
	return err
}
fmt.Println(nav.IsLogin, nav.Uname)
```

## 服务总览

| 服务字段 | 主要用途 |
| --- | --- |
| `Login` | 登录状态、导航栏信息、关注/粉丝/动态计数。 |
| `User` | 用户空间、名片、投稿、关系统计、空间设置。 |
| `Video` | 视频详情、播放地址、统计、点赞投币收藏、稍后再看相关操作。 |
| `Search` | 默认搜索词、热搜、建议、综合搜索、分类搜索。 |
| `Comment` | 评论列表、楼中楼、评论操作、举报。 |
| `Danmaku` | XML 弹幕、分段弹幕、历史弹幕、弹幕发送/撤回/点赞。 |
| `Fav` | 收藏夹信息、收藏夹内容、收藏夹创建编辑、资源移动复制删除。 |
| `History` | 历史记录、历史开关、稍后再看。 |
| `ClientInfo` | IP、地区和客户端网络信息。 |
| `Misc` | 时间戳、buvid、短链、MathJax、服务端日期脚本等杂项接口。 |

下文表格中的认证标记：

| 标记 | 含义 |
| --- | --- |
| Public | 不要求登录态。 |
| Login | 需要 `SESSDATA`。 |
| CSRF | 需要 `SESSDATA` 和 `BiliJCT`，会修改账号状态。 |
| WBI | 调用时自动进行 WBI 签名。 |

## Login API

| 方法 | 认证 | 返回 | 说明 |
| --- | --- | --- | --- |
| `Nav(ctx)` | Public | `*NavInfo` | 获取导航栏账号信息。未登录时允许 B 站返回 `-101`，可用 `IsLogin` 判断。返回中包含 `WBIImg`，内部也用它更新 WBI key。 |
| `NavStat(ctx)` | Login | `*NavStat` | 获取关注数、粉丝数、动态数。 |

## ClientInfo API

| 方法 | 认证 | 返回 | 说明 |
| --- | --- | --- | --- |
| `Zone(ctx)` | Public | `*ZoneInfo` | 当前请求 IP 的地区信息。 |
| `LiveIPInfo(ctx)` | Public | `*LiveIPInfo` | 直播接口侧的 IP 信息。 |
| `AppIP(ctx)` | Public | `*LiveIPInfo` | App 资源接口侧的 IP 信息。 |
| `IPInfo(ctx, IPInfoParams)` | Public | `*LiveIPInfo` | 查询指定 IP；`IPInfoParams.IP` 为空时按接口默认行为处理。 |

## Misc API

| 方法 | 认证 | 返回 | 说明 |
| --- | --- | --- | --- |
| `Now(ctx)` | Public | `*Timestamp` | 点击接口当前时间。 |
| `ReportNow(ctx)` | Public | `*ReportTimestamp` | 上报接口当前时间。 |
| `Buvid(ctx)` | Public | `*BuvidInfo` | 获取 `b_3`、`b_4`、`b_nut`。 |
| `Buvid3(ctx)` | Public | `*Buvid3Info` | 获取 `buvid`。 |
| `RTCTimestamp(ctx)` | Public | `*RTCTimestamp` | 直播 RTC 时间戳。 |
| `ServerDateJS(ctx)` | Public | `[]byte` | 获取 `serverdate.js` 原始内容。 |
| `MathJax(ctx, MathJaxParams)` | Public | `*MathJaxResult` | 将 TeX 转成 SVG；参数字段为 `Tex`。 |
| `ShortLink(ctx, ShortLinkParams)` | Public | `*ShortLinkInfo` | 生成或解析分享短链；参数可带 `Build`、`BVID`、`AID`、`URL`。 |

## Search API

| 方法 | 认证 | 返回 | 说明 |
| --- | --- | --- | --- |
| `Default(ctx)` | Public, WBI | `*SearchDefault` | 默认搜索词。 |
| `Square(ctx, params ...SearchSquareParams)` | Public, WBI | `*SearchSquare` | 热搜面板。无参数时默认 `Limit: 10`、`Platform: "web"`。 |
| `Hotword(ctx)` | Public | `*SearchHotword` | 搜索热词接口，返回结构不是标准 B 站 envelope。 |
| `Suggest(ctx, SearchSuggestParams)` | Public | `*SearchSuggest` | 搜索建议；常用字段 `Term`、`MainVer`、`Highlight`。 |
| `All(ctx, SearchAllParams)` | Public, WBI | `*SearchAll` | 综合搜索；常用字段 `Keyword`、`Page`、`PageSize`。 |
| `Type(ctx, SearchTypeParams)` | Public, WBI | `*SearchTypeResult` | 分类搜索；`SearchType` 可使用包内常量。 |

`SearchType` 常量：

```go
biliapi.SearchTypeVideo
biliapi.SearchTypeBangumi
biliapi.SearchTypeMovie
biliapi.SearchTypeLive
biliapi.SearchTypeArticle
biliapi.SearchTypeTopic
biliapi.SearchTypeUser
```

## User API

| 方法 | 认证 | 返回 | 说明 |
| --- | --- | --- | --- |
| `SpaceInfo(ctx, UserSpaceInfoParams)` | Public, WBI | `*UserSpaceInfo` | 用户空间资料；`Mid` 必须为正数。 |
| `Card(ctx, UserCardParams)` | Public | `*UserCard` | 用户名片；`Mid` 必须为正数，可带 `Photo`、`Relation`。 |
| `MyInfo(ctx)` | Login | `*UserSpaceInfo` | 当前登录用户资料。 |
| `Cards(ctx, UserCardsParams)` | Public | `*UserCards` | 批量用户名片；`Mids` 会编码为逗号分隔。 |
| `ArcSearch(ctx, UserArcSearchParams)` | Public, WBI | `*UserArcSearch` | 用户投稿列表；`Mid` 必须为正数。 |
| `UpStat(ctx, mid)` | Public | `*UserStat` | UP 主稿件和点赞统计。 |
| `RelationStat(ctx, vmid)` | Public | `*RelationStat` | 关注、粉丝等关系统计。 |
| `NavNum(ctx, mid)` | Public | `*NavNum` | 用户空间导航数量。 |
| `SetNotice(ctx, notice)` | CSRF | `error` | 修改个人公告。 |
| `SetTags(ctx, tags)` | CSRF | `error` | 修改个人标签。 |
| `SetTopArchive(ctx, aid)` | CSRF | `error` | 设置置顶稿件。 |
| `CancelTopArchive(ctx)` | CSRF | `error` | 取消置顶稿件。 |

## Video API

| 方法 | 认证 | 返回 | 说明 |
| --- | --- | --- | --- |
| `View(ctx, VideoViewParams)` | Public, WBI | `*VideoView` | 视频基础详情。`AID` 和 `BVID` 二选一。 |
| `ViewDetail(ctx, VideoViewParams)` | Public, WBI | `*VideoDetail` | 视频详情聚合数据，包含 `View`、相关推荐等字段。 |
| `Desc(ctx, VideoViewParams)` | Public | `string` | 视频简介。 |
| `PageList(ctx, VideoViewParams)` | Public | `[]Page` | 视频分 P 列表。 |
| `PlayURL(ctx, PlayURLParams)` | Public, WBI | `*PlayURL` | 播放地址。需要 `AID`/`BVID` 二选一和正数 `CID`。 |
| `PlayerInfo(ctx, PlayerInfoParams)` | Public, WBI | `*PlayerInfo` | 播放器信息。需要 `AID`/`BVID` 二选一和正数 `CID`。 |
| `Stat(ctx, VideoViewParams)` | Public | `*ArchiveStat` | 视频统计。 |
| `HasLike(ctx, VideoViewParams)` | Login | `*VideoHasLike` | 当前账号是否点赞。 |
| `Coins(ctx, VideoViewParams)` | Login | `*VideoCoins` | 当前账号投币数量。 |
| `Favoured(ctx, aid)` | Login | `*VideoFavoured` | 当前账号是否收藏指定 `aid`。 |
| `Related(ctx, VideoViewParams)` | Public | `[]VideoView` | 相关推荐。 |
| `OnlineTotal(ctx, aid, cid, bvid)` | Public | `*OnlineTotal` | 在线人数。`aid`/`bvid` 二选一，`cid` 必须为正数。 |
| `Shot(ctx, aid, cid, bvid)` | Public | `*VideoShot` | 视频快照索引。 |
| `AIConclusion(ctx, bvid, cid, upMid)` | Public | `*AIConclusion` | AI 总结结果。 |
| `Like(ctx, VideoLikeParams)` | CSRF | `error` | 点赞或取消点赞；`Like` 为 `true` 表示点赞，`false` 表示取消。 |
| `Coin(ctx, VideoCoinParams)` | CSRF | `error` | 投币；`Multiply` 为投币数量，`SelectLike` 表示同时点赞。 |
| `Favorite(ctx, VideoFavoriteParams)` | CSRF | `error` | 添加或移除收藏夹；使用 `AddMediaIDs`、`DelMediaIDs`。 |
| `Triple(ctx, VideoViewParams)` | CSRF | `error` | 一键三连。 |

`PlayURLParams` 常用字段：

```go
biliapi.PlayURLParams{
	BVID:        "BV...",
	CID:         123,
	QN:          80,
	FnVal:       16,
	FourK:       true,
	HighQuality: true,
}
```

## Comment API

| 方法 | 认证 | 返回 | 说明 |
| --- | --- | --- | --- |
| `List(ctx, CommentListParams)` | Public | `*CommentList` | 传统评论列表。 |
| `Main(ctx, CommentMainParams)` | Public, WBI | `*CommentMain` | 主评论列表，支持 cursor。 |
| `Replies(ctx, CommentRepliesParams)` | Public | `*CommentReplies` | 楼中楼回复列表。 |
| `Dialog(ctx, CommentDialogParams)` | Public | `*CommentReplies` | 对话串回复列表。 |
| `Hot(ctx, CommentListParams)` | Public | `*CommentList` | 热门评论。 |
| `Info(ctx, oid, typ, rpid)` | Public | `*CommentReply` | 单条评论信息。该接口可能返回包裹结构或空载荷。 |
| `Count(ctx, oid, typ)` | Public | `*CommentCount` | 评论数。 |
| `Add(ctx, CommentAddParams)` | CSRF | `*CommentReply` | 发表评论或回复。 |
| `Action(ctx, CommentActionParams)` | CSRF | `error` | 评论点赞等 action。 |
| `Hate(ctx, CommentActionParams)` | CSRF | `error` | 评论点踩或取消点踩。 |
| `Delete(ctx, oid, typ, rpid)` | CSRF | `error` | 删除评论。 |
| `Top(ctx, oid, typ, rpid, action)` | CSRF | `error` | 置顶或取消置顶评论。 |
| `Report(ctx, oid, typ, rpid, reason, content)` | CSRF | `error` | 举报评论。 |

`CommentType` 常量：

```go
biliapi.CommentTypeVideo
biliapi.CommentTypeArticle
biliapi.CommentTypeDynamic
```

## Danmaku API

| 方法 | 认证 | 返回 | 说明 |
| --- | --- | --- | --- |
| `XML(ctx, DanmakuXMLParams)` | Public | `*DanmakuXML` | 通过 XML list 接口读取弹幕；`CID` 必须为正数。 |
| `XMLByCID(ctx, cid)` | Public | `*DanmakuXML` | 通过 `comment.bilibili.com/{cid}.xml` 读取弹幕。 |
| `Segment(ctx, DanmakuSegmentParams)` | Public, WBI | `[]byte` | 分段弹幕原始字节，调用方自行解析 protobuf。 |
| `View(ctx, DanmakuViewParams)` | Public | `*DanmakuView` | 弹幕 view 信息。 |
| `HistoryIndex(ctx, DanmakuHistoryIndexParams)` | Login | `*DanmakuHistoryIndex` | 可查询历史弹幕的月份。 |
| `HistorySegment(ctx, DanmakuHistoryParams)` | Login | `[]byte` | 历史分段弹幕原始字节。 |
| `HistoryXML(ctx, DanmakuHistoryParams)` | Login | `*DanmakuXML` | 历史 XML 弹幕。 |
| `ThumbStats(ctx, DanmakuThumbStatsParams)` | Public | `*DanmakuThumbStats` | 弹幕点赞统计。 |
| `Post(ctx, DanmakuPostParams)` | CSRF | `error` | 发送弹幕。 |
| `Recall(ctx, cid, dmid)` | CSRF | `error` | 撤回弹幕。 |
| `ThumbUp(ctx, oid, dmid, up)` | CSRF | `error` | 点赞或取消点赞弹幕。 |

XML 弹幕会先尝试按 XML 直接解析，失败后再尝试 GB18030 解码。

## Fav API

| 方法 | 认证 | 返回 | 说明 |
| --- | --- | --- | --- |
| `FolderInfo(ctx, FavFolderInfoParams)` | Public | `*FavFolder` | 收藏夹信息。 |
| `CreatedListAll(ctx, FavCreatedListAllParams)` | Public | `*FavCreatedListAll` | 指定用户创建的收藏夹。 |
| `CollectedList(ctx, FavCollectedListParams)` | Public | `*FavCollectedList` | 指定用户收藏的收藏夹。 |
| `ResourceList(ctx, FavResourceListParams)` | Public | `*FavResourceList` | 收藏夹内容列表。 |
| `ResourceIDs(ctx, mediaID)` | Public | `[]FavResourceRef` | 收藏夹内资源 ID 列表。 |
| `ResourceInfos(ctx, FavResourceInfosParams)` | Public | `[]FavResource` | 批量资源详情。 |
| `AddFolder(ctx, FavFolderAddParams)` | CSRF | `*FavFolder` | 创建收藏夹。 |
| `EditFolder(ctx, mediaID, FavFolderAddParams)` | CSRF | `error` | 编辑收藏夹。 |
| `DeleteFolder(ctx, mediaID)` | CSRF | `error` | 删除收藏夹。 |
| `MoveResources(ctx, FavResourceMoveParams)` | CSRF | `error` | 移动收藏夹资源。 |
| `CopyResources(ctx, FavResourceMoveParams)` | CSRF | `error` | 复制收藏夹资源。 |
| `DeleteResources(ctx, mediaID, resources)` | CSRF | `error` | 批量删除收藏夹资源。 |
| `CleanInvalidResources(ctx, mediaID)` | CSRF | `error` | 清理失效收藏。 |

`FavResourceID` 用于编码资源引用：

```go
biliapi.FavResourceID{ID: 1700001, Type: 2}
```

## History API

| 方法 | 认证 | 返回 | 说明 |
| --- | --- | --- | --- |
| `Cursor(ctx, HistoryCursorParams)` | Login | `*HistoryCursor` | cursor 版历史记录。 |
| `Legacy(ctx, HistoryLegacyParams)` | Login | `*HistoryLegacy` | 旧版历史记录。返回可能是数组或对象，包内会兼容。 |
| `Shadow(ctx)` | Login | `*ShadowStatus` | 查询历史记录暂停状态。 |
| `ToView(ctx)` | Login | `*ToViewList` | 稍后再看列表。 |
| `Delete(ctx, []HistoryDeleteItem)` | CSRF | `error` | 删除历史记录。 |
| `Clear(ctx)` | CSRF | `error` | 清空历史记录。 |
| `SetShadow(ctx, shadow)` | CSRF | `error` | 开启或关闭历史记录暂停。 |
| `AddToView(ctx, aid)` | CSRF | `error` | 加入稍后再看。 |
| `DeleteToView(ctx, aid)` | CSRF | `error` | 从稍后再看删除。 |
| `ClearToView(ctx)` | CSRF | `error` | 清空稍后再看。 |

## 返回模型与原始字段

常用公共模型集中在 `models.go`，包括：

| 类型 | 说明 |
| --- | --- |
| `Owner` | 视频或收藏资源作者。 |
| `OfficialInfo` | 认证信息。 |
| `VipInfo`、`VipLabel` | 大会员信息。 |
| `ArchiveStat`、`ArchiveRights` | 视频统计和权限。 |
| `Page`、`Dimension` | 分 P 和尺寸信息。 |
| `Cursor`、`CursorString` | 部分接口的 cursor 数据。 |
| `FlexibleString` | 兼容线上返回字符串、数字、`null` 的字段。 |
| `MixedStringInt` | `json.RawMessage` 别名，用于保留可能是字符串或数字的字段。 |

如果结构体字段类型是 `json.RawMessage`，说明该字段的线上结构可能较复杂或不稳定。调用方可以按业务场景二次解析：

```go
var payload map[string]any
if len(view.UGCSeason) > 0 {
	_ = json.Unmarshal(view.UGCSeason, &payload)
}
```

## 测试与验证

文档对应的公共 API 可用基础测试验证：

```powershell
go test ./...
go vet ./...
```

真实接口集成测试会访问 B 站线上接口：

```powershell
op run --env-file=.env -- go test ./... -run Integration -count=1 -v
```

凭据相关测试读取以下环境变量之一：

| 环境变量 | 说明 |
| --- | --- |
| `BILIAPI_CREDENTIALS_JSON` | JSON 格式的 `Credentials`。 |
| `BILIAPI_SESSDATA` | `SESSDATA` Cookie。 |
| `BILIAPI_BILI_JCT` | `bili_jct` CSRF token。 |
| `BILIAPI_DEDE_USER_ID` | `DedeUserID` Cookie。 |
| `BILIAPI_DEDE_USER_ID_CKMD5` | `DedeUserID__ckMd5` Cookie。 |
| `BILIAPI_SID` | `sid` Cookie。 |
| `BILIAPI_REFRESH_TOKEN` | refresh token；当前只会读入 `Credentials.RefreshToken`，不会自动刷新 Cookie。 |

集成测试还可用 `BILIAPI_TEST_BVID`、`BILIAPI_TEST_AIDS`、`BILIAPI_TEST_FAV_ID`、`BILIAPI_TEST_SEARCH_KEYWORD` 指定真实测试对象。
