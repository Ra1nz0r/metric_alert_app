package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ra1nz0r/metric_alert_app/internal/storage"
)

type HandlerService struct {
	sMS storage.MetricService
}

func NewHandlers(sMS storage.MetricService) *HandlerService {
	return &HandlerService{sMS: sMS}
}

// Выводит все метрики из локального хранилища при GET запросе.
// Принимает интерфейс, с реализованным методом чтения всех метрик из
// локального хранилища и объединения в одно целое.
// Формат JSON, в виде {"Alloc":146464,"Frees":10,...}.
func (hs *HandlerService) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	res, errJSON := json.Marshal(hs.sMS.AllMetricsFromStorage())
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

// Выводит значение метрики при GET запросе по типу и имени.
// Принимает интерфейс с реализованным методом получения
// указанной метрики из хранилища.
// Формат text, вид 12345.
func (hs *HandlerService) GetMetricByName(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "type")
	mName := chi.URLParam(r, "name")

	res, errGMV := hs.sMS.GetMetricVal(mType, mName)
	if errGMV != nil {
		ErrReturn(errGMV, http.StatusNotFound, w)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write([]byte(fmt.Sprintf("%v", res))); errWrite != nil {
		log.Print("failed attempt WRITE response")
		return
	}
}

// Обновляет значение метрик в зависимости от типа и имени метрики.
// Тип gauge, float64 — новое значение должно замещет предыдущее.
// Тип counter, int64 — новое значение должно добавляется к предыдущему, если какое-то значение уже было известно серверу.
// Принимает метрики по протоколу HTTP методом POST.
// При успешном приёме возвращает http.StatusOK.
// При попытке передать запрос без имени метрики возвращет http.StatusNotFound.
// При попытке передать запрос с некорректным типом метрики или значением возвращет http.StatusBadRequest.
// Принимает интерфейс, с созданным новым и инициализированным хранилищем,
// где реализованы методы для работы с ним.
func (hs *HandlerService) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
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
		hs.sMS.UpdateGauge(mName, v)
	case mType == "counter" && codeStatus != 404:
		v, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			ErrReturn(err, http.StatusBadRequest, w)
			return
		}
		hs.sMS.UpdateCounter(mName, v)
	default:
		codeStatus = http.StatusBadRequest
	}

	w.WriteHeader(codeStatus)
}
