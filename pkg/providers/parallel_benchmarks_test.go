package providers

import (
	"conintracker-hiring/pkg/models"
	"context"
	"testing"
)

// BenchmarkNormalizeTransactionsParallel benchmarks parallel normalization
func BenchmarkNormalizeTransactionsParallel(b *testing.B) {
	fixtures := GetMediumFixture()
	normalizer := NewEtherscanNormalizer()
	parallelNormalizer := NewParallelNormalizer(normalizer)

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
}

// BenchmarkStreamNormalizeResults benchmarks streaming normalization
func BenchmarkStreamNormalizeResults(b *testing.B) {
	fixtures := GetMediumFixture()
	normalizer := NewEtherscanNormalizer()
	parallelNormalizer := NewParallelNormalizer(normalizer)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		resultChan := parallelNormalizer.StreamNormalizeResults(
			ctx,
			fixtures.NormalTxs,
			fixtures.InternalTxs,
			fixtures.TokenTxs,
			fixtures.NFTTxs,
			fixtures.ERC1155Txs,
		)
		// Drain the channel
		for range resultChan {
		}
	}
}

// BenchmarkNormalizationWorkerCounts benchmarks normalization with different worker counts
func BenchmarkNormalizationWorkerCounts(b *testing.B) {
	fixtures := GetMediumFixture()
	normalizer := NewEtherscanNormalizer()

	for _, workerCount := range []int{1, 2, 4, 8} {
		b.Run("Workers"+string(rune(48+workerCount)), func(b *testing.B) {
			parallelNormalizer := NewParallelNormalizer(normalizer)
			parallelNormalizer.SetWorkerCount(workerCount)

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
}

// BenchmarkParallelVsSequentialNormalization compares performance
func BenchmarkParallelVsSequentialNormalization(b *testing.B) {
	fixtures := GetMediumFixture()
	normalizer := NewEtherscanNormalizer()

	// Sequential normalization
	b.Run("Sequential", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result []*models.Transaction

			// Normalize normal transactions
			for _, tx := range fixtures.NormalTxs {
				if norm, err := normalizer.NormalizeNormalTx(tx); err == nil {
					result = append(result, norm)
				}
			}
			// Normalize internal transactions
			for _, tx := range fixtures.InternalTxs {
				if norm, err := normalizer.NormalizeInternalTx(tx); err == nil {
					result = append(result, norm)
				}
			}
			// Normalize token transfers
			for _, tx := range fixtures.TokenTxs {
				if norm, err := normalizer.NormalizeERC20Tx(tx); err == nil {
					result = append(result, norm)
				}
			}
			// Normalize NFTs
			for _, tx := range fixtures.NFTTxs {
				if norm, err := normalizer.NormalizeERC721Tx(tx); err == nil {
					result = append(result, norm)
				}
			}
			// Normalize ERC-1155
			for _, tx := range fixtures.ERC1155Txs {
				if norm, err := normalizer.NormalizeERC1155Tx(tx); err == nil {
					result = append(result, norm)
				}
			}
			_ = result
		}
	})

	// Parallel normalization
	b.Run("Parallel", func(b *testing.B) {
		parallelNormalizer := NewParallelNormalizer(normalizer)
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
