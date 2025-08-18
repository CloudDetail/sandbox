package config

import (
	"os"
	"strconv"
	"time"

	"github.com/CloudDetail/apo-sandbox/logging"
	"github.com/joho/godotenv"
)

var config *Config

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Faults   FaultsConfig   `json:"faults"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MySQL MySQLConfig `json:"mysql"`
	Redis RedisConfig `json:"redis"`
}

// MySQLConfig MySQL配置
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

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Database int    `json:"database"`
}

// FaultsConfig 故障配置
type FaultsConfig struct {
	CPU     CPUFaultConfig     `json:"cpu"`
	Latency LatencyFaultConfig `json:"latency"`
	Redis   RedisFaultConfig   `json:"redis"`
}

// CPUFaultConfig CPU故障配置
type CPUFaultConfig struct {
	DefaultDuration int `json:"default_duration"`
}

// LatencyFaultConfig 延迟故障配置
type LatencyFaultConfig struct {
	DefaultDelay int `json:"default_delay"`
}

// RedisFaultConfig Redis故障配置
type RedisFaultConfig struct {
	DefaultDelay int `json:"default_delay"`
}

// LoadConfig 加载配置
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
		},
		Database: DatabaseConfig{
			MySQL: MySQLConfig{
				Host:           getEnv("MYSQL_HOST", "localhost"),
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
				Host:     getEnv("REDIS_HOST", "localhost"),
				Port:     getEnvInt("REDIS_PORT", 6379),
				Password: getEnv("REDIS_PASSWORD", ""),
				Database: getEnvInt("REDIS_DATABASE", 0),
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

// 辅助函数：获取环境变量
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// 辅助函数：获取环境变量并转换为整数
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// 辅助函数：获取环境变量并转换为浮点数
func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}

// 辅助函数：获取环境变量并转换为时间间隔
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
