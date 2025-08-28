package api

import (
	"encoding/json"
	"net/http"

	"github.com/CloudDetail/apo-sandbox/service"
)

type BusinessAPI struct {
	Service *service.BusinessService
}

func (b *BusinessAPI) GetUsers1(w http.ResponseWriter, r *http.Request) {
	active := r.URL.Query().Get("mode")
	chaos := ""
	if active == "1" {
		chaos = "latency"
	}
	result, err := b.Service.GetUsers(chaos, 0)
	if err != nil {
		http.Error(w, "get users failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{
		"data": result,
	})
}

func (b *BusinessAPI) GetUsers2(w http.ResponseWriter, r *http.Request) {
	active := r.URL.Query().Get("mode")
	chaos := ""
	if active == "1" {
		chaos = "cpu"
	}
	result, err := b.Service.GetUsers(chaos, 0)
	if err != nil {
		http.Error(w, "get users failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{
		"data": result,
	})
}

func (b *BusinessAPI) GetUsers3(w http.ResponseWriter, r *http.Request) {
	active := r.URL.Query().Get("mode")
	chaos := ""
	if active == "1" {
		chaos = "redis_latency"
	}
	result, err := b.Service.GetUsers(chaos, 0)
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
