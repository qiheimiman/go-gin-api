#!/bin/bash

SERVICE_NAME="go-gin-api"
PIDFILE="./${SERVICE_NAME}.pid"
BINARY="./${SERVICE_NAME}"
PORT=9999
ENV="dev"

echo "🔄 开始部署 $SERVICE_NAME..."

# 确保在项目根目录
cd "$(dirname "$0")/.." || exit 1

# 拉取代码
git checkout go.mod go.sum
git pull || {
    echo "❌ 拉取代码失败"
    exit 1
}

# 整理依赖
go mod tidy || {
    echo "❌ go mod tidy 失败"
    exit 1
}

# 编译
echo "📦 编译中..."
go build -o "$BINARY" main.go || {
    echo "❌ 编译失败"
    exit 1
}

# 生成 Swagger（可选）
if command -v swag >/dev/null 2>&1; then
    echo "📝 生成 Swagger 文档..."
    swag init || echo "⚠️ swag init 失败，跳过"
fi

# 获取旧进程 PID（通过 PID 文件或精确匹配）
PID=""
if [ -f "$PIDFILE" ]; then
    PID=$(cat "$PIDFILE")
    if ! ps -p "$PID" > /dev/null 2>&1; then
        echo "🧹 旧 PID 文件无效，进程已不存在"
        PID=""
    fi
else
    # 备用方案：精确匹配启动命令
    PID=$(pgrep -f "^$BINARY.*-port $PORT")
fi

# 重启或启动
if [ -n "$PID" ]; then
    echo "🔄 正在重启服务（PID: $PID）..."
    kill -USR1 "$PID" || {
        echo "❌ 无法发送信号，尝试强制重启..."
        kill "$PID" && sleep 1
    }
    # 更新 PID 文件（假设程序会重新写入）
else
    echo "🚀 服务未运行，正在启动..."
fi

# 启动新服务（如果没运行）
if ! pgrep -f "^$BINARY.*-port $PORT" > /dev/null; then
    nohup "$BINARY" -port "$PORT" -env "$ENV" > ./start_api.log 2>&1 &
    NEW_PID=$!
    echo "$NEW_PID" > "$PIDFILE"
    echo "✅ 服务已启动，PID: $NEW_PID"
fi

echo "🎉 部署完成"