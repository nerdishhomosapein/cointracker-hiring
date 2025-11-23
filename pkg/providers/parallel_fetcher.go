package providers

import (
	"conintracker-hiring/pkg/models"
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

// ParallelFetcher orchestrates concurrent fetching of different transaction types
// while respecting rate limits and maintaining error handling
type ParallelFetcher struct {
	provider      Provider
	normalizer    Normalizer
	maxConcurrent int // Max concurrent fetch operations (default 3 for Etherscan)
	timeout       time.Duration // Per-fetch timeout
}

// FetchTypeResult holds the result of fetching a specific transaction type
type FetchTypeResult struct {
	TxType TransactionType
	Txs    []*models.Transaction
	Err    error
	Count  int
}

// TransactionType enum for identifying fetch type
type TransactionType int

const (
	TxTypeNormal TransactionType = iota
	TxTypeInternal
	TxTypeToken
	TxTypeNFT
	TxTypeERC1155
)

func (t TransactionType) String() string {
	switch t {
	case TxTypeNormal:
		return "Normal"
	case TxTypeInternal:
		return "Internal"
	case TxTypeToken:
		return "ERC-20"
	case TxTypeNFT:
		return "ERC-721"
	case TxTypeERC1155:
		return "ERC-1155"
	default:
		return "Unknown"
	}
}

// NewParallelFetcher creates a new parallel fetcher with sensible defaults
func NewParallelFetcher(provider Provider, normalizer Normalizer) *ParallelFetcher {
	return &ParallelFetcher{
		provider:      provider,
		normalizer:    normalizer,
		maxConcurrent: 3, // Etherscan allows ~5 req/sec, so 3 concurrent is safe
		timeout:       30 * time.Second,
	}
}

// SetMaxConcurrent sets the maximum number of concurrent fetch operations
func (pf *ParallelFetcher) SetMaxConcurrent(max int) {
	if max > 0 && max <= 10 {
		pf.maxConcurrent = max
	}
}

// SetTimeout sets the timeout for individual fetch operations
func (pf *ParallelFetcher) SetTimeout(timeout time.Duration) {
	if timeout > 0 {
		pf.timeout = timeout
	}
}

// FetchAllTransactionsParallel fetches all transaction types concurrently
func (pf *ParallelFetcher) FetchAllTransactionsParallel(
	ctx context.Context,
	address string,
	startPage, endPage int,
) ([]*models.Transaction, error) {
	// Create a semaphore to limit concurrent operations
	sem := make(chan struct{}, pf.maxConcurrent)
	defer close(sem)

	// Result channel to collect all results
	resultChan := make(chan *FetchTypeResult, 5) // 5 fetch types
	var wg sync.WaitGroup

	// Helper function to wrap fetch operations with semaphore
	fetchWithSemaphore := func(fetchFunc func() (*FetchTypeResult), txType TransactionType) {
		defer wg.Done()

		// Acquire semaphore slot
		sem <- struct{}{}
		defer func() { <-sem }()

		// Create context with timeout
		fetchCtx, cancel := context.WithTimeout(ctx, pf.timeout)
		defer cancel()

		// Execute fetch in goroutine
		resultChan <- pf.executeFetch(fetchCtx, fetchFunc, txType)
	}

	// Launch all fetch operations
	wg.Add(5)
	go fetchWithSemaphore(func() *FetchTypeResult {
		return pf.fetchNormalTransactionsConcurrent(ctx, address, startPage, endPage)
	}, TxTypeNormal)

	go fetchWithSemaphore(func() *FetchTypeResult {
		return pf.fetchInternalTransactionsConcurrent(ctx, address, startPage, endPage)
	}, TxTypeInternal)

	go fetchWithSemaphore(func() *FetchTypeResult {
		return pf.fetchTokenTransfersConcurrent(ctx, address, startPage, endPage)
	}, TxTypeToken)

	go fetchWithSemaphore(func() *FetchTypeResult {
		return pf.fetchNFTTransfersConcurrent(ctx, address, startPage, endPage)
	}, TxTypeNFT)

	go fetchWithSemaphore(func() *FetchTypeResult {
		return pf.fetchERC1155TransfersConcurrent(ctx, address, startPage, endPage)
	}, TxTypeERC1155)

	// Close result channel when all operations complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect all results
	var allTransactions []*models.Transaction
	var errors []error
	var fetchStats map[TransactionType]int = make(map[TransactionType]int)

	for result := range resultChan {
		if result.Err != nil {
			errors = append(errors, fmt.Errorf("%s fetch failed: %w", result.TxType.String(), result.Err))
		} else if result.Txs != nil {
			allTransactions = append(allTransactions, result.Txs...)
			fetchStats[result.TxType] = result.Count
		}
	}

	// If all fetches failed, return error
	if len(errors) == 5 {
		return nil, fmt.Errorf("all transaction fetches failed: %v", errors)
	}

	// Sort all transactions
	if len(allTransactions) > 0 {
		sort.Sort(models.TransactionList(allTransactions))
	}

	return allTransactions, nil
}

// executeFetch safely executes a fetch operation with timeout handling
func (pf *ParallelFetcher) executeFetch(
	ctx context.Context,
	fetchFunc func() *FetchTypeResult,
	txType TransactionType,
) *FetchTypeResult {
	// Use a channel to capture the result with timeout
	resultChan := make(chan *FetchTypeResult, 1)

	go func() {
		resultChan <- fetchFunc()
	}()

	select {
	case result := <-resultChan:
		return result
	case <-ctx.Done():
		return &FetchTypeResult{
			TxType: txType,
			Err:    fmt.Errorf("fetch timeout for %s transactions", txType.String()),
		}
	}
}

// fetchNormalTransactionsConcurrent fetches normal transactions (internal parallel helper)
func (pf *ParallelFetcher) fetchNormalTransactionsConcurrent(
	ctx context.Context,
	address string,
	startPage, endPage int,
) *FetchTypeResult {
	rawTxs, err := pf.provider.FetchNormalTransactions(ctx, address, startPage, endPage)
	if err != nil {
		return &FetchTypeResult{TxType: TxTypeNormal, Err: err}
	}

	var normalized []*models.Transaction
	for _, tx := range rawTxs {
		if norm, err := pf.normalizer.NormalizeNormalTx(tx); err == nil {
			normalized = append(normalized, norm)
		}
	}

	return &FetchTypeResult{
		TxType: TxTypeNormal,
		Txs:    normalized,
		Count:  len(normalized),
	}
}

// fetchInternalTransactionsConcurrent fetches internal transactions
func (pf *ParallelFetcher) fetchInternalTransactionsConcurrent(
	ctx context.Context,
	address string,
	startPage, endPage int,
) *FetchTypeResult {
	rawTxs, err := pf.provider.FetchInternalTransactions(ctx, address, startPage, endPage)
	if err != nil {
		return &FetchTypeResult{TxType: TxTypeInternal, Err: err}
	}

	var normalized []*models.Transaction
	for _, tx := range rawTxs {
		if norm, err := pf.normalizer.NormalizeInternalTx(tx); err == nil {
			normalized = append(normalized, norm)
		}
	}

	return &FetchTypeResult{
		TxType: TxTypeInternal,
		Txs:    normalized,
		Count:  len(normalized),
	}
}

// fetchTokenTransfersConcurrent fetches token transfers
func (pf *ParallelFetcher) fetchTokenTransfersConcurrent(
	ctx context.Context,
	address string,
	startPage, endPage int,
) *FetchTypeResult {
	rawTxs, err := pf.provider.FetchTokenTransfers(ctx, address, startPage, endPage)
	if err != nil {
		return &FetchTypeResult{TxType: TxTypeToken, Err: err}
	}

	var normalized []*models.Transaction
	for _, tx := range rawTxs {
		if norm, err := pf.normalizer.NormalizeERC20Tx(tx); err == nil {
			normalized = append(normalized, norm)
		}
	}

	return &FetchTypeResult{
		TxType: TxTypeToken,
		Txs:    normalized,
		Count:  len(normalized),
	}
}

// fetchNFTTransfersConcurrent fetches NFT transfers
func (pf *ParallelFetcher) fetchNFTTransfersConcurrent(
	ctx context.Context,
	address string,
	startPage, endPage int,
) *FetchTypeResult {
	rawTxs, err := pf.provider.FetchNFTTransfers(ctx, address, startPage, endPage)
	if err != nil {
		return &FetchTypeResult{TxType: TxTypeNFT, Err: err}
	}

	var normalized []*models.Transaction
	for _, tx := range rawTxs {
		if norm, err := pf.normalizer.NormalizeERC721Tx(tx); err == nil {
			normalized = append(normalized, norm)
		}
	}

	return &FetchTypeResult{
		TxType: TxTypeNFT,
		Txs:    normalized,
		Count:  len(normalized),
	}
}

// fetchERC1155TransfersConcurrent fetches ERC-1155 transfers
func (pf *ParallelFetcher) fetchERC1155TransfersConcurrent(
	ctx context.Context,
	address string,
	startPage, endPage int,
) *FetchTypeResult {
	rawTxs, err := pf.provider.FetchERC1155Transfers(ctx, address, startPage, endPage)
	if err != nil {
		return &FetchTypeResult{TxType: TxTypeERC1155, Err: err}
	}

	var normalized []*models.Transaction
	for _, tx := range rawTxs {
		if norm, err := pf.normalizer.NormalizeERC1155Tx(tx); err == nil {
			normalized = append(normalized, norm)
		}
	}

	return &FetchTypeResult{
		TxType: TxTypeERC1155,
		Txs:    normalized,
		Count:  len(normalized),
	}
}
