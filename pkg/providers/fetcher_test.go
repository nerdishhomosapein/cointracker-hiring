package providers

import (
	"conintracker-hiring/pkg/models"
	"context"
	"testing"
)

// MockProvider implements Provider interface for testing
type MockProvider struct {
	normalTxs     []EtherscanNormalTx
	internalTxs   []EtherscanInternalTx
	tokenTxs      []EtherscanTokenTx
	nftTxs        []EtherscanTokenTx
	erc1155Txs    []EtherscanTokenTx
	shouldError   bool
}

func (mp *MockProvider) FetchNormalTransactions(ctx context.Context, address string, startPage, endPage int) ([]EtherscanNormalTx, error) {
	if mp.shouldError {
		return nil, errMock
	}
	return mp.normalTxs, nil
}

func (mp *MockProvider) FetchInternalTransactions(ctx context.Context, address string, startPage, endPage int) ([]EtherscanInternalTx, error) {
	if mp.shouldError {
		return nil, errMock
	}
	return mp.internalTxs, nil
}

func (mp *MockProvider) FetchTokenTransfers(ctx context.Context, address string, startPage, endPage int) ([]EtherscanTokenTx, error) {
	if mp.shouldError {
		return nil, errMock
	}
	return mp.tokenTxs, nil
}

func (mp *MockProvider) FetchNFTTransfers(ctx context.Context, address string, startPage, endPage int) ([]EtherscanTokenTx, error) {
	if mp.shouldError {
		return nil, errMock
	}
	return mp.nftTxs, nil
}

func (mp *MockProvider) FetchERC1155Transfers(ctx context.Context, address string, startPage, endPage int) ([]EtherscanTokenTx, error) {
	if mp.shouldError {
		return nil, errMock
	}
	return mp.erc1155Txs, nil
}

var errMock = testError("mock error")

type testError string

func (e testError) Error() string {
	return string(e)
}

func TestFetchAllTransactions(t *testing.T) {
	mockProvider := &MockProvider{
		normalTxs: []EtherscanNormalTx{
			{
				Hash:     "0x1234",
				From:     "0xfrom",
				To:       "0xto",
				Value:    "1000000000000000000",
				GasUsed:  "21000",
				GasPrice: "50000000000",
				BlockNumber: "100",
				TimeStamp: "1000",
			},
		},
		internalTxs: []EtherscanInternalTx{
			{
				Hash:     "0x5678",
				From:     "0xfrom",
				To:       "0xto",
				Value:    "500000000000000000",
				GasUsed:  "40000",
				BlockNumber: "99",
				TimeStamp: "999",
			},
		},
		tokenTxs: []EtherscanTokenTx{
			{
				Hash:            "0x9012",
				From:            "0xfrom",
				To:              "0xto",
				ContractAddress: "0xtoken",
				TokenSymbol:     "USDC",
				TokenDecimal:    "6",
				Value:           "1000000000",
				GasUsed:         "80000",
				GasPrice:        "55000000000",
				BlockNumber:     "98",
				TimeStamp:       "998",
			},
		},
	}

	normalizer := NewEtherscanNormalizer()
	fetcher := NewTransactionFetcher(mockProvider, normalizer)

	txs, err := fetcher.FetchAllTransactions(context.Background(), "0xtest", 1, 1)
	if err != nil {
		t.Fatalf("FetchAllTransactions() error = %v", err)
	}

	if len(txs) != 3 {
		t.Errorf("Expected 3 transactions, got %d", len(txs))
	}

	// Verify sorting by block number (ascending)
	if txs[0].BlockNumber != 98 || txs[1].BlockNumber != 99 || txs[2].BlockNumber != 100 {
		t.Errorf("Transactions not sorted correctly by block number")
	}

	// Verify types are correct
	typeCount := map[models.TransactionType]int{}
	for _, tx := range txs {
		typeCount[tx.Type]++
	}

	if typeCount[models.TypeEthTransfer] != 1 {
		t.Errorf("Expected 1 ETH transfer, got %d", typeCount[models.TypeEthTransfer])
	}
	if typeCount[models.TypeInternal] != 1 {
		t.Errorf("Expected 1 internal transfer, got %d", typeCount[models.TypeInternal])
	}
	if typeCount[models.TypeERC20Transfer] != 1 {
		t.Errorf("Expected 1 ERC-20 transfer, got %d", typeCount[models.TypeERC20Transfer])
	}
}

func TestFetchAllTransactionsWithError(t *testing.T) {
	mockProvider := &MockProvider{
		shouldError: true,
	}

	normalizer := NewEtherscanNormalizer()
	fetcher := NewTransactionFetcher(mockProvider, normalizer)

	_, err := fetcher.FetchAllTransactions(context.Background(), "0xtest", 1, 1)
	if err == nil {
		t.Error("Expected error, got none")
	}
}

func TestFetchAllTransactionsEmpty(t *testing.T) {
	mockProvider := &MockProvider{
		normalTxs:   []EtherscanNormalTx{},
		internalTxs: []EtherscanInternalTx{},
		tokenTxs:    []EtherscanTokenTx{},
		nftTxs:      []EtherscanTokenTx{},
		erc1155Txs:  []EtherscanTokenTx{},
	}

	normalizer := NewEtherscanNormalizer()
	fetcher := NewTransactionFetcher(mockProvider, normalizer)

	txs, err := fetcher.FetchAllTransactions(context.Background(), "0xtest", 1, 1)
	if err != nil {
		t.Fatalf("FetchAllTransactions() error = %v", err)
	}

	if len(txs) != 0 {
		t.Errorf("Expected 0 transactions, got %d", len(txs))
	}
}

func TestFetchAllTransactionsMixedTypes(t *testing.T) {
	mockProvider := &MockProvider{
		normalTxs: []EtherscanNormalTx{
			{
				Hash:        "0x1",
				From:        "0xfrom",
				To:          "0xto",
				Value:       "1000000000000000000",
				GasUsed:     "21000",
				GasPrice:    "50000000000",
				BlockNumber: "100",
				TimeStamp:   "1000",
			},
		},
		nftTxs: []EtherscanTokenTx{
			{
				Hash:            "0x2",
				From:            "0xfrom",
				To:              "0xto",
				ContractAddress: "0xnft",
				TokenSymbol:     "BAYC",
				TokenID:         "1337",
				GasUsed:         "150000",
				GasPrice:        "60000000000",
				BlockNumber:     "101",
				TimeStamp:       "1001",
			},
		},
		erc1155Txs: []EtherscanTokenTx{
			{
				Hash:            "0x3",
				From:            "0xfrom",
				To:              "0xto",
				ContractAddress: "0xerc1155",
				TokenSymbol:     "POLY",
				TokenID:         "999",
				TokenValue:      "50",
				GasUsed:         "150000",
				GasPrice:        "65000000000",
				BlockNumber:     "102",
				TimeStamp:       "1002",
			},
		},
	}

	normalizer := NewEtherscanNormalizer()
	fetcher := NewTransactionFetcher(mockProvider, normalizer)

	txs, err := fetcher.FetchAllTransactions(context.Background(), "0xtest", 1, 1)
	if err != nil {
		t.Fatalf("FetchAllTransactions() error = %v", err)
	}

	if len(txs) != 3 {
		t.Errorf("Expected 3 transactions, got %d", len(txs))
	}

	// Check types
	if txs[0].Type != models.TypeEthTransfer {
		t.Errorf("First should be ETH transfer, got %s", txs[0].Type)
	}
	if txs[1].Type != models.TypeERC721Transfer {
		t.Errorf("Second should be ERC-721, got %s", txs[1].Type)
	}
	if txs[2].Type != models.TypeERC1155Transfer {
		t.Errorf("Third should be ERC-1155, got %s", txs[2].Type)
	}

	// Check NFT fields
	if txs[1].TokenID != "1337" {
		t.Errorf("NFT TokenID mismatch")
	}
	if txs[2].Amount != "50" {
		t.Errorf("ERC-1155 Amount mismatch, expected 50 got %s", txs[2].Amount)
	}
}
