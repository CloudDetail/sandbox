package com.apo.sandbox.service;

import com.apo.sandbox.dao.Store;
import com.apo.sandbox.fault.FaultManager;
import com.apo.sandbox.model.User;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;

@Service
public class BusinessService {
    private static final Logger log = LoggerFactory.getLogger(BusinessService.class);
    private final FaultManager faultManager;
    private final Store store;

    public BusinessService(FaultManager faultManager, Store store) {
        this.faultManager = faultManager;
        this.store = store;
    }

    public List<User> getUsersCached(Optional<String> chaosType, Optional<Integer> duration) {
        if (chaosType.isPresent() && !chaosType.get().isEmpty()) {
            Map<String, Object> params = new HashMap<>();

            if (duration.isPresent() && duration.get() > 0) {
                params.put("duration", duration.get());
            }
            try {
                faultManager.startFault(chaosType.get(), params);
            } catch (Exception e) {
            }
        } else {
            faultManager.stopAllFaults();
        }

        try {
            store.queryUsersCached();
            return store.queryUsersFromDatabase();
        } catch (Exception e) {
            log.error("Failed to get users: {}", e.getMessage());
            return Collections.emptyList();
        }
    }
}