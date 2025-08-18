package com.apo.sandbox;

import com.apo.sandbox.config.AppProperties;
import com.apo.sandbox.dao.IRedisClient;
import com.apo.sandbox.dao.MockRedisClient;
import com.apo.sandbox.dao.RedisClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import redis.clients.jedis.JedisPool;
import redis.clients.jedis.JedisPoolConfig;
import redis.clients.jedis.exceptions.JedisConnectionException;

import java.time.Duration;

@SpringBootApplication
public class SandboxApplication {

    private static final Logger log = LoggerFactory.getLogger(SandboxApplication.class);

    public static void main(String[] args) {
        SpringApplication.run(SandboxApplication.class, args);
    }

    @Bean
    public IRedisClient redisClient(AppProperties props) {
        try {
            final JedisPoolConfig poolConfig = new JedisPoolConfig();
            poolConfig.setMaxTotal(10);
            poolConfig.setBlockWhenExhausted(true);
            poolConfig.setMaxWait(Duration.ofMillis(3000));

            JedisPool jedisPool;
            String password = props.getRedisPassword();
            if (password != null && !password.isEmpty()) {
                jedisPool = new JedisPool(poolConfig, props.getRedisHost(), props.getRedisPort(), 2000, password);
            } else {
                jedisPool = new JedisPool(poolConfig, props.getRedisHost(), props.getRedisPort(), 2000);
            }

            // Test connection
            jedisPool.getResource().close();
            log.info("Successfully connected to Redis at {}:{}.", props.getRedisHost(), props.getRedisPort());
            return new RedisClient(jedisPool);
        } catch (JedisConnectionException e) {
            log.error("Could not connect to Redis at {}:{}. Using mock client. Error: {}",
                    props.getRedisHost(), props.getRedisPort(), e.getMessage());
            return new MockRedisClient();
        }
    }
}