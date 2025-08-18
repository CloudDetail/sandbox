package com.apo.sandbox.config;

import com.apo.sandbox.dao.DatabaseClient;
import com.apo.sandbox.dao.IDatabaseClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class DatabaseConfig {
    private static final Logger log = LoggerFactory.getLogger(DatabaseConfig.class);

    private final DatabaseProperties databaseProperties;

    @Autowired
    public DatabaseConfig(DatabaseProperties databaseProperties) {
        this.databaseProperties = databaseProperties;
    }

    @Bean
    public IDatabaseClient databaseClient() {
        log.info("Creating database client with properties: {}", databaseProperties);

        // Create and return DatabaseClient instance
        return new DatabaseClient(
                databaseProperties.getHost(),
                databaseProperties.getPort(),
                databaseProperties.getUsername(),
                databaseProperties.getPassword(),
                databaseProperties.getDatabase(),
                databaseProperties.getMaxConnections(),
                databaseProperties.getConnTimeout(),
                databaseProperties.getReadTimeout(),
                databaseProperties.getWriteTimeout()
        );
    }
}