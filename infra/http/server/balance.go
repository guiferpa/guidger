package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/guiferpa/guidger/domain/ledger"
	"github.com/guiferpa/guidger/httputil"
)

type GetBalanceResponseBody struct {
	User     string            `json:"user"`
	Balances map[string]string `json:"balances"`
}

func GetBalance(ldgr ledger.Ledger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := chi.URLParam(r, "user")
		if user == "" {
			httputil.Error(w, http.StatusBadRequest, fmt.Errorf("user parameter is required"))
			return
		}

		balances, err := ldgr.GetBalanceByUser(user)
		if err != nil {
			httputil.Error(w, http.StatusInternalServerError, err)
			return
		}

		response := GetBalanceResponseBody{
			User:     user,
			Balances: balances,
		}

		httputil.Respond(w, http.StatusOK, response)
	}
}
