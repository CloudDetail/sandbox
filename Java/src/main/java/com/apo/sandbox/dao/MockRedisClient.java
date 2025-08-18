package com.apo.sandbox.dao;

import com.apo.sandbox.model.User;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Collections;
import java.util.List;

public class MockRedisClient implements IRedisClient {
    private static final Logger log = LoggerFactory.getLogger(MockRedisClient.class);

    public MockRedisClient() {
        log.warn("Using MockRedisClient. All Redis operations will be simulated.");
    }

    @Override
    public List<String> getUserIDs() {
        log.info("MOCK: Getting user IDs.");
        return Collections.emptyList(); // Return empty to force fallback logic
    }

    @Override
    public User getUser(String userId) {
        log.info("MOCK: Getting user for ID: {}", userId);
        return null;
    }

    @Override
    public void setUser(User user) {
        log.info("MOCK: Setting user: {}", user.getId());
    }

    @Override
    public void setUserIDs(List<String> userIds) {
        log.info("MOCK: Setting user IDs list.");
    }

    @Override
    public void startFault(int delay) {
        log.warn("MOCK: Cannot start Redis fault. Redis is not connected.");
    }

    @Override
    public void stopFault() {
        log.warn("MOCK: Cannot stop Redis fault. Redis is not connected.");
    }
}