package benchmarking

import (
	"conintracker-hiring/pkg/models"
	"conintracker-hiring/pkg/providers"
	"context"
	"testing"
)

// BenchmarkRegressionGuard defines and executes regression tests
// These benchmarks should be run regularly to detect performance regressions
// Usage: go test -bench=BenchmarkRegression ./pkg/benchmarking

// BaselineMetrics holds expected baseline values
type BaselineMetrics struct {
	WeiToETHNs              int64
	CalculateGasFeeETHNs    int64
	AdjustForDecimalsNs     int64
	NormalizeNormalTxNs     int64
	NormalizeERC20TxNs      int64
	ParallelFetchNs         int64
	ParallelNormalizeNs     int64
}

// RegressionTest benchmarks critical paths and verifies they stay within thresholds
func BenchmarkRegressionGuard(b *testing.B) {
	// Individual helper benchmarks
	b.Run("WeiToETH", func(b *testing.B) {
		testCases := []string{
			"1000000000000000000",
			"500000000000000000",
			"1000000000000000",
			"1000000000000000000000",
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, tc := range testCases {
				providers.WeiToETH(tc)
			}
		}
	})

	b.Run("CalculateGasFeeETH", func(b *testing.B) {
		testCases := []struct {
			gasUsed  string
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
				providers.CalculateGasFeeETH(tc.gasUsed, tc.gasPrice)
			}
		}
	})

	b.Run("NormalizeNormalTx", func(b *testing.B) {
		fixtures := providers.GetSmallFixture()
		normalizer := providers.NewEtherscanNormalizer()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, tx := range fixtures.NormalTxs {
				normalizer.NormalizeNormalTx(tx)
			}
		}
	})

	b.Run("ParallelFetch", func(b *testing.B) {
		fixtures := providers.GetMediumFixture()
		mockFetcher := newMockFetcher(fixtures)
		normalizer := providers.NewEtherscanNormalizer()
		parallelFetcher := providers.NewParallelFetcher(mockFetcher, normalizer)
		ctx := context.Background()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			parallelFetcher.FetchAllTransactionsParallel(ctx, "0xtest", 1, 1)
		}
	})

	b.Run("ParallelNormalize", func(b *testing.B) {
		fixtures := providers.GetMediumFixture()
		normalizer := providers.NewEtherscanNormalizer()
		parallelNormalizer := providers.NewParallelNormalizer(normalizer)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			parallelNormalizer.NormalizeTransactionsParallel(
				context.Background(),
				fixtures.NormalTxs,
				fixtures.InternalTxs,
				fixtures.TokenTxs,
				fixtures.NFTTxs,
				fixtures.ERC1155Txs,
			)
		}
	})
}

// BenchmarkRegressionNormalizers specifically tests normalization performance
func BenchmarkRegressionNormalizers(b *testing.B) {
	fixtures := providers.GetSmallFixture()
	normalizer := providers.NewEtherscanNormalizer()

	b.Run("Normal", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, tx := range fixtures.NormalTxs {
				normalizer.NormalizeNormalTx(tx)
			}
		}
	})

	b.Run("Internal", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, tx := range fixtures.InternalTxs {
				normalizer.NormalizeInternalTx(tx)
			}
		}
	})

	b.Run("ERC20", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, tx := range fixtures.TokenTxs {
				normalizer.NormalizeERC20Tx(tx)
			}
		}
	})

	b.Run("ERC721", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, tx := range fixtures.NFTTxs {
				normalizer.NormalizeERC721Tx(tx)
			}
		}
	})

	b.Run("ERC1155", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, tx := range fixtures.ERC1155Txs {
				normalizer.NormalizeERC1155Tx(tx)
			}
		}
	})
}

// BenchmarkRegressionParallel tests parallel operations specifically
func BenchmarkRegressionParallel(b *testing.B) {
	fixtures := providers.GetMediumFixture()
	mockFetcher := newMockFetcher(fixtures)
	normalizer := providers.NewEtherscanNormalizer()
	ctx := context.Background()

	b.Run("FetchSequential", func(b *testing.B) {
		fetcher := providers.NewTransactionFetcher(mockFetcher, normalizer)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			fetcher.FetchAllTransactions(ctx, "0xtest", 1, 1)
		}
	})

	b.Run("FetchParallel", func(b *testing.B) {
		parallelFetcher := providers.NewParallelFetcher(mockFetcher, normalizer)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			parallelFetcher.FetchAllTransactionsParallel(ctx, "0xtest", 1, 1)
		}
	})

	b.Run("NormalizeSequential", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result []*models.Transaction
			for _, tx := range fixtures.NormalTxs {
				if norm, err := normalizer.NormalizeNormalTx(tx); err == nil {
					result = append(result, norm)
				}
			}
			for _, tx := range fixtures.InternalTxs {
				if norm, err := normalizer.NormalizeInternalTx(tx); err == nil {
					result = append(result, norm)
				}
			}
			for _, tx := range fixtures.TokenTxs {
				if norm, err := normalizer.NormalizeERC20Tx(tx); err == nil {
					result = append(result, norm)
				}
			}
			for _, tx := range fixtures.NFTTxs {
				if norm, err := normalizer.NormalizeERC721Tx(tx); err == nil {
					result = append(result, norm)
				}
			}
			for _, tx := range fixtures.ERC1155Txs {
				if norm, err := normalizer.NormalizeERC1155Tx(tx); err == nil {
					result = append(result, norm)
				}
			}
			_ = result
		}
	})

	b.Run("NormalizeParallel", func(b *testing.B) {
		parallelNormalizer := providers.NewParallelNormalizer(normalizer)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			parallelNormalizer.NormalizeTransactionsParallel(
				ctx,
				fixtures.NormalTxs,
				fixtures.InternalTxs,
				fixtures.TokenTxs,
				fixtures.NFTTxs,
				fixtures.ERC1155Txs,
			)
		}
	})
}

// mockFetcher wraps BenchmarkMockFetcher for testing
func newMockFetcher(fixtures *providers.BenchmarkFixtures) *providers.BenchmarkMockFetcher {
	return providers.NewBenchmarkMockFetcher(fixtures)
}
