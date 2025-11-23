package providers

import (
	"conintracker-hiring/pkg/models"
	"context"
	"fmt"
	"sync"
)

// ParallelNormalizer processes multiple transactions concurrently
type ParallelNormalizer struct {
	normalizer  Normalizer
	workerCount int
	bufferSize  int
}

// NormalizationStats tracks statistics about the normalization process
type NormalizationStats struct {
	TotalProcessed int
	SuccessCount   int
	ErrorCount     int
	Errors         []error
}

// NewParallelNormalizer creates a new parallel normalizer
func NewParallelNormalizer(normalizer Normalizer) *ParallelNormalizer {
	return &ParallelNormalizer{
		normalizer:  normalizer,
		workerCount: 4, // Default to 4 workers (CPU-bound)
		bufferSize:  1000,
	}
}

// SetWorkerCount sets the number of normalization workers
func (pn *ParallelNormalizer) SetWorkerCount(count int) {
	if count > 0 && count <= 16 {
		pn.workerCount = count
	}
}

// SetBufferSize sets the size of the result buffer
func (pn *ParallelNormalizer) SetBufferSize(size int) {
	if size > 0 && size <= 10000 {
		pn.bufferSize = size
	}
}

// NormalizationResult holds both successful transactions and error information
type NormalizationResult struct {
	Transactions []*models.Transaction
	Stats        NormalizationStats
}

// normalizeWorkerPoolTyped is a type-safe worker pool using generics
func normalizeWorkerPoolTyped[T any](
	ctx context.Context,
	items []T,
	normalizeFunc func(T) (*models.Transaction, error),
	workerCount int,
	resultChan chan<- *models.Transaction,
	statsChan chan<- NormalizationStats,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	workQueue := make(chan T, len(items))

	// Populate work queue
	go func() {
		defer close(workQueue)
		for _, item := range items {
			select {
			case workQueue <- item:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Spawn worker goroutines
	var workerWg sync.WaitGroup
	var statsMutex sync.Mutex
	stats := NormalizationStats{}

	for i := 0; i < workerCount; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for item := range workQueue {
				select {
				case <-ctx.Done():
					return
				default:
					result, err := normalizeFunc(item)
					
					statsMutex.Lock()
					stats.TotalProcessed++
					if err != nil {
						stats.ErrorCount++
						stats.Errors = append(stats.Errors, fmt.Errorf("normalization failed: %w", err))
					} else if result != nil {
						stats.SuccessCount++
						select {
						case resultChan <- result:
						case <-ctx.Done():
							statsMutex.Unlock()
							return
						}
					}
					statsMutex.Unlock()
				}
			}
		}()
	}

	// Wait for all workers to complete
	workerWg.Wait()
	
	// Send stats
	select {
	case statsChan <- stats:
	case <-ctx.Done():
	}
}

// NormalizeTransactionsParallel normalizes transactions in parallel with error tracking
func (pn *ParallelNormalizer) NormalizeTransactionsParallel(
	ctx context.Context,
	normalTxs []EtherscanNormalTx,
	internalTxs []EtherscanInternalTx,
	tokenTxs []EtherscanTokenTx,
	nftTxs []EtherscanTokenTx,
	erc1155Txs []EtherscanTokenTx,
) *NormalizationResult {
	// Total work items
	totalWork := len(normalTxs) + len(internalTxs) + len(tokenTxs) + len(nftTxs) + len(erc1155Txs)

	// Result channel with buffering
	resultChan := make(chan *models.Transaction, pn.bufferSize)
	statsChan := make(chan NormalizationStats, 5) // 5 transaction types

	// WaitGroup to track goroutine completion
	var wg sync.WaitGroup

	// Process each transaction type with type-safe workers
	if len(normalTxs) > 0 {
		wg.Add(1)
		go normalizeWorkerPoolTyped(ctx, normalTxs, pn.normalizer.NormalizeNormalTx, 
			pn.workerCount, resultChan, statsChan, &wg)
	}

	if len(internalTxs) > 0 {
		wg.Add(1)
		go normalizeWorkerPoolTyped(ctx, internalTxs, pn.normalizer.NormalizeInternalTx, 
			pn.workerCount, resultChan, statsChan, &wg)
	}

	if len(tokenTxs) > 0 {
		wg.Add(1)
		go normalizeWorkerPoolTyped(ctx, tokenTxs, pn.normalizer.NormalizeERC20Tx, 
			pn.workerCount, resultChan, statsChan, &wg)
	}

	if len(nftTxs) > 0 {
		wg.Add(1)
		go normalizeWorkerPoolTyped(ctx, nftTxs, pn.normalizer.NormalizeERC721Tx, 
			pn.workerCount, resultChan, statsChan, &wg)
	}

	if len(erc1155Txs) > 0 {
		wg.Add(1)
		go normalizeWorkerPoolTyped(ctx, erc1155Txs, pn.normalizer.NormalizeERC1155Tx, 
			pn.workerCount, resultChan, statsChan, &wg)
	}

	// Close channels when all workers complete
	go func() {
		wg.Wait()
		close(resultChan)
		close(statsChan)
	}()

	// Collect results and stats
	result := make([]*models.Transaction, 0, totalWork)
	aggregateStats := NormalizationStats{}

	done := false
	for !done {
		select {
		case tx, ok := <-resultChan:
			if !ok {
				resultChan = nil
			} else if tx != nil {
				result = append(result, tx)
			}
		case stats, ok := <-statsChan:
			if !ok {
				statsChan = nil
			} else {
				aggregateStats.TotalProcessed += stats.TotalProcessed
				aggregateStats.SuccessCount += stats.SuccessCount
				aggregateStats.ErrorCount += stats.ErrorCount
				aggregateStats.Errors = append(aggregateStats.Errors, stats.Errors...)
			}
		}
		
		if resultChan == nil && statsChan == nil {
			done = true
		}
	}

	return &NormalizationResult{
		Transactions: result,
		Stats:        aggregateStats,
	}
}


// StreamNormalizeResults returns a channel of normalized transactions for streaming processing
func (pn *ParallelNormalizer) StreamNormalizeResults(
	ctx context.Context,
	normalTxs []EtherscanNormalTx,
	internalTxs []EtherscanInternalTx,
	tokenTxs []EtherscanTokenTx,
	nftTxs []EtherscanTokenTx,
	erc1155Txs []EtherscanTokenTx,
) chan *models.Transaction {
	resultChan := make(chan *models.Transaction, pn.bufferSize)

	go func() {
		defer close(resultChan)

		// Reuse the new type-safe implementation but discard error stats for streaming
		result := pn.NormalizeTransactionsParallel(ctx, normalTxs, internalTxs, tokenTxs, nftTxs, erc1155Txs)
		
		// Stream the results
		for _, tx := range result.Transactions {
			select {
			case resultChan <- tx:
			case <-ctx.Done():
				return
			}
		}
	}()

	return resultChan
}
