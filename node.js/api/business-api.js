const BusinessService = require('../service/business-service');
const logger = require('../logging');

class BusinessAPI {
    constructor(service) {
        this.service = service;
    }

    async getUsersCached(req, res) {
        try {
            const chaos = req.query.chaos;
            const durationParam = req.query.duration;
            const duration = durationParam ? parseInt(durationParam) : 0;

            const result = await this.service.getUsersCached(chaos, duration);

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
