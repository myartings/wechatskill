---
name: wechat
description: |
  微信公众号文章搜索助手。查询公众号文章和公众号信息，获取文章全文内容。
  当用户提到微信公众号、公众号文章、微信文章、WeChat、搜公众号、看公众号文章等任何与微信公众号相关的查询操作时使用此 skill。
  支持两种搜索方式：web search 搜索最新文章（推荐），以及搜狗微信搜索。
---

# 微信公众号查询助手

## 首次使用

首次使用前需要初始化：

```shell
cd ~/.openclaw/skills/wechat && bash scripts/setup.sh
```

这会自动编译服务器并启动。之后无需重复执行。

## 搜索文章的两种方式

### 方式一：Web Search（推荐，结果更新更全）

直接使用 web search 搜索微信公众号文章，搜索关键词格式：

- 搜索特定公众号最新文章：`site:mp.weixin.qq.com "公众号名" 最新文章`
- 搜索特定话题的文章：`site:mp.weixin.qq.com "公众号名" "关键词"`

Web search 返回的 mp.weixin.qq.com 链接可以直接用于获取文章全文（见下方"获取文章全文"）。

**当用户要求搜索某个公众号的最新文章时，优先使用此方式。**

### 方式二：搜狗微信搜索（备选）

通过 Python 脚本调用搜狗微信搜索。搜狗的索引可能不包含最新文章，适合搜索较热门的内容。

脚本路径：`~/.openclaw/skills/wechat/scripts/wechat_client.py`

```shell
# 检查服务器运行状态
python3 ~/.openclaw/skills/wechat/scripts/wechat_client.py status

# 按关键词搜索公众号文章
python3 ~/.openclaw/skills/wechat/scripts/wechat_client.py search-articles "关键词"

# 按名称搜索公众号
python3 ~/.openclaw/skills/wechat/scripts/wechat_client.py search-accounts "公众号名称"
```

## 获取文章全文

使用搜索结果中的 URL 获取文章完整内容：

```shell
# URL 可以来自 web search 或搜狗搜索的结果
python3 ~/.openclaw/skills/wechat/scripts/wechat_client.py article "文章URL"
```

## 操作流程

1. **每次操作前**先运行 `status` 检查服务器状态
2. 如服务器未启动：运行 `bash ~/.openclaw/skills/wechat/scripts/setup.sh`
3. **搜索最新文章**：优先用 web search（方式一），搜狗搜索不到时也用 web search
4. **搜索热门/通用文章**：可以用搜狗搜索（方式二）
5. **查看全文**：使用搜索到的 URL，运行 `article "URL"` 获取文章全文

## 重要规则

- 搜索最新文章时**优先使用 web search**，搜狗索引经常滞后
- Python 脚本命令用于搜狗搜索和获取文章全文
- 如遇到验证码或反爬限制，建议稍后重试
