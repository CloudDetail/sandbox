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

func (b *BusinessAPI) GetUsersCached(w http.ResponseWriter, r *http.Request) {
	chaos := r.URL.Query().Get("chaos")
	durationParam := r.URL.Query().Get("duration")
	duration, err := strconv.Atoi(durationParam)
	result, err := b.Service.GetUsersCached(chaos, duration)
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
