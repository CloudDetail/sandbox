# python/main.py
from flask import Flask, request, g
import time
import logging

from config import load_config
from logging_config import setup_logging
from storage import MySQLClient, RedisClient, Store
from service import BusinessService
from api import BusinessAPI, business_api_bp

# 创建Flask应用实例
app = Flask(__name__)
logger = logging.getLogger(__name__)

fault_manager = None

@app.before_request
def logging_middleware():
    g.start_time = time.time()
    logger.info(f"Request started: {request.method} {request.path}")

@app.after_request
def after_request(response):
    duration = time.time() - g.start_time
    logger.info(f"Request completed: {request.method} {request.path} - {duration:.4f}s")
    return response

def create_app():
    global fault_manager

    # 加载配置
    app_config = load_config()

    # 设置日志
    setup_logging()

    if app_config.server.deploy_proxy:
        import requests
        proxy_json = {
                "name": "redis",
                "listen": "localhost:6379",
                "upstream": "redis-service:6379"
            }
            # Ignore error if proxy doesn't exist
        try:
            requests.post("http://localhost:8474/proxies", json=proxy_json)
        except requests.RequestException:
            pass  # Ignore errors during proxy creation
    # 初始化存储层
    mysql_client = MySQLClient(app_config.database.mysql)
    redis_client = RedisClient(app_config.database.redis)
    store = Store(mysql_client, redis_client)

    # 初始化业务服务
    business_service = BusinessService(store)

    # 先初始化API层并注册路由到Blueprint
    business_api = BusinessAPI(business_service)
    business_api.register_routes()

    # 然后注册Blueprint到Flask应用
    app.register_blueprint(business_api_bp)

    # Debugging: Print all registered URL rules
    with app.app_context():
        logger.info("Registered URL rules:")
        for rule in app.url_map.iter_rules():
            logger.info(f"  Rule: {rule.endpoint} | Methods: {rule.methods} | Path: {rule.rule}")

    return app, app_config

# 当使用flask run命令时，Flask会调用这个函数
def create_app_for_flask():
    return create_app()[0]

def wait_for_shutdown():
    # Placeholder for graceful shutdown.
    # In a real Flask app, you'd handle signals.
    import signal
    import sys

    def signal_handler(sig, frame):
        logger.info("Shutdown signal received")
        if fault_manager:
            fault_manager.stop_all_faults()
            logger.info("All faults stopped.")
        logger.info("Server stopped gracefully")
        sys.exit(0)

    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)
    logger.info("Press Ctrl+C to shut down gracefully")
    # Keep the main thread alive, waiting for signals
    while True:
        time.sleep(1)

if __name__ == "__main__":
    flask_app, config = create_app()
    port = int(config.server.port)
    logger.info(f"Starting server on port {port}")
    try:
        flask_app.run(host='0.0.0.0', port=port, debug=False)
    except Exception as e:
        logger.error(f"Server error: {e}")
    finally:
        wait_for_shutdown()
