package api

import (
	"encoding/json"
	"net/http"

	"github.com/CloudDetail/apo-sandbox/config"
	"github.com/CloudDetail/apo-sandbox/service"
)

type BusinessAPI struct {
	Service *service.BusinessService
}

func (b *BusinessAPI) GetUsers1(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	delayMs := config.LoadConfig().Faults.Latency.DefaultDelay
	result, err := b.Service.GetUsers1(mode, delayMs)
	if err != nil {
		http.Error(w, "get users failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{
		"data": result,
	})
}

func (b *BusinessAPI) GetUsers2(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	faultConfig := config.LoadConfig().Faults.CPU
	duration := faultConfig.DefaultDuration

	result, err := b.Service.GetUsers2(mode, duration)
	if err != nil {
		http.Error(w, "get users failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{
		"data": result,
	})
}

func (b *BusinessAPI) GetUsers3(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	faultConfig := config.LoadConfig().Faults.Redis
	duration := faultConfig.DefaultDelay
	result, err := b.Service.GetUsers3(mode, duration)
	if err != nil {
		http.Error(w, "get users failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{
		"data": result,
	})
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}
