const logger = require('../logging');

class FaultManager {
    constructor() {
        this.faults = new Map();
    }

    register(fault) {
        this.faults.set(fault.getName(), fault);
        logger.info(`Fault registered: ${fault.getName()}`);
    }

    async startFault(chaosType, params = {}) {
        const fault = this.faults.get(chaosType);
        if (!fault) {
            throw new Error(`Fault ${chaosType} not found`);
        }

        try {
            await fault.start(params);
            logger.info(`Fault ${chaosType} started successfully`);
            return true;
        } catch (error) {
            logger.error(`Failed to start fault ${chaosType}: ${error.message}`);
            throw error;
        }
    }

    async stopFault(chaosType) {
        const fault = this.faults.get(chaosType);
        if (!fault) {
            throw new Error(`Fault ${chaosType} not found`);
        }

        try {
            await fault.stop();
            logger.info(`Fault ${chaosType} stopped successfully`);
            return true;
        } catch (error) {
            logger.error(`Failed to stop fault ${chaosType}: ${error.message}`);
            throw error;
        }
    }

    async stopAllFaults() {
        const stopPromises = Array.from(this.faults.values()).map(fault => fault.stop());
        try {
            await Promise.all(stopPromises);
        } catch (error) {
            logger.error(`Error stopping all faults: ${error.message}`);
        }
    }

    getStatus() {
        const status = {};
        for (const [name, fault] of this.faults) {
            status[name] = {
                active: fault.isActive(),
                name: fault.getName()
            };
        }
        return status;
    }

    listActive() {
        const active = [];
        for (const [name, fault] of this.faults) {
            if (fault.isActive()) {
                active.push(name);
            }
        }
        return active;
    }
}

module.exports = FaultManager;
