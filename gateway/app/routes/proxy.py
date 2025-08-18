from fastapi import APIRouter, Request, HTTPException
from fastapi.responses import Response
from typing import Dict, Any
import httpx
from urllib.parse import urljoin

from app.config import settings
from app.utils.http_client import http_client
from app.utils.logger import logger

router = APIRouter()


@router.api_route("/api/users", methods=["GET", "POST", "PUT", "DELETE", "PATCH"])
@router.api_route("/api/users/{path:path}", methods=["GET", "POST", "PUT", "DELETE", "PATCH"])
async def proxy_users_api(request: Request, path: str = ""):
    """
    拦截并转发 /api/users 相关的所有请求到目标服务
    
    Args:
        request: FastAPI请求对象
        path: 可选的路径参数 (如 /api/users/123 中的 '123')
    
    Returns:
        Response: 目标服务的响应
    """
    
    try:
        # 构建目标URL
        if path:
            target_path = f"/api/users/{path}"
        else:
            target_path = "/api/users"
            
        target_url = urljoin(settings.target_service_url.rstrip('/') + '/', target_path.lstrip('/'))
        
        # 获取请求体
        body = await request.body()
        json_data = None
        
        # 尝试解析JSON数据
        if body:
            content_type = request.headers.get('content-type', '')
            if 'application/json' in content_type:
                try:
                    import json
                    json_data = json.loads(body)
                except json.JSONDecodeError:
                    logger.warning("Failed to parse JSON body, using raw content")
        
        response = await http_client.forward_request(
            method=request.method,
            target_url=target_url,
            headers=dict(request.headers),
            params=dict(request.query_params),
            json_data=json_data,
            content=body if not json_data else None
        )
        

        
        response_headers = {}
        for key, value in response.headers.items():
            if key.lower() not in ['content-length', 'transfer-encoding', 'connection']:
                response_headers[key] = value
        
        return Response(
            content=response.content,
            status_code=response.status_code,
            headers=response_headers,
            media_type=response.headers.get('content-type')
        )
        
    except httpx.TimeoutException:
        logger.error("Request to target service timed out")
        raise HTTPException(status_code=504, detail="Gateway timeout")
        
    except httpx.ConnectError:
        logger.error("Failed to connect to target service")
        raise HTTPException(status_code=502, detail="Bad gateway - unable to connect to target service")
        
    except httpx.RequestError as e:
        logger.error(f"Request error: {e}")
        raise HTTPException(status_code=502, detail="Bad gateway - request failed")
        
    except Exception as e:
        logger.error(f"Unexpected error in proxy: {e}")
        raise HTTPException(status_code=500, detail="Internal server error")


@router.get("/health")
async def health_check():
    """健康检查端点"""
    return {
        "status": "healthy",
        "service": "api-gateway",
        "target_service": settings.target_service_url
    }


@router.get("/")
async def root():
    """根路径信息"""
    return {
        "message": "API Gateway is running",
        "version": "1.0.0",
        "target_service": settings.target_service_url,
        "endpoints": [
            "GET,POST,PUT,DELETE,PATCH /api/users",
            "GET /health"
        ]
    }