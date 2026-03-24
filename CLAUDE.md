# WeChat MCP Server

微信公众号查询 MCP 服务器 + OpenClaw Skill。

## 构建与运行

```bash
export PATH=$HOME/go-sdk/go/bin:$HOME/go/bin:$PATH
go build -o wechat-mcp .
./wechat-mcp --port 8090
```

Go 版本要求 >= 1.24（MCP SDK 会自动下载 Go 1.25 toolchain）。

## 架构

- 通过搜狗微信搜索（weixin.sogou.com）抓取公众号文章和账号信息
- 通过 goquery 解析 HTML 提取结构化数据
- 文章全文通过直接访问 mp.weixin.qq.com 文章页面提取
- 双接口：MCP (`/mcp`) + REST API (`/api/v1/*`)
- 无需认证，纯只读查询

## 搜狗微信搜索注意事项

- `type=2` 搜索文章，`type=1` 搜索公众号
- **公众号搜索（type=1）反爬较严**，无 cookie 时可能返回空结果，文章搜索（type=2）通常正常
- 搜狗有反爬机制，频繁请求可能触发验证码
- HTML 结构可能变化，选择器需要定期维护
- 搜索结果中的文章 URL 是搜狗跳转链接（https://weixin.sogou.com/link?...）
- 日期通过 JS `timeConvert()` 传递 unix 时间戳，需要解析

## 项目结构

```
├── main.go              # 入口，CLI 参数
├── app_server.go        # AppServer 管理 service + MCP server
├── mcp_server.go        # 3 个 MCP tool 注册
├── mcp_handlers.go      # MCP handler 实现
├── routes.go            # REST API 路由
├── service.go           # 业务逻辑层
├── wechat/              # 微信搜索客户端
│   ├── client.go        # HTTP 客户端（带 cookie jar）
│   ├── search.go        # 搜狗微信搜索（文章 + 公众号）
│   ├── article.go       # 文章内容提取
│   └── types.go         # 数据类型
├── SKILL.md             # OpenClaw skill 定义
└── scripts/
    ├── setup.sh         # 初始化脚本（编译 + 启动）
    └── wechat_client.py # Python CLI 客户端
```

## 添加新功能的流程

1. 在 `wechat/` 中添加客户端方法
2. 在 `service.go` 中添加代理方法
3. 在 `mcp_server.go` 中注册 tool
4. 在 `mcp_handlers.go` 中实现 handler
5. 可选：在 `routes.go` 中添加 REST API endpoint

## 测试

```bash
# 启动
./wechat-mcp --port 8090

# 检查状态
curl -s http://localhost:8090/api/v1/status

# 搜索文章
curl -s -X POST http://localhost:8090/api/v1/search/articles -H "Content-Type: application/json" -d '{"keyword":"AI"}'

# 搜索公众号
curl -s -X POST http://localhost:8090/api/v1/search/accounts -H "Content-Type: application/json" -d '{"keyword":"人民日报"}'
```
