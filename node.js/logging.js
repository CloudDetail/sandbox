const winston = require('winston');

class Logger {
    constructor() {
        this.logger = winston.createLogger({
            level: 'info',
            format: winston.format.combine(
                winston.format.timestamp(),
                winston.format.errors({ stack: true }),
                winston.format.json()
            ),
            defaultMeta: { service: 'apo-sandbox' },
            transports: [
                new winston.transports.Console({
                    format: winston.format.combine(
                        winston.format.colorize(),
                        winston.format.simple()
                    )
                })
            ]
        });
    }

    setLevel(level) {
        this.logger.level = level;
    }

    info(message, ...args) {
        if (args.length > 0) {
            this.logger.info(message, ...args);
        } else {
            this.logger.info(message);
        }
    }

    warn(message, ...args) {
        if (args.length > 0) {
            this.logger.warn(message, ...args);
        } else {
            this.logger.warn(message);
        }
    }

    error(message, ...args) {
        if (args.length > 0) {
            this.logger.error(message, ...args);
        } else {
            this.logger.error(message);
        }
    }

    debug(message, ...args) {
        if (args.length > 0) {
            this.logger.debug(message, ...args);
        } else {
            this.logger.debug(message);
        }
    }
}

const logger = new Logger();

module.exports = logger;
