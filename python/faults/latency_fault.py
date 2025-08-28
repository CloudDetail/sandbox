import time
import threading
import subprocess
import os

from faults.fault_manager import ChaosFault
from config import load_config

class LatencyFault(ChaosFault):
    def __init__(self, iface: str = "eth0"):
        self._active = False
        self._lock = threading.Lock()
        self._iface = iface
        self._delay_ms = 0

    def name(self) -> str:
        return "latency"

    def is_active(self) -> bool:
        with self._lock:
            return self._active

    def start(self, params: dict) -> None:
        with self._lock:
            if self._active:
                return

            config = load_config()
            delay_ms = params.get("duration", config.faults.latency.default_delay)
            if delay_ms < 1:
                delay_ms = 100

            try:
                # First, safely clear any existing tc rules
                self._clear_tc()
                
                # Add new delay rule
                cmd = ["tc", "qdisc", "add", "dev", self._iface, "root", "netem", "delay", f"{delay_ms}ms"]
                
                result = subprocess.run(cmd, check=True, capture_output=True, text=True)
                
            except subprocess.CalledProcessError as e:
                error_msg = e.stderr.decode().strip() if e.stderr else str(e)
                raise RuntimeError(f"Failed to add delay on {self._iface}: {error_msg}")
            except Exception as e:
                raise RuntimeError(f"Unexpected error while adding delay on {self._iface}: {e}")

            self._delay_ms = delay_ms
            self._active = True

    def stop(self) -> None:
        with self._lock:
            if not self._active:
                return
            try:
                # Safely clear tc rules
                self._clear_tc()
            except Exception as e:
                return
            finally:
                self._active = False
                self._delay_ms = 0

    def _clear_tc(self) -> None:
        """
        Safely clear existing tc qdisc rules.
        Handles various error cases gracefully.
        """
        cmd = ["tc", "qdisc", "del", "dev", self._iface, "root"]
        try:
            subprocess.run(cmd, check=True, capture_output=True)
        except subprocess.CalledProcessError as e:
            return


