package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ra1nz0r/metric_alert_app/internal/storage"
)

func UpdateMetrics(h storage.MetricService, w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "type")
	mName := chi.URLParam(r, "name")
	mValue := chi.URLParam(r, "value")

	codeStatus := http.StatusOK

	switch {
	case strings.TrimSpace(mName) == "":
		codeStatus = http.StatusNotFound
	case mType == "gauge" && codeStatus != 404:
		v, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			ErrReturn(err, http.StatusBadRequest, w)
			return
		}
		h.UpdateGauge(mName, v)
	case mType == "counter" && codeStatus != 404:
		v, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			ErrReturn(err, http.StatusBadRequest, w)
			return
		}
		h.UpdateCounter(mName, v)
	default:
		codeStatus = http.StatusBadRequest
	}

	w.WriteHeader(codeStatus)
}
