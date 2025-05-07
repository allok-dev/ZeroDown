#!/bin/bash

# 定义颜色
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo -e "${GREEN}[测试脚本] 启动测试脚本${NC}"

# 检查端口是否被占用
if lsof -i:8080 > /dev/null 2>&1; then
    echo -e "${GREEN}[测试脚本] 错误: 端口 8080 已被占用${NC}"
    exit 1
fi

# 启动服务器
echo -e "${GREEN}[测试脚本] 启动服务器...${NC}"
go run main.go > >(sed 's/^/[服务器] /') 2>&1 &
sleep 2

# 获取服务器进程 PID
SERVER_PID=$(lsof -i:8080 -t)
if [ -z "$SERVER_PID" ]; then
    echo -e "${GREEN}[测试脚本] 错误: 无法获取服务器 PID${NC}"
    exit 1
fi
echo -e "${GREEN}[测试脚本] 服务器 PID: $SERVER_PID${NC}"

# 并行执行 HTTP 请求和重启信号
echo -e "${GREEN}[测试脚本] 发送HTTP请求到旧服务器...${NC}"
{
  response=$(curl -s http://localhost:8080)
  echo -e "${GREEN}[测试脚本] 旧请求收到响应: $response${NC}" 
} 2>&1 &

sleep 3
echo -e "${GREEN}[测试脚本] 发送重启信号...${NC}"
kill -INT $SERVER_PID
echo -e "${GREEN}[测试脚本] 等待新服务器启动...${NC}"

# 再次测试HTTP请求
echo -e "${GREEN}[测试脚本] 发送HTTP请求到新服务器...${NC}"
response=$(curl -s http://localhost:8080)
echo -e "${GREEN}[测试脚本] 新请求收到响应: $response${NC}"

# 获取新服务器的 PID
NEW_PID=$(lsof -i:8080 -t)
if [ -z "$NEW_PID" ]; then
    echo -e "${GREEN}[测试脚本] 错误: 无法获取新服务器 PID${NC}"
    exit 1
fi
echo -e "${GREEN}[测试脚本] 新服务器 PID: $NEW_PID${NC}"

# 并行执行 HTTP 请求和重启信号
echo -e "${GREEN}[测试脚本] 发送HTTP请求到服务器...${NC}"
{
  response=$(curl -s http://localhost:8080)
  echo -e "${GREEN}[测试脚本] 请求收到响应: $response${NC}" 
} 2>&1 &

sleep 3

# 关闭服务器
echo -e "${GREEN}[测试脚本] 发送关闭信号...${NC}"
kill -QUIT $NEW_PID

sleep 60
echo -e "${GREEN}[测试脚本] 测试完成${NC}"