package com.apo.sandbox.service;

import com.apo.sandbox.dao.Store;
import com.apo.sandbox.fault.FaultManager;
import com.apo.sandbox.model.User;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import eu.rekawek.toxiproxy.Proxy;
import eu.rekawek.toxiproxy.model.ToxicDirection;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.atomic.AtomicBoolean;

@Service
public class BusinessService {
    private static final Logger log = LoggerFactory.getLogger(BusinessService.class);
    private final Store store;
    private final AtomicBoolean latencyActive = new AtomicBoolean(false);
    private final AtomicBoolean redisLatencyActive = new AtomicBoolean(false);

    public final Proxy toxiProxy;

    public BusinessService(FaultManager faultManager, Store store, Proxy toxiProxy) {
        this.store = store;
        this.toxiProxy = toxiProxy;
    }

    public List<User> getUsersWithLatency(Optional<String> mode, int duration) {
        if ("1".equals(mode.orElse(""))) {
            if (!latencyActive.get()) {
                try {
                    // Clear any existing traffic control rules to ensure a clean state
                    clearTc();

                    // Inject network latency using Linux tc (traffic control) to simulate network
                    // delays
                    // This adds a delay to all packets on the eth0 interface
                    String command = String.format("tc qdisc add dev %s root netem delay %dms", "eth0", duration);
                    executeCommand(command.split(" "));

                    latencyActive.set(true);
                } catch (Exception e) {
                    log.error("type 1 failed");
                }
            }
        } else {
            stopFaults();
        }

        try {
            store.queryUserFromRedis();
            return store.queryUserFromMySQL();
        } catch (Exception e) {
            log.error("Failed to get users: {}", e.getMessage());
            return Collections.emptyList();
        }
    }

    public List<User> getUsersWithCPUBurn(Optional<String> mode, int duration) {
        if ("1".equals(mode.orElse(""))) {
            // Calculate target duration in nanoseconds for precise CPU burning
            long targetDurationNanos = duration * 1_000_000L;
            long startTime = System.nanoTime();

            // Burn CPU cycles by performing intensive mathematical calculations
            // This creates artificial CPU load without side effects
            while (System.nanoTime() - startTime < targetDurationNanos) {
                // Compute Fibonacci sequence recursively to maximize CPU usage
                fibonacci(18);
            }
        } else {
            stopFaults();
        }

        try {
            store.queryUserFromRedis();
            return store.queryUserFromMySQL();
        } catch (Exception e) {
            log.error("Failed to get users: {}", e.getMessage());
            return Collections.emptyList();
        }
    }

    public List<User> getUsersWithRedisLatency(Optional<String> mode, int duration) {
        if ("1".equals(mode.orElse(""))) {
            if (!redisLatencyActive.get()) {
                try {
                    // Use Toxiproxy to simulate Redis latency
                    // This simulates slow Redis responses without affecting actual Redis server
                    toxiProxy.toxics().latency("redis_latency", ToxicDirection.DOWNSTREAM, duration);

                    redisLatencyActive.set(true);
                } catch (Exception e) {
                    log.error("type 3 failed");
                }
            }
        } else {
            stopFaults();
        }

        try {
            store.queryUserFromRedis();
            return store.queryUserFromMySQL();
        } catch (Exception e) {
            log.error("Failed to get users: {}", e.getMessage());
            return Collections.emptyList();
        }
    }

    private void clearTc() throws Exception {
        String command = String.format("tc qdisc del dev %s root", "eth0");
        try {
            executeCommand(command.split(" "));
        } catch (Exception e) {
            if (e.getMessage().contains("No such file or directory") ||
                    e.getMessage().contains("No qdisc") ||
                    e.getMessage().contains("Cannot delete qdisc with handle of zero")) {
                log.info("No existing tc rule to clear, which is fine.");
            } else {
                throw e;
            }
        }
    }

    private void executeCommand(String... command) throws Exception {
        ProcessBuilder pb = new ProcessBuilder(command);
        Process process = pb.start();
        int exitCode = process.waitFor();

        if (exitCode != 0) {
            BufferedReader reader = new BufferedReader(new InputStreamReader(process.getErrorStream()));
            StringBuilder output = new StringBuilder();
            String line;
            while ((line = reader.readLine()) != null) {
                output.append(line).append("\n");
            }
            throw new RuntimeException(
                    String.format("Command failed with exit code %d: %s", exitCode, output.toString()));
        }
    }

    private int fibonacci(int n) {
        if (n <= 1) {
            return n;
        }
        return fibonacci(n - 1) + fibonacci(n - 2);
    }

    private void stopFaults() {
        if (latencyActive.get()) {
            try {
                clearTc();
                latencyActive.set(false);
            } catch (Exception e) {
                log.error("Failed to clear tc rules: {}", e.getMessage());
            }
        }
        if (redisLatencyActive.get()) {
            try {
                var redisLatencyToxic = toxiProxy.toxics().get("redis_latency");

                redisLatencyToxic.remove();

                redisLatencyActive.set(false);
            } catch (IOException e) {
                redisLatencyActive.set(false);
            } catch (Exception e) {
                log.info("stop failed");
            }
        }
    }
}