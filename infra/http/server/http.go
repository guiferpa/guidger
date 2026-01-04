package server

import (
	"net/http"

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

func NewServerHTTP(opts NewServerHTTPOptions) *http.Server {
	if opts.Port == "" {
		opts.Port = ":8080"
	}

	router := chi.NewRouter()

	router.Route("/webhook", func(r chi.Router) {
		r.Post("/", ProcessWebhook(opts.Domain.Webhook, opts.Domain.Ledger))
	})

	router.Route("/balance", func(r chi.Router) {
		r.Get("/{user}", GetBalance(opts.Domain.Ledger))
	})

	return &http.Server{
		Addr:    opts.Port,
		Handler: router,
	}
}
