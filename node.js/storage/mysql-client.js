const mysql = require('mysql2/promise');
const logger = require('../logging');

class MySQLClient {
    constructor() {
        this.connection = null;
        this.connected = false;
    }

    async connect(config) {
        try {
            this.connection = await mysql.createConnection({
                host: config.host,
                port: config.port,
                user: config.username,
                password: config.password,
                database: config.database,
                connectTimeout: config.connTimeout,
                acquireTimeout: config.connTimeout,
                timeout: config.readTimeout,
                charset: 'utf8mb4'
            });

            this.connected = true;
            logger.info('MySQL Client Connected');
            return true;
        } catch (error) {
            logger.warn('Failed to connect to MySQL:', error.message);
            this.connected = false;
            return false;
        }
    }

    async query(sql, params = []) {
        if (!this.connection || !this.connected) {
            throw new Error('MySQL client not connected');
        }

        try {
            const [rows] = await this.connection.execute(sql, params);
            return rows;
        } catch (error) {
            logger.error('MySQL query failed:', error.message);
            throw error;
        }
    }

    async disconnect() {
        if (this.connection) {
            await this.connection.end();
            this.connected = false;
            logger.info('MySQL client disconnected');
        }
    }
}

module.exports = MySQLClient;
