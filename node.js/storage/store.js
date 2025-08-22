const MySQLClient = require('./mysql-client');
const RedisClient = require('./redis-client');
const logger = require('../logging');
const { User } = require('../models/user');

class Store {
    constructor() {
        this.mysql = null;
        this.redis = null;
    }

    async initMySQL(config) {
        this.mysql = new MySQLClient();
        const connected = await this.mysql.connect(config);
        if (!connected) {
            logger.warn('MySQL connection failed, using mock client');
            this.mysql = new MockMySQLClient();
        }
    }

    async initRedis(config) {
        try {
            this.redis = new RedisClient(config.host, config.port, config.password, config.database);
            const connected = await this.redis.connect();
            if (!connected) {
                logger.warn('Redis connection failed, using mock client');
                this.redis = new MockRedisClient();
            }
        } catch (err) {
            logger.warn(`Redis init error (${err.code || err.message}), using mock client`);
            this.redis = new MockRedisClient();
        }
    }

    async queryUsersCached() {
        try {
            // 如果Redis客户端为nil，模拟HTTP操作
            if (!this.redis || !this.redis.connected) {
                logger.info('Redis client is nil. Simulating HTTP operation to fetch users.');
                await new Promise(resolve => setTimeout(resolve, 10));

                const users = [];
                for (let i = 0; i < 10; i++) {
                    const user = new User(
                        i,
                        `Mock User HTTP ${i + 1}`,
                        `mock_http${i + 1}@example.com`
                    );
                    users.push(user);
                }
                logger.info('Successfully simulated HTTP operation for 10 users.');
                return users;
            }

            // 首先尝试从Redis缓存获取用户ID列表
            const userIDs = await this.redis.getUserIDs();
            if (userIDs && userIDs.length > 0) {
                const users = [];
                for (const userID of userIDs) {
                    const user = await this.redis.getUser(userID);
                    if (user) {
                        users.push(user);
                    } else {
                        logger.warn(`Failed to get user ${userID} from Redis`);
                    }
                }

                if (users.length === userIDs.length) {
                    logger.info('All users retrieved from Redis cache by individual IDs.');
                    return users;
                }
                logger.warn('Incomplete users retrieved from Redis cache. Re-fetching from MySQL.');
            }

            // 如果Redis中没有，模拟10个用户并缓存
            const users = [];
            const newUserIDs = [];

            for (let i = 0; i < 10; i++) {
                const user = new User(
                    i,
                    `Mock User ${i + 1}`,
                    `mock${i + 1}@example.com`
                );
                users.push(user);
                newUserIDs.push(user.id);

                // 缓存单个用户
                try {
                    await this.redis.setUser(user, 0);
                } catch (error) {
                    logger.warn(`Failed to cache user ${user.id} in Redis: ${error.message}`);
                }
            }

            // 缓存所有用户ID
            try {
                await this.redis.setUserIDs(newUserIDs, 0);
            } catch (error) {
                logger.warn(`Failed to cache all user IDs in Redis: ${error.message}`);
            }

            logger.info('Mocked 10 users and cached in Redis (individual users and IDs).');
            return users;
        } catch (error) {
            logger.error('Failed to query users:', error.message);
            throw error;
        }
    }

    async disconnect() {
        if (this.mysql) {
            await this.mysql.disconnect();
        }
        if (this.redis) {
            await this.redis.disconnect();
        }
    }

    /**
     * 从数据库获取用户
     * 如果数据库连接没有初始化，则mock 10个用户返回
     * 如果数据库连接初始化但是没有查到用户，则mock 10个用户并写回数据库
     * @returns {Promise<User[]>} 用户列表
     */
    async queryUsersFromDatabase() {
        try {
            logger.info('Querying users from database.');
            // 检查MySQL连接是否初始化
            // 注：在当前实现中，我们通过检查mysql实例是否为MockMySQLClient来判断连接状态
            if (!this.mysql || this.mysql instanceof MockMySQLClient) {
                logger.info('MySQL connection is not initialized. Returning mock users.');
                return this._generateMockUsers(10);
            }

            // 从数据库查询用户
            const users = await this.mysql.query('SELECT * FROM users');

            if (users && users.length > 0) {
                logger.info(`Retrieved ${users.length} users from database.`);
                return users.map(user => new User(user.id, user.name, user.email));
            }

            // 如果没有查到用户，生成mock用户并写入数据库
            logger.info('No users found in database. Generating and writing mock users.');
            const mockUsers = this._generateMockUsers(10);

            // 写入数据库
            for (const user of mockUsers) {
                await this.mysql.query(
                    'INSERT INTO users (id, name, email) VALUES (?, ?, ?)',
                    [user.id, user.name, user.email]
                );
            }

            logger.info(`Successfully wrote ${mockUsers.length} mock users to database.`);
            return mockUsers;
        } catch (error) {
            logger.error('Failed to query users from database:', error.message);
            throw error;
        }
    }

    /**
     * 生成指定数量的mock用户
     * @param {number} count 用户数量
     * @returns {User[]} mock用户列表
     */
    _generateMockUsers(count) {
        const users = [];
        for (let i = 0; i < count; i++) {
            const user = new User(
                Date.now() + i, // 使用时间戳+索引作为唯一ID
                `Mock User ${i + 1}`,
                `mock${i + 1}@example.com`
            );
            users.push(user);
        }
        return users;
    }
}

// 模拟MySQL客户端
class MockMySQLClient {
    async connect() { return false; }
    async query() { return []; }
    async disconnect() { }
}

// 模拟Redis客户端
class MockRedisClient {
    constructor() {
        this.connected = false;
    }
    async connect() { return false; }
    async startFault() { logger.info('Mock Redis fault started'); }
    async stopFault() { logger.info('Mock Redis fault stopped'); }
    async disconnect() { }
}

module.exports = Store;
