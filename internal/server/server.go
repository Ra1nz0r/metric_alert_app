package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ra1nz0r/metric_alert_app/internal/config"
	hd "github.com/ra1nz0r/metric_alert_app/internal/handlers"
	"github.com/ra1nz0r/metric_alert_app/internal/logger"
	"github.com/ra1nz0r/metric_alert_app/internal/storage"
)

// Запускает агент, который будет принимать метрики от агента.
func Run() {
	r := chi.NewRouter()

	hs := hd.NewHandlers(storage.New())

	if errLog := logger.Initialize(config.DefLogLevel); errLog != nil {
		log.Fatal(errLog)
	}

	logger.Log.Info("Running handlers.")

	log.Println("Running handlers.")

	r.Use(WithLogging)
	//r.Handle("/", nil)

	r.Post("/update/{type}/{name}/{value}", hs.UpdateMetrics)

	r.Get("/", hs.GetAllMetrics)
	r.Get("/value/{type}/{name}", hs.GetMetricByName)

	config.ServerFlags()
	log.Printf("Starting server on: '%s'", config.DefServerHost)

	srv := http.Server{
		Addr:         config.DefServerHost,
		Handler:      r,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}

	go func() {
		if errListn := srv.ListenAndServe(); !errors.Is(errListn, http.ErrServerClosed) {
			log.Fatal("HTTP server error ", errListn)
		}
		log.Println("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if errShut := srv.Shutdown(shutdownCtx); errShut != nil {
		log.Fatal("HTTP shutdown error", errShut)
	}
	log.Println("Graceful shutdown complete.")
}

func WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := logginResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		logger.Log.Sugar().Infoln(
			"URI:", r.RequestURI,
			"Method:", r.Method,
			"Status:", responseData.status,
			"Duration:", time.Since(start),
			"Size:", responseData.size,
		)
	})
}

type responseData struct {
	status int
	size   int
}

type logginResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *logginResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *logginResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
