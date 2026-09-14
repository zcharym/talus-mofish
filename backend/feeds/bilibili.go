package feeds

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type biliEnvelope[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type biliNavData struct {
	IsLogin bool   `json:"isLogin"`
	Uname   string `json:"uname"`
	Mid     int64  `json:"mid"`
}

type biliToViewData struct {
	Count int              `json:"count"`
	List  []biliToViewItem `json:"list"`
}

type biliToViewItem struct {
	BVID    string `json:"bvid"`
	AID     int64  `json:"aid"`
	Title   string `json:"title"`
	Pic     string `json:"pic"`
	Desc    string `json:"desc"`
	Pubdate int64  `json:"pubdate"`
	AddAt   int64  `json:"add_at"`
	Owner   struct {
		Name string `json:"name"`
	} `json:"owner"`
}

type biliSpaceData struct {
	List struct {
		VList []biliSpaceVideo `json:"vlist"`
	} `json:"list"`
}

type biliSpaceVideo struct {
	BVID        string `json:"bvid"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Pic         string `json:"pic"`
	Description string `json:"description"`
	Created     int64  `json:"created"`
}

type biliFavData struct {
	Medias []biliFavMedia `json:"medias"`
}

type biliFavMedia struct {
	ID    int64  `json:"id"`
	BVID  string `json:"bvid"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	Intro string `json:"intro"`
	Ctime int64  `json:"ctime"`
	Upper struct {
		Name string `json:"name"`
	} `json:"upper"`
}

func (c *Client) fetchBilibili(ctx context.Context, spec Spec, sessdata string) (string, []FetchedItem, error) {
	switch {
	case spec.RemoteID == RemoteBilibiliToView:
		if strings.TrimSpace(sessdata) == "" {
			return "", nil, fmt.Errorf("Bilibili SESSDATA is required for watch later (Config → Feeds)")
		}
		return c.fetchBilibiliToView(ctx, sessdata)
	case strings.HasPrefix(spec.RemoteID, "fav:"):
		return c.fetchBilibiliFav(ctx, strings.TrimPrefix(spec.RemoteID, "fav:"), sessdata)
	case strings.HasPrefix(spec.RemoteID, "mid:"):
		return c.fetchBilibiliSpace(ctx, strings.TrimPrefix(spec.RemoteID, "mid:"), sessdata)
	default:
		return "", nil, fmt.Errorf("unsupported Bilibili source %q", spec.RemoteID)
	}
}

func (c *Client) fetchBilibiliToView(ctx context.Context, sessdata string) (string, []FetchedItem, error) {
	body, err := c.getJSON(ctx, "https://api.bilibili.com/x/v2/history/toview", biliCookie(sessdata), map[string]string{
		"Referer": "https://www.bilibili.com/watchlater/",
	})
	if err != nil {
		return "", nil, err
	}
	return parseBiliToView(body)
}

func parseBiliToView(body []byte) (string, []FetchedItem, error) {
	var env biliEnvelope[biliToViewData]
	if err := json.Unmarshal(body, &env); err != nil {
		return "", nil, fmt.Errorf("bilibili watch later: %w", err)
	}
	if env.Code != 0 {
		return "", nil, fmt.Errorf("bilibili watch later: %s (%d)", biliMessage(env.Message, env.Code), env.Code)
	}
	items := make([]FetchedItem, 0, len(env.Data.List))
	for _, entry := range env.Data.List {
		items = append(items, biliItem(entry.BVID, entry.Title, entry.Owner.Name, entry.Desc, entry.Pic, entry.Pubdate, entry.AddAt))
	}
	return "Bilibili watch later", items, nil
}

func (c *Client) fetchBilibiliSpace(ctx context.Context, mid, sessdata string) (string, []FetchedItem, error) {
	endpoint := "https://api.bilibili.com/x/space/arc/search?mid=" + url.QueryEscape(mid) + "&pn=1&ps=30"
	body, err := c.getJSON(ctx, endpoint, biliCookie(sessdata), map[string]string{
		"Referer": "https://space.bilibili.com/" + mid,
	})
	if err != nil {
		return "", nil, err
	}
	var env biliEnvelope[biliSpaceData]
	if err := json.Unmarshal(body, &env); err != nil {
		return "", nil, fmt.Errorf("bilibili space: %w", err)
	}
	if env.Code != 0 {
		return "", nil, fmt.Errorf("bilibili uploads: %s (%d). Paste an RSS URL instead if this API is blocked", biliMessage(env.Message, env.Code), env.Code)
	}
	items := make([]FetchedItem, 0, len(env.Data.List.VList))
	for _, entry := range env.Data.List.VList {
		items = append(items, biliItem(entry.BVID, entry.Title, entry.Author, entry.Description, entry.Pic, entry.Created, 0))
	}
	return "Bilibili " + mid, items, nil
}

func (c *Client) fetchBilibiliFav(ctx context.Context, mediaID, sessdata string) (string, []FetchedItem, error) {
	endpoint := "https://api.bilibili.com/x/v3/fav/resource/list?media_id=" + url.QueryEscape(mediaID) + "&ps=20&pn=1"
	body, err := c.getJSON(ctx, endpoint, biliCookie(sessdata), map[string]string{
		"Referer": "https://www.bilibili.com/",
	})
	if err != nil {
		return "", nil, err
	}
	var env biliEnvelope[biliFavData]
	if err := json.Unmarshal(body, &env); err != nil {
		return "", nil, fmt.Errorf("bilibili favorites: %w", err)
	}
	if env.Code != 0 {
		return "", nil, fmt.Errorf("bilibili favorites: %s (%d)", biliMessage(env.Message, env.Code), env.Code)
	}
	items := make([]FetchedItem, 0, len(env.Data.Medias))
	for _, entry := range env.Data.Medias {
		items = append(items, biliItem(entry.BVID, entry.Title, entry.Upper.Name, entry.Intro, entry.Cover, entry.Ctime, 0))
	}
	return "Bilibili favorites", items, nil
}

func (c *Client) pingBilibili(ctx context.Context, sessdata string) (PingResult, error) {
	if strings.TrimSpace(sessdata) == "" {
		return PingResult{}, fmt.Errorf("Bilibili SESSDATA cookie is not set")
	}
	body, err := c.getJSON(ctx, "https://api.bilibili.com/x/web-interface/nav", biliCookie(sessdata), map[string]string{
		"Referer": "https://www.bilibili.com/",
	})
	if err != nil {
		return PingResult{}, err
	}
	var env biliEnvelope[biliNavData]
	if err := json.Unmarshal(body, &env); err != nil {
		return PingResult{}, err
	}
	if env.Code != 0 || !env.Data.IsLogin {
		return PingResult{}, fmt.Errorf("Bilibili login failed: %s", biliMessage(env.Message, env.Code))
	}
	detail := strings.TrimSpace(env.Data.Uname)
	if detail == "" {
		detail = strconv.FormatInt(env.Data.Mid, 10)
	}
	return PingResult{OK: true, Service: "bilibili", Detail: "Signed in as " + detail}, nil
}

func biliCookie(sessdata string) string {
	sessdata = strings.TrimSpace(sessdata)
	if sessdata == "" {
		return ""
	}
	if strings.Contains(strings.ToLower(sessdata), "sessdata=") {
		return sessdata
	}
	return "SESSDATA=" + sessdata
}

func biliMessage(message string, code int) string {
	message = strings.TrimSpace(message)
	if message == "" || message == "0" {
		return fmt.Sprintf("error %d", code)
	}
	return message
}

func biliItem(bvid, title, author, desc, pic string, published, added int64) FetchedItem {
	bvid = strings.TrimSpace(bvid)
	ts := published
	if ts == 0 {
		ts = added
	}
	publishedAt := time.Time{}
	if ts > 0 {
		publishedAt = time.Unix(ts, 0).UTC()
	}
	id := bvid
	if id == "" {
		id = title
	}
	link := ""
	if bvid != "" {
		link = "https://www.bilibili.com/video/" + bvid
	}
	return FetchedItem{
		ExternalID:   id,
		Title:        strings.TrimSpace(title),
		URL:          link,
		Author:       strings.TrimSpace(author),
		Summary:      sanitizeSummary(desc),
		ThumbnailURL: strings.TrimSpace(pic),
		PublishedAt:  publishedAt,
	}
}
