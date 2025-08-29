const express = require('express');
const Config = require('./config');
const logger = require('./logging');
const Store = require('./storage/store');
const BusinessService = require('./service/business-service');
const BusinessAPI = require('./api/business-api');

let faultManager;
let server;

async function main() {
    try {
        // 加载配置
        const appConfig = Config.loadConfig();

        // 设置日志级别
        logger.setLevel('info');

        // 初始化存储层
        const store = new Store();
        await store.initMySQL(appConfig.database.mysql);
        await store.initRedis(appConfig.database.redis);

        // 初始化故障管理器
        await initFaultManager(store.redis);

        // 初始化业务服务
        const businessService = new BusinessService(store, faultManager);

        // 初始化API层
        const businessAPI = new BusinessAPI(businessService);

        // 创建Express应用
        const app = express();

        // 中间件
        app.use(express.json());
        app.use(loggingMiddleware);

        // 路由
        app.get('/api/users/1', (req, res) => businessAPI.getUsers1(req, res));
        app.get('/api/users/2', (req, res) => businessAPI.getUsers2(req, res));
        app.get('/api/users/3', (req, res) => businessAPI.getUsers3(req, res));

        // 健康检查
        app.get('/health', (req, res) => {
            res.json({ status: 'ok', timestamp: new Date().toISOString() });
        });

        // 故障状态
        app.get('/faults/status', (req, res) => {
            res.json(faultManager.getStatus());
        });

        // 启动服务器
        server = app.listen(appConfig.server.port, () => {
            logger.info(`server is listening ${appConfig.server.port}`);
        });

        // 优雅关闭
        setupGracefulShutdown();

    } catch (error) {
        logger.error('start failed:', error.message);
    }
}

function loggingMiddleware(req, res, next) {
    const start = Date.now();

    logger.info(`request start: ${req.method} ${req.path}`);

    res.on('finish', () => {
        const duration = Date.now() - start;
        logger.info(`request finish: ${req.method} ${req.path} - ${duration}ms`);
    });

    next();
}

function setupGracefulShutdown() {
    const shutdown = async (signal) => {
        logger.info(`close signal: ${signal}`);

        if (server) {
            server.close(() => {
                logger.info('server shut down');
                process.exit(0);
            });

            setTimeout(() => {
                logger.error('shut down server timeout');
                process.exit(1);
            }, 30000);
        }
    };

    process.on('SIGINT', () => shutdown('SIGINT'));
    process.on('SIGTERM', () => shutdown('SIGTERM'));
}

main().catch(error => {
    logger.error('start server failed:', error);
    process.exit(1);
});
