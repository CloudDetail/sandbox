const BusinessService = require('../service/business-service');
const logger = require('../logging');

class BusinessAPI {
    constructor(service) {
        this.service = service;
    }

    async getUsers1(req,res) {
        try {
            const mode = parseInt(req.query.mode) || 0;
            const result = await this.service.getUsersWithLatency(mode);
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

    async getUsers2(req,res) {
        try {
            const mode = parseInt(req.query.mode) || 0;
            const result = await this.service.getUsersWithCPUBurn(mode);
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

    async getUsers3(req,res) {
        try {
            const mode = parseInt(req.query.mode) || 0;
            const result = await this.service.getUsersWithRedisLatency(mode);
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
