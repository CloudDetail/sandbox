# APO Sandbox Node.js 版本

这是一个基于Node.js的故障注入沙箱应用，用于模拟应用操作并可以为自身注入故障。

## 功能特性

- **CPU故障**: 阻塞当前请求，消耗CPU时间
- **Redis故障**: 发送故障命令触发redis proxy的故障处理
- **网络延迟**: 通过tc命令实现网络接口延迟

## 安装和运行

### 前置要求

- Node.js 18+
- Redis服务器
- MySQL服务器（可选）
- Linux系统（用于tc命令）

### 安装依赖

```bash
npm install
```

### 配置环境变量

复制 `env.example` 到 `.env` 并修改配置：

```bash
cp env.example .env
```

### 启动应用

```bash
# 开发模式
npm run dev

# 生产模式
npm start
```

## API接口

### 获取用户数据（带故障注入）

```
GET /api/users?chaos=cpu&duration=200
```

**查询参数:**
- `chaos`: 故障类型 (`cpu`, `redis_latency`, `latency`)
- `duration`: 故障持续时间（毫秒）

**故障类型说明:**
- `cpu`: CPU故障，阻塞当前请求
- `redis_latency`: Redis延迟故障，触发redis proxy故障处理
- `latency`: 网络延迟故障，通过tc命令实现

### 健康检查

```
GET /health
```

### 故障状态

```
GET /faults/status
```

## 故障注入示例

### CPU故障
```bash
curl "http://localhost:3500/api/users?chaos=cpu&duration=500"
```

### Redis故障
```bash
curl "http://localhost:3500/api/users?chaos=redis_latency&duration=100"
```

### 网络延迟
```bash
curl "http://localhost:3500/api/users?chaos=latency&duration=200"
```

### 停止所有故障
```bash
curl "http://localhost:3500/api/users"
```

## Docker部署

```bash
# 构建镜像
docker build -t apo-sandbox-nodejs .

# 运行容器
docker run -p 3500:3500 --env-file .env apo-sandbox-nodejs
```

## 注意事项

1. **tc命令权限**: 网络延迟故障需要root权限执行tc命令
2. **Redis连接**: 确保Redis服务器可访问
3. **网络接口**: 默认使用eth0接口，可根据实际情况修改

## 故障处理机制

- **CPU故障**: 通过斐波那契计算阻塞当前请求
- **Redis故障**: 向Redis发送特定命令，触发proxy故障处理
- **网络延迟**: 使用Linux tc命令在网络接口上添加延迟规则

## 日志

应用使用Winston日志库，支持不同级别的日志输出。日志包含时间戳和结构化信息，便于故障排查。
