package server

import (
	"encoding/json"
	"net/http"

	"github.com/guiferpa/guidger/domain/ledger"
	"github.com/guiferpa/guidger/domain/webhook"
	"github.com/guiferpa/guidger/httputil"
)

type ProcessWebhookRequestBody struct {
	Asset  string `json:"asset"`
	Amount string `json:"amount"`
	User   string `json:"user"`
}

func ProcessWebhook(whk webhook.Webhook, ldgr ledger.Ledger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timestamp, err := httputil.GetHeaderInt64(r, "X-Timestamp")
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, err)
			return
		}
		nonce, err := httputil.GetHeaderString(r, "X-Nonce")
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, err)
			return
		}
		signature, err := httputil.GetHeaderString(r, "X-Signature")
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, err)
			return
		}
		var body ProcessWebhookRequestBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			httputil.Error(w, http.StatusBadRequest, err)
			return
		}
		payload, err := json.Marshal(body)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, err)
			return
		}
		if err := whk.ValidateSignature(timestamp, nonce, payload, signature); err != nil {
			httputil.Error(w, http.StatusUnauthorized, err)
			return
		}
		if err := ldgr.Apply(body.User, body.Asset, body.Amount); err != nil {
			httputil.Error(w, http.StatusInternalServerError, err)
			return
		}
		httputil.Respond(w, http.StatusOK, map[string]string{"message": "Webhook received"})
	}
}
