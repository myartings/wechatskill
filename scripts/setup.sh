#!/bin/bash
# 微信公众号 MCP 服务器初始化脚本
set -e

SKILL_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$SKILL_DIR"

export PATH="$HOME/go-sdk/go/bin:$HOME/go/bin:$PATH"

# 创建 ~/.openclaw/skills/wechat 软链接（兼容 SKILL.md 中的路径）
LINK_PATH="$HOME/.openclaw/skills/wechat"
if [ ! -e "$LINK_PATH" ]; then
    mkdir -p "$(dirname "$LINK_PATH")"
    ln -s "$SKILL_DIR" "$LINK_PATH"
    echo "创建软链接: $LINK_PATH -> $SKILL_DIR"
elif [ -L "$LINK_PATH" ] && [ "$(readlink "$LINK_PATH")" != "$SKILL_DIR" ]; then
    rm "$LINK_PATH"
    ln -s "$SKILL_DIR" "$LINK_PATH"
    echo "更新软链接: $LINK_PATH -> $SKILL_DIR"
fi

# 编译（如果没有二进制或源码更新了）
if [ ! -f "$SKILL_DIR/wechat-mcp" ] || [ "$SKILL_DIR/main.go" -nt "$SKILL_DIR/wechat-mcp" ]; then
    echo "编译 wechat-mcp..."
    go build -o "$SKILL_DIR/wechat-mcp" .
fi

# 启动（如果没在运行）
if ! curl -s http://localhost:8090/api/v1/status > /dev/null 2>&1; then
    echo "启动 wechat-mcp 服务器..."
    "$SKILL_DIR/wechat-mcp" --port 8090 &
    sleep 2
fi

echo "微信公众号服务器就绪"
python3 "$SKILL_DIR/scripts/wechat_client.py" status
