package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/guiferpa/guidger/domain/ledger"
	"github.com/guiferpa/guidger/domain/webhook"
	"github.com/guiferpa/guidger/infra/signer/hmac"
	"github.com/guiferpa/guidger/infra/storage/memory"
)

func setupTestServer() (*http.Server, *memory.MemoryStorage, *hmac.HMACSigner) {
	sm := memory.NewMemoryStorage()
	sgr := hmac.NewHMACSigner()

	s := NewServerHTTP(NewServerHTTPOptions{
		Port: ":0", // Use :0 to get a random port for testing
		Domain: DomainOptions{
			Ledger: ledger.New(ledger.NewOptions{
				Storage: sm,
			}),
			Webhook: webhook.New(webhook.NewOptions{
				SignatureTolerance: 5 * time.Minute,
				SignatureSecret:    "test-secret-key",
				Storage:            sm,
				Signer:             sgr,
			}),
		},
	})

	return s, sm, sgr
}

func closeTestServer(s *http.Server) {
	_ = s.Close() // In tests, we can safely ignore Close errors
}

func TestProcessWebhook_Success(t *testing.T) {
	s, _, sgr := setupTestServer()
	defer closeTestServer(s)

	timestamp := time.Now().Unix()
	nonce := "test-nonce-1"
	payload := []byte(`{"user":"user1","asset":"BTC","amount":"1.5"}`)
	signature := sgr.GenerateSignature(timestamp, nonce, payload, "test-secret-key")

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-Nonce", nonce)
	req.Header.Set("X-Signature", signature)

	w := httptest.NewRecorder()
	s.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["message"] != "Webhook received" {
		t.Fatalf("expected message 'Webhook received', got %s", response["message"])
	}
}

func TestProcessWebhook_InvalidSignature(t *testing.T) {
	s, _, _ := setupTestServer()
	defer closeTestServer(s)

	timestamp := time.Now().Unix()
	nonce := "test-nonce-2"
	payload := []byte(`{"user":"user1","asset":"BTC","amount":"1.5"}`)
	signature := "invalid-signature"

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-Nonce", nonce)
	req.Header.Set("X-Signature", signature)

	w := httptest.NewRecorder()
	s.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProcessWebhook_TimestampTooOld(t *testing.T) {
	s, _, sgr := setupTestServer()
	defer closeTestServer(s)

	// Timestamp 10 minutes ago (beyond 5 minute tolerance)
	timestamp := time.Now().Add(-10 * time.Minute).Unix()
	nonce := "test-nonce-3"
	payload := []byte(`{"user":"user1","asset":"BTC","amount":"1.5"}`)
	signature := sgr.GenerateSignature(timestamp, nonce, payload, "test-secret-key")

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-Nonce", nonce)
	req.Header.Set("X-Signature", signature)

	w := httptest.NewRecorder()
	s.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProcessWebhook_DuplicateNonce(t *testing.T) {
	s, _, sgr := setupTestServer()
	defer closeTestServer(s)

	timestamp := time.Now().Unix()
	nonce := "test-nonce-4"
	payload := []byte(`{"user":"user1","asset":"BTC","amount":"1.5"}`)
	signature := sgr.GenerateSignature(timestamp, nonce, payload, "test-secret-key")

	// First request should succeed
	req1 := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
	req1.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestamp))
	req1.Header.Set("X-Nonce", nonce)
	req1.Header.Set("X-Signature", signature)

	w1 := httptest.NewRecorder()
	s.Handler.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("first request expected status 200, got %d: %s", w1.Code, w1.Body.String())
	}

	// Second request with same nonce should fail
	timestamp2 := time.Now().Unix()
	signature2 := sgr.GenerateSignature(timestamp2, nonce, payload, "test-secret-key")

	req2 := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
	req2.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestamp2))
	req2.Header.Set("X-Nonce", nonce)
	req2.Header.Set("X-Signature", signature2)

	w2 := httptest.NewRecorder()
	s.Handler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusUnauthorized {
		t.Fatalf("second request expected status 401 (replay attack), got %d: %s", w2.Code, w2.Body.String())
	}
}

func TestProcessWebhook_MissingHeaders(t *testing.T) {
	s, _, _ := setupTestServer()
	defer closeTestServer(s)

	payload := []byte(`{"user":"user1","asset":"BTC","amount":"1.5"}`)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
	// Missing headers

	w := httptest.NewRecorder()
	s.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetBalance_UserWithBalances(t *testing.T) {
	s, _, sgr := setupTestServer()
	defer closeTestServer(s)

	// First, add some balances via webhook
	timestamp := time.Now().Unix()
	nonce1 := "test-nonce-balance-1"
	payload1 := []byte(`{"user":"user2","asset":"BTC","amount":"1.5"}`)
	signature1 := sgr.GenerateSignature(timestamp, nonce1, payload1, "test-secret-key")

	req1 := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload1))
	req1.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestamp))
	req1.Header.Set("X-Nonce", nonce1)
	req1.Header.Set("X-Signature", signature1)

	w1 := httptest.NewRecorder()
	s.Handler.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("webhook request expected status 200, got %d: %s", w1.Code, w1.Body.String())
	}

	// Add another asset
	timestamp2 := time.Now().Unix()
	nonce2 := "test-nonce-balance-2"
	payload2 := []byte(`{"user":"user2","asset":"ETH","amount":"2.25"}`)
	signature2 := sgr.GenerateSignature(timestamp2, nonce2, payload2, "test-secret-key")

	req2 := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload2))
	req2.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestamp2))
	req2.Header.Set("X-Nonce", nonce2)
	req2.Header.Set("X-Signature", signature2)

	w2 := httptest.NewRecorder()
	s.Handler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("webhook request expected status 200, got %d: %s", w2.Code, w2.Body.String())
	}

	// Now get the balance
	req3 := httptest.NewRequest(http.MethodGet, "/balance/user2", nil)
	w3 := httptest.NewRecorder()
	s.Handler.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w3.Code, w3.Body.String())
	}

	var response GetBalanceResponseBody
	if err := json.NewDecoder(w3.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.User != "user2" {
		t.Fatalf("expected user 'user2', got %s", response.User)
	}

	if len(response.Balances) != 2 {
		t.Fatalf("expected 2 balances, got %d", len(response.Balances))
	}

	if response.Balances["BTC"] != "1.5" {
		t.Fatalf("expected BTC balance '1.5', got %s", response.Balances["BTC"])
	}

	if response.Balances["ETH"] != "2.25" {
		t.Fatalf("expected ETH balance '2.25', got %s", response.Balances["ETH"])
	}
}

func TestGetBalance_UserWithoutBalances(t *testing.T) {
	s, _, _ := setupTestServer()
	defer closeTestServer(s)

	req := httptest.NewRequest(http.MethodGet, "/balance/user3", nil)
	w := httptest.NewRecorder()
	s.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response GetBalanceResponseBody
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.User != "user3" {
		t.Fatalf("expected user 'user3', got %s", response.User)
	}

	if len(response.Balances) != 0 {
		t.Fatalf("expected 0 balances, got %d", len(response.Balances))
	}
}

func TestGetBalance_MissingUser(t *testing.T) {
	s, _, _ := setupTestServer()
	defer closeTestServer(s)

	req := httptest.NewRequest(http.MethodGet, "/balance/", nil)
	w := httptest.NewRecorder()
	s.Handler.ServeHTTP(w, req)

	// chi router will return 404 for missing path parameter
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestWebhookToBalance_IntegrationFlow(t *testing.T) {
	s, _, sgr := setupTestServer()
	defer closeTestServer(s)

	user := "integration-user"
	asset := "USDT"
	amount := "100.50"

	// Step 1: Send webhook
	timestamp := time.Now().Unix()
	nonce := "test-nonce-integration"
	payload := []byte(fmt.Sprintf(`{"user":"%s","asset":"%s","amount":"%s"}`, user, asset, amount))
	signature := sgr.GenerateSignature(timestamp, nonce, payload, "test-secret-key")

	req1 := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
	req1.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestamp))
	req1.Header.Set("X-Nonce", nonce)
	req1.Header.Set("X-Signature", signature)

	w1 := httptest.NewRecorder()
	s.Handler.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("webhook expected status 200, got %d: %s", w1.Code, w1.Body.String())
	}

	// Step 2: Add more to the same asset
	timestamp2 := time.Now().Unix()
	nonce2 := "test-nonce-integration-2"
	amount2 := "50.25"
	payload2 := []byte(fmt.Sprintf(`{"user":"%s","asset":"%s","amount":"%s"}`, user, asset, amount2))
	signature2 := sgr.GenerateSignature(timestamp2, nonce2, payload2, "test-secret-key")

	req2 := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload2))
	req2.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestamp2))
	req2.Header.Set("X-Nonce", nonce2)
	req2.Header.Set("X-Signature", signature2)

	w2 := httptest.NewRecorder()
	s.Handler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("webhook expected status 200, got %d: %s", w2.Code, w2.Body.String())
	}

	// Step 3: Get balance and verify it's the sum
	req3 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/balance/%s", user), nil)
	w3 := httptest.NewRecorder()
	s.Handler.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Fatalf("balance expected status 200, got %d: %s", w3.Code, w3.Body.String())
	}

	var response GetBalanceResponseBody
	if err := json.NewDecoder(w3.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Balances[asset] != "150.75" {
		t.Fatalf("expected balance '150.75' (100.50 + 50.25), got %s", response.Balances[asset])
	}
}
