const ChaosFault = require('./fault');
const Config = require('../config');
const { exec } = require('child_process');
const { promisify } = require('util');

const execAsync = promisify(exec);

class LatencyFault extends ChaosFault {
    constructor(iface = 'eth0') {
        super();
        this.active = false;
        this.iface = iface;
        this.delay = 0;
    }

    getName() {
        return 'latency';
    }

    async start(params = {}) {
        const config = Config.loadConfig();
        let delayMs = config.faults.latency.defaultDelay;

        if (params.duration && params.duration > 0) {
            delayMs = parseInt(params.duration);
        }

        if (this.active) {
            return
        }

        try {
            await this.clearTC();

            const cmd = `tc qdisc add dev ${this.iface} root netem delay ${delayMs}ms`;
            const { stdout, stderr } = await execAsync(cmd);


            this.delay = delayMs;
            this.active = true;
        } catch (error) {
            throw error;
        }
    }

    async stop() {
        if (!this.active) {
            return;
        }

        try {
            await this.clearTC();
            this.active = false;
        } catch (error) {
            throw error;
        }
    }

    async clearTC() {
        try {
            const cmd = `tc qdisc del dev ${this.iface} root`;
            const { stdout, stderr } = await execAsync(cmd);

            if (stderr && !stderr.includes('No such file or directory') &&
                !stderr.includes('No qdisc')) {
            }
        } catch (error) {
            if (!error.message.includes('No such file or directory') &&
                !error.message.includes('No qdisc')) {
            }
        }
    }

    isActive() {
        return this.active;
    }
}

module.exports = LatencyFault;
