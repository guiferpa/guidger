package main

import (
	"fmt"
	"log"
	"time"

	"github.com/guiferpa/guidger/domain/ledger"
	"github.com/guiferpa/guidger/domain/webhook"
	"github.com/guiferpa/guidger/envutil"
	"github.com/guiferpa/guidger/infra/http/server"
	"github.com/guiferpa/guidger/infra/signer/hmac"
	"github.com/guiferpa/guidger/infra/storage/memory"
)

func main() {
	port, err := envutil.GetInt("PORT", 8080)
	if err != nil {
		log.Fatalf("failed to get port: %v", err)
	}

	secret, err := envutil.GetStringRequired("WEBHOOK_SECRET")
	if err != nil {
		log.Fatalf("failed to get webhook secret: %v", err)
	}

	sm := memory.NewMemoryStorage()
	sgr := hmac.NewHMACSigner()

	s := server.NewServerHTTP(server.NewServerHTTPOptions{
		Port: fmt.Sprintf(":%d", port),
		Domain: server.DomainOptions{
			Ledger: ledger.New(ledger.NewOptions{
				Storage: sm,
			}),
			Webhook: webhook.New(webhook.NewOptions{
				SignatureTolerance: 5 * time.Minute,
				SignatureSecret:    secret,
				Storage:            sm,
				Signer:             sgr,
			}),
		},
	})

	log.Printf("Server starting on port :%d", port)
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
