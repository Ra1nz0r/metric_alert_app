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
	config.ServerFlags()

	r := chi.NewRouter()

	hs := hd.NewHandlers(storage.New())

	if errLog := logger.Initialize(config.DefLogLevel); errLog != nil {
		log.Fatal(errLog)
	}

	logger.Log.Info("Running handlers.")

	r.Use(hs.WithLogging)

	r.Post("/update/{type}/{name}/{value}", hs.UpdateMetrics)

	r.Get("/", hs.GetAllMetrics)
	r.Get("/value/{type}/{name}", hs.GetMetricByName)

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
