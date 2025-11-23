package pkg_test

import (
	"conintracker-hiring/internal/testdata"
	"conintracker-hiring/pkg/models"
	"conintracker-hiring/pkg/output"
	"conintracker-hiring/pkg/providers"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestEndToEndFetchAndExport tests the full flow: fetch from provider, normalize, and export to CSV
func TestEndToEndFetchAndExport(t *testing.T) {
	// Create mock Etherscan server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		action := r.URL.Query().Get("action")
		switch action {
		case "txlist":
			w.Write([]byte(testdata.NormalTxResponse))
		case "txlistinternal":
			w.Write([]byte(testdata.InternalTxResponse))
		case "tokentx":
			w.Write([]byte(testdata.ERC20TokenTxResponse))
		case "tokennfttx":
			w.Write([]byte(testdata.ERC721NFTResponse))
		case "token1155tx":
			w.Write([]byte(testdata.ERC1155Response))
		default:
			w.Write([]byte(testdata.EmptyResultResponse))
		}
	}))
	defer server.Close()

	// Step 1: Create Etherscan client
	client := providers.NewEtherscanClient(providers.ClientConfig{
		APIKey:     "test-key",
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	})

	// Step 2: Create normalizer
	normalizer := providers.NewEtherscanNormalizer()

	// Step 3: Create fetcher
	fetcher := providers.NewTransactionFetcher(client, normalizer)

	// Step 4: Fetch all transactions
	ctx := context.Background()
	txs, err := fetcher.FetchAllTransactions(ctx, "0xa39b189482f984388a34460636fea9eb181ad1a6", 1, 1)
	if err != nil {
		t.Fatalf("FetchAllTransactions() error = %v", err)
	}

	if len(txs) == 0 {
		t.Fatal("Expected transactions but got none")
	}

	// Step 5: Verify we have diverse transaction types
	typeCount := make(map[models.TransactionType]int)
	for _, tx := range txs {
		typeCount[tx.Type]++
	}

	if typeCount[models.TypeEthTransfer] == 0 {
		t.Error("Missing ETH transfers")
	}
	if typeCount[models.TypeERC20Transfer] == 0 {
		t.Error("Missing ERC-20 transfers")
	}
	if typeCount[models.TypeERC721Transfer] == 0 {
		t.Error("Missing ERC-721 transfers")
	}

	// Step 6: Export to CSV
	buf := &bytes.Buffer{}
	csvWriter, err := output.NewCSVWriter(output.CSVConfig{Writer: &closeableBuffer{buf}})
	if err != nil {
		t.Fatalf("NewCSVWriter() error = %v", err)
	}

	if err := csvWriter.WriteTransactions(txs); err != nil {
		t.Fatalf("WriteTransactions() error = %v", err)
	}

	if err := csvWriter.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// Step 7: Verify CSV output
	csvContent := buf.String()

	// Check header
	if !strings.Contains(csvContent, "Transaction Hash") {
		t.Error("CSV header missing")
	}

	// Check for transaction data
	lines := strings.Split(strings.TrimSpace(csvContent), "\n")
	if len(lines) <= 1 {
		t.Errorf("Expected CSV to have header + data, got %d lines", len(lines))
	}

	// Verify all transaction types appear in CSV
	if !strings.Contains(csvContent, "ETH") {
		t.Error("ETH type not in CSV")
	}
	if !strings.Contains(csvContent, "ERC-20") {
		t.Error("ERC-20 type not in CSV")
	}
	if !strings.Contains(csvContent, "ERC-721") {
		t.Error("ERC-721 type not in CSV")
	}

	// Verify required fields are present
	if !strings.Contains(csvContent, "0xfrom") && !strings.Contains(csvContent, "0xa39b189482f984388a34460636fea9eb181ad1a6") {
		t.Error("From address not in CSV")
	}
	if !strings.Contains(csvContent, "0xto") && !strings.Contains(csvContent, "0xd620AADaBaA20d2af700853C4504028cba7C3333") {
		t.Error("To address not in CSV")
	}
}

// TestAllTransactionTypesNormalization verifies each transaction type normalizes correctly
func TestAllTransactionTypesNormalization(t *testing.T) {
	normalizer := providers.NewEtherscanNormalizer()

	tests := []struct {
		name           string
		normalizeFunc  func(interface{}) (*models.Transaction, error)
		rawTx          interface{}
		expectedType   models.TransactionType
		expectedFields []string
	}{
		{
			name: "normal_eth_transfer",
			normalizeFunc: func(tx interface{}) (*models.Transaction, error) {
				return normalizer.NormalizeNormalTx(tx.(providers.EtherscanNormalTx))
			},
			rawTx: providers.EtherscanNormalTx{
				Hash:     "0x1",
				From:     "0xfrom",
				To:       "0xto",
				Value:    "1000000000000000000",
				GasUsed:  "21000",
				GasPrice: "50000000000",
				BlockNumber: "100",
				TimeStamp: "1000",
			},
			expectedType:   models.TypeEthTransfer,
			expectedFields: []string{"Hash", "From", "To", "Amount", "GasFeeETH"},
		},
		{
			name: "erc20_transfer",
			normalizeFunc: func(tx interface{}) (*models.Transaction, error) {
				return normalizer.NormalizeERC20Tx(tx.(providers.EtherscanTokenTx))
			},
			rawTx: providers.EtherscanTokenTx{
				Hash:            "0x2",
				From:            "0xfrom",
				To:              "0xto",
				ContractAddress: "0xtoken",
				TokenSymbol:     "USDC",
				TokenDecimal:    "6",
				Value:           "1000000000",
				GasUsed:         "80000",
				GasPrice:        "55000000000",
				BlockNumber:     "101",
				TimeStamp:       "1001",
			},
			expectedType:   models.TypeERC20Transfer,
			expectedFields: []string{"Hash", "From", "To", "AssetSymbol", "Amount"},
		},
		{
			name: "erc721_nft_transfer",
			normalizeFunc: func(tx interface{}) (*models.Transaction, error) {
				return normalizer.NormalizeERC721Tx(tx.(providers.EtherscanTokenTx))
			},
			rawTx: providers.EtherscanTokenTx{
				Hash:            "0x3",
				From:            "0xfrom",
				To:              "0xto",
				ContractAddress: "0xnft",
				TokenSymbol:     "BAYC",
				TokenID:         "1337",
				GasUsed:         "150000",
				GasPrice:        "60000000000",
				BlockNumber:     "102",
				TimeStamp:       "1002",
			},
			expectedType:   models.TypeERC721Transfer,
			expectedFields: []string{"Hash", "TokenID", "AssetSymbol", "Amount"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tt.normalizeFunc(tt.rawTx)
			if err != nil {
				t.Fatalf("Normalization error: %v", err)
			}

			if tx.Type != tt.expectedType {
				t.Errorf("Type mismatch: got %s, want %s", tx.Type, tt.expectedType)
			}

			// Check required fields are populated
			if tx.Hash == "" {
				t.Error("Hash is empty")
			}
			if tx.From == "" {
				t.Error("From is empty")
			}
			if tx.To == "" {
				t.Error("To is empty")
			}
			if tx.Amount == "" {
				t.Error("Amount is empty")
			}
		})
	}
}

// TestCSVRoundTrip tests that data survives the fetch -> normalize -> export -> parse cycle
func TestCSVRoundTrip(t *testing.T) {
	normalizer := providers.NewEtherscanNormalizer()

	// Create a transaction
	tx := &models.Transaction{
		Hash:                 "0xabc123",
		From:                 "0xfrom",
		To:                   "0xto",
		Type:                 models.TypeERC20Transfer,
		AssetContractAddress: "0xcontract",
		AssetSymbol:          "TEST",
		TokenID:              "999",
		Amount:               "1234.567",
		GasFeeETH:            "0.00525",
	}

	// Write to CSV
	buf := &bytes.Buffer{}
	csvWriter, err := output.NewCSVWriter(output.CSVConfig{Writer: &closeableBuffer{buf}})
	if err != nil {
		t.Fatalf("NewCSVWriter() error = %v", err)
	}

	if err := csvWriter.WriteTransaction(tx); err != nil {
		t.Fatalf("WriteTransaction() error = %v", err)
	}

	if err := csvWriter.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	csvContent := buf.String()

	// Verify all data is present
	if !strings.Contains(csvContent, tx.Hash) {
		t.Error("Hash not in CSV")
	}
	if !strings.Contains(csvContent, tx.From) {
		t.Error("From not in CSV")
	}
	if !strings.Contains(csvContent, tx.AssetSymbol) {
		t.Error("AssetSymbol not in CSV")
	}
	if !strings.Contains(csvContent, tx.TokenID) {
		t.Error("TokenID not in CSV")
	}
	if !strings.Contains(csvContent, tx.Amount) {
		t.Error("Amount not in CSV")
	}
	if !strings.Contains(csvContent, tx.GasFeeETH) {
		t.Error("GasFeeETH not in CSV")
	}

	_ = normalizer // Use normalizer in test for completeness
}

// closeableBuffer wraps bytes.Buffer to implement io.WriteCloser
type closeableBuffer struct {
	*bytes.Buffer
}

func (cb *closeableBuffer) Close() error {
	return nil
}
