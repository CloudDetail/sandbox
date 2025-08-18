from pydantic_settings import BaseSettings
from typing import Optional


class Settings(BaseSettings):
    """应用程序配置"""
    
    # Gateway配置
    gateway_host: str = "0.0.0.0"
    gateway_port: int = 8000
    
    # 目标服务配置
    target_service_url: str
    
    # 日志配置
    log_level: str = "INFO"
    
    # HTTP客户端配置
    request_timeout: int = 30
    
    class Config:
        env_file = ".env"
        case_sensitive = False


# 全局配置实例
settings = Settings()