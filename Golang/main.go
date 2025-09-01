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
	toxiproxy "github.com/Shopify/toxiproxy/v2/client"
	"github.com/gorilla/mux"
)

var (
	faultManager *fault.Manager
)

func main() {
	// load config
	appConfig := config.LoadConfig()

	// set log level
	logging.SetLevel(logging.LevelInfo)

	// init storage layer
	store, err := initStorage()
	if err != nil {
		panic(err)
	}

	// init fault manager
	// initFaultManager(store.Redis)

	// init business service
	businessService := service.NewBusinessService(store)

	// init business api
	businessAPI := &api.BusinessAPI{Service: businessService}

	// setup router
	router := setupRouter(businessAPI)

	// setup http server
	server := &http.Server{
		Addr:         ":" + appConfig.Server.Port,
		Handler:      router,
		ReadTimeout:  appConfig.Server.ReadTimeout,
		WriteTimeout: appConfig.Server.WriteTimeout,
	}

	// start http server
	go func() {
		logging.Info("Starting server on port %s", appConfig.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("Server error: %v", err)
		}
	}()

	// wait for shutdown signal
	waitForShutdown(server)
}

func initFaultManager(redisClient *storage.RedisClient) {

	faultManager = fault.NewManager()
	cpuFault := fault.NewCPUFault()
	faultManager.Register(cpuFault)
	latencyFault := fault.NewLatencyFault("eth0")
	faultManager.Register(latencyFault)
	redisFault := fault.NewRedisLatencyFault(redisClient)
	faultManager.Register(redisFault)

	logging.Info("%s", "Fault manager initialized")
}

func initStorage() (*storage.Store, error) {
	appConfig := config.LoadConfig()
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
		mysqlClient = &storage.MySQLClient{}
	}

	store := &storage.Store{}
	if appConfig.Server.DeployProxy {
		client := toxiproxy.NewClient(appConfig.Database.Proxy.Addr)
		proxies, err := client.Populate([]toxiproxy.Proxy{{
			Name:     "redis",
			Listen:   appConfig.Database.Proxy.ListenAddr,
			Upstream: "redis-service:6379",
			Enabled:  true,
		}})
		if err != nil {
			return nil, err
		}
		store.Proxy = proxies[0]
	}
	redisClient, err := storage.NewRedis(fmt.Sprintf("%s:%d", appConfig.Database.Redis.Host, appConfig.Database.Redis.Port))
	if err != nil {
		return nil, err
	}
	store.Redis = redisClient
	store.MySQL = mysqlClient

	logging.Info("%s", "Storage layer initialized")
	return store, nil
}

func setupRouter(businessAPI *api.BusinessAPI) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/users/1", businessAPI.GetUsers1).Methods("GET")
	router.HandleFunc("/api/users/2", businessAPI.GetUsers2).Methods("GET")
	router.HandleFunc("/api/users/3", businessAPI.GetUsers3).Methods("GET")

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
