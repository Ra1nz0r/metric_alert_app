package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ra1nz0r/metric_alert_app/internal/storage"
)

func GetAllMetrics(h storage.MetricService, w http.ResponseWriter, r *http.Request) {
	res, errJSON := json.Marshal(h.AllMetricsFromStorage())
	if errJSON != nil {
		http.Error(w, errJSON.Error(), http.StatusInternalServerError)
		//logerr.ErrEvent("failed attempt json-marshal response", errJSON)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write([]byte(res)); errWrite != nil {
		log.Print("failed attempt WRITE response")
		return
	}

}
