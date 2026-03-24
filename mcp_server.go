package main

import (
	"encoding/json"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func newMCPServer(app *AppServer) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "wechat-mcp",
		Version: "1.0.0",
	}, nil)

	inputSchema := func(props map[string]any, required []string) json.RawMessage {
		schema := map[string]any{
			"type":       "object",
			"properties": props,
		}
		if len(required) > 0 {
			schema["required"] = required
		}
		raw, _ := json.Marshal(schema)
		return raw
	}

	s.AddTool(&mcp.Tool{
		Name:        "search_articles",
		Description: "Search WeChat official account articles by keyword. Returns title, author, account, summary, URL, and publish date.",
		InputSchema: inputSchema(map[string]any{
			"keyword": map[string]any{"type": "string", "description": "Search keyword"},
			"page":    map[string]any{"type": "integer", "description": "Page number (starts from 1, default 1)"},
		}, []string{"keyword"}),
	}, app.handleSearchArticles)

	s.AddTool(&mcp.Tool{
		Name:        "search_accounts",
		Description: "Search WeChat official accounts by name. Returns account name, wechat_id, description, and recent article.",
		InputSchema: inputSchema(map[string]any{
			"keyword": map[string]any{"type": "string", "description": "Account name or keyword to search"},
			"page":    map[string]any{"type": "integer", "description": "Page number (starts from 1, default 1)"},
		}, []string{"keyword"}),
	}, app.handleSearchAccounts)

	s.AddTool(&mcp.Tool{
		Name:        "get_article_content",
		Description: "Get the full content of a WeChat article by its URL. Extracts title, author, publish date, and full text content.",
		InputSchema: inputSchema(map[string]any{
			"url": map[string]any{"type": "string", "description": "The WeChat article URL (mp.weixin.qq.com)"},
		}, []string{"url"}),
	}, app.handleGetArticleContent)

	return s
}

func newStreamableHandler(s *mcp.Server) http.Handler {
	return mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return s
	}, nil)
}
