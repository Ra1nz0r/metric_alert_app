package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ra1nz0r/metric_alert_app/internal/storage"
)

type UpdateMetricsResult struct {
	val interface{}
	err error
}

func UpdateMetrics(w http.ResponseWriter, r *http.Request, h storage.MetricService) {
	mType := chi.URLParam(r, "type")
	mName := chi.URLParam(r, "name")
	mValue := chi.URLParam(r, "value")

	var umr UpdateMetricsResult

	codeStatus := http.StatusOK

	switch {
	case strings.TrimSpace(mName) == "":
		codeStatus = http.StatusNotFound
	case mType == "gauge" && codeStatus != 404:
		v, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			log.Println("Error from strconv: ", err)
			return
		}
		umr.val, umr.err = h.UpdateGauge(mName, v)
	case mType == "counter" && codeStatus != 404:
		v, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			log.Println("Error from strconv: ", err)
			return
		}
		umr.val, umr.err = h.UpdateCounter(mName, v)
	default:
		codeStatus = http.StatusBadRequest
	}

	if umr.err != nil {
		ErrReturn(umr.err, w)
		return
	}

	res, errJSON := json.Marshal(umr.val)
	if errJSON != nil {
		http.Error(w, errJSON.Error(), http.StatusInternalServerError)
		//logerr.ErrEvent("failed attempt json-marshal response", errJSON)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	w.WriteHeader(codeStatus)

	if _, errWrite := w.Write([]byte(res)); errWrite != nil {
		log.Print("failed attempt WRITE response")
		return
	}
}
