package com.apo.sandbox.dao;

import com.apo.sandbox.model.User;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import redis.clients.jedis.Jedis;
import redis.clients.jedis.JedisPool;

import java.util.List;

public class RedisClient implements IRedisClient {
    private final JedisPool jedisPool;
    private final ObjectMapper objectMapper = new ObjectMapper();
    private static final String USER_IDS_KEY = "user_ids";

    public RedisClient(JedisPool jedisPool) {
        this.jedisPool = jedisPool;
    }

    private String userKey(String userId) {
        return "user:" + userId;
    }

    @Override
    public List<String> getUserIDs() throws JsonProcessingException {
        try (Jedis jedis = jedisPool.getResource()) {
            String json = jedis.get(USER_IDS_KEY);
            if (json == null || json.isEmpty()) {
                return null;
            }
            return objectMapper.readValue(json, new TypeReference<>() {
            });
        }
    }

    @Override
    public User getUser(String userId) throws JsonProcessingException {
        try (Jedis jedis = jedisPool.getResource()) {
            String json = jedis.get(userKey(userId));
            if (json == null || json.isEmpty()) {
                return null;
            }
            return objectMapper.readValue(json, User.class);
        }
    }

    @Override
    public void setUser(User user) throws JsonProcessingException {
        try (Jedis jedis = jedisPool.getResource()) {
            String json = objectMapper.writeValueAsString(user);
            jedis.set(userKey(user.getId()), json);
        }
    }

    @Override
    public void setUserIDs(List<String> userIds) throws JsonProcessingException {
        try (Jedis jedis = jedisPool.getResource()) {
            String json = objectMapper.writeValueAsString(userIds);
            jedis.set(USER_IDS_KEY, json);
        }
    }

    @Override
    public void startFault(int delay) {
        try (Jedis jedis = jedisPool.getResource()) {
            // This sends a custom command to the Redis proxy
            jedis.sendCommand(() -> "FAULT.START".getBytes(), String.valueOf(delay).getBytes());
        } catch (Exception e) {
            throw new RuntimeException("Failed to start Redis fault", e);
        }
    }

    @Override
    public void stopFault() {
        try (Jedis jedis = jedisPool.getResource()) {
            jedis.sendCommand(() -> "FAULT.STOP".getBytes());
        } catch (Exception e) {
            throw new RuntimeException("Failed to stop Redis fault", e);
        }
    }
}