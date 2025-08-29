package com.apo.sandbox.dao;

import com.apo.sandbox.model.User;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Repository;

import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

@Repository
public class Store {
    private static final Logger log = LoggerFactory.getLogger(Store.class);
    private final IRedisClient redisClient;
    private final IDatabaseClient dbClient;

    public Store(IRedisClient redisClient, IDatabaseClient dbClient) {
        this.redisClient = redisClient;
        this.dbClient = dbClient;
    }

    public List<User> queryUserFromMySQL() throws Exception {
        // Check if database is connected
        if (!dbClient.isConnected()) {
            log.info("Database is not connected. Returning mocked users.");
            return mockUsers("DB_Mock_", 10);
        }

        // Try to get users from database
        List<User> users = dbClient.getUsers();
        if (users != null && !users.isEmpty()) {
            log.info("Successfully fetched {} users from database.", users.size());
            return users;
        }

        // If database has no data, mock users and save to database
        log.info("Database has no data. Mocking users and saving to database.");
        users = mockUsers("DB_Saved_", 10);

        try {
            dbClient.saveUsers(users);
            log.info("Successfully saved 10 mocked users to database.");
        } catch (Exception e) {
            log.error("Failed to save mocked users to database: {}", e.getMessage());
            // Continue even if save fails
        }

        return users;
    }

    // Helper method to mock users
    private List<User> mockUsers(String prefix, int count) {
        List<User> users = new ArrayList<>();
        for (int i = 0; i < count; i++) {
            users.add(new User(
                    UUID.randomUUID().toString(),
                    String.format("%sUser %d", prefix, i + 1),
                    String.format("%suser%d@apo.com", prefix.toLowerCase(), i + 1)));
        }
        return users;
    }

    public List<User> queryUserFromRedis() throws Exception {
        // Try to get from Redis cache
        List<String> userIDs = redisClient.getUserIDs();
        if (userIDs != null && !userIDs.isEmpty()) {
            List<User> users = new ArrayList<>();
            for (String userId : userIDs) {
                User user = redisClient.getUser(userId);
                if (user != null) {
                    users.add(user);
                } else {
                    log.warn("Failed to get user {} from Redis.", userId);
                }
            }
            if (users.size() == userIDs.size()) {
                log.info("All users retrieved from Redis cache by individual IDs.");
                return users;
            }
            log.warn("Incomplete users retrieved from Redis cache. Re-fetching and caching.");
        }

        // If not in Redis or incomplete, mock 10 users and cache them
        log.info("Mocking 10 users and caching in Redis.");
        List<User> users = new ArrayList<>();
        List<String> newUserIDs = new ArrayList<>();
        for (int i = 0; i < 10; i++) {
            User user = new User(
                    UUID.randomUUID().toString(),
                    String.format("Mock User %d", i + 1),
                    String.format("mock%d@apo.com", i + 1));
            users.add(user);
            newUserIDs.add(user.getId());
            try {
                redisClient.setUser(user);
            } catch (Exception e) {
                log.warn("Failed to cache user {} in Redis: {}", user.getId(), e.getMessage());
            }
        }

        try {
            redisClient.setUserIDs(newUserIDs);
        } catch (Exception e) {
            log.warn("Failed to cache all user IDs in Redis: {}", e.getMessage());
        }

        log.info("Mocked 10 users and cached in Redis (individual users and IDs).");
        return users;
    }
}