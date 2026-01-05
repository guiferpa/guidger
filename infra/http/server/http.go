package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/guiferpa/guidger/domain/ledger"
	"github.com/guiferpa/guidger/domain/webhook"
)

type DomainOptions struct {
	Ledger  ledger.Ledger
	Webhook webhook.Webhook
}

type NewServerHTTPOptions struct {
	Port   string
	Domain DomainOptions
}

// NewServerHTTP creates a new HTTP server with structured logging, telemetry, and routing configured
func NewServerHTTP(opts NewServerHTTPOptions) *http.Server {
	if opts.Port == "" {
		opts.Port = ":8080"
	}

	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Initialize telemetry
	telemetry := NewTelemetry()

	router := chi.NewRouter()

	// Apply middleware (order matters)
	router.Use(RequestIDMiddleware)
	router.Use(LoggingMiddleware(logger))
	router.Use(ErrorLoggingMiddleware(logger))
	router.Use(TelemetryMiddleware(telemetry))

	// Routes
	router.Route("/webhook", func(r chi.Router) {
		r.Post("/", ProcessWebhook(opts.Domain.Webhook, opts.Domain.Ledger))
	})

	router.Route("/balance", func(r chi.Router) {
		r.Get("/{user}", GetBalance(opts.Domain.Ledger))
	})

	// Metrics endpoint
	router.Get("/metrics", GetMetrics(telemetry))

	return &http.Server{
		Addr:    opts.Port,
		Handler: router,
	}
}
