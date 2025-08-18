class ChaosFault {
    constructor() {
        if (this.constructor === ChaosFault) {
            throw new Error('ChaosFault is an abstract class');
        }
    }

    async start(params = {}) {
        throw new Error('start method must be implemented');
    }

    async stop() {
        throw new Error('stop method must be implemented');
    }

    isActive() {
        throw new Error('isActive method must be implemented');
    }

    getName() {
        throw new Error('getName method must be implemented');
    }
}

module.exports = ChaosFault;
