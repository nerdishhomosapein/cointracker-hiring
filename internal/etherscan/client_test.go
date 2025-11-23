package etherscan

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fixture loads a JSON file from internal/etherscan/testdata.
func fixture(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", name, err)
	}
	return string(data)
}

func TestClientFetchNormalTx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("module") != "account" || q.Get("action") != "txlist" {
			t.Fatalf("unexpected module/action: %s/%s", q.Get("module"), q.Get("action"))
		}
		if got := q.Get("page"); got != "1" {
			t.Fatalf("expected page=1 got=%s", got)
		}
		if got := q.Get("offset"); got != "2" {
			t.Fatalf("expected offset=2 got=%s", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(fixture(t, "normal_tx.json")))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "dummy-key", srv.Client())
	txs, err := client.GetNormalTx(context.Background(), "0xabc", 1, 2, "asc")
	if err != nil {
		t.Fatalf("GetNormalTx returned error: %v", err)
	}
	if len(txs) != 2 {
		t.Fatalf("expected 2 txs, got %d", len(txs))
	}
	if txs[0].Hash != "0xhash1" || txs[1].Hash != "0xhash2" {
		t.Fatalf("unexpected hashes: %+v", txs)
	}
	if txs[0].GasPrice != "1000000000" || txs[0].GasUsed != "21000" {
		t.Fatalf("unexpected gas fields: %+v", txs[0])
	}
}

func TestClientFetchInternalTx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("module") != "account" || q.Get("action") != "txlistinternal" {
			t.Fatalf("unexpected module/action: %s/%s", q.Get("module"), q.Get("action"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(fixture(t, "internal_tx.json")))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "dummy-key", srv.Client())
	txs, err := client.GetInternalTx(context.Background(), "0xabc", 1, 10, "asc")
	if err != nil {
		t.Fatalf("GetInternalTx returned error: %v", err)
	}
	if len(txs) != 1 {
		t.Fatalf("expected 1 tx, got %d", len(txs))
	}
	if txs[0].ContractAddress != "0xcontract1" || txs[0].Value != "5000000000000000" {
		t.Fatalf("unexpected tx: %+v", txs[0])
	}
}

func TestClientFetchTokenTx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("module") != "account" || q.Get("action") != "tokentx" {
			t.Fatalf("unexpected module/action: %s/%s", q.Get("module"), q.Get("action"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(fixture(t, "erc20_tx.json")))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "dummy-key", srv.Client())
	txs, err := client.GetTokenTx(context.Background(), "0xabc", 1, 10, "asc")
	if err != nil {
		t.Fatalf("GetTokenTx returned error: %v", err)
	}
	if len(txs) != 1 {
		t.Fatalf("expected 1 token tx, got %d", len(txs))
	}
	if txs[0].TokenSymbol != "USDC" || txs[0].TokenDecimal != "6" {
		t.Fatalf("unexpected token fields: %+v", txs[0])
	}
}

func TestClientFetchERC721Tx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("module") != "account" || q.Get("action") != "tokennfttx" {
			t.Fatalf("unexpected module/action: %s/%s", q.Get("module"), q.Get("action"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(fixture(t, "erc721_tx.json")))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "dummy-key", srv.Client())
	txs, err := client.GetERC721Tx(context.Background(), "0xabc", 1, 10, "asc")
	if err != nil {
		t.Fatalf("GetERC721Tx returned error: %v", err)
	}
	if len(txs) != 1 {
		t.Fatalf("expected 1 nft tx, got %d", len(txs))
	}
	if txs[0].TokenID != "12345" || txs[0].TokenName != "CoolNFTs" {
		t.Fatalf("unexpected nft fields: %+v", txs[0])
	}
}

func TestClientFetchERC1155Tx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("module") != "account" || q.Get("action") != "token1155tx" {
			t.Fatalf("unexpected module/action: %s/%s", q.Get("module"), q.Get("action"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(fixture(t, "erc1155_tx.json")))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "dummy-key", srv.Client())
	txs, err := client.GetERC1155Tx(context.Background(), "0xabc", 1, 10, "asc")
	if err != nil {
		t.Fatalf("GetERC1155Tx returned error: %v", err)
	}
	if len(txs) != 1 {
		t.Fatalf("expected 1 erc1155 tx, got %d", len(txs))
	}
	if txs[0].TokenID != "777" || txs[0].TokenValue != "3" {
		t.Fatalf("unexpected erc1155 fields: %+v", txs[0])
	}
}

func TestClientHandlesEtherscanError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fixture(t, "error_status.json")))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "dummy-key", srv.Client())
	_, err := client.GetNormalTx(context.Background(), "0xabc", 1, 10, "asc")
	if err == nil {
		t.Fatalf("expected error on status=0")
	}
	if !strings.Contains(err.Error(), "NOTOK") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClientHandlesHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "dummy-key", srv.Client())
	_, err := client.GetNormalTx(context.Background(), "0xabc", 1, 10, "asc")
	if err == nil {
		t.Fatalf("expected http error")
	}
}
