---
name: wechat
description: |
  微信公众号文章搜索助手。查询公众号文章和公众号信息，获取文章全文内容。
  当用户提到微信公众号、公众号文章、微信文章、WeChat、搜公众号、看公众号文章等任何与微信公众号相关的查询操作时使用此 skill。
  支持两种搜索方式：web search 搜索最新文章（推荐），以及搜狗微信搜索。
---

# 规则

1. **只用下面的 python3 命令。禁止用 curl、wget、httpie 或任何其他方式。**
2. Skill 目录（`SKILL_DIR`）：`~/.claude/skills/wechat` 或 `~/.openclaw/skills/wechat`，取实际存在的路径
3. 首次使用先运行初始化：`cd <SKILL_DIR> && bash scripts/setup.sh`（自动编译服务器并启动）
4. 每次操作前先运行 `P status` 检查服务器状态，如未启动则重新运行 setup.sh

# 命令

以下是全部可用命令，`P` 代表 `python3 <SKILL_DIR>/scripts/wechat_client.py`。

| 功能 | 命令 |
|------|------|
| 检查状态 | `P status` |
| 搜索文章 | `P search-articles "关键词"` |
| 搜索公众号 | `P search-accounts "公众号名称"` |
| 获取文章全文 | `P article "文章URL"` |

## 搜狗搜不到时的备选方案

如果搜狗搜索返回的结果不相关或为空，可以用 web search 补充搜索：

- 搜索格式：`site:mp.weixin.qq.com "公众号名" "关键词"`
- Web search 返回的 mp.weixin.qq.com 链接可以直接用于获取文章全文

# 示例

```shell
P status
P search-articles "AI"
P article "https://mp.weixin.qq.com/s/xxx"
```
