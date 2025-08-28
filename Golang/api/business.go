package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CloudDetail/apo-sandbox/service"
)

type BusinessAPI struct {
	Service *service.BusinessService
}

func (b *BusinessAPI) GetUsers1(w http.ResponseWriter, r *http.Request) {
	active := r.URL.Query().Get("mode")
	durationParam := r.URL.Query().Get("duration")
	duration, err := strconv.Atoi(durationParam)
	if err != nil {
		http.Error(w, "invalid duration parameter", http.StatusBadRequest)
		return
	}
	chaos := "latency"
	if active == "0" {
		chaos = ""
	}
	result, err := b.Service.GetUsers(chaos, duration)
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
	durationParam := r.URL.Query().Get("duration")
	duration, err := strconv.Atoi(durationParam)
	if err != nil {
		http.Error(w, "invalid duration parameter", http.StatusBadRequest)
		return
	}
	chaos := "cpu"
	if active == "0" {
		chaos = ""
	}
	result, err := b.Service.GetUsers(chaos, duration)
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
	durationParam := r.URL.Query().Get("duration")
	duration, err := strconv.Atoi(durationParam)
	if err != nil {
		http.Error(w, "invalid duration parameter", http.StatusBadRequest)
		return
	}
	chaos := "redis_latency"
	if active == "0" {
		chaos = ""
	}
	result, err := b.Service.GetUsers(chaos, duration)
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
