package main

import (
	"github.com/guiferpa/guidger/domain/account"
	"github.com/guiferpa/guidger/domain/webhook"
	"github.com/guiferpa/guidger/infra/http/server"
	"github.com/guiferpa/guidger/infra/ledger/memory"
)

func main() {
	l := memory.NewLedgerMemory()

	s := server.NewHTTPServer(server.NewHTTPServerOptions{
		Port: ":8080",
		Domain: server.DomainOptions{
			Account: account.New(l),
			Webhook: webhook.New(),
		},
	})

	s.ListenAndServe()
}
