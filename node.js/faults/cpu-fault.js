const ChaosFault = require('./fault');
const Config = require('../config');

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
            return;
        }

        const config = Config.loadConfig();
        let targetDuration = config.faults.cpu.defaultDuration;

        if (params.duration) {
            targetDuration = parseInt(params.duration);
        }

        this.active = true;

        const startTime = Date.now();
        while (Date.now() - startTime < targetDuration) {
            this.fibonacci(38);
        }

        const actualDuration = Date.now() - startTime;
        this.active = false;
    }

    async stop() {
        this.active = false;
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
