# WeChat Official Account MCP Server

微信公众号文章搜索 MCP 服务器 + OpenClaw Skill。

通过搜狗微信搜索引擎查询公众号文章，获取文章全文内容。

## 功能

- **搜索文章** — 按关键词搜索微信公众号文章，返回标题、摘要、日期、URL
- **搜索公众号** — 按名称搜索微信公众号
- **获取文章全文** — 提取微信文章的完整内容（支持搜狗跳转链接和直接微信链接）

## 快速开始

### 构建

```bash
go build -o wechat-mcp .
```

### 运行

```bash
./wechat-mcp --port 8090
```

### 作为 OpenClaw Skill 使用

```bash
cd ~/.openclaw/skills/wechat && bash scripts/setup.sh
```

## 接口

### MCP (Model Context Protocol)

端点：`http://localhost:8090/mcp`

提供 3 个 MCP Tool：

| Tool | 说明 |
|------|------|
| `search_articles` | 按关键词搜索公众号文章 |
| `search_accounts` | 按名称搜索公众号 |
| `get_article_content` | 获取文章全文内容 |

### REST API

```bash
# 检查状态
curl http://localhost:8090/api/v1/status

# 搜索文章
curl -X POST http://localhost:8090/api/v1/search/articles \
  -H "Content-Type: application/json" \
  -d '{"keyword":"AI"}'

# 搜索公众号
curl -X POST http://localhost:8090/api/v1/search/accounts \
  -H "Content-Type: application/json" \
  -d '{"keyword":"人民日报"}'

# 获取文章全文
curl -X POST http://localhost:8090/api/v1/article/content \
  -H "Content-Type: application/json" \
  -d '{"url":"https://weixin.sogou.com/link?url=..."}'
```

### Python CLI

```bash
python3 scripts/wechat_client.py status
python3 scripts/wechat_client.py search-articles "AI"
python3 scripts/wechat_client.py search-accounts "人民日报"
python3 scripts/wechat_client.py article "文章URL"
```

## 技术架构

- **Go** + MCP SDK + Gin
- 通过搜狗微信搜索（weixin.sogou.com）抓取文章和公众号信息
- goquery 解析 HTML
- 自动解析搜狗 JS 跳转链接获取真实微信文章 URL
- 双接口：MCP + REST API
- 无需认证，纯只读查询

## 已知限制

- 公众号搜索受搜狗反爬限制，无 cookie 时可能返回空结果
- 频繁请求可能触发搜狗验证码
- 搜索结果可能不包含最新发布的文章

## License

MIT
