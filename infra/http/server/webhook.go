package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/guiferpa/guidger/domain/ledger"
	"github.com/guiferpa/guidger/domain/webhook"
	"github.com/guiferpa/guidger/httputil"
)

// ProcessWebhookRequestBody represents the JSON body of a webhook request
type ProcessWebhookRequestBody struct {
	Asset  string `json:"asset"`  // Asset identifier (e.g., "BTC", "ETH")
	Amount string `json:"amount"` // Decimal amount as string (e.g., "1.5", "0.0001")
	User   string `json:"user"`   // User identifier
}

// ProcessWebhook handles POST /webhook requests
// It validates the signature, timestamp, and nonce, then applies the balance update
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
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, err)
			return
		}

		var body ProcessWebhookRequestBody
		if err := json.Unmarshal(payload, &body); err != nil {
			httputil.Error(w, http.StatusBadRequest, err)
			return
		}

		requestID := GetRequestID(r)
		logger := slog.Default()

		if err := whk.ValidateSignature(timestamp, nonce, payload, signature); err != nil {
			logger.Warn("webhook_validation_failed",
				"request_id", requestID,
				"error", err.Error(),
				"user", body.User,
			)
			httputil.Error(w, http.StatusUnauthorized, err)
			return
		}
		if err := ldgr.Apply(body.User, body.Asset, body.Amount); err != nil {
			logger.Error("ledger_apply_failed",
				"request_id", requestID,
				"error", err.Error(),
				"user", body.User,
				"asset", body.Asset,
				"amount", body.Amount,
			)
			httputil.Error(w, http.StatusInternalServerError, err)
			return
		}

		logger.Info("webhook_processed",
			"request_id", requestID,
			"user", body.User,
			"asset", body.Asset,
			"amount", body.Amount,
		)
		httputil.Respond(w, http.StatusOK, map[string]string{"message": "Webhook received"})
	}
}
