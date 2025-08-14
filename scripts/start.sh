#!/bin/bash

# 任务名
SERVICE_NAME="go-gin-api"

# 日志文件（覆盖写入）
LOG_FILE="./deploy.log"

# 清空或创建日志文件
> "$LOG_FILE"

# 将后续所有输出重定向到日志（可选：同时保留屏幕输出）
# 如果你希望：屏幕也显示 + 写入日志，使用下面这行
exec > >(tee -a "$LOG_FILE") 2>&1
# 注意：tee -a 是追加，但我们前面已清空，所以等效于覆盖

# 如果你希望：只写日志，不打印到屏幕，用这行代替：
# exec > "$LOG_FILE" 2>&1

echo "======================================"
echo "部署开始 $(date '+%Y-%m-%d %H:%M:%S')"
echo "======================================"

git checkout go.mod go.sum

# 拉取代码
echo "📥 正在拉取代码..."
git pull || {
    echo "❌ 拉取代码失败"
    exit 1
}

# 执行 go mod tidy
echo "📦 执行 go mod tidy..."
go mod tidy || {
    echo "❌ go mod tidy 失败"
    exit 1
}

# 编译主程序
# echo "🔨 正在编译..."
# go build -o go-gin-api main.go || {
#     echo "❌ 编译失败"
#     exit 1
# }

# 使用 pgrep 获取服务进程的 PID
PID=$(pgrep -f "^\.\/go run main.*-env")

# 检查 PID 是否存在
if [ -z "$PID" ]; then
    echo "✅ 开始启动服务..."
    nohup ./go run main -env dev > ./start_api.log 2>&1 &
else
    echo "🛑 正在停止服务（进程ID: $PID）..."
    kill -15 "$PID"  # 优雅关闭（SIGTERM）
    sleep 2

    # 如果进程仍未退出，强制杀死
    if pgrep -f "^\.\/go run main.*-env" > /dev/null; then
        echo "💥 进程未退出，强制终止..."
        kill -9 "$PID"
    fi

    echo "🚀 启动新服务..."
    nohup ./go run main -env dev > ./start_api.log 2>&1 &
fi

echo "🎉 基础部署完成"