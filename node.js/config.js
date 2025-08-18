require('dotenv').config();

class Config {
    constructor() {
        this.server = {
            port: process.env.PORT || '3500',
            readTimeout: parseInt(process.env.READ_TIMEOUT) || 30000,
            writeTimeout: parseInt(process.env.WRITE_TIMEOUT) || 30000
        };

        this.database = {
            mysql: {
                host: process.env.MYSQL_HOST || 'localhost',
                port: parseInt(process.env.MYSQL_PORT) || 3306,
                username: process.env.MYSQL_USERNAME || 'root',
                password: process.env.MYSQL_PASSWORD || '',
                database: process.env.MYSQL_DATABASE || 'sandbox',
                maxConnections: parseInt(process.env.MYSQL_MAX_CONNECTIONS) || 10,
                connTimeout: parseInt(process.env.MYSQL_CONN_TIMEOUT) || 30000,
                readTimeout: parseInt(process.env.MYSQL_READ_TIMEOUT) || 10000,
                writeTimeout: parseInt(process.env.MYSQL_WRITE_TIMEOUT) || 10000
            },
            redis: {
                host: process.env.REDIS_HOST || 'localhost',
                port: parseInt(process.env.REDIS_PORT) || 6379,
                password: process.env.REDIS_PASSWORD || '',
                database: parseInt(process.env.REDIS_DATABASE) || 0
            }
        };

        this.faults = {
            cpu: {
                defaultDuration: parseInt(process.env.CPU_FAULT_DEFAULT_DURATION) || 200
            },
            latency: {
                defaultDelay: parseInt(process.env.LATENCY_FAULT_DEFAULT_DELAY) || 200,
                maxDelay: parseInt(process.env.LATENCY_FAULT_MAX_DELAY) || 5000
            },
            redis: {
                defaultDelay: parseInt(process.env.REDIS_FAULT_DEFAULT_DELAY) || 100,
                maxDelay: parseInt(process.env.REDIS_FAULT_MAX_DELAY) || 2000
            }
        };
    }

    static loadConfig() {
        if (!Config.instance) {
            Config.instance = new Config();
        }
        return Config.instance;
    }
}

module.exports = Config;
