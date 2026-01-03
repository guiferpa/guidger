package server

import (
	"net/http"

	"github.com/guiferpa/guidger/domain/account"
	"github.com/guiferpa/guidger/domain/webhook"
)

func ProcessWebhook(whk webhook.Webhook, acc account.Account) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}
}
