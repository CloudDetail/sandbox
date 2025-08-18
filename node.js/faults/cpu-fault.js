const ChaosFault = require('./fault');
const Config = require('../config');
const logger = require('../logging');

class CPUFault extends ChaosFault {
    constructor() {
        super();
        this.active = false;
    }

    getName() {
        return 'cpu';
    }

    async start(params = {}) {
        if (this.active) {
            logger.info('CPU fault already active');
            return;
        }

        const config = Config.loadConfig();
        let targetDuration = config.faults.cpu.defaultDuration;

        if (params.duration) {
            targetDuration = parseInt(params.duration);
        }

        this.active = true;
        logger.info(`CPU fault started, will consume CPU for ${targetDuration}ms`);

        const startTime = Date.now();
        while (Date.now() - startTime < targetDuration) {
            this.fibonacci(38);
        }

        const actualDuration = Date.now() - startTime;
        logger.info(`CPU fault completed, consumed ${actualDuration}ms CPU time`);
        this.active = false;
    }

    async stop() {
        this.active = false;
        logger.info('CPU fault stopped');
    }

    isActive() {
        return this.active;
    }

    fibonacci(n) {
        if (n <= 1) {
            return n;
        }
        return this.fibonacci(n - 1) + this.fibonacci(n - 2);
    }
}

module.exports = CPUFault;
