import time
from fastapi import Request, Response
from starlette.middleware.base import BaseHTTPMiddleware
from app.utils.logger import logger


class LoggingMiddleware(BaseHTTPMiddleware):
    """日志记录中间件"""
    
    async def dispatch(self, request: Request, call_next):
        start_time = time.time()
        
        # 记录请求信息
        logger.info(
            f"Incoming request: {request.method} {request.url} "
            f"from {request.client.host if request.client else 'unknown'}"
        )
        
        try:
            response = await call_next(request)
            
            # 计算处理时间
            process_time = time.time() - start_time
            
            # 记录响应信息
            logger.info(
                f"Request completed: {request.method} {request.url} "
                f"- Status: {response.status_code} "
                f"- Time: {process_time:.3f}s"
            )
            
            # 添加处理时间到响应头
            response.headers["X-Process-Time"] = str(process_time)
            
            return response
            
        except Exception as e:
            # 记录错误信息
            process_time = time.time() - start_time
            logger.error(
                f"Request failed: {request.method} {request.url} "
                f"- Error: {str(e)} "
                f"- Time: {process_time:.3f}s"
            )
            raise