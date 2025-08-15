package storage

import (
	"database/sql"
	"fmt"

	"github.com/CloudDetail/apo-sandbox/logging"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLClient struct {
	DB *sql.DB
}

func NewMySQL(dsn string) (*MySQLClient, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	client := &MySQLClient{DB: db}

	// Initialize schema
	err = client.InitSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MySQL schema: %w", err)
	}

	return client, nil
}

func (c *MySQLClient) InitSchema() error {
	createUserTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) NOT NULL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE
	);
	`
	_, err := c.DB.Exec(createUserTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	logging.Info("%s", "Users table checked/created successfully.")

	createOrderTableSQL := `
	CREATE TABLE IF NOT EXISTS orders (
		id VARCHAR(36) NOT NULL PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL,
		item VARCHAR(255) NOT NULL,
		amount DECIMAL(10, 2) NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`
	_, err = c.DB.Exec(createOrderTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create orders table: %w", err)
	}
	logging.Info("%s", "Orders table checked/created successfully.")

	return nil
}

func (c *MySQLClient) QueryRow(query string, args ...interface{}) *sql.Row {
	if c.DB == nil {
		logging.Warn("%s", "MySQL not implement for QueryRow")
		return nil // Or handle error appropriately
	}
	return c.DB.QueryRow(query, args...)
}

func (c *MySQLClient) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if c.DB == nil {
		logging.Warn("%s", "MySQL not implement")
		return nil, fmt.Errorf("%s", "MySQL not implement")
	}
	return c.DB.Query(query, args...)
}

func (c *MySQLClient) Exec(query string, args ...interface{}) (sql.Result, error) {
	if c.DB == nil {
		logging.Warn("%s", "MySQL not implement")
		return nil, fmt.Errorf("%s", "MySQL not implement")
	}
	return c.DB.Exec(query, args...)
}
