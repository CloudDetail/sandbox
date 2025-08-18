const ChaosFault = require('./fault');
const Config = require('../config');
const logger = require('../logging');

class RedisLatencyFault extends ChaosFault {
    constructor(redisClient) {
        super();
        this.active = false;
        this.redisClient = redisClient;
        this.delay = 0;
    }

    getName() {
        return 'redis_latency';
    }

    async start(params = {}) {
        const config = Config.loadConfig();
        let delay = config.faults.redis.defaultDelay;

        if (params.duration) {
            delay = parseInt(params.duration);
        }

        if (this.active && this.delay == delay) {
            logger.info('Redis fault already active');
            return;
        }

        try {
            if (this.redisClient && this.redisClient.startFault) {
                await this.redisClient.startFault(delay);
                this.active = true;
                this.delay = delay;
                logger.info(`Redis latency fault started with delay: ${delay}ms`);
            } else {
                logger.warn('Redis client not available, simulating fault behavior');
                this.active = true;
            }
        } catch (error) {
            logger.error(`Failed to start redis latency fault: ${error.message}`);
            throw error;
        }
    }

    async stop() {
        if (!this.active) {
            return;
        }

        try {
            if (this.redisClient && this.redisClient.stopFault) {
                await this.redisClient.stopFault();
            }
            this.active = false;
            logger.info('Redis latency fault stopped');
        } catch (error) {
            logger.error(`Failed to stop redis latency fault: ${error.message}`);
            throw error;
        }
    }

    isActive() {
        return this.active;
    }
}

module.exports = RedisLatencyFault;
