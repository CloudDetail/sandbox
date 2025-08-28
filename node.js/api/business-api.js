const BusinessService = require('../service/business-service');
const logger = require('../logging');

class BusinessAPI {
    constructor(service) {
        this.service = service;
    }

    async getUsersCached(req, res) {
        try {
            const mode = req.query.mode;
            let chaos = "";
            switch (mode) {
                case '1':
                    chaos = 'latency';
                    break;
                case '2':
                    chaos = 'cpu';
                    break;
                case '3':
                    chaos = 'redis_latency';
                    break;
                default:
                    break;
            }

            const result = await this.service.getUsersCached(chaos, 0);

            res.json({
                data: result
            });
        } catch (error) {
            logger.error('API error:', error.message);
            res.status(500).json({
                error: error.message
            });
        }
    }
}

module.exports = BusinessAPI;
