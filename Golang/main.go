package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CloudDetail/apo-sandbox/api"
	"github.com/CloudDetail/apo-sandbox/config"
	"github.com/CloudDetail/apo-sandbox/fault"
	"github.com/CloudDetail/apo-sandbox/logging"
	"github.com/CloudDetail/apo-sandbox/service"
	"github.com/CloudDetail/apo-sandbox/storage"
	"github.com/gorilla/mux"
)

var (
	faultManager *fault.Manager
)

func main() {
	// 加载配置
	appConfig := config.LoadConfig()

	// 设置日志级别
	logging.SetLevel(logging.LevelInfo)

	// 初始化存储层
	store := initStorage()

	// 初始化故障管理器
	initFaultManager(store.Redis)

	// 初始化业务服务
	businessService := service.NewBusinessService(store, faultManager)

	// 初始化API层
	businessAPI := &api.BusinessAPI{Service: businessService}

	// 设置路由
	router := setupRouter(businessAPI)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         ":" + appConfig.Server.Port,
		Handler:      router,
		ReadTimeout:  appConfig.Server.ReadTimeout,
		WriteTimeout: appConfig.Server.WriteTimeout,
	}

	// 启动服务器
	go func() {
		logging.Info("Starting server on port %s", appConfig.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("Server error: %v", err)
		}
	}()

	// 等待中断信号
	waitForShutdown(server)
}

// 初始化故障管理器
func initFaultManager(redisClient *storage.RedisClient) {
	faultManager = fault.NewManager()
	cpuFault := fault.NewCPUFault()
	faultManager.Register(cpuFault)
	latencyFault := fault.NewLatencyFault("eth0")
	faultManager.Register(latencyFault)
	redisFault := fault.NewRedisLatencyFault(redisClient)
	faultManager.Register(redisFault)
}

// 初始化存储层
func initStorage() *storage.Store {
	appConfig := config.LoadConfig()
	// 初始化MySQL连接
	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&timeout=%s&readTimeout=%s&writeTimeout=%s",
		appConfig.Database.MySQL.Username,
		appConfig.Database.MySQL.Password,
		appConfig.Database.MySQL.Host,
		appConfig.Database.MySQL.Port,
		appConfig.Database.MySQL.Database,
		appConfig.Database.MySQL.ConnTimeout,
		appConfig.Database.MySQL.ReadTimeout,
		appConfig.Database.MySQL.WriteTimeout,
	)

	mysqlClient, err := storage.NewMySQL(mysqlDSN)
	if err != nil {
		logging.Warn("Failed to connect to MySQL: %v", err)
		// 创建模拟客户端
		mysqlClient = &storage.MySQLClient{}
	}

	// 初始化Redis连接
	redisAddr := fmt.Sprintf("%s:%d", appConfig.Database.Redis.Host, appConfig.Database.Redis.Port)
	redisClient := storage.NewRedis(redisAddr)
	store := &storage.Store{
		MySQL: mysqlClient,
		Redis: redisClient,
	}

	logging.Info("%s", "Storage layer initialized")
	return store
}

func setupRouter(businessAPI *api.BusinessAPI) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/users", businessAPI.GetUsersCached).Methods("GET")

	router.Use(loggingMiddleware)

	return router
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logging.Info("Request started: %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		logging.Info("Request completed: %s %s - %v", r.Method, r.URL.Path, duration)
	})
}

func waitForShutdown(server *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logging.Info("%s", "Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if faultManager != nil {
		faultManager.StopAllFaults()
		logging.Info("%s", "All faults stopped")
	}

	if err := server.Shutdown(ctx); err != nil {
		logging.Error("Server shutdown error: %v", err)
	}

	logging.Info("%s", "Server stopped gracefully")
}
