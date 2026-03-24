---
name: wechat
description: |
  微信公众号文章搜索助手。通过搜狗微信搜索引擎查询公众号文章和公众号信息，获取文章全文内容。
  当用户提到微信公众号、公众号文章、微信文章、WeChat、搜公众号、看公众号文章等任何与微信公众号相关的查询操作时使用此 skill。
  必须且只能通过 scripts/wechat_client.py 脚本来执行所有操作，不要尝试其他方式。
---

# 微信公众号查询助手

通过 Python 脚本查询微信公众号内容。**所有操作必须且只能通过下面的命令执行，禁止使用 curl 或其他方式直接调用 API。**

## 首次使用

首次使用前需要初始化：

```shell
cd ~/.openclaw/skills/wechat && bash scripts/setup.sh
```

这会自动编译服务器并启动。之后无需重复执行。

## 所有可用命令

脚本路径：`~/.openclaw/skills/wechat/scripts/wechat_client.py`

### 检查状态

```shell
# 检查服务器运行状态
python3 ~/.openclaw/skills/wechat/scripts/wechat_client.py status
```

### 搜索文章

```shell
# 按关键词搜索公众号文章
python3 ~/.openclaw/skills/wechat/scripts/wechat_client.py search-articles "关键词"
```

### 搜索公众号

```shell
# 按名称搜索公众号
python3 ~/.openclaw/skills/wechat/scripts/wechat_client.py search-accounts "公众号名称"
```

### 获取文章全文

```shell
# 获取文章的完整内容（URL 来自搜索结果）
python3 ~/.openclaw/skills/wechat/scripts/wechat_client.py article "文章URL"
```

## 操作流程

1. **每次操作前**先运行 `status` 检查服务器状态
2. 如服务器未启动：运行 `bash ~/.openclaw/skills/wechat/scripts/setup.sh`
3. 搜索文章：使用 `search-articles "关键词"` 搜索，结果包含标题、公众号名、摘要和 URL
4. 搜索公众号：使用 `search-accounts "名称"` 搜索公众号
5. 查看全文：使用搜索结果中的 URL，运行 `article "URL"` 获取文章全文

## 重要规则

- **只使用上面列出的 python3 命令，不要用 curl 或其他方式**
- 搜索结果来自搜狗微信搜索，可能不包含最新发布的文章
- 如遇到验证码或反爬限制，建议稍后重试
