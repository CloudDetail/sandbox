package com.apo.sandbox.fault;

import com.apo.sandbox.config.AppProperties;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

import java.util.Map;

@Component
public class CpuFault implements Fault {
    private final AppProperties appProperties;

    public CpuFault(AppProperties appProperties) {
        this.appProperties = appProperties;
    }

    @Override
    public String getName() {
        return "cpu";
    }

    @Override
    public void start(Map<String, Object> params) {
        int durationMs = (int) params.getOrDefault("duration", appProperties.getCpuFaultDefaultDuration());
        long targetDurationNanos = durationMs * 1_000_000L;

        long startTime = System.nanoTime();

        while (System.nanoTime() - startTime < targetDurationNanos) {
            fibonacci(18); // A reasonably expensive computation
        }
    }

    private int fibonacci(int n) {
        if (n <= 1) {
            return n;
        }
        return fibonacci(n - 1) + fibonacci(n - 2);
    }

    @Override
    public void stop() {
        // This is a one-shot fault, no stop action needed.
    }

    @Override
    public boolean isActive() {
        // As per the Go logic, this fault is instantaneous and doesn't have a
        // persistent active state.
        return false;
    }
}