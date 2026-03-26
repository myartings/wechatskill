package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func parseArgs(req *mcp.CallToolRequest) map[string]any {
	var args map[string]any
	if req.Params.Arguments != nil {
		json.Unmarshal(req.Params.Arguments, &args)
	}
	if args == nil {
		args = map[string]any{}
	}
	return args
}

func getStringArg(args map[string]any, key string) string {
	v, ok := args[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

func getIntArg(args map[string]any, key string, defaultVal int) int {
	v, ok := args[key]
	if !ok {
		return defaultVal
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case json.Number:
		i, _ := n.Int64()
		return int(i)
	}
	return defaultVal
}

func toJSON(v any) string {
	raw, _ := json.MarshalIndent(v, "", "  ")
	return string(raw)
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}

func errorResult(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error: %s", err.Error())}},
		IsError: true,
	}
}

func (a *AppServer) handleSearchArticles(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	keyword := getStringArg(args, "keyword")
	if keyword == "" {
		return errorResult(fmt.Errorf("keyword is required")), nil
	}
	page := getIntArg(args, "page", 1)

	articles, err := a.service.SearchArticles(ctx, keyword, page)
	if err != nil {
		return errorResult(err), nil
	}
	if len(articles) == 0 {
		return textResult("No articles found."), nil
	}
	return textResult(toJSON(articles)), nil
}

func (a *AppServer) handleSearchAccounts(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	keyword := getStringArg(args, "keyword")
	if keyword == "" {
		return errorResult(fmt.Errorf("keyword is required")), nil
	}
	page := getIntArg(args, "page", 1)

	accounts, err := a.service.SearchAccounts(ctx, keyword, page)
	if err != nil {
		return errorResult(err), nil
	}
	if len(accounts) == 0 {
		return textResult("No accounts found."), nil
	}
	return textResult(toJSON(accounts)), nil
}

func (a *AppServer) handleGetAccountArticles(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	accountName := getStringArg(args, "account_name")
	if accountName == "" {
		return errorResult(fmt.Errorf("account_name is required")), nil
	}

	articles, err := a.service.GetAccountArticles(ctx, accountName)
	if err != nil {
		return errorResult(err), nil
	}
	if len(articles) == 0 {
		return textResult("No articles found for this account."), nil
	}
	return textResult(toJSON(articles)), nil
}

func (a *AppServer) handleGetArticleContent(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	url := getStringArg(args, "url")
	if url == "" {
		return errorResult(fmt.Errorf("url is required")), nil
	}

	detail, err := a.service.GetArticleContent(ctx, url)
	if err != nil {
		return errorResult(err), nil
	}
	return textResult(toJSON(detail)), nil
}
