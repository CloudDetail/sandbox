package com.apo.sandbox.fault;

import com.apo.sandbox.config.AppProperties;
import com.apo.sandbox.dao.IRedisClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

import java.util.Map;
import java.util.concurrent.atomic.AtomicBoolean;

@Component
public class RedisLatencyFault implements Fault {
    private static final Logger log = LoggerFactory.getLogger(RedisLatencyFault.class);
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
            log.info("Redis fault is already active.");
            return;
        }
        int delay = (int) params.getOrDefault("duration", appProperties.getRedisFaultDefaultDelay());
        try {
            redisClient.startFault(delay);
            active.set(true);
            log.info("Redis latency fault started with delay: {}ms", delay);
        } catch (Exception e) {
            log.error("Failed to start Redis latency fault: {}", e.getMessage());
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
            log.info("Redis latency fault stopped.");
        } catch (Exception e) {
            log.error("Failed to stop Redis latency fault: {}", e.getMessage());
        }
    }

    @Override
    public boolean isActive() {
        return active.get();
    }
}