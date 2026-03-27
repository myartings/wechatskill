---
name: wechat
description: |
  微信公众号文章搜索助手。查询公众号文章和公众号信息，获取文章全文内容。
  当用户提到微信公众号、公众号文章、微信文章、WeChat、搜公众号、看公众号文章等任何与微信公众号相关的查询操作时使用此 skill。
  支持两种搜索方式：web search 搜索最新文章（推荐），以及搜狗微信搜索。
---

# 微信公众号查询助手

## Skill 目录

`SKILL_DIR`：`${CLAUDE_PLUGIN_ROOT}` 或 `~/.openclaw/skills/wechat` 或 `~/.claude/skills/wechat`，取实际存在的路径。

`P` 代表 `python3 <SKILL_DIR>/scripts/wechat_client.py`。

## 首次使用

首次使用前需要初始化：

```shell
cd <SKILL_DIR> && bash scripts/setup.sh
```

这会自动编译服务器并启动。之后无需重复执行。

## 搜索文章

通过 Python 脚本调用搜狗微信搜索，速度快，适合大部分公众号。

```shell
P status
P search-articles "关键词"
P search-accounts "公众号名称"
```

### 搜狗搜不到时的备选方案

如果搜狗搜索返回的结果不相关或为空，可以用 web search 补充搜索：

- 搜索格式：`site:mp.weixin.qq.com "公众号名" "关键词"`
- Web search 返回的 mp.weixin.qq.com 链接可以直接用于获取文章全文

## 获取文章全文

使用搜索结果中的 URL 获取文章完整内容：

```shell
P article "文章URL"
```

## 操作流程

1. **每次操作前**先运行 `P status` 检查服务器状态
2. 如服务器未启动：运行 `cd <SKILL_DIR> && bash scripts/setup.sh`
3. **搜索文章**：使用搜狗搜索（`P search-articles`），速度快
4. **搜狗搜不到时**：用 web search 搜 `site:mp.weixin.qq.com "公众号名"` 作为备选
5. **查看全文**：使用搜索到的 URL，运行 `P article "URL"` 获取文章全文

## 重要规则

- **优先使用搜狗搜索**，速度快，大部分公众号都能搜到
- 搜狗搜不到时再用 web search 补充
- 如遇到验证码或反爬限制，建议稍后重试
