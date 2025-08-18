import time
import logging
import threading
import subprocess
import os

from faults.fault_manager import ChaosFault
from config import load_config

logger = logging.getLogger(__name__)

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
                logger.info("Latency fault already active.")
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
                logger.debug(f"Executing tc command: {' '.join(cmd)}")
                
                result = subprocess.run(cmd, check=True, capture_output=True, text=True)
                logger.debug(f"tc command output: {result.stdout}")
                
            except subprocess.CalledProcessError as e:
                error_msg = e.stderr.decode().strip() if e.stderr else str(e)
                logger.error(f"Failed to add delay on {self._iface}: {error_msg}")
                if e.stdout:
                    logger.debug(f"Command stdout: {e.stdout.decode().strip()}")
                raise RuntimeError(f"Failed to add delay on {self._iface}: {error_msg}")
            except Exception as e:
                logger.error(f"Unexpected error while adding delay on {self._iface}: {e}")
                raise RuntimeError(f"Unexpected error while adding delay on {self._iface}: {e}")

            self._delay_ms = delay_ms
            self._active = True
            logger.info(f"Successfully set simulated {self._delay_ms}ms delay on {self._iface}.")

    def stop(self) -> None:
        with self._lock:
            if not self._active:
                logger.debug(f"Latency fault on {self._iface} not active, nothing to stop.")
                return
            try:
                # Safely clear tc rules
                self._clear_tc()
            except Exception as e:
                logger.error(f"Failed to clear tc on {self._iface}: {e}")
                # Even if clearing fails, mark as inactive
                # Don't raise exception to avoid blocking stop operation
            finally:
                self._active = False
                self._delay_ms = 0
                logger.info(f"Latency fault stopped on {self._iface}.")

    def _clear_tc(self) -> None:
        """
        Safely clear existing tc qdisc rules.
        Handles various error cases gracefully.
        """
        cmd = ["tc", "qdisc", "del", "dev", self._iface, "root"]
        try:
            subprocess.run(cmd, check=True, capture_output=True)
            logger.debug(f"Successfully cleared tc qdisc for {self._iface}")
        except subprocess.CalledProcessError as e:
            error_msg = e.stderr.decode().strip()
            # Handle various error cases that indicate no qdisc exists
            if any(msg in error_msg.lower() for msg in [
                "no such file or directory",
                "no qdisc",
                "cannot delete qdisc with handle of zero",
                "no qdisc with this handle",
                "invalid qdisc handle"
            ]):
                logger.debug(f"No existing tc qdisc to delete on {self._iface}: {error_msg}")
            else:
                # For other unexpected errors, log warning but don't fail
                logger.warning(f"Unexpected error while clearing tc qdisc on {self._iface}: {error_msg}")
                # Don't raise the exception, just log it

