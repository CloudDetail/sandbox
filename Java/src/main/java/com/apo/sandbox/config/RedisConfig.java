package com.apo.sandbox.config;

import com.apo.sandbox.dao.IRedisClient;
import com.apo.sandbox.dao.MockRedisClient;
import com.apo.sandbox.dao.RedisClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.DependsOn;

import redis.clients.jedis.JedisPool;
import redis.clients.jedis.JedisPoolConfig;
import redis.clients.jedis.exceptions.JedisConnectionException;

import java.time.Duration;

@Configuration
@DependsOn("toxiProxy")
public class RedisConfig {

    private static final Logger log = LoggerFactory.getLogger(RedisConfig.class);

    @Bean
    public IRedisClient redisClient(AppProperties props) {
        String proxyHost = props.getProxyListenAddr().split(":")[0];
        int proxyPort = Integer.parseInt(props.getProxyListenAddr().split(":")[1]);
        try {
            final JedisPoolConfig poolConfig = new JedisPoolConfig();
            poolConfig.setMaxTotal(10);
            poolConfig.setBlockWhenExhausted(true);
            poolConfig.setMaxWait(Duration.ofMillis(3000));

            JedisPool jedisPool;
            String password = props.getRedisPassword();
            if (password != null && !password.isEmpty()) {
                jedisPool = new JedisPool(poolConfig, proxyHost, proxyPort, 2000, password);
            } else {
                jedisPool = new JedisPool(poolConfig, proxyHost, proxyPort, 2000);
            }

            // Test connection
            jedisPool.getResource().close();
            log.info("Successfully connected to Redis at {}:{}.", proxyHost, proxyPort);
            return new RedisClient(jedisPool);
        } catch (JedisConnectionException e) {
            log.error("Could not connect to Redis at {}:{}. Using mock client. Error: {}",
                    proxyHost, proxyPort, e.getMessage());
            return new MockRedisClient();
        }
    }
}