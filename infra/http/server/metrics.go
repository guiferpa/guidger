package server

import (
	"net/http"

	"github.com/guiferpa/guidger/httputil"
)

// GetMetrics returns the current telemetry metrics
func GetMetrics(telemetry *Telemetry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := telemetry.GetMetrics()
		httputil.Respond(w, http.StatusOK, metrics)
	}
}
