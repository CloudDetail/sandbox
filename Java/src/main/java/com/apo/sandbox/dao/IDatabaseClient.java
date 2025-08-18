package com.apo.sandbox.dao;

import com.apo.sandbox.model.User;
import java.util.List;

public interface IDatabaseClient {
    boolean isConnected();
    List<User> getUsers();
    void saveUsers(List<User> users);
}