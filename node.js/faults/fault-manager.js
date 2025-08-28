class FaultManager {
    constructor() {
        this.faults = new Map();
    }

    register(fault) {
        this.faults.set(fault.getName(), fault);
    }

    async startFault(chaosType, params = {}) {
        const fault = this.faults.get(chaosType);
        if (!fault) {
            throw new Error(`Fault ${chaosType} not found`);
        }

        try {
            await fault.start(params);
            return true;
        } catch (error) {
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
            return true;
        } catch (error) {
            throw error;
        }
    }

    async stopAllFaults() {
        const stopPromises = Array.from(this.faults.values()).map(fault => fault.stop());
        try {
            await Promise.all(stopPromises);
        } catch (error) {
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
