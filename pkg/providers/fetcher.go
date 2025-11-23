package providers

import (
	"conintracker-hiring/pkg/models"
	"context"
	"fmt"
	"sort"
)

// TransactionFetcher orchestrates fetching and normalizing transactions from a provider
type TransactionFetcher struct {
	provider   Provider
	normalizer Normalizer
}

// FetchResult holds the result of fetching a specific transaction type
type FetchResult struct {
	Transactions []*models.Transaction
	Err          error
}

// NewTransactionFetcher creates a new transaction fetcher
func NewTransactionFetcher(provider Provider, normalizer Normalizer) *TransactionFetcher {
	return &TransactionFetcher{
		provider:   provider,
		normalizer: normalizer,
	}
}

// FetchAllTransactions fetches all transaction types for an address and returns normalized transactions
func (tf *TransactionFetcher) FetchAllTransactions(ctx context.Context, address string, startPage, endPage int) ([]*models.Transaction, error) {
	// Fetch all transaction types sequentially to respect rate limits
	var allTransactions []*models.Transaction

	// Fetch normal transactions
	normalTxs, err := tf.fetchNormalTransactions(ctx, address, startPage, endPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch normal transactions: %w", err)
	}
	if normalTxs != nil {
		allTransactions = append(allTransactions, normalTxs...)
	}

	// Fetch internal transactions
	internalTxs, err := tf.fetchInternalTransactions(ctx, address, startPage, endPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch internal transactions: %w", err)
	}
	if internalTxs != nil {
		allTransactions = append(allTransactions, internalTxs...)
	}

	// Fetch ERC-20 token transfers
	tokenTxs, err := tf.fetchTokenTransfers(ctx, address, startPage, endPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch token transfers: %w", err)
	}
	if tokenTxs != nil {
		allTransactions = append(allTransactions, tokenTxs...)
	}

	// Fetch ERC-721 NFT transfers
	nftTxs, err := tf.fetchNFTTransfers(ctx, address, startPage, endPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch NFT transfers: %w", err)
	}
	if nftTxs != nil {
		allTransactions = append(allTransactions, nftTxs...)
	}

	// Fetch ERC-1155 token transfers
	erc1155Txs, err := tf.fetchERC1155Transfers(ctx, address, startPage, endPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ERC-1155 transfers: %w", err)
	}
	if erc1155Txs != nil {
		allTransactions = append(allTransactions, erc1155Txs...)
	}

	// Sort by block number and timestamp
	sort.Sort(models.TransactionList(allTransactions))

	return allTransactions, nil
}

// fetchNormalTransactions fetches and normalizes normal ETH transfers
func (tf *TransactionFetcher) fetchNormalTransactions(ctx context.Context, address string, startPage, endPage int) ([]*models.Transaction, error) {
	rawTxs, err := tf.provider.FetchNormalTransactions(ctx, address, startPage, endPage)
	if err != nil {
		return nil, err
	}

	var normalized []*models.Transaction
	for _, tx := range rawTxs {
		norm, err := tf.normalizer.NormalizeNormalTx(tx)
		if err != nil {
			// Log and skip invalid transactions
			continue
		}
		normalized = append(normalized, norm)
	}

	return normalized, nil
}

// fetchInternalTransactions fetches and normalizes internal transfers
func (tf *TransactionFetcher) fetchInternalTransactions(ctx context.Context, address string, startPage, endPage int) ([]*models.Transaction, error) {
	rawTxs, err := tf.provider.FetchInternalTransactions(ctx, address, startPage, endPage)
	if err != nil {
		return nil, err
	}

	var normalized []*models.Transaction
	for _, tx := range rawTxs {
		norm, err := tf.normalizer.NormalizeInternalTx(tx)
		if err != nil {
			continue
		}
		normalized = append(normalized, norm)
	}

	return normalized, nil
}

// fetchTokenTransfers fetches and normalizes ERC-20 token transfers
func (tf *TransactionFetcher) fetchTokenTransfers(ctx context.Context, address string, startPage, endPage int) ([]*models.Transaction, error) {
	rawTxs, err := tf.provider.FetchTokenTransfers(ctx, address, startPage, endPage)
	if err != nil {
		return nil, err
	}

	var normalized []*models.Transaction
	for _, tx := range rawTxs {
		norm, err := tf.normalizer.NormalizeERC20Tx(tx)
		if err != nil {
			continue
		}
		normalized = append(normalized, norm)
	}

	return normalized, nil
}

// fetchNFTTransfers fetches and normalizes ERC-721 NFT transfers
func (tf *TransactionFetcher) fetchNFTTransfers(ctx context.Context, address string, startPage, endPage int) ([]*models.Transaction, error) {
	rawTxs, err := tf.provider.FetchNFTTransfers(ctx, address, startPage, endPage)
	if err != nil {
		return nil, err
	}

	var normalized []*models.Transaction
	for _, tx := range rawTxs {
		norm, err := tf.normalizer.NormalizeERC721Tx(tx)
		if err != nil {
			continue
		}
		normalized = append(normalized, norm)
	}

	return normalized, nil
}

// fetchERC1155Transfers fetches and normalizes ERC-1155 multi-token transfers
func (tf *TransactionFetcher) fetchERC1155Transfers(ctx context.Context, address string, startPage, endPage int) ([]*models.Transaction, error) {
	rawTxs, err := tf.provider.FetchERC1155Transfers(ctx, address, startPage, endPage)
	if err != nil {
		return nil, err
	}

	var normalized []*models.Transaction
	for _, tx := range rawTxs {
		norm, err := tf.normalizer.NormalizeERC1155Tx(tx)
		if err != nil {
			continue
		}
		normalized = append(normalized, norm)
	}

	return normalized, nil
}
