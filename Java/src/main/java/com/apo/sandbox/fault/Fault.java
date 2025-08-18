package com.apo.sandbox.fault;

import java.util.Map;

public interface Fault {
    String getName();

    void start(Map<String, Object> params) throws Exception;

    void stop() throws Exception;

    boolean isActive();
}