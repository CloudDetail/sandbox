import logging
import sys
from typing import Optional
from app.config import settings


def setup_logger(name: Optional[str] = None) -> logging.Logger:
    """设置日志记录器"""
    
    logger = logging.getLogger(name or __name__)
    
    if not logger.handlers:
        # 设置日志级别
        log_level = getattr(logging, settings.log_level.upper(), logging.INFO)
        logger.setLevel(log_level)
        
        # 创建控制台处理器
        handler = logging.StreamHandler(sys.stdout)
        handler.setLevel(log_level)
        
        # 设置日志格式
        formatter = logging.Formatter(
            '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
        )
        handler.setFormatter(formatter)
        
        logger.addHandler(handler)
    
    return logger


# 创建默认日志记录器
logger = setup_logger("gateway")