#!/bin/bash

# 任务名
SERVICE_NAME="go-gin-api"

git checkout go.mod go.sum

# 拉取代码
git pull

# 执行 go mod tidy
go mod tidy

# 编译主程序
go build -o go-gin-api main.go

# 使用 pgrep 获取服务进程的 PID
PID=$(pgrep -f "^\.\/go-gin-api.*-env")

# 检查 PID 是否存在
if [ -z "$PID" ]; then
    echo "开始启动服务..."
    nohup ./go-gin-api -env dev  > ./start_api.log 2>&1 &
else
    echo "正在停止服务（进程ID: $PID）..."
    kill -15 "$PID"  # 优雅关闭（SIGTERM）
    sleep 2

    # 如果进程仍未退出，强制杀死
    if pgrep -f "^\.\/go-gin-api.*-env" > /dev/null; then
        echo "进程未退出，强制终止..."
        kill -9 "$PID"
    fi

    echo "启动新服务..."
    nohup ./go-gin-api -env dev > ./start_api.log 2>&1 &
fi

# 使用 Supervisor 重启任务
#supervisorctl restart "$supervisor_task"

