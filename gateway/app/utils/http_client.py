import httpx
from typing import Dict, Any, Optional
from app.config import settings
from app.utils.logger import logger


class HTTPClient:
    """HTTP客户端工具类"""
    
    def __init__(self):
        self.timeout = httpx.Timeout(settings.request_timeout)
        
    async def forward_request(
        self,
        method: str,
        target_url: str,
        headers: Optional[Dict[str, str]] = None,
        params: Optional[Dict[str, Any]] = None,
        json_data: Optional[Dict[str, Any]] = None,
        content: Optional[bytes] = None
    ) -> httpx.Response:
        """
        转发HTTP请求到目标服务
        
        Args:
            method: HTTP方法 (GET, POST, PUT, DELETE等)
            target_url: 目标URL
            headers: 请求头
            params: 查询参数
            json_data: JSON数据
            content: 原始请求体
            
        Returns:
            httpx.Response: 目标服务的响应
        """
        
        # 过滤掉可能导致问题的头部，并添加标准头部
        filtered_headers = {}
        if headers:
            skip_headers = {'host', 'content-length', 'transfer-encoding', 'connection'}
            filtered_headers = {
                k: v for k, v in headers.items() 
                if k.lower() not in skip_headers
            }
        
        # 如果没有User-Agent，添加一个标准的
        if 'user-agent' not in [k.lower() for k in filtered_headers.keys()]:
            filtered_headers['User-Agent'] = 'Gateway-Client/1.0'
        
        logger.info(f"Forwarding {method} request to: {target_url}")
        logger.info(f"Headers: {filtered_headers}")
        logger.info(f"Params: {params}")
        
        async with httpx.AsyncClient(timeout=self.timeout) as client:
            try:
                response = await client.request(
                    method=method,
                    url=target_url,
                    headers=filtered_headers,
                    params=params,
                    json=json_data,
                    content=content
                )
                
                return response
                
            except httpx.TimeoutException:
                logger.error(f"Request to {target_url} timed out")
                raise
            except httpx.RequestError as e:
                logger.error(f"Request to {target_url} failed: {e}")
                raise
            except Exception as e:
                logger.error(f"Unexpected error forwarding request: {e}")
                raise


http_client = HTTPClient()