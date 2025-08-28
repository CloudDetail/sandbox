package api

import (
	"encoding/json"
	"net/http"

	"github.com/CloudDetail/apo-sandbox/service"
)

type BusinessAPI struct {
	Service *service.BusinessService
}

func (b *BusinessAPI) GetUsersCached(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	var chaos string
	switch mode {
	case "1":
		chaos = "latency"
	case "2":
		chaos = "cpu"
	case "3":
		chaos = "redis_latency"
	}
	result, err := b.Service.GetUsersCached(chaos, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
