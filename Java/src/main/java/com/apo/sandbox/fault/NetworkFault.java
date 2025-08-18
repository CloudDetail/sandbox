package com.apo.sandbox.fault;

import com.apo.sandbox.config.AppProperties;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.util.Map;
import java.util.concurrent.atomic.AtomicBoolean;

@Component
public class NetworkFault implements Fault {
    private static final Logger log = LoggerFactory.getLogger(NetworkFault.class);
    private final AtomicBoolean active = new AtomicBoolean(false);
    private final String networkInterface;
    private final AppProperties appProperties;
    private int currentDelay = 0;

    public NetworkFault(AppProperties appProperties) {
        this.appProperties = appProperties;
        this.networkInterface = appProperties.getNetworkFaultInterface();
    }

    @Override
    public String getName() {
        return "latency";
    }

    @Override
    public synchronized void start(Map<String, Object> params) throws Exception {
        if (active.get()) {
            log.warn("Network fault is already active.");
            return;
        }
        int delayMs = (int) params.getOrDefault("duration", appProperties.getLatencyFaultDefaultDelay());
        if (delayMs < 1) {
            delayMs = 100;
        }

        clearTc(); // Clear any existing rules first

        String command = String.format("tc qdisc add dev %s root netem delay %dms", networkInterface, delayMs);
        log.info("Executing command: {}", command);

        executeCommand(command.split(" "));

        this.currentDelay = delayMs;
        active.set(true);
        log.info("Successfully added {}ms delay on interface {}.", delayMs, networkInterface);
    }

    // ... (rest of the file is unchanged, including stop(), clearTc(),
    // executeCommand(), isActive())
    @Override
    public synchronized void stop() throws Exception {
        if (!active.get()) {
            return;
        }
        clearTc();
        active.set(false);
        log.info("Successfully cleared network delay on interface {}.", networkInterface);
    }

    private void clearTc() throws Exception {
        log.info("Attempting to clear tc rules on interface {}.", networkInterface);
        String command = String.format("tc qdisc del dev %s root", networkInterface);
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
            throw new RuntimeException(String.format("Command failed with exit code %d: %s", exitCode, output));
        }
    }

    @Override
    public boolean isActive() {
        return active.get();
    }
}