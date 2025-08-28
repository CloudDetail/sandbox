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
    private final Map<String, Fault> faults;

    public FaultManager(List<Fault> faultList) {
        this.faults = faultList.stream().collect(
                Collectors.toConcurrentMap(Fault::getName, Function.identity()));
    }

    public void startFault(String faultType, Map<String, Object> params) throws Exception {
        Fault fault = faults.get(faultType);
        if (fault != null) {
            fault.start(params);
        }
    }

    public void stopAllFaults() {
        faults.values().forEach(fault -> {
            if (fault.isActive()) {
                try {
                    fault.stop();
                } catch (Exception e) {

                }
            }
        });
    }
}