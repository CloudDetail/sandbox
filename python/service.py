# python/service.py
import logging
import json

from storage import Store
from faults.fault_manager import Manager

logger = logging.getLogger(__name__)

class BusinessService:
    def __init__(self, store: Store, fault_manager: Manager):
        self.store = store
        self.fault_manager = fault_manager

    def get_users_cached(self, chaos_type: str = None, duration: int = 0) -> tuple[str, None]:
        if chaos_type:
            params = {}
            if duration > 0:
                params = {"duration": duration}
            self.fault_manager.start_fault(chaos_type, params)
        else:
            self.fault_manager.stop_all_faults()

        if chaos_type == "redis_latency":
            users, err = self.store.query_or_create_users()
            if err:
                return "", err
        else:
            users, err = self.store.query_or_create_users()
            if err:
                return "", err
        
        return json.dumps([user.__dict__ for user in users]), None
