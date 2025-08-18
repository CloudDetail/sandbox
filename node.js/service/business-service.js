const logger = require('../logging');

class BusinessService {
    constructor(store, faultManager) {
        this.store = store;
        this.faultManager = faultManager;
    }

    async getUsersCached(chaosType, duration) {
        try {
            if (chaosType && chaosType.length > 0) {
                const params = { duration: duration || 0 };
                await this.faultManager.startFault(chaosType, params);
            } else {
                await this.faultManager.stopAllFaults();
                logger.info('All faults stopped');
            }

            if (chaosType == "redis_latency") {
                const users = await this.store.queryUsersCached();
                return JSON.stringify(users);
            } else {
                const users = await this.store.queryUsersFromDatabase();
                return JSON.stringify(users);
            }
        } catch (error) {
            logger.error('GetUsers failed:', error.message);
            throw error;
        }
    }
}

module.exports = BusinessService;
