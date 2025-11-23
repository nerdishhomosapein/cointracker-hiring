package providers

import (
	"conintracker-hiring/pkg/models"
	"context"
	"sync"
)

// ParallelNormalizer processes multiple transactions concurrently
type ParallelNormalizer struct {
	normalizer  Normalizer
	workerCount int
	bufferSize  int
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

// NormalizeTransactionsParallel normalizes transactions in parallel
func (pn *ParallelNormalizer) NormalizeTransactionsParallel(
	ctx context.Context,
	normalTxs []EtherscanNormalTx,
	internalTxs []EtherscanInternalTx,
	tokenTxs []EtherscanTokenTx,
	nftTxs []EtherscanTokenTx,
	erc1155Txs []EtherscanTokenTx,
) []*models.Transaction {
	// Total work items
	totalWork := len(normalTxs) + len(internalTxs) + len(tokenTxs) + len(nftTxs) + len(erc1155Txs)

	// Result channel with buffering
	resultChan := make(chan *models.Transaction, pn.bufferSize)

	// WaitGroup to track goroutine completion
	var wg sync.WaitGroup

	// Helper function to normalize a slice with worker pool
	normalizeSlice := func(
		items interface{},
		normalizeFunc func(interface{}) *models.Transaction,
		count int,
	) {
		if count == 0 {
			return
		}

		wg.Add(1)
		go pn.normalizeWorkerPool(ctx, items, normalizeFunc, count, resultChan, &wg)
	}

	// Spawn workers for each transaction type
	normalizeSlice(normalTxs, func(item interface{}) *models.Transaction {
		if tx, ok := item.(EtherscanNormalTx); ok {
			result, _ := pn.normalizer.NormalizeNormalTx(tx)
			return result
		}
		return nil
	}, len(normalTxs))

	normalizeSlice(internalTxs, func(item interface{}) *models.Transaction {
		if tx, ok := item.(EtherscanInternalTx); ok {
			result, _ := pn.normalizer.NormalizeInternalTx(tx)
			return result
		}
		return nil
	}, len(internalTxs))

	normalizeSlice(tokenTxs, func(item interface{}) *models.Transaction {
		if tx, ok := item.(EtherscanTokenTx); ok {
			result, _ := pn.normalizer.NormalizeERC20Tx(tx)
			return result
		}
		return nil
	}, len(tokenTxs))

	normalizeSlice(nftTxs, func(item interface{}) *models.Transaction {
		if tx, ok := item.(EtherscanTokenTx); ok {
			result, _ := pn.normalizer.NormalizeERC721Tx(tx)
			return result
		}
		return nil
	}, len(nftTxs))

	normalizeSlice(erc1155Txs, func(item interface{}) *models.Transaction {
		if tx, ok := item.(EtherscanTokenTx); ok {
			result, _ := pn.normalizer.NormalizeERC1155Tx(tx)
			return result
		}
		return nil
	}, len(erc1155Txs))

	// Close result channel when all workers complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	result := make([]*models.Transaction, 0, totalWork)
	for tx := range resultChan {
		if tx != nil {
			result = append(result, tx)
		}
	}

	return result
}

// normalizeWorkerPool processes items with a pool of workers
func (pn *ParallelNormalizer) normalizeWorkerPool(
	ctx context.Context,
	items interface{},
	normalizeFunc func(interface{}) *models.Transaction,
	count int,
	resultChan chan *models.Transaction,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	// Create work queue
	workQueue := make(chan interface{}, count)

	// Populate work queue based on type
	go func() {
		switch v := items.(type) {
		case []EtherscanNormalTx:
			for _, item := range v {
				select {
				case workQueue <- item:
				case <-ctx.Done():
					return
				}
			}
		case []EtherscanInternalTx:
			for _, item := range v {
				select {
				case workQueue <- item:
				case <-ctx.Done():
					return
				}
			}
		case []EtherscanTokenTx:
			for _, item := range v {
				select {
				case workQueue <- item:
				case <-ctx.Done():
					return
				}
			}
		}
		close(workQueue)
	}()

	// Spawn worker goroutines
	var workerWg sync.WaitGroup
	for i := 0; i < pn.workerCount; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for item := range workQueue {
				select {
				case <-ctx.Done():
					return
				default:
					if result := normalizeFunc(item); result != nil {
						select {
						case resultChan <- result:
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}()
	}

	// Wait for all workers to complete
	workerWg.Wait()
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

		var wg sync.WaitGroup

		// Helper to normalize slice and stream results
		streamSlice := func(items interface{}, normalizeFunc func(interface{}) *models.Transaction, count int) {
			if count == 0 {
				return
			}
			wg.Add(1)
			go pn.normalizeWorkerPool(ctx, items, normalizeFunc, count, resultChan, &wg)
		}

		// Spawn workers
		streamSlice(normalTxs, func(item interface{}) *models.Transaction {
			if tx, ok := item.(EtherscanNormalTx); ok {
				result, _ := pn.normalizer.NormalizeNormalTx(tx)
				return result
			}
			return nil
		}, len(normalTxs))

		streamSlice(internalTxs, func(item interface{}) *models.Transaction {
			if tx, ok := item.(EtherscanInternalTx); ok {
				result, _ := pn.normalizer.NormalizeInternalTx(tx)
				return result
			}
			return nil
		}, len(internalTxs))

		streamSlice(tokenTxs, func(item interface{}) *models.Transaction {
			if tx, ok := item.(EtherscanTokenTx); ok {
				result, _ := pn.normalizer.NormalizeERC20Tx(tx)
				return result
			}
			return nil
		}, len(tokenTxs))

		streamSlice(nftTxs, func(item interface{}) *models.Transaction {
			if tx, ok := item.(EtherscanTokenTx); ok {
				result, _ := pn.normalizer.NormalizeERC721Tx(tx)
				return result
			}
			return nil
		}, len(nftTxs))

		streamSlice(erc1155Txs, func(item interface{}) *models.Transaction {
			if tx, ok := item.(EtherscanTokenTx); ok {
				result, _ := pn.normalizer.NormalizeERC1155Tx(tx)
				return result
			}
			return nil
		}, len(erc1155Txs))

		wg.Wait()
	}()

	return resultChan
}
