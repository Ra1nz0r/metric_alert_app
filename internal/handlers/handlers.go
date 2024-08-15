package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ra1nz0r/metric_alert_app/internal/logger"
	"github.com/ra1nz0r/metric_alert_app/internal/storage"
)

// Интерфейс для взаимодействия хендлеров хранилищем
type HandlerService struct {
	sMS storage.MetricService
}

func NewHandlers(sMS storage.MetricService) *HandlerService {
	return &HandlerService{sMS: sMS}
}

// Собирает все метрики метрики из локального хранилища и выводит их в
// результирующей карте при получении GET запроса.
// Вызывает метод интерфейса, который возвращает копию локального хранилища.
// Формат JSON, в виде {"Alloc":146464,"Frees":10,...}.
func (hs *HandlerService) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]any)

	g, c := hs.sMS.MakeStorageCopy()

	for k, v := range *g {
		res[k] = v
	}

	for k, v := range *c {
		res[k] = v
	}

	ans, errJSON := json.Marshal(res)
	if errJSON != nil {
		logger.Zap.Error(fmt.Errorf("failed attempt json-marshal response: %w", errJSON))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write([]byte(ans)); errWrite != nil {
		logger.Zap.Error("failed attempt WRITE response")
		return
	}

}

// При получении GET запроса вида "/value/{type}/{name}", берёт тип с названием метрики
// и в ответе выводит её значение в "text/plain" формате, вид - 12345.
// При успешном запросе возвращает http.StatusOK.
// Выводит ошибку если тип указан неправильно или имя отсутствует в хранилище.
// Вызывает метод интерфейса, который возвращает копию локального хранилища.
func (hs *HandlerService) GetMetricByName(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "type")
	mName := chi.URLParam(r, "name")

	g, c := hs.sMS.MakeStorageCopy()

	var resVal any

	switch mType {
	case "gauge":
		gVal, ok := (*g)[mName]
		if !ok {
			ErrReturn(fmt.Errorf("metric not found"), http.StatusNotFound, w)
			return
		}
		resVal = gVal
	case "counter":
		cVal, ok := (*c)[mName]
		if !ok {
			ErrReturn(fmt.Errorf("metric not found"), http.StatusNotFound, w)
			return
		}
		resVal = cVal
	default:
		ErrReturn(fmt.Errorf("type not found"), http.StatusNotFound, w)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write([]byte(fmt.Sprintf("%v", resVal))); errWrite != nil {
		logger.Zap.Error("failed attempt WRITE response")
		return
	}
}

// Обновляет значение метрик в зависимости от типа и имени метрики.
// Тип gauge, float64 — новое значение замещает предыдущее.
// Тип counter, int64 — новое значение добавляется к предыдущему, если какое-то значение уже было известно серверу.
// Принимает метрики по протоколу HTTP методом POST.
// При успешном приёме возвращает http.StatusOK.
// При попытке передать запрос без имени метрики возвращет http.StatusNotFound.
// При попытке передать запрос с некорректным типом метрики или значением возвращет http.StatusBadRequest.
// Вызывает методы интерфейса хранилища, где реализовано взаимодействие и работа с ним.
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

func (hs *HandlerService) WithRequestDetails(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h.ServeHTTP(w, r)

		logger.Zap.Info(
			"URI:", r.RequestURI,
			"Method:", r.Method,
			"Duration:", time.Since(start),
		)
	})
}

func (hs *HandlerService) WithResponseDetails(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := logginResponseWriter{
			ResponseWriter: w,
			status:         0,
			size:           0,
		}

		h.ServeHTTP(&lw, r)

		logger.Zap.Info(
			"Status:", lw.status,
			"Size:", lw.size,
		)
	})
}

type logginResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (r *logginResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}

func (r *logginResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.status = statusCode
}
