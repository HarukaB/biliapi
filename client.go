package biliapi

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const (
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0 Safari/537.36"
	defaultReferer   = "https://www.bilibili.com/"
	defaultTimeout   = 15 * time.Second
)

type Client struct {
	httpClient *http.Client
	jar        http.CookieJar
	userAgent  string
	referer    string
	creds      Credentials

	wbi *wbiState

	Login      *LoginService
	User       *UserService
	Video      *VideoService
	Search     *SearchService
	Comment    *CommentService
	Danmaku    *DanmakuService
	Fav        *FavService
	History    *HistoryService
	ClientInfo *ClientInfoService
	Misc       *MiscService
}

type Option func(*Client)

func NewClient(opts ...Option) *Client {
	jar, _ := cookiejar.New(nil)
	c := &Client{
		jar:       jar,
		userAgent: defaultUserAgent,
		referer:   defaultReferer,
		wbi:       newWBIState(time.Now),
	}
	c.httpClient = &http.Client{
		Timeout: defaultTimeout,
		Jar:     jar,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{Timeout: defaultTimeout}
	}
	if c.jar == nil {
		if c.httpClient.Jar != nil {
			c.jar = c.httpClient.Jar
		} else {
			jar, _ := cookiejar.New(nil)
			c.jar = jar
			c.httpClient.Jar = jar
		}
	}
	c.bindServices()
	if !c.creds.IsZero() {
		c.SetCredentials(c.creds)
	}
	return c
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient == nil {
			return
		}
		c.httpClient = httpClient
		if httpClient.Jar != nil {
			c.jar = httpClient.Jar
			return
		}
		if c.jar != nil {
			httpClient.Jar = c.jar
		}
	}
}

func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		if userAgent != "" {
			c.userAgent = userAgent
		}
	}
}

func WithReferer(referer string) Option {
	return func(c *Client) {
		if referer != "" {
			c.referer = referer
		}
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if timeout <= 0 {
			return
		}
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
	}
}

func WithCredentials(creds Credentials) Option {
	return func(c *Client) {
		c.creds = creds
	}
}

func (c *Client) SetCredentials(creds Credentials) {
	c.creds = creds
	if c.jar == nil {
		return
	}
	hosts := []string{
		"https://bilibili.com/",
		"https://www.bilibili.com/",
		"https://api.bilibili.com/",
		"https://passport.bilibili.com/",
		"https://space.bilibili.com/",
		"https://s.search.bilibili.com/",
	}
	for _, raw := range hosts {
		u, err := url.Parse(raw)
		if err != nil {
			continue
		}
		c.jar.SetCookies(u, creds.Cookies())
	}
}

func (c *Client) Credentials() Credentials {
	return c.creds
}

func (c *Client) bindServices() {
	c.Login = &LoginService{client: c}
	c.User = &UserService{client: c}
	c.Video = &VideoService{client: c}
	c.Search = &SearchService{client: c}
	c.Comment = &CommentService{client: c}
	c.Danmaku = &DanmakuService{client: c}
	c.Fav = &FavService{client: c}
	c.History = &HistoryService{client: c}
	c.ClientInfo = &ClientInfoService{client: c}
	c.Misc = &MiscService{client: c}
}
