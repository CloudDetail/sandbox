package com.apo.sandbox.dao;

import com.apo.sandbox.model.User;
import com.fasterxml.jackson.core.JsonProcessingException;

import java.util.List;

public interface IRedisClient {
    List<String> getUserIDs() throws JsonProcessingException;

    User getUser(String userId) throws JsonProcessingException;

    void setUser(User user) throws JsonProcessingException;

    void setUserIDs(List<String> userIds) throws JsonProcessingException;

    void startFault(int delay);

    void stopFault();
}