package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/myartings/wechatskill/wechat"
)

type AppServer struct {
	service   *WechatService
	mcpServer *mcp.Server
	port      string
}

func NewAppServer(port string) *AppServer {
	app := &AppServer{
		service: NewWechatService(wechat.NewClient()),
		port:    port,
	}
	app.mcpServer = newMCPServer(app)
	return app
}
