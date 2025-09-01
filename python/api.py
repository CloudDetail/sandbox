# python/api.py
from flask import Blueprint, request, jsonify
import logging

from service import BusinessService

business_api_bp = Blueprint('business_api', __name__)
logger = logging.getLogger(__name__)

class BusinessAPI:
    def __init__(self, service: BusinessService):
        self.service = service

    def register_routes(self):
        @business_api_bp.route("/api/users/1", methods=["GET"])
        def get_users_1():
            mode_str = request.args.get("mode", "0")
            try:
                mode = int(mode_str)
            except ValueError:
                mode = 0  # Default to 0 if conversion fails

            result, err = self.service.get_users_latency(mode)
            if err:
                logger.error(f"Error in GetUsersCached: {err}")
                return jsonify({"error": str(err)}), 500

            return result, 200, {'Content-Type': 'application/json'}

        @business_api_bp.route("/api/users/2", methods=["GET"])
        def get_users_2():
            mode_str = request.args.get("mode", "0")
            try:
                mode = int(mode_str)
            except ValueError:
                mode = 0  # Default to 0 if conversion fails

            result, err = self.service.get_users_cpu_burn(mode)
            if err:
                logger.error(f"Error in GetUsersCached: {err}")
                return jsonify({"error": str(err)}), 500

            return result, 200, {'Content-Type': 'application/json'}

        @business_api_bp.route("/api/users/3", methods=["GET"])
        def get_users_3():
            mode_str = request.args.get("mode", "0")
            try:
                mode = int(mode_str)
            except ValueError:
                mode = 0  # Default to 0 if conversion fails

            result, err = self.service.get_users_redis_latency(mode)
            if err:
                logger.error(f"Error in GetUsersCached: {err}")
                return jsonify({"error": str(err)}), 500

            return result, 200, {'Content-Type': 'application/json'}
