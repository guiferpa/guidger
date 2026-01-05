package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/guiferpa/guidger/domain/ledger"
	"github.com/guiferpa/guidger/httputil"
)

// GetBalanceResponseBody represents the JSON response for a balance query
type GetBalanceResponseBody struct {
	User     string            `json:"user"`     // User identifier
	Balances map[string]string `json:"balances"` // Map of asset names to balance strings
}

// GetBalance handles GET /balance/{user} requests
// Returns all balances for the specified user
func GetBalance(ldgr ledger.Ledger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := chi.URLParam(r, "user")
		if user == "" {
			httputil.Error(w, http.StatusBadRequest, fmt.Errorf("user parameter is required"))
			return
		}

		requestID := GetRequestID(r)
		logger := slog.Default()

		balances, err := ldgr.GetBalanceByUser(user)
		if err != nil {
			logger.Error("balance_query_failed",
				"request_id", requestID,
				"error", err.Error(),
				"user", user,
			)
			httputil.Error(w, http.StatusInternalServerError, err)
			return
		}

		logger.Info("balance_queried",
			"request_id", requestID,
			"user", user,
			"assets_count", len(balances),
		)

		response := GetBalanceResponseBody{
			User:     user,
			Balances: balances,
		}

		httputil.Respond(w, http.StatusOK, response)
	}
}
