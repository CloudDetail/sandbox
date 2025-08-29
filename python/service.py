# python/service.py
import logging
import json
import math
import requests
import subprocess
import threading
import time

from storage import Store

logger = logging.getLogger(__name__)

class BusinessService:
    """
    BusinessService class provides methods for fault injection testing
    including network latency, CPU burn, and Redis latency simulation
    """
    def __init__(self, store: Store, iface: str = "eth0"):
        """
        Initialize BusinessService with store and network interface

        Args:
            store (Store): Data store object
            iface (str): Network interface for traffic control (default: 'eth0')
        """
        self._active = False
        self._lock = threading.Lock()
        self._iface = iface
        self.store = store

        proxy_json = {
            "name": "basicProxy",
            "listen": "localhost:6379",
            "upstream": "redis-service:6379"
        }
        # Ignore error if proxy doesn't exist
        try:
            requests.post("http://localhost:8474/proxies", json=proxy_json)
        except requests.RequestException:
            pass  # Ignore errors during proxy creation

    def get_users_latency(self, mode: int = 0) -> tuple[str, None]:
        """
        Simulate network latency by adding traffic control rules

        Args:
            mode (int): 1 to enable latency, 0 to disable

        Returns:
            tuple[str, None]: JSON stringified users data and error (if any)
        """
        if mode == 1:
            # Use the `_active` variable to record the fault injection status
            # keep it enabled during the injection process.
            if not self._active:
                try:
                    self._clear_tc()
                    cmd = ["tc", "qdisc", "add", "dev", self._iface, "root", "netem", "delay", "200ms"]
                    result = subprocess.run(cmd, check=True, capture_output=True, text=True)
                except subprocess.CalledProcessError as e:
                    error_msg = e.stderr.decode() if e.stderr else str(e)
                    raise RuntimeError(f"Failed to add delay on {self._iface}: {error_msg}")
                except Exception as e:
                    raise RuntimeError(f"Unexpected error while adding delay on {self._iface}: {e}")
                self._active = True
        else:
            if self._active:
                self._clear_tc()
                self._active = False

        self.store.query_users_from_redis()
        users, err = self.store.query_or_create_users_from_mysql()
        if err:
            return "", err

        return json.dumps([user.__dict__ for user in users]), None

    def _clear_tc(self) -> None:
        """
        Safely clear existing tc qdisc rules.
        Handles various error cases gracefully.
        """
        cmd = ["tc", "qdisc", "del", "dev", self._iface, "root"]
        try:
            subprocess.run(cmd, check=True, capture_output=True)
        except subprocess.CalledProcessError as e:
            error_msg = e.stderr.decode().strip()
            logger.debug(f"No existing tc qdisc to delete on {self._iface}: {error_msg}")

    def get_users_cpu_burn(self, mode: int = 0) -> tuple[str, None]:
        """
        Simulate high CPU usage by performing intensive mathematical calculations

        Args:
            mode (int): 1 to enable CPU burn, 0 to disable

        Returns:
            tuple[str, None]: JSON stringified users data and error (if any)
        """
        if mode == 1:
            start_time = time.time()
            while (time.time() - start_time) < 0.2:  # 200ms
                result = 0
                for i in range(100000):
                    x = i + 1
                    sqrt_x = math.sqrt(x * math.pi)
                    sin_x = math.sin(x)
                    cos_x = math.cos(x)
                    result += sqrt_x * sin_x + cos_x * math.log(x + 1)

        self.store.query_users_from_redis()
        users, err = self.store.query_or_create_users_from_mysql()
        if err:
            return "", err
        return json.dumps([user.__dict__ for user in users]), None

    def get_users_redis_latency(self, mode: int = 0) -> tuple[str, None]:
        """
        Simulate Redis latency using Toxiproxy

        Args:
            mode (int): 1 to enable Redis latency, 0 to disable

        Returns:
            tuple[str, None]: JSON stringified users data and error (if any)
        """
        if mode == 1:
            # Use the `_active` variable to record the fault injection status
            # keep it enabled during the injection process.
            if not self._active:
                # Toxiproxy is a framework for simulating network conditions.
                # https://github.com/shopify/toxiproxy
                try:
                    requests.post("http://localhost:8474/proxies/basicProxy/toxics", json={
                        "name": "latency",
                        "type": "latency",
                        "stream": "all",
                        "attributes": {"latency": 200}
                    })
                    self._active = True
                except requests.RequestException as e:
                    logger.warning(f"Failed to add Redis latency: {e}")
        else:
            if self._active:
                try:
                    requests.delete("http://localhost:8474/proxies/basicProxy/toxics/latency")
                    self._active = False
                except requests.RequestException as e:
                    logger.warning(f"Failed to remove Redis latency: {e}")

        self.store.query_users_from_redis()
        users, err = self.store.query_or_create_users_from_mysql()
        if err:
            return "", err
        return json.dumps([user.__dict__ for user in users]), None