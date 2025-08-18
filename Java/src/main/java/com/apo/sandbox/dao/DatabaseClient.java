package com.apo.sandbox.dao;

import com.apo.sandbox.model.User;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import java.sql.*;
import java.time.Duration;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

public class DatabaseClient implements IDatabaseClient {
    private static final Logger log = LoggerFactory.getLogger(DatabaseClient.class);
    private final String host;
    private final int port;
    private final String username;
    private final String password;
    private final String database;
    private final int maxConnections;
    private final Duration connTimeout;
    private final Duration readTimeout;
    private final Duration writeTimeout;
    private Connection connection;
    private boolean initialized = false;

    public DatabaseClient(String host, int port, String username, String password, String database,
                         int maxConnections, Duration connTimeout, Duration readTimeout, Duration writeTimeout) {
        this.host = host;
        this.port = port;
        this.username = username;
        this.password = password;
        this.database = database;
        this.maxConnections = maxConnections;
        this.connTimeout = connTimeout;
        this.readTimeout = readTimeout;
        this.writeTimeout = writeTimeout;

        // Initialize connection
        initialize();
    }

    private void initialize() {
        try {
            // Register JDBC driver
            Class.forName("com.mysql.cj.jdbc.Driver");

            // Create connection URL
            String url = String.format("jdbc:mysql://%s:%d/%s?serverTimezone=UTC&connectTimeout=%d&socketTimeout=%d",
                    host, port, database, (int)connTimeout.toMillis(), (int)readTimeout.toMillis());

            // Establish connection
            log.info("Attempting to connect to database at {}", url);
            connection = DriverManager.getConnection(url, username, password);
            initialized = true;
            log.info("Successfully connected to database");

            // Create users table if it doesn't exist
            createUsersTable();
        } catch (ClassNotFoundException e) {
            log.error("MySQL JDBC driver not found: {}", e.getMessage());
        } catch (SQLException e) {
            log.error("Failed to connect to database: {}", e.getMessage());
            // Don't set initialized to true
        }
    }

    private void createUsersTable() {
        String createTableSQL = "CREATE TABLE IF NOT EXISTS users (" +
                "id VARCHAR(36) PRIMARY KEY, " +
                "name VARCHAR(100) NOT NULL, " +
                "email VARCHAR(100) NOT NULL UNIQUE" +
                ")";

        try (Statement stmt = connection.createStatement()) {
            stmt.execute(createTableSQL);
            log.info("Users table created or already exists");
        } catch (SQLException e) {
            log.error("Failed to create users table: {}", e.getMessage());
        }
    }

    @Override
    public boolean isConnected() {
        if (!initialized) {
            return false;
        }

        try {
            if (connection == null || connection.isClosed()) {
                log.warn("Database connection is closed. Reconnecting...");
                initialize();
            }
            return !connection.isClosed();
        } catch (SQLException e) {
            log.error("Error checking connection status: {}", e.getMessage());
            return false;
        }
    }

    @Override
    public List<User> getUsers() {
        if (!isConnected()) {
            log.warn("Cannot get users: database not connected");
            return null;
        }

        List<User> users = new ArrayList<>();
        String query = "SELECT id, name, email FROM users";

        try (Statement stmt = connection.createStatement();
             ResultSet rs = stmt.executeQuery(query)) {

            while (rs.next()) {
                User user = new User(
                        rs.getString("id"),
                        rs.getString("name"),
                        rs.getString("email")
                );
                users.add(user);
            }

            log.info("Retrieved {} users from database", users.size());
            return users;
        } catch (SQLException e) {
            log.error("Failed to retrieve users from database: {}", e.getMessage());
            return null;
        }
    }

    @Override
    public void saveUsers(List<User> users) {
        if (!isConnected()) {
            log.warn("Cannot save users: database not connected");
            return;
        }

        String insertSQL = "INSERT INTO users (id, name, email) VALUES (?, ?, ?)";

        try (PreparedStatement pstmt = connection.prepareStatement(insertSQL)) {
            connection.setAutoCommit(false);

            for (User user : users) {
                pstmt.setString(1, user.getId());
                pstmt.setString(2, user.getName());
                pstmt.setString(3, user.getEmail());
                pstmt.addBatch();
            }

            int[] result = pstmt.executeBatch();
            connection.commit();
            log.info("Successfully saved {} users to database", result.length);
        } catch (SQLException e) {
            try {
                if (connection != null) {
                    connection.rollback();
                }
            } catch (SQLException rollbackEx) {
                log.error("Failed to rollback transaction: {}", rollbackEx.getMessage());
            }
            log.error("Failed to save users to database: {}", e.getMessage());
        } finally {
            try {
                if (connection != null) {
                    connection.setAutoCommit(true);
                }
            } catch (SQLException autoCommitEx) {
                log.error("Failed to set auto-commit: {}", autoCommitEx.getMessage());
            }
        }
    }

    // Close connection when not needed
    public void close() {
        if (connection != null) {
            try {
                connection.close();
                log.info("Database connection closed");
            } catch (SQLException e) {
                log.error("Failed to close database connection: {}", e.getMessage());
            }
        }
    }
}