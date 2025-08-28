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
        @business_api_bp.route("/api/users", methods=["GET"])
        def get_users_cached():
            mode = request.args.get("mode")
            chaos = ""
            if mode == "1":
                chaos = "latency"
            elif mode == "2":
                chaos = "cpu"
            elif mode == "3":
                chaos = "redis_latency"
            if chaos and chaos != "none": # Check if chaos is explicitly provided and not "none"
                result, err = self.service.get_users_cached(chaos, 0)
            else:
                result, err = self.service.get_users_cached(None, 0) # Stop all faults if no chaos or "none"

            if err:
                logger.error(f"Error in GetUsersCached: {err}")
                return jsonify({"error": str(err)}), 500

            return result, 200, {'Content-Type': 'application/json'}
