package wechat

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var timeConvertRe = regexp.MustCompile(`timeConvert\('(\d+)'\)`)

// SearchArticles searches for WeChat articles via Sogou WeChat search.
// type=2 searches articles.
func (c *Client) SearchArticles(ctx context.Context, keyword string, page int) ([]Article, error) {
	if page < 1 {
		page = 1
	}
	searchURL := fmt.Sprintf("%s/weixin?type=2&query=%s&page=%d",
		SogouBaseURL, url.QueryEscape(keyword), page)

	body, err := c.get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("search articles: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("parse HTML: %w", err)
	}

	var articles []Article

	doc.Find("div.txt-box").Each(func(i int, s *goquery.Selection) {
		article := Article{}

		// Title and URL
		titleEl := s.Find("h3 a")
		article.Title = strings.TrimSpace(titleEl.Text())
		if href, exists := titleEl.Attr("href"); exists {
			article.URL = href
		}

		// Account name
		accountEl := s.Find("div.s-p a[data-z]")
		if accountEl.Length() == 0 {
			accountEl = s.Find("div.s-p a")
		}
		article.AccountName = strings.TrimSpace(accountEl.Text())

		// Summary
		summaryEl := s.Find("p.txt-info")
		article.Summary = strings.TrimSpace(summaryEl.Text())

		// Publish date
		dateEl := s.Find("span.s2")
		if dateEl.Length() == 0 {
			dateEl = s.Find("div.s-p span")
		}
		article.PublishDate = strings.TrimSpace(dateEl.Text())

		if article.Title != "" {
			articles = append(articles, article)
		}
	})

	// Fallback: try alternative selector for news-box layout
	if len(articles) == 0 {
		doc.Find("ul.news-list li").Each(func(i int, s *goquery.Selection) {
			article := Article{}
			titleEl := s.Find("h3 a")
			article.Title = strings.TrimSpace(titleEl.Text())
			if href, exists := titleEl.Attr("href"); exists {
				article.URL = href
			}
			article.Summary = strings.TrimSpace(s.Find("p.txt-info").Text())
			article.AccountName = strings.TrimSpace(s.Find("a.account").Text())

			if article.Title != "" {
				articles = append(articles, article)
			}
		})
	}

	// Post-process: fix relative URLs and JS date strings
	for i := range articles {
		articles[i].URL = fixURL(articles[i].URL)
		articles[i].PublishDate = fixDate(articles[i].PublishDate)
	}

	return articles, nil
}

// fixURL converts relative Sogou URLs to absolute.
func fixURL(href string) string {
	if href == "" {
		return ""
	}
	if strings.HasPrefix(href, "/") {
		return SogouBaseURL + href
	}
	return href
}

// fixDate extracts unix timestamp from JS code like "document.write(timeConvert('1774368920'))"
// and converts it to a human-readable date.
func fixDate(raw string) string {
	if raw == "" {
		return ""
	}
	matches := timeConvertRe.FindStringSubmatch(raw)
	if len(matches) > 1 {
		ts, err := strconv.ParseInt(matches[1], 10, 64)
		if err == nil {
			return time.Unix(ts, 0).Format("2006-01-02 15:04")
		}
	}
	return raw
}

// GetAccountArticles searches for recent articles from a specific WeChat official account.
// It uses Sogou article search (type=2) with the account name as keyword,
// then filters results to only include articles from the matching account.
func (c *Client) GetAccountArticles(ctx context.Context, accountName string) ([]Article, error) {
	// Search articles using the account name as keyword — this returns
	// more recent results than scraping the Sogou profile page.
	allArticles, err := c.SearchArticles(ctx, accountName, 1)
	if err != nil {
		return nil, fmt.Errorf("search articles for account: %w", err)
	}

	// Filter to only keep articles from the target account
	var matched []Article
	target := strings.ToLower(accountName)
	for _, a := range allArticles {
		if strings.ToLower(a.AccountName) == target ||
			strings.Contains(strings.ToLower(a.AccountName), target) ||
			strings.Contains(target, strings.ToLower(a.AccountName)) {
			matched = append(matched, a)
		}
	}

	// If strict matching found nothing, return all results (the keyword
	// search already scoped them to be relevant to the account name).
	if len(matched) == 0 {
		return allArticles, nil
	}
	return matched, nil
}

// SearchAccounts searches for WeChat official accounts via Sogou WeChat search.
// type=1 searches accounts.
func (c *Client) SearchAccounts(ctx context.Context, keyword string, page int) ([]Account, error) {
	if page < 1 {
		page = 1
	}
	searchURL := fmt.Sprintf("%s/weixin?type=1&query=%s&page=%d",
		SogouBaseURL, url.QueryEscape(keyword), page)

	body, err := c.get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("search accounts: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("parse HTML: %w", err)
	}

	var accounts []Account

	doc.Find("div.txt-box").Each(func(i int, s *goquery.Selection) {
		account := Account{}

		// Account name
		nameEl := s.Find("h3 a")
		account.Name = strings.TrimSpace(nameEl.Text())
		if href, exists := nameEl.Attr("href"); exists {
			if strings.HasPrefix(href, "/") {
				account.ProfileURL = SogouBaseURL + href
			} else {
				account.ProfileURL = href
			}
		}

		// WeChat ID
		idEl := s.Find("label[name='em_weixinhao']")
		account.WechatID = strings.TrimSpace(idEl.Text())

		// Description
		descEl := s.Find("dl:nth-child(1) dd")
		account.Description = strings.TrimSpace(descEl.Text())

		// Avatar
		imgEl := s.Parent().Find("div.img-box img")
		if src, exists := imgEl.Attr("src"); exists {
			account.AvatarURL = src
		}

		// Recent article
		recentEl := s.Find("dl:last-child dd a")
		account.RecentArticle = strings.TrimSpace(recentEl.Text())

		if account.Name != "" {
			accounts = append(accounts, account)
		}
	})

	// Fallback selector
	if len(accounts) == 0 {
		doc.Find("div.gzh-box2").Each(func(i int, s *goquery.Selection) {
			account := Account{}
			account.Name = strings.TrimSpace(s.Find("p.tit a").Text())
			account.WechatID = strings.TrimSpace(s.Find("label[name='em_weixinhao']").Text())
			account.Description = strings.TrimSpace(s.Find("dl:first-of-type dd").Text())
			if href, exists := s.Find("p.tit a").Attr("href"); exists {
				if strings.HasPrefix(href, "/") {
					account.ProfileURL = SogouBaseURL + href
				} else {
					account.ProfileURL = href
				}
			}

			if account.Name != "" {
				accounts = append(accounts, account)
			}
		})
	}

	return accounts, nil
}
