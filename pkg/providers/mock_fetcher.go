package providers

import (
	"context"
)

// BenchmarkMockFetcher is a mock implementation for benchmarking without network calls
type BenchmarkMockFetcher struct {
	fixtures *BenchmarkFixtures
}

// NewBenchmarkMockFetcher creates a new mock fetcher with fixtures
func NewBenchmarkMockFetcher(fixtures *BenchmarkFixtures) *BenchmarkMockFetcher {
	return &BenchmarkMockFetcher{fixtures: fixtures}
}

// FetchNormalTransactions returns mock normal transactions
func (b *BenchmarkMockFetcher) FetchNormalTransactions(ctx context.Context, address string, startPage, endPage int) ([]EtherscanNormalTx, error) {
	return b.fixtures.NormalTxs, nil
}

// FetchInternalTransactions returns mock internal transactions
func (b *BenchmarkMockFetcher) FetchInternalTransactions(ctx context.Context, address string, startPage, endPage int) ([]EtherscanInternalTx, error) {
	return b.fixtures.InternalTxs, nil
}

// FetchTokenTransfers returns mock token transfers
func (b *BenchmarkMockFetcher) FetchTokenTransfers(ctx context.Context, address string, startPage, endPage int) ([]EtherscanTokenTx, error) {
	return b.fixtures.TokenTxs, nil
}

// FetchNFTTransfers returns mock NFT transfers
func (b *BenchmarkMockFetcher) FetchNFTTransfers(ctx context.Context, address string, startPage, endPage int) ([]EtherscanTokenTx, error) {
	return b.fixtures.NFTTxs, nil
}

// FetchERC1155Transfers returns mock ERC-1155 transfers
func (b *BenchmarkMockFetcher) FetchERC1155Transfers(ctx context.Context, address string, startPage, endPage int) ([]EtherscanTokenTx, error) {
	return b.fixtures.ERC1155Txs, nil
}
