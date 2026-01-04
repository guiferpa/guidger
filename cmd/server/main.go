package main

import (
	"time"

	"github.com/guiferpa/guidger/domain/ledger"
	"github.com/guiferpa/guidger/domain/webhook"
	"github.com/guiferpa/guidger/infra/http/server"
	"github.com/guiferpa/guidger/infra/storage/memory"
)

func main() {
	sm := memory.NewStorageMemory()

	s := server.NewHTTPServer(server.NewHTTPServerOptions{
		Port: ":8080",
		Domain: server.DomainOptions{
			Ledger: ledger.New(ledger.NewOptions{
				Storage: sm,
			}),
			Webhook: webhook.New(webhook.NewOptions{
				SignatureTolerance: 5 * time.Minute,
				Storage:            sm,
			}),
		},
	})

	s.ListenAndServe()
}
