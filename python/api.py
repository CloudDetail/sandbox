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
            chaos = request.args.get("chaos")
            duration_param = request.args.get("duration")
            duration = 0

            if duration_param:
                try:
                    duration = int(duration_param)
                except ValueError:
                    logger.error(f"Invalid duration parameter: {duration_param}")
                    return jsonify({"error": "Invalid duration parameter"}), 400

            if chaos and chaos != "none": # Check if chaos is explicitly provided and not "none"
                result, err = self.service.get_users_cached(chaos, duration)
            else:
                result, err = self.service.get_users_cached(None, 0) # Stop all faults if no chaos or "none"

            if err:
                logger.error(f"Error in GetUsersCached: {err}")
                return jsonify({"error": str(err)}), 500

            return result, 200, {'Content-Type': 'application/json'}
