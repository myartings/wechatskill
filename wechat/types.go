package wechat

// Article represents a WeChat article from search results.
type Article struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	AccountName string `json:"account_name"`
	Summary     string `json:"summary"`
	PublishDate string `json:"publish_date"`
	ImageURL    string `json:"image_url,omitempty"`
}

// Account represents a WeChat official account from search results.
type Account struct {
	Name          string `json:"name"`
	WechatID      string `json:"wechat_id"`
	Description   string `json:"description"`
	AvatarURL     string `json:"avatar_url,omitempty"`
	RecentArticle string `json:"recent_article,omitempty"`
	ProfileURL    string `json:"profile_url,omitempty"`
}

// ArticleDetail represents the full content of a WeChat article.
type ArticleDetail struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	AccountName string `json:"account_name"`
	PublishDate string `json:"publish_date"`
	Content     string `json:"content"`
	URL         string `json:"url"`
}
