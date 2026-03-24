package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *AppServer) setupRoutes(r *gin.Engine) {
	mcpHandler := newStreamableHandler(a.mcpServer)
	r.Any("/mcp", gin.WrapH(mcpHandler))

	api := r.Group("/api/v1")
	{
		api.GET("/status", a.apiStatus)
		api.POST("/search/articles", a.apiSearchArticles)
		api.POST("/search/accounts", a.apiSearchAccounts)
		api.POST("/article/content", a.apiGetArticleContent)
	}
}

func (a *AppServer) apiStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "wechat-mcp"})
}

func (a *AppServer) apiSearchArticles(c *gin.Context) {
	var req struct {
		Keyword string `json:"keyword"`
		Page    int    `json:"page"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keyword is required"})
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	articles, err := a.service.SearchArticles(c.Request.Context(), req.Keyword, req.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": articles})
}

func (a *AppServer) apiSearchAccounts(c *gin.Context) {
	var req struct {
		Keyword string `json:"keyword"`
		Page    int    `json:"page"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keyword is required"})
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	accounts, err := a.service.SearchAccounts(c.Request.Context(), req.Keyword, req.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": accounts})
}

func (a *AppServer) apiGetArticleContent(c *gin.Context) {
	var req struct {
		URL string `json:"url"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}
	detail, err := a.service.GetArticleContent(c.Request.Context(), req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, detail)
}
