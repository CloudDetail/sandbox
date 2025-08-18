package com.apo.sandbox.fault;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Function;
import java.util.stream.Collectors;

@Service
public class FaultManager {
    private static final Logger log = LoggerFactory.getLogger(FaultManager.class);
    private final Map<String, Fault> faults;

    public FaultManager(List<Fault> faultList) {
        this.faults = faultList.stream().collect(
                Collectors.toConcurrentMap(Fault::getName, Function.identity()));
    }

    public void startFault(String faultType, Map<String, Object> params) throws Exception {
        Fault fault = faults.get(faultType);
        if (fault != null) {
            fault.start(params);
        } else {
            log.error("Unknown fault type: {}", faultType);
        }
    }

    public void stopAllFaults() {
        log.info("Stopping all active faults...");
        faults.values().forEach(fault -> {
            if (fault.isActive()) {
                try {
                    fault.stop();
                } catch (Exception e) {
                    log.error("Failed to stop fault '{}': {}", fault.getName(), e.getMessage());
                }
            }
        });
    }
}