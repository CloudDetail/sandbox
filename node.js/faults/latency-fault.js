const ChaosFault = require('./fault');
const Config = require('../config');
const logger = require('../logging');
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
            logger.info('Latency fault already active');
            return
        }

        try {
            await this.clearTC();

            const cmd = `tc qdisc add dev ${this.iface} root netem delay ${delayMs}ms`;
            const { stdout, stderr } = await execAsync(cmd);

            if (stderr && !stderr.includes('RTNETLINK answers: File exists')) {
                logger.warn(`TC command warning: ${stderr}`);
            }

            this.delay = delayMs;
            this.active = true;
            logger.info(`Successfully added ${delayMs}ms delay on ${this.iface}`);
        } catch (error) {
            logger.error(`Failed to add delay: ${error.message}`);
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
            logger.info('Latency fault stopped');
        } catch (error) {
            logger.error(`Failed to stop latency fault: ${error.message}`);
            throw error;
        }
    }

    async clearTC() {
        try {
            const cmd = `tc qdisc del dev ${this.iface} root`;
            const { stdout, stderr } = await execAsync(cmd);

            if (stderr && !stderr.includes('No such file or directory') &&
                !stderr.includes('No qdisc')) {
                logger.warn(`TC clear warning: ${stderr}`);
            }
        } catch (error) {
            if (!error.message.includes('No such file or directory') &&
                !error.message.includes('No qdisc')) {
                logger.warn(`TC clear error: ${error.message}`);
            }
        }
    }

    isActive() {
        return this.active;
    }
}

module.exports = LatencyFault;
