#!/bin/bash
echo "启动测试脚本"

# 检查端口是否被占用
if lsof -i:8080 > /dev/null 2>&1; then
    echo "错误: 端口 8080 已被占用"
    exit 1
fi

# 启动服务器
go run main.go &
sleep 2

# 获取服务器进程 PID（通过端口号）
SERVER_PID=$(lsof -i:8080 -t)
if [ -z "$SERVER_PID" ]; then
    echo "错误: 无法获取服务器 PID"
    exit 1
fi
echo "服务器 PID: $SERVER_PID"
echo "等待3秒..."
sleep 3

# 测试HTTP请求
echo "发送HTTP请求..."
curl http://localhost:8080
echo

# 测试重启
echo "发送重启信号..."
kill -INT $SERVER_PID
echo "等待新服务器启动..."
sleep 5

# 获取新服务器的 PID
NEW_PID=$(lsof -i:8080 -t)
if [ -z "$NEW_PID" ]; then
    echo "错误: 无法获取新服务器 PID"
    exit 1
fi
echo "新服务器 PID: $NEW_PID"

# 再次测试HTTP请求
echo "发送HTTP请求到新服务器..."
curl http://localhost:8080
echo

# 关闭服务器
echo "发送关闭信号..."
kill -QUIT $NEW_PID
echo "测试完成"