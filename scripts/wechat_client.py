#!/usr/bin/env python3
"""WeChat Official Account MCP Server CLI Client.

Usage:
    python wechat_client.py status                     # Check server status
    python wechat_client.py search-articles <keyword>  # Search articles by keyword
    python wechat_client.py search-accounts <keyword>  # Search official accounts
    python wechat_client.py article <url>              # Get full article content
"""

import sys
import json
import urllib.request
import urllib.error

BASE_URL = "http://localhost:8090"


def api(method, path, data=None):
    url = BASE_URL + path
    headers = {"Content-Type": "application/json", "Accept": "application/json"}
    body = json.dumps(data).encode() if data else None
    req = urllib.request.Request(url, data=body, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            return json.loads(resp.read())
    except urllib.error.HTTPError as e:
        return {"error": f"HTTP {e.code}: {e.read().decode()[:200]}"}
    except Exception as e:
        return {"error": str(e)}


def pp(obj):
    print(json.dumps(obj, indent=2, ensure_ascii=False))


def cmd_status():
    r = api("GET", "/api/v1/status")
    if r.get("status") == "ok":
        print("服务器运行中 ✓")
    else:
        print(f"服务器状态异常: {r}")


def cmd_search_articles(keyword):
    r = api("POST", "/api/v1/search/articles", {"keyword": keyword})
    if "error" in r:
        print(f"错误: {r['error']}")
        return
    articles = r.get("data", [])
    if not articles:
        print("未找到相关文章")
        return
    print(f"找到 {len(articles)} 篇文章:\n")
    for a in articles:
        print(f"  [{a.get('account_name', '?')}] {a.get('title', '')}")
        if a.get("summary"):
            print(f"    摘要: {a['summary'][:80]}")
        if a.get("publish_date"):
            print(f"    日期: {a['publish_date']}")
        if a.get("url"):
            print(f"    URL: {a['url']}")
        print()


def cmd_search_accounts(keyword):
    r = api("POST", "/api/v1/search/accounts", {"keyword": keyword})
    if "error" in r:
        print(f"错误: {r['error']}")
        return
    accounts = r.get("data", [])
    if not accounts:
        print("未找到相关公众号")
        return
    print(f"找到 {len(accounts)} 个公众号:\n")
    for a in accounts:
        print(f"  {a.get('name', '?')} (ID: {a.get('wechat_id', '?')})")
        if a.get("description"):
            print(f"    简介: {a['description'][:80]}")
        if a.get("recent_article"):
            print(f"    最近文章: {a['recent_article']}")
        print()


def cmd_article(url):
    r = api("POST", "/api/v1/article/content", {"url": url})
    if "error" in r:
        print(f"错误: {r['error']}")
        return
    print(f"标题: {r.get('title', '?')}")
    print(f"作者: {r.get('author', '?')}")
    print(f"公众号: {r.get('account_name', '?')}")
    print(f"日期: {r.get('publish_date', '?')}")
    print(f"\n{'='*60}\n")
    print(r.get("content", "(无内容)"))


def main():
    if len(sys.argv) < 2:
        print(__doc__)
        sys.exit(1)

    cmd = sys.argv[1]
    args = sys.argv[2:]

    commands = {
        "status": (cmd_status, 0),
        "search-articles": (cmd_search_articles, 1),
        "search-accounts": (cmd_search_accounts, 1),
        "article": (cmd_article, 1),
    }

    if cmd not in commands:
        print(f"未知命令: {cmd}")
        print(__doc__)
        sys.exit(1)

    func, nargs = commands[cmd]
    if len(args) < nargs:
        print(f"命令 '{cmd}' 需要 {nargs} 个参数")
        sys.exit(1)

    func(*args[:nargs])


if __name__ == "__main__":
    main()
