package com.apo.sandbox.fault;

import com.apo.sandbox.config.AppProperties;
import com.apo.sandbox.dao.IRedisClient;
import org.springframework.stereotype.Component;

import java.util.Map;
import java.util.concurrent.atomic.AtomicBoolean;

@Component
public class RedisLatencyFault implements Fault {
    private final IRedisClient redisClient;
    private final AppProperties appProperties;
    private final AtomicBoolean active = new AtomicBoolean(false);

    public RedisLatencyFault(IRedisClient redisClient, AppProperties appProperties) {
        this.redisClient = redisClient;
        this.appProperties = appProperties;
    }

    @Override
    public String getName() {
        // The Go code seems to use "redis_latency" as the name, but the config key is
        // REDIS_FAULT
        // Let's align the fault name with the Go version for clarity in the API
        return "redis_latency";
    }

    @Override
    public synchronized void start(Map<String, Object> params) {
        if (active.get()) {
            return;
        }
        int delay = (int) params.getOrDefault("duration", appProperties.getRedisFaultDefaultDelay());
        try {
            redisClient.startFault(delay);
            active.set(true);
        } catch (Exception e) {
        }
    }

    // ... (rest of the file is unchanged, including stop() and isActive())
    @Override
    public synchronized void stop() {
        if (!active.get()) {
            return;
        }
        try {
            redisClient.stopFault();
            active.set(false);
        } catch (Exception e) {
        }
    }

    @Override
    public boolean isActive() {
        return active.get();
    }
}