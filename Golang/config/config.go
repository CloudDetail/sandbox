package config

import (
	"os"
	"strconv"
	"time"

	"github.com/CloudDetail/apo-sandbox/logging"
	"github.com/joho/godotenv"
)

var config *Config

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Faults   FaultsConfig   `json:"faults"`
}

type ServerConfig struct {
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	DeployProxy  bool          `json:"deploy_proxy"`
}

type DatabaseConfig struct {
	MySQL MySQLConfig `json:"mysql"`
	Redis RedisConfig `json:"redis"`
	Proxy ProxyConfig `json:"proxy"`
}

type MySQLConfig struct {
	Host           string        `json:"host"`
	Port           int           `json:"port"`
	Username       string        `json:"username"`
	Password       string        `json:"password"`
	Database       string        `json:"database"`
	MaxConnections int           `json:"max_connections"`
	ConnTimeout    time.Duration `json:"connection_timeout"`
	ReadTimeout    time.Duration `json:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Database int    `json:"database"`
}

type ProxyConfig struct {
	Addr       string `json:"addr"`        // toxiproxy server address
	ListenAddr string `json:"listen_addr"` // toxiproxy redis proxy listen address
}

type FaultsConfig struct {
	CPU     CPUFaultConfig     `json:"cpu"`
	Latency LatencyFaultConfig `json:"latency"`
	Redis   RedisFaultConfig   `json:"redis"`
}

type CPUFaultConfig struct {
	DefaultDuration int `json:"default_duration"`
}

type LatencyFaultConfig struct {
	DefaultDelay int `json:"default_delay"`
}

type RedisFaultConfig struct {
	DefaultDelay int `json:"default_delay"`
}

func LoadConfig() *Config {
	if config != nil {
		return config
	}
	if err := godotenv.Load(); err != nil {
		logging.Info("未找到.env文件，将使用默认配置: %v", err)
	}
	config = &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "3500"),
			ReadTimeout:  getEnvDuration("READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getEnvDuration("WRITE_TIMEOUT", 30*time.Second),
			DeployProxy:  getEnvBool("DEPLOY_PROXY", false),
		},
		Database: DatabaseConfig{
			MySQL: MySQLConfig{
				Host:           getEnv("MYSQL_HOST", "mysql-service"),
				Port:           getEnvInt("MYSQL_PORT", 3306),
				Username:       getEnv("MYSQL_USERNAME", "root"),
				Password:       getEnv("MYSQL_PASSWORD", ""),
				Database:       getEnv("MYSQL_DATABASE", "sandbox"),
				MaxConnections: getEnvInt("MYSQL_MAX_CONNECTIONS", 10),
				ConnTimeout:    getEnvDuration("MYSQL_CONN_TIMEOUT", 30*time.Second),
				ReadTimeout:    getEnvDuration("MYSQL_READ_TIMEOUT", 10*time.Second),
				WriteTimeout:   getEnvDuration("MYSQL_WRITE_TIMEOUT", 10*time.Second),
			},
			Redis: RedisConfig{
				Host:     getEnv("REDIS_HOST", "redis-service"),
				Port:     getEnvInt("REDIS_PORT", 6379),
				Password: getEnv("REDIS_PASSWORD", ""),
				Database: getEnvInt("REDIS_DATABASE", 0),
			},
			Proxy: ProxyConfig{
				Addr:       getEnv("PROXY_ADDR", "localhost:8474"),
				ListenAddr: getEnv("PROXY_LISTEN_ADDR", "localhost:6379"),
			},
		},
		Faults: FaultsConfig{
			CPU: CPUFaultConfig{
				DefaultDuration: getEnvInt("CPU_FAULT_DEFAULT_DURATION", 200),
			},
			Latency: LatencyFaultConfig{
				DefaultDelay: getEnvInt("LATENCY_FAULT_DEFAULT_DELAY", 200),
			},
			Redis: RedisFaultConfig{
				DefaultDelay: getEnvInt("REDIS_FAULT_DEFAULT_DELAY", 100),
			},
		},
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
