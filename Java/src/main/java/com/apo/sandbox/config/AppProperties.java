package com.apo.sandbox.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

@Component
public class AppProperties {

    // --- Redis Configuration ---
    @Value("${REDIS_HOST:localhost}")
    private String redisHost;

    @Value("${REDIS_PORT:6379}")
    private int redisPort;

    @Value("${REDIS_PASSWORD:}")
    private String redisPassword;

    // --- Faults Configuration ---
    @Value("${CPU_FAULT_DEFAULT_DURATION:200}")
    private int cpuFaultDefaultDuration;

    @Value("${LATENCY_FAULT_DEFAULT_DELAY:200}")
    private int latencyFaultDefaultDelay;

    @Value("${NETWORK_FAULT_INTERFACE:eth0}")
    private String networkFaultInterface;

    @Value("${REDIS_FAULT_DEFAULT_DELAY:100}")
    private int redisFaultDefaultDelay;

    // --- Getters ---
    public String getRedisHost() {
        return redisHost;
    }

    public int getRedisPort() {
        return redisPort;
    }

    public String getRedisPassword() {
        return redisPassword;
    }

    public int getCpuFaultDefaultDuration() {
        return cpuFaultDefaultDuration;
    }

    public int getLatencyFaultDefaultDelay() {
        return latencyFaultDefaultDelay;
    }

    public String getNetworkFaultInterface() {
        return networkFaultInterface;
    }

    public int getRedisFaultDefaultDelay() {
        return redisFaultDefaultDelay;
    }
}