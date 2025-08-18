const redis = require('redis');
const logger = require('../logging');
const { User, Order } = require('../models/user');

class RedisClient {
    constructor(host = 'localhost', port = 6379, password = '', database = 0) {
        this.host = host;
        this.port = port;
        this.password = password;
        this.database = database;
        this.client = null;
        this.connected = false;
    }

    async connect() {
        try {
            this.client = redis.createClient({
                socket: {
                    host: this.host,
                    port: this.port,
                    reconnectStrategy: false, // 禁用重连
                    connectTimeout: 1000      // 1 秒超时
                },
                password: this.password,
                database: this.database
            });
    
            this.client.on('error', () => {
                this.connected = false;
            });
    
            this.client.on('connect', () => {
                logger.info('Redis Client Connected');
                this.connected = true;
            });
    
            await this.client.connect();
            return true;
        } catch (error) {
            this.connected = false;
            return false;
        }
    }    

    async ping() {
        if (!this.client || !this.connected) {
            return false;
        }
        try {
            const result = await this.client.ping();
            return result === 'PONG';
        } catch (error) {
            logger.error('Redis ping failed:', error.message);
            return false;
        }
    }

    async get(key) {
        if (!this.client || !this.connected) {
            throw new Error('Redis client not connected');
        }
        return await this.client.get(key);
    }

    async set(key, value, expiration = 0) {
        if (!this.client || !this.connected) {
            throw new Error('Redis client not connected');
        }
        if (expiration > 0) {
            return await this.client.setEx(key, expiration, value);
        } else {
            return await this.client.set(key, value);
        }
    }

    // 用户相关操作
    async setUser(user, expiration = 0) {
        const key = `user:${user.id}`;
        const userData = JSON.stringify(user.toJSON());
        return await this.set(key, userData, expiration);
    }

    async getUser(id) {
        const key = `user:${id}`;
        try {
            const userData = await this.get(key);
            if (userData === null) {
                return null; // User not found in Redis
            }
            return User.fromJSON(JSON.parse(userData));
        } catch (error) {
            logger.error(`Failed to get user ${id}:`, error.message);
            return null;
        }
    }

    // 用户ID列表操作
    async setUserIDs(userIDs, expiration = 0) {
        const userIDsJSON = JSON.stringify(userIDs);
        return await this.set('all_user_ids', userIDsJSON, expiration);
    }

    async getUserIDs() {
        try {
            const userIDsData = await this.get('all_user_ids');
            if (userIDsData === null) {
                return null; // User IDs not found in Redis
            }
            return JSON.parse(userIDsData);
        } catch (error) {
            logger.error('Failed to get user IDs:', error.message);
            return null;
        }
    }

    async startFault(delay) {
        if (!this.client || !this.connected) {
            throw new Error('Redis client not connected');
        }

        try {
            const result = await this.client.sendCommand(['FAULT.START', delay.toString()]);
            logger.info(`Start redis latency ${delay}.`);
            return result;
        } catch (error) {
            logger.error(`Failed to send fault command to Redis proxy: ${error.message}`);
            throw error;
        }
    }

    async stopFault() {
        if (!this.client || !this.connected) {
            throw new Error('Redis client not connected');
        }

        try {
            const result = await this.client.sendCommand(['FAULT.STOP']);
            logger.info('Redis latency stopped');
            return result;
        } catch (error) {
            logger.error(`Failed to send stop fault command to Redis proxy: ${error.message}`);
            throw error;
        }
    }

    async disconnect() {
        if (this.client) {
            await this.client.quit();
            this.connected = false;
            logger.info('Redis client disconnected');
        }
    }
}

module.exports = RedisClient;
