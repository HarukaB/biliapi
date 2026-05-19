# biliapi

Go client for Bilibili Web REST APIs

```go
client := biliapi.NewClient(biliapi.WithCredentials(biliapi.Credentials{
	SESSDATA: "your SESSDATA",
	BiliJCT:  "your bili_jct",
}))

view, err := client.Video.View(ctx, biliapi.VideoViewParams{BVID: "BV..."})
```

## Scope

- Web REST APIs first.
- Typed services: `Login`, `User`, `Video`, `Search`, `Comment`, `Danmaku`, `Fav`, `History`, `ClientInfo`, `Misc`.
- WBI signing is handled automatically for WBI endpoints.
- Cookies and CSRF are represented by `Credentials`; callers should not pass raw Cookie strings through business APIs.
- Mutating methods are present as request skeletons and require `Credentials.BiliJCT`, but tests do not call live mutating endpoints.

## Attribution

API endpoint documentation and examples are based on [rinnein/bilibili-API-collect](https://github.com/rinnein/bilibili-API-collect.git).

## Validation

```powershell
go test ./...
go vet ./...
```

Real API credential test:

```powershell
op run --env-file=.env -- go test ./... -run Integration -count=1 -v
```

Alternatively set `BILIAPI_CREDENTIALS_JSON` to a JSON object matching `Credentials`.
Set `BILIAPI_TEST_AIDS` to a comma-separated list such as `1700001,av1700002`
to exercise real av ids in the video integration chain.
Set `BILIAPI_TEST_FAV_ID` to exercise a real favorite folder media id; the
default is `2519150301`.

The integration suite always exercises public real APIs for misc, client info,
search, video, comments, danmaku, and favorites. Credential-only checks are
skipped unless credentials are provided through the environment.
