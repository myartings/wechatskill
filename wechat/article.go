package wechat

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	publishDateRe = regexp.MustCompile(`var\s+publish_time\s*=\s*"([^"]+)"`)
	// Extract metadata from og: meta tags and JS variables
	ogTitleRe   = regexp.MustCompile(`property="og:title"\s+content="([^"]+)"`)
	ogDescRe    = regexp.MustCompile(`property="og:description"\s+content="([^"]+)"`)
	msgTitleRe  = regexp.MustCompile(`var\s+msg_title\s*=\s*'([^']*)'`)
	msgDescRe   = regexp.MustCompile(`var\s+msg_desc\s*=\s*'([^']*)'`)
	nicknameRe  = regexp.MustCompile(`var\s+nickname\s*=\s*(?:html_decode\()?['"]([^'"]+)['"]`)
	oriHeadRe   = regexp.MustCompile(`var\s+ori_head_img_url\s*=\s*"([^"]*)"`)
	msgSourceRe = regexp.MustCompile(`var\s+msg_source_url\s*=\s*'([^']*)'`)
)

// GetArticleContent fetches a WeChat article and extracts its content.
// If the URL is a Sogou redirect link, it will first resolve it to the real mp.weixin.qq.com URL.
func (c *Client) GetArticleContent(ctx context.Context, articleURL string) (*ArticleDetail, error) {
	// Auto-resolve Sogou redirect links
	if strings.Contains(articleURL, "weixin.sogou.com/link") {
		realURL, err := c.ResolveSogouLink(articleURL)
		if err != nil {
			return nil, fmt.Errorf("resolve sogou link: %w", err)
		}
		articleURL = realURL
	}

	body, err := c.get(articleURL)
	if err != nil {
		return nil, fmt.Errorf("fetch article: %w", err)
	}

	htmlStr := string(body)

	// Detect WeChat anti-crawler verification page
	if strings.Contains(htmlStr, "环境异常") || strings.Contains(htmlStr, "操作频繁") {
		return nil, fmt.Errorf("微信环境异常：当前 IP 被微信限制，需要完成人机验证才能继续访问。建议：1) 稍后重试 2) 通过代理换 IP 3) 让用户直接复制文章内容")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("parse article HTML: %w", err)
	}

	detail := &ArticleDetail{
		URL: articleURL,
	}

	// === Title ===
	// Try DOM elements first
	detail.Title = strings.TrimSpace(doc.Find("#activity-name").Text())
	if detail.Title == "" {
		detail.Title = strings.TrimSpace(doc.Find("h1.rich_media_title").Text())
	}
	// Fallback: og:title meta tag
	if detail.Title == "" {
		if m := ogTitleRe.FindStringSubmatch(htmlStr); len(m) > 1 {
			detail.Title = htmlUnescape(m[1])
		}
	}
	// Fallback: JS variable
	if detail.Title == "" {
		if m := msgTitleRe.FindStringSubmatch(htmlStr); len(m) > 1 {
			detail.Title = htmlUnescape(m[1])
		}
	}

	// === Author / Account name ===
	detail.Author = strings.TrimSpace(doc.Find("#js_author_name").Text())
	if detail.Author == "" {
		detail.Author = strings.TrimSpace(doc.Find("span.rich_media_meta_text").First().Text())
	}

	detail.AccountName = strings.TrimSpace(doc.Find("#js_name").Text())
	if detail.AccountName == "" {
		detail.AccountName = strings.TrimSpace(doc.Find("a.rich_media_meta_link").Text())
	}
	if detail.AccountName == "" {
		if m := nicknameRe.FindStringSubmatch(htmlStr); len(m) > 1 {
			detail.AccountName = htmlUnescape(m[1])
		}
	}

	// === Publish date ===
	if m := publishDateRe.FindStringSubmatch(htmlStr); len(m) > 1 {
		detail.PublishDate = m[1]
	}
	if detail.PublishDate == "" {
		if m := createTimeRe.FindStringSubmatch(htmlStr); len(m) > 1 {
			detail.PublishDate = m[1]
		}
	}
	if detail.PublishDate == "" {
		if m := varCreateTimeRe.FindStringSubmatch(htmlStr); len(m) > 1 {
			ts, err := strconv.ParseInt(m[1], 10, 64)
			if err == nil {
				detail.PublishDate = time.Unix(ts, 0).Format("2006-01-02 15:04")
			}
		}
	}
	if detail.PublishDate == "" {
		detail.PublishDate = strings.TrimSpace(doc.Find("#publish_time").Text())
	}

	// === Content ===
	// Try #js_content first (standard WeChat article body)
	contentEl := doc.Find("#js_content")
	if contentEl.Length() == 0 {
		contentEl = doc.Find("div.rich_media_content")
	}

	if contentEl.Length() > 0 {
		contentEl.Find("script, style").Remove()

		var parts []string
		contentEl.Find("p, section, h1, h2, h3, h4, h5, h6, blockquote, li").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				parts = append(parts, text)
			}
		})

		if len(parts) > 0 {
			detail.Content = strings.Join(parts, "\n\n")
		} else {
			detail.Content = strings.TrimSpace(contentEl.Text())
		}
	}

	// Fallback: use og:description if no content extracted
	if detail.Content == "" {
		if m := ogDescRe.FindStringSubmatch(htmlStr); len(m) > 1 {
			detail.Content = "(摘要) " + htmlUnescape(m[1])
		}
	}
	if detail.Content == "" {
		if m := msgDescRe.FindStringSubmatch(htmlStr); len(m) > 1 {
			detail.Content = "(摘要) " + htmlUnescape(m[1])
		}
	}

	if detail.Title == "" && detail.Content == "" {
		return nil, fmt.Errorf("could not extract article content (article may be expired, require JavaScript rendering, or blocked by verification)")
	}

	return detail, nil
}

// htmlUnescape handles basic HTML entity unescaping.
func htmlUnescape(s string) string {
	r := strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&quot;", `"`,
		"&#39;", "'",
		"&nbsp;", " ",
	)
	return r.Replace(s)
}
