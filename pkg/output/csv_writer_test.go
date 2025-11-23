package output

import (
	"conintracker-hiring/pkg/models"
	"bytes"
	"strings"
	"testing"
	"time"
)

type WriteCloserBuffer struct {
	*bytes.Buffer
}

func (wcb *WriteCloserBuffer) Close() error {
	return nil
}

func TestNewCSVWriter(t *testing.T) {
	buf := &WriteCloserBuffer{Buffer: &bytes.Buffer{}}
	writer, err := NewCSVWriter(CSVConfig{Writer: buf})
	if err != nil {
		t.Fatalf("NewCSVWriter() error = %v", err)
	}
	
	// Close to flush the header
	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// Check that header was written
	content := buf.String()
	if !strings.Contains(content, "Transaction Hash") {
		t.Errorf("Header not written correctly. Content: %s", content)
	}
}

func TestWriteTransaction(t *testing.T) {
	buf := &WriteCloserBuffer{Buffer: &bytes.Buffer{}}
	writer, err := NewCSVWriter(CSVConfig{Writer: buf})
	if err != nil {
		t.Fatalf("NewCSVWriter() error = %v", err)
	}

	tx := &models.Transaction{
		Hash:      "0x1234",
		Timestamp: time.Unix(1700000000, 0),
		From:      "0xfrom",
		To:        "0xto",
		Type:      models.TypeEthTransfer,
		Amount:    "1.5",
		GasFeeETH: "0.001",
	}

	if err := writer.WriteTransaction(tx); err != nil {
		t.Fatalf("WriteTransaction() error = %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	content := buf.String()
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines (header + transaction), got %d", len(lines))
	}

	// Check that transaction data is present
	if !strings.Contains(content, "0x1234") {
		t.Error("Transaction hash not found in CSV")
	}
	if !strings.Contains(content, "ETH") {
		t.Error("Transaction type not found in CSV")
	}
}

func TestWriteMultipleTransactions(t *testing.T) {
	buf := &WriteCloserBuffer{Buffer: &bytes.Buffer{}}
	writer, err := NewCSVWriter(CSVConfig{Writer: buf})
	if err != nil {
		t.Fatalf("NewCSVWriter() error = %v", err)
	}

	txs := []*models.Transaction{
		{
			Hash:      "0x1111",
			Timestamp: time.Unix(1700000000, 0),
			From:      "0xfrom1",
			To:        "0xto1",
			Type:      models.TypeEthTransfer,
			Amount:    "1.0",
			GasFeeETH: "0.001",
		},
		{
			Hash:      "0x2222",
			Timestamp: time.Unix(1700000001, 0),
			From:      "0xfrom2",
			To:        "0xto2",
			Type:      models.TypeERC20Transfer,
			AssetSymbol: "USDC",
			Amount:    "100.0",
			GasFeeETH: "0.002",
		},
		{
			Hash:      "0x3333",
			Timestamp: time.Unix(1700000002, 0),
			From:      "0xfrom3",
			To:        "0xto3",
			Type:      models.TypeERC721Transfer,
			TokenID:   "1337",
			AssetSymbol: "BAYC",
			Amount:    "1",
			GasFeeETH: "0.003",
		},
	}

	if err := writer.WriteTransactions(txs); err != nil {
		t.Fatalf("WriteTransactions() error = %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	content := buf.String()
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines (header + 3 transactions), got %d", len(lines))
	}

	// Verify all transactions are present
	if !strings.Contains(content, "0x1111") {
		t.Error("First transaction hash not found")
	}
	if !strings.Contains(content, "0x2222") {
		t.Error("Second transaction hash not found")
	}
	if !strings.Contains(content, "0x3333") {
		t.Error("Third transaction hash not found")
	}

	// Verify types
	if !strings.Contains(content, "ERC-20") {
		t.Error("ERC-20 type not found")
	}
	if !strings.Contains(content, "ERC-721") {
		t.Error("ERC-721 type not found")
	}

	// Verify NFT token ID
	if !strings.Contains(content, "1337") {
		t.Error("Token ID 1337 not found")
	}
}

func TestCSVFormatting(t *testing.T) {
	buf := &WriteCloserBuffer{Buffer: &bytes.Buffer{}}
	writer, err := NewCSVWriter(CSVConfig{Writer: buf})
	if err != nil {
		t.Fatalf("NewCSVWriter() error = %v", err)
	}

	tx := &models.Transaction{
		Hash:                 "0xabc123",
		Timestamp:            time.Date(2023, 11, 15, 10, 30, 45, 0, time.UTC),
		From:                 "0xfrom",
		To:                   "0xto",
		Type:                 models.TypeERC20Transfer,
		AssetContractAddress: "0xcontract",
		AssetSymbol:          "USDC",
		Amount:               "1234.567",
		GasFeeETH:            "0.00525",
	}

	if err := writer.WriteTransaction(tx); err != nil {
		t.Fatalf("WriteTransaction() error = %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	content := buf.String()

	// Check for proper RFC3339 timestamp format
	if !strings.Contains(content, "2023-11-15T10:30:45Z") {
		t.Errorf("Timestamp not in RFC3339 format. Content: %s", content)
	}

	// Check for all fields
	if !strings.Contains(content, "0xcontract") {
		t.Error("Asset contract address not found")
	}
	if !strings.Contains(content, "USDC") {
		t.Error("Token symbol not found")
	}
	if !strings.Contains(content, "1234.567") {
		t.Error("Amount not found")
	}
	if !strings.Contains(content, "0.00525") {
		t.Error("Gas fee not found")
	}
}

func TestCSVWithSpecialCharacters(t *testing.T) {
	buf := &WriteCloserBuffer{Buffer: &bytes.Buffer{}}
	writer, err := NewCSVWriter(CSVConfig{Writer: buf})
	if err != nil {
		t.Fatalf("NewCSVWriter() error = %v", err)
	}

	tx := &models.Transaction{
		Hash:      "0x1234",
		Timestamp: time.Unix(1700000000, 0),
		From:      "0xfrom",
		To:        "0xto",
		Type:      models.TypeEthTransfer,
		AssetSymbol: "TEST,SYMBOL", // Contains comma
		Amount:    "1.0",
		GasFeeETH: "0.001",
	}

	if err := writer.WriteTransaction(tx); err != nil {
		t.Fatalf("WriteTransaction() error = %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// CSV writer should properly escape the comma in the symbol
	content := buf.String()
	if !strings.Contains(content, `"TEST,SYMBOL"`) {
		t.Error("Special characters not properly escaped in CSV")
	}
}

func TestEmptyTransactions(t *testing.T) {
	buf := &WriteCloserBuffer{Buffer: &bytes.Buffer{}}
	writer, err := NewCSVWriter(CSVConfig{Writer: buf})
	if err != nil {
		t.Fatalf("NewCSVWriter() error = %v", err)
	}

	if err := writer.WriteTransactions([]*models.Transaction{}); err != nil {
		t.Fatalf("WriteTransactions() error = %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	content := buf.String()
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) != 1 {
		t.Errorf("Expected only header line, got %d lines", len(lines))
	}
}
