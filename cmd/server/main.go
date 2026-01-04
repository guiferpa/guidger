package main

import (
	"log"
	"time"

	"github.com/guiferpa/guidger/domain/ledger"
	"github.com/guiferpa/guidger/domain/webhook"
	"github.com/guiferpa/guidger/infra/http/server"
	"github.com/guiferpa/guidger/infra/signer/hmac"
	"github.com/guiferpa/guidger/infra/storage/memory"
)

func main() {
	sm := memory.NewMemoryStorage()
	sgr := hmac.NewHMACSigner()

	s := server.NewServerHTTP(server.NewServerHTTPOptions{
		Port: ":8080",
		Domain: server.DomainOptions{
			Ledger: ledger.New(ledger.NewOptions{
				Storage: sm,
			}),
			Webhook: webhook.New(webhook.NewOptions{
				SignatureTolerance: 5 * time.Minute,
				Storage:            sm,
				Signer:             sgr,
			}),
		},
	})

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
