const logger = require('../logging');
const toxiproxyClient = require("toxiproxy-node-client");
const toxiproxy = new toxiproxyClient.Toxiproxy("http://localhost:8474");
const proxyBody = {
    listen: "localhost:6379",
    name: "basicProxy",
    upstream: "redis-service:6379"
};
const proxy = toxiproxy.createProxy(proxyBody)

/**
 * BusinessService class provides methods for fault injection testing
 * including network latency, CPU burn, and Redis latency simulation
 */
class BusinessService {
    /**
     * Initialize BusinessService with store and network interface
     * @param {Object} store - Data store object
     * @param {string} iface - Network interface for traffic control (default: 'eth0')
     */
    constructor(store,iface = 'eth0') {
        this.store = store;
        this.active = false;
        this.iface = iface;
        this.toxic = null;
    }

    /**
     * Simulate network latency by adding traffic control rules
     * @param {number} mode - 1 to enable latency, 0 to disable
     * @returns {Promise<string>} JSON stringified users data
     */
    async getUserWithLatency(mode) {
        if (mode == 1) {
            // Use the `active` variable to record the fault injection status
            // keep it enabled during the injection process.
            if (!this.active) {
                await this.clearTC();
                const cmd = `tc qdisc add dev ${this.iface} root netem delay 200ms`;
                const { stdout, stderr } = await execAsync(cmd);
                if (stderr && !stderr.includes('RTNETLINK answers: File exists')) {
                    logger.warn(`TC command warning: ${stderr}`);
                }
                this.active = true;
            }
        }else{
            if (this.active) {
                await this.clearTC();
                this.active = false;
            }
        }

        await this.store.queryUsersCachedFromRedis();
        const users = await this.store.queryUsersFromDatabase();
        return JSON.stringify(users);
    }

    /**
     * Clear traffic control rules from network interface
     * @returns {Promise<void>}
     */
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

    /**
     * Simulate high CPU usage by calculating Fibonacci numbers
     * @param {number} mode - 1 to enable CPU burn, 0 to disable
     * @returns {Promise<string>} JSON stringified users data
     */
    async getUsersWithCpuBurn(mode) {
        if (mode == 1) {
            const startTime = Date.now();
            while (Date.now() - startTime < 200) {
                this.fibonacci(38);
            }
        }

        await this.store.queryUsersCachedFromRedis();
        const users = await this.store.queryUsersFromDatabase();
        return JSON.stringify(users);
    }

    /**
     * Calculate Fibonacci number recursively (CPU intensive)
     * @param {number} n - Number to calculate Fibonacci for
     * @returns {number} Fibonacci result
     */
    fibonacci(n) {
        if (n <= 1) {
            return n;
        }
        return this.fibonacci(n - 1) + this.fibonacci(n - 2);
    }

    /**
     * Simulate Redis latency using Toxiproxy
     * @param {number} mode - 1 to enable Redis latency, 0 to disable
     * @returns {Promise<string>} JSON stringified users data
     */
    async getUsersWithRedisLatency(mode) {
        if (mode == 1) {
            // Use the `active` variable to record the fault injection status
            // keep it enabled during the injection process.
            if (!this.active) {
                const toxicBody = {
                    name: "latency",
                    attributes: { latency: 200 },
                    type: "latency"
                };
                // Toxiproxy is a framework for simulating network conditions.
                // https://github.com/shopify/toxiproxy
                this.toxic = await proxy.addToxic(toxicBody);
            }
        }else{
            if (this.active && this.toxic) {
                await this.toxic.remove();
                this.toxic = null;
                this.active = false;
            }
        }

        await this.store.queryUsersCachedFromRedis();
        const users = await this.store.queryUsersFromDatabase();
        return JSON.stringify(users);
    }
}

module.exports = BusinessService;
