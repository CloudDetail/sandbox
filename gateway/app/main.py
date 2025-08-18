from contextlib import asynccontextmanager
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.config import settings
from app.middleware.logging import LoggingMiddleware
from app.routes import proxy
from app.utils.logger import logger


@asynccontextmanager
async def lifespan(app: FastAPI):
    """应用生命周期管理"""
    logger.info("API Gateway starting up...")
    logger.info(f"Target service URL: {settings.target_service_url}")
    logger.info(f"Gateway listening on: {settings.gateway_host}:{settings.gateway_port}")
    
    yield
    
    logger.info("API Gateway shutting down...")


def create_app() -> FastAPI:
    """创建FastAPI应用实例"""
    
    app = FastAPI(
        title="API Gateway",
        description="A simple API gateway for forwarding requests to target services",
        version="1.0.0",
        lifespan=lifespan
    )
    
    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )
    
    app.add_middleware(LoggingMiddleware)
    
    app.include_router(proxy.router)
    
    return app


app = create_app()


if __name__ == "__main__":
    import uvicorn
    
    uvicorn.run(
        "app.main:app",
        host=settings.gateway_host,
        port=settings.gateway_port,
        reload=True,
        log_level=settings.log_level.lower()
    )