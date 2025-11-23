package providers

import (
	"conintracker-hiring/internal/testdata"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestEtherscanClientFetchNormalTransactions(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testdata.NormalTxResponse))
	}))
	defer server.Close()

	cfg := ClientConfig{
		APIKey:     "test-key",
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	}
	client := NewEtherscanClient(cfg)

	txs, err := client.FetchNormalTransactions(context.Background(), "0xa39b189482f984388a34460636fea9eb181ad1a6", 1, 1)
	if err != nil {
		t.Fatalf("FetchNormalTransactions() error = %v", err)
	}

	if len(txs) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(txs))
	}

	if txs[0].Hash != "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef" {
		t.Errorf("First transaction hash mismatch")
	}
}

func TestEtherscanClientFetchInternalTransactions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testdata.InternalTxResponse))
	}))
	defer server.Close()

	cfg := ClientConfig{
		APIKey:     "test-key",
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	}
	client := NewEtherscanClient(cfg)

	txs, err := client.FetchInternalTransactions(context.Background(), "0xa39b189482f984388a34460636fea9eb181ad1a6", 1, 1)
	if err != nil {
		t.Fatalf("FetchInternalTransactions() error = %v", err)
	}

	if len(txs) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(txs))
	}

	if txs[0].Type != "call" {
		t.Errorf("Expected call type, got %s", txs[0].Type)
	}
}

func TestEtherscanClientFetchTokenTransfers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testdata.ERC20TokenTxResponse))
	}))
	defer server.Close()

	cfg := ClientConfig{
		APIKey:     "test-key",
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	}
	client := NewEtherscanClient(cfg)

	txs, err := client.FetchTokenTransfers(context.Background(), "0xa39b189482f984388a34460636fea9eb181ad1a6", 1, 1)
	if err != nil {
		t.Fatalf("FetchTokenTransfers() error = %v", err)
	}

	if len(txs) != 2 {
		t.Errorf("Expected 2 token transactions, got %d", len(txs))
	}

	if txs[0].TokenSymbol != "USDC" {
		t.Errorf("Expected USDC token, got %s", txs[0].TokenSymbol)
	}
}

func TestEtherscanClientFetchNFTTransfers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testdata.ERC721NFTResponse))
	}))
	defer server.Close()

	cfg := ClientConfig{
		APIKey:     "test-key",
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	}
	client := NewEtherscanClient(cfg)

	txs, err := client.FetchNFTTransfers(context.Background(), "0xa39b189482f984388a34460636fea9eb181ad1a6", 1, 1)
	if err != nil {
		t.Fatalf("FetchNFTTransfers() error = %v", err)
	}

	if len(txs) != 1 {
		t.Errorf("Expected 1 NFT transaction, got %d", len(txs))
	}

	if txs[0].TokenID != "1337" {
		t.Errorf("Expected token ID 1337, got %s", txs[0].TokenID)
	}
}

func TestEtherscanClientFetchERC1155Transfers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testdata.ERC1155Response))
	}))
	defer server.Close()

	cfg := ClientConfig{
		APIKey:     "test-key",
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	}
	client := NewEtherscanClient(cfg)

	txs, err := client.FetchERC1155Transfers(context.Background(), "0xa39b189482f984388a34460636fea9eb181ad1a6", 1, 1)
	if err != nil {
		t.Fatalf("FetchERC1155Transfers() error = %v", err)
	}

	if len(txs) != 1 {
		t.Errorf("Expected 1 ERC-1155 transaction, got %d", len(txs))
	}

	if txs[0].TokenValue != "50" {
		t.Errorf("Expected token value 50, got %s", txs[0].TokenValue)
	}
}

func TestEtherscanClientErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testdata.ErrorResponse))
	}))
	defer server.Close()

	cfg := ClientConfig{
		APIKey:     "invalid-key",
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	}
	client := NewEtherscanClient(cfg)

	_, err := client.FetchNormalTransactions(context.Background(), "0xa39b189482f984388a34460636fea9eb181ad1a6", 1, 1)
	if err == nil {
		t.Error("Expected error for invalid API key, got none")
	}
}

func TestEtherscanClientRateLimiting(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testdata.EmptyResultResponse))
	}))
	defer server.Close()

	cfg := ClientConfig{
		APIKey:     "test-key",
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	}
	client := NewEtherscanClient(cfg)

	start := time.Now()
	
	// Make two requests in quick succession
	client.FetchNormalTransactions(context.Background(), "0xa39b189482f984388a34460636fea9eb181ad1a6", 1, 1)
	client.FetchNormalTransactions(context.Background(), "0xa39b189482f984388a34460636fea9eb181ad1a6", 1, 1)
	
	elapsed := time.Since(start)

	// Should have rate limited the second request
	if elapsed < RateLimitDelay {
		t.Errorf("Expected rate limiting delay, but elapsed time %v is less than %v", elapsed, RateLimitDelay)
	}

	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

func TestEtherscanClientEmptyResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testdata.EmptyResultResponse))
	}))
	defer server.Close()

	cfg := ClientConfig{
		APIKey:     "test-key",
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	}
	client := NewEtherscanClient(cfg)

	txs, err := client.FetchNormalTransactions(context.Background(), "0x0000000000000000000000000000000000000000", 1, 1)
	if err != nil {
		t.Fatalf("FetchNormalTransactions() error = %v", err)
	}

	if len(txs) != 0 {
		t.Errorf("Expected 0 transactions, got %d", len(txs))
	}
}

func TestNewEtherscanClient(t *testing.T) {
	tests := []struct {
		name string
		cfg  ClientConfig
	}{
		{
			name: "with_custom_config",
			cfg: ClientConfig{
				APIKey:   "test-key",
				BaseURL:  "http://custom-url",
			},
		},
		{
			name: "with_defaults",
			cfg: ClientConfig{
				APIKey: "test-key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewEtherscanClient(tt.cfg)
			if client.apiKey != tt.cfg.APIKey {
				t.Errorf("API key mismatch")
			}
			if tt.cfg.BaseURL != "" && client.baseURL != tt.cfg.BaseURL {
				t.Errorf("Base URL mismatch")
			}
		})
	}
}
