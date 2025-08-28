const ChaosFault = require('./fault');
const Config = require('../config');

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
            return;
        }

        try {
            if (this.redisClient && this.redisClient.startFault) {
                await this.redisClient.startFault(delay);
                this.active = true;
                this.delay = delay;
            } else {
                this.active = true;
            }
        } catch (error) {
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
        } catch (error) {
            throw error;
        }
    }

    isActive() {
        return this.active;
    }
}

module.exports = RedisLatencyFault;
