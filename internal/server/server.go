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
	hd "github.com/ra1nz0r/metric_alert_app/internal/handlers"
	"github.com/ra1nz0r/metric_alert_app/internal/storage"
)

func Run() {
	var h storage.MetricService = storage.New()

	r := chi.NewRouter()

	log.Println("Running handlers.")
	r.Handle("/", nil)

	r.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		hd.UpdateMetrics(w, r, h)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		hd.GetAllMetrics(w, r, h)
	})

	serverLink := `0.0.0.0:8080`
	log.Printf("Starting server on: '%s'", serverLink)

	srv := http.Server{
		Addr:         serverLink,
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
