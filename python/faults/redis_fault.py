# python/faults/redis_fault.py
import logging
import threading

from faults.fault_manager import ChaosFault
from config import load_config
from storage import RedisClient

logger = logging.getLogger(__name__)

class RedisLatencyFault(ChaosFault):
    def __init__(self, redis_client: RedisClient):
        self._active = False
        self._lock = threading.Lock()
        self._redis_client = redis_client
        self._original_redis_get = None
        self._original_redis_set = None
        self._fault_delay = 0

    def name(self) -> str:
        return "redis_latency"

    def is_active(self) -> bool:
        with self._lock:
            return self._active

    def start(self, params: dict) -> None:
        with self._lock:
            if self._active:
                logger.info("Redis latency fault already active.")
                return

            config = load_config()
            delay_ms = params.get("duration", config.faults.redis.default_delay)

            self._fault_delay = delay_ms

            if self._redis_client:
                self._redis_client.start_fault(delay_ms)

            self._active = True

    def stop(self) -> None:
        with self._lock:
            if not self._active:
                return

            if self._redis_client:
                self._redis_client.stop_fault()

            self._active = False
            self._fault_delay = 0
            logger.info("Redis latency fault stopped.")
