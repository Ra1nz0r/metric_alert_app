package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ra1nz0r/metric_alert_app/internal/storage"
)

func GetMetricByName(h storage.MetricService, w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "type")
	mName := chi.URLParam(r, "name")

	codeStatus := http.StatusOK

	ans, err := h.GetMetricVal(mType, mName)
	if err != nil {
		ErrReturn(err, http.StatusNotFound, w)
		return
	}

	res, errJSON := json.Marshal(ans)
	if errJSON != nil {
		http.Error(w, errJSON.Error(), http.StatusInternalServerError)
		//logerr.ErrEvent("failed attempt json-marshal response", errJSON)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(codeStatus)

	if _, errWrite := w.Write([]byte(res)); errWrite != nil {
		log.Print("failed attempt WRITE response")
		return
	}
}
