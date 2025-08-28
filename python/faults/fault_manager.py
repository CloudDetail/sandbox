# python/faults/fault_manager.py
import logging
from typing import Dict, Any, List
import threading

logger = logging.getLogger(__name__)

class ChaosFault:
    def start(self, params: Dict[str, Any]) -> None:
        raise NotImplementedError

    def stop(self) -> None:
        raise NotImplementedError

    def is_active(self) -> bool:
        raise NotImplementedError

    def name(self) -> str:
        raise NotImplementedError

class Manager:
    def __init__(self):
        self._lock = threading.Lock()
        self.faults: Dict[str, ChaosFault] = {}

    def register(self, fault: ChaosFault) -> None:
        with self._lock:
            self.faults[fault.name()] = fault

    def start_fault(self, chaos_type: str, params: Dict[str, Any]) -> None:
        with self._lock:
            f = self.faults.get(chaos_type)
        if not f:
            return
        try:
            # Start fault asynchronously - don't wait for completion
            f.start(params)
        except Exception as e:
            return

    def stop_fault(self, chaos_type: str) -> None:
        with self._lock:
            f = self.faults.get(chaos_type)
        if not f:
            return
        try:
            f.stop()
        except Exception as e:
            return

    def stop_all_faults(self) -> None:
        with self._lock:
            faults_to_stop = list(self.faults.values())
        for f in faults_to_stop:
            if f.is_active():
                self.stop_fault(f.name())

    def status(self) -> Dict[str, Any]:
        with self._lock:
            status_data = {
                name: {"active": f.is_active(), "name": f.name()}
                for name, f in self.faults.items()
            }
        return status_data

    def list_active(self) -> List[str]:
        with self._lock:
            active = [name for name, f in self.faults.items() if f.is_active()]
        return active
