package main

import (
	"context"

	"github.com/myartings/wechatskill/wechat"
)

type WechatService struct {
	client *wechat.Client
}

func NewWechatService(client *wechat.Client) *WechatService {
	return &WechatService{client: client}
}

// SearchArticles searches WeChat articles by keyword via Sogou.
func (s *WechatService) SearchArticles(ctx context.Context, keyword string, page int) ([]wechat.Article, error) {
	return s.client.SearchArticles(ctx, keyword, page)
}

// SearchAccounts searches WeChat official accounts by name via Sogou.
func (s *WechatService) SearchAccounts(ctx context.Context, keyword string, page int) ([]wechat.Account, error) {
	return s.client.SearchAccounts(ctx, keyword, page)
}

// GetAccountArticles fetches recent articles from a specific WeChat official account.
func (s *WechatService) GetAccountArticles(ctx context.Context, accountName string) ([]wechat.Article, error) {
	return s.client.GetAccountArticles(ctx, accountName)
}

// GetArticleContent fetches and extracts the full content of a WeChat article.
func (s *WechatService) GetArticleContent(ctx context.Context, url string) (*wechat.ArticleDetail, error) {
	return s.client.GetArticleContent(ctx, url)
}
