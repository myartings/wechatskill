package wechat

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
	"time"
)

var sogouURLConcatRe = regexp.MustCompile(`url\s*\+=\s*'([^']*)'`)

const (
	SogouBaseURL = "https://weixin.sogou.com"
	UserAgent    = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	jar, _ := cookiejar.New(nil)
	c := &Client{
		httpClient: &http.Client{
			Jar:     jar,
			Timeout: 15 * time.Second,
		},
	}
	// Warm up: visit Sogou WeChat search homepage to collect initial cookies
	// (SNUID, SUID, etc.) that are needed for account searches (type=1).
	c.warmup()
	return c
}

// warmup visits the Sogou WeChat search homepage to populate the cookie jar.
func (c *Client) warmup() {
	req, err := http.NewRequest("GET", SogouBaseURL, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Drain body so the connection can be reused; cookies are already stored in the jar.
	io.Copy(io.Discard, resp.Body)
}

func (c *Client) get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Referer", SogouBaseURL)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	return body, nil
}

// ResolveSogouLink follows a Sogou redirect link (weixin.sogou.com/link?url=...)
// and extracts the real mp.weixin.qq.com URL from the JavaScript redirect page.
func (c *Client) ResolveSogouLink(sogouURL string) (string, error) {
	body, err := c.get(sogouURL)
	if err != nil {
		return "", fmt.Errorf("fetch sogou link: %w", err)
	}

	content := string(body)

	// Sogou uses: var url = ''; url += 'https://mp.'; url += 'weixin.qq.c'; ...
	matches := sogouURLConcatRe.FindAllStringSubmatch(content, -1)
	if len(matches) > 0 {
		var sb strings.Builder
		for _, m := range matches {
			sb.WriteString(m[1])
		}
		realURL := sb.String()
		if strings.Contains(realURL, "mp.weixin.qq.com") {
			return realURL, nil
		}
	}

	return "", fmt.Errorf("could not resolve Sogou redirect URL")
}
