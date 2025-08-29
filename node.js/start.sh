#!/bin/bash

echo "启动 APO Sandbox Node.js 版本..."

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo "错误: Node.js 未安装"
    exit 1
fi

# 检查npm是否安装
if ! command -v npm &> /dev/null; then
    echo "错误: npm 未安装"
    exit 1
fi

# 安装依赖
echo "安装依赖..."
npm install

# 检查.env文件
if [ ! -f .env ]; then
    echo "警告: .env 文件不存在，使用默认配置"
    cp env.example .env 2>/dev/null || echo "无法复制 env.example"
fi

# 启动应用
echo "启动应用..."
npm run start:otel
