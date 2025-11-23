package providers

import (
	"context"
	"testing"
)

// BenchmarkWeiToETH benchmarks the wei to ETH conversion
func BenchmarkWeiToETH(b *testing.B) {
	testCases := []string{
		"1000000000000000000",    // 1 ETH
		"500000000000000000",     // 0.5 ETH
		"1000000000000000",       // 0.001 ETH
		"1000000000000000000000", // 1000 ETH
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			weiToETH(tc)
		}
	}
}

// BenchmarkCalculateGasFeeETH benchmarks gas fee calculation
func BenchmarkCalculateGasFeeETH(b *testing.B) {
	testCases := []struct {
		gasUsed string
		gasPrice string
	}{
		{"21000", "20000000000"},
		{"65000", "30000000000"},
		{"150000", "50000000000"},
		{"200000", "100000000000"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			calculateGasFeeETH(tc.gasUsed, tc.gasPrice)
		}
	}
}

// BenchmarkAdjustForDecimals benchmarks token decimal adjustment
func BenchmarkAdjustForDecimals(b *testing.B) {
	testCases := []struct {
		value    string
		decimals int
	}{
		{"1000000000000000000", 18}, // USDC-like (18 decimals)
		{"1000000", 6},               // USDC (6 decimals)
		{"1000", 8},                  // Other token (8 decimals)
		{"1000000000000000000000", 18}, // Large value
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			adjustForDecimals(tc.value, tc.decimals)
		}
	}
}

// BenchmarkParseUint64 benchmarks uint64 parsing
func BenchmarkParseUint64(b *testing.B) {
	testCases := []string{
		"19000000",
		"1700000000",
		"21000",
		"18446744073709551615", // max uint64
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			parseUint64(tc)
		}
	}
}

// BenchmarkParseTimestamp benchmarks timestamp parsing
func BenchmarkParseTimestamp(b *testing.B) {
	testCases := []string{
		"1700000000",
		"1600000000",
		"1500000000",
		"1400000000",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			parseTimestamp(tc)
		}
	}
}

// BenchmarkNormalizeNormalTx benchmarks normal transaction normalization
func BenchmarkNormalizeNormalTx(b *testing.B) {
	fixtures := GetSmallFixture()
	normalizer := NewEtherscanNormalizer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tx := range fixtures.NormalTxs {
			normalizer.NormalizeNormalTx(tx)
		}
	}
}

// BenchmarkNormalizeInternalTx benchmarks internal transaction normalization
func BenchmarkNormalizeInternalTx(b *testing.B) {
	fixtures := GetSmallFixture()
	normalizer := NewEtherscanNormalizer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tx := range fixtures.InternalTxs {
			normalizer.NormalizeInternalTx(tx)
		}
	}
}

// BenchmarkNormalizeERC20Tx benchmarks ERC-20 token normalization
func BenchmarkNormalizeERC20Tx(b *testing.B) {
	fixtures := GetSmallFixture()
	normalizer := NewEtherscanNormalizer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tx := range fixtures.TokenTxs {
			normalizer.NormalizeERC20Tx(tx)
		}
	}
}

// BenchmarkNormalizeERC721Tx benchmarks ERC-721 NFT normalization
func BenchmarkNormalizeERC721Tx(b *testing.B) {
	fixtures := GetSmallFixture()
	normalizer := NewEtherscanNormalizer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tx := range fixtures.NFTTxs {
			normalizer.NormalizeERC721Tx(tx)
		}
	}
}

// BenchmarkNormalizeERC1155Tx benchmarks ERC-1155 token normalization
func BenchmarkNormalizeERC1155Tx(b *testing.B) {
	fixtures := GetSmallFixture()
	normalizer := NewEtherscanNormalizer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tx := range fixtures.ERC1155Txs {
			normalizer.NormalizeERC1155Tx(tx)
		}
	}
}

// BenchmarkNormalizationPipeline benchmarks the full normalization pipeline
func BenchmarkNormalizationPipeline(b *testing.B) {
	fixtures := GetMediumFixture()
	normalizer := NewEtherscanNormalizer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Normalize all transaction types
		for _, tx := range fixtures.NormalTxs {
			normalizer.NormalizeNormalTx(tx)
		}
		for _, tx := range fixtures.InternalTxs {
			normalizer.NormalizeInternalTx(tx)
		}
		for _, tx := range fixtures.TokenTxs {
			normalizer.NormalizeERC20Tx(tx)
		}
		for _, tx := range fixtures.NFTTxs {
			normalizer.NormalizeERC721Tx(tx)
		}
		for _, tx := range fixtures.ERC1155Txs {
			normalizer.NormalizeERC1155Tx(tx)
		}
	}
}

// BenchmarkFetchAllTransactions benchmarks the fetch orchestration
func BenchmarkFetchAllTransactions(b *testing.B) {
	fixtures := GetMediumFixture()
	mockFetcher := NewBenchmarkMockFetcher(fixtures)
	normalizer := NewEtherscanNormalizer()
	fetcher := NewTransactionFetcher(mockFetcher, normalizer)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fetcher.FetchAllTransactions(ctx, "0x1234567890123456789012345678901234567890", 1, 1)
	}
}

// BenchmarkParallelFetchAllTransactions benchmarks parallel fetch orchestration
func BenchmarkParallelFetchAllTransactions(b *testing.B) {
	fixtures := GetMediumFixture()
	mockFetcher := NewBenchmarkMockFetcher(fixtures)
	normalizer := NewEtherscanNormalizer()
	parallelFetcher := NewParallelFetcher(mockFetcher, normalizer)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parallelFetcher.FetchAllTransactionsParallel(ctx, "0x1234567890123456789012345678901234567890", 1, 1)
	}
}

// BenchmarkParallelFetchVsSequential compares parallel vs sequential fetch performance
func BenchmarkParallelFetchVsSequential(b *testing.B) {
	fixtures := GetMediumFixture()
	mockFetcher := NewBenchmarkMockFetcher(fixtures)
	normalizer := NewEtherscanNormalizer()
	ctx := context.Background()

	b.Run("Sequential", func(b *testing.B) {
		fetcher := NewTransactionFetcher(mockFetcher, normalizer)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			fetcher.FetchAllTransactions(ctx, "0x1234567890123456789012345678901234567890", 1, 1)
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		parallelFetcher := NewParallelFetcher(mockFetcher, normalizer)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			parallelFetcher.FetchAllTransactionsParallel(ctx, "0x1234567890123456789012345678901234567890", 1, 1)
		}
	})
}
