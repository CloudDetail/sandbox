# python/config.py
from dotenv import load_dotenv
import os
from datetime import timedelta

class ServerConfig:
    def __init__(self):
        self.port = os.getenv("PORT", "3500")
        self.read_timeout = timedelta(seconds=int(os.getenv("READ_TIMEOUT_SEC", 30)))
        self.write_timeout = timedelta(seconds=int(os.getenv("WRITE_TIMEOUT_SEC", 30)))
        self.deploy_proxy = os.getenv("DEPLOY_PROXY", "false") == "true"

class MySQLConfig:
    def __init__(self):
        self.host = os.getenv("MYSQL_HOST", "mysql-service")
        self.port = int(os.getenv("MYSQL_PORT", 3306))
        self.username = os.getenv("MYSQL_USERNAME", "root")
        self.password = os.getenv("MYSQL_PASSWORD", "")
        self.database = os.getenv("MYSQL_DATABASE", "sandbox")
        self.max_connections = int(os.getenv("MYSQL_MAX_CONNECTIONS", 10))
        self.conn_timeout = timedelta(seconds=int(os.getenv("MYSQL_CONN_TIMEOUT_SEC", 30)))
        self.read_timeout = timedelta(seconds=int(os.getenv("MYSQL_READ_TIMEOUT_SEC", 10)))
        self.write_timeout = timedelta(seconds=int(os.getenv("MYSQL_WRITE_TIMEOUT_SEC", 10)))

class RedisConfig:
    def __init__(self):
        self.host = os.getenv("REDIS_HOST", "localhost")
        self.port = int(os.getenv("REDIS_PORT", 6379))
        self.password = os.getenv("REDIS_PASSWORD", "")
        self.database = int(os.getenv("REDIS_DATABASE", 0))

class CPUFaultConfig:
    def __init__(self):
        self.default_duration = int(os.getenv("CPU_FAULT_DEFAULT_DURATION", 200))

class LatencyFaultConfig:
    def __init__(self):
        self.default_delay = int(os.getenv("LATENCY_FAULT_DEFAULT_DELAY", 200))
        self.max_delay = int(os.getenv("LATENCY_FAULT_MAX_DELAY", 5000))

class RedisFaultConfig:
    def __init__(self):
        self.default_delay = int(os.getenv("REDIS_FAULT_DEFAULT_DELAY", 100))
        self.max_delay = int(os.getenv("REDIS_FAULT_MAX_DELAY", 2000))

class FaultsConfig:
    def __init__(self):
        self.cpu = CPUFaultConfig()
        self.latency = LatencyFaultConfig()
        self.redis = RedisFaultConfig()

class DatabaseConfig:
    def __init__(self):
        self.mysql = MySQLConfig()
        self.redis = RedisConfig()

class Config:
    _instance = None

    def __new__(cls):
        if cls._instance is None:
            cls._instance = super(Config, cls).__new__(cls)
            load_dotenv()
            cls._instance.server = ServerConfig()
            cls._instance.database = DatabaseConfig()
            cls._instance.faults = FaultsConfig()
        return cls._instance

def load_config():
    return Config()
