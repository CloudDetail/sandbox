import multiprocessing
import os

# 绑定地址和端口
bind = f"0.0.0.0:{os.getenv('PORT', '3500')}"

# 工作进程数量（根据CPU核心数调整）
workers = multiprocessing.cpu_count() * 2 + 1

# 每个工作进程的线程数
threads = 4

# 工作模式：sync（同步）或 gevent（异步）
worker_class = "sync"

# 最大并发连接数
worker_connections = 1000

# 超时时间
timeout = 30

# 优雅关闭超时
graceful_timeout = 30

# 保持连接超时
keepalive = 2

# 日志级别
loglevel = "info"

# 访问日志格式
access_log_format = '%(h)s %(l)s %(u)s %(t)s "%(r)s" %(s)s %(b)s "%(f)s" "%(a)s" %(D)s'

# 预加载应用
preload_app = True