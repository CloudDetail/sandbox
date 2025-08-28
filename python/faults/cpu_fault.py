import time
import threading

from faults.fault_manager import ChaosFault
from config import load_config

class CPUFault(ChaosFault):
    """
    A class to simulate a CPU-intensive fault.
    """

    def __init__(self):
        """
        Initializes the CPUFault class.
        """
        self._name = "cpu"
        self._active = False
        self._lock = threading.Lock()
        self._stop_event = threading.Event()

    def name(self):
        """
        Returns the name of the fault.
        """
        return self._name

    def start(self, params=None):
        """
        Starts the CPU fault for a specified duration.
        :param params: A dictionary containing parameters for the fault.
                       Expected key: "duration" (int, milliseconds).
        """
        if params is None:
            params = {}

        with self._lock:
            if self._active:
                return
            self._active = True
            self._stop_event.clear()

        # Default duration if not provided in params.
        config = load_config()
        target_duration = params.get("duration", config.faults.cpu.default_duration)
        target_duration_s = target_duration / 1000.0

        start_time = time.time()

        iteration_count = 0
        while (time.time() - start_time) < target_duration_s and not self._stop_event.is_set():
            self._cpu_intensive_work()
            iteration_count += 1

        end_time = time.time()
        elapsed_time = end_time - start_time

        with self._lock:
            self._active = False

        return None

    def _cpu_intensive_work(self):
        import math
        result = 0
        for i in range(100000):
            x = i + 1
            sqrt_x = math.sqrt(x*math.pi)
            sin_x = math.sin(x)
            cos_x = math.cos(x)
            result += sqrt_x * sin_x + cos_x * math.log(x + 1)
        return result

    def _fibonacci(self, n):
        """
        A recursive function to calculate the nth Fibonacci number.
        This is a CPU-intensive operation.
        """
        if n <= 1:
            return n
        return self._fibonacci(n - 1) + self._fibonacci(n - 2)

    def stop(self):
        """
        Stops the CPU fault.
        """
        with self._lock:
            if not self._active:
                return
            self._stop_event.set()
            self._active = False
        return None

    def is_active(self):
        """
        Checks if the fault is currently active.
        """
        with self._lock:
            return self._active