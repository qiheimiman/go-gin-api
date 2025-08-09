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
PID=$(pgrep -f "^\.\/$SERVICE_NAME.*-port 9999")

# 检查 PID 是否存在
if [ -z "$PID" ]; then
    echo "开始启动服务..."
    nohup ./go-gin-api -port 9999 -env dev  > ./start_api.log 2>&1 &
else
    # 发送 SIGUSR1 信号到进程，告诉 endless 重启服务
    echo "正在重启服务（进程ID: $PID）..."
    kill -1 $PID
    echo "重启信号已发送。"
fi

# 使用 Supervisor 重启任务
#supervisorctl restart "$supervisor_task"

