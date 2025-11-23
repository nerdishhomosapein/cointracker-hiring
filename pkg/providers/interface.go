package providers

import (
	"context"
	"conintracker-hiring/pkg/models"
)

// Provider defines the interface for blockchain data providers
type Provider interface {
	// FetchNormalTransactions fetches normal ETH transfers for an address
	FetchNormalTransactions(ctx context.Context, address string, startPage, endPage int) ([]EtherscanNormalTx, error)
	
	// FetchInternalTransactions fetches internal contract interactions
	FetchInternalTransactions(ctx context.Context, address string, startPage, endPage int) ([]EtherscanInternalTx, error)
	
	// FetchTokenTransfers fetches ERC-20 token transfers
	FetchTokenTransfers(ctx context.Context, address string, startPage, endPage int) ([]EtherscanTokenTx, error)
	
	// FetchNFTTransfers fetches ERC-721 NFT transfers
	FetchNFTTransfers(ctx context.Context, address string, startPage, endPage int) ([]EtherscanTokenTx, error)
	
	// FetchERC1155Transfers fetches ERC-1155 multi-token transfers
	FetchERC1155Transfers(ctx context.Context, address string, startPage, endPage int) ([]EtherscanTokenTx, error)
}

// Normalizer defines the interface for converting provider responses to normalized transactions
type Normalizer interface {
	// NormalizeNormalTx converts Etherscan normal tx to normalized transaction
	NormalizeNormalTx(tx EtherscanNormalTx) (*models.Transaction, error)
	
	// NormalizeInternalTx converts Etherscan internal tx to normalized transaction
	NormalizeInternalTx(tx EtherscanInternalTx) (*models.Transaction, error)
	
	// NormalizeERC20Tx converts Etherscan ERC-20 tx to normalized transaction
	NormalizeERC20Tx(tx EtherscanTokenTx) (*models.Transaction, error)
	
	// NormalizeERC721Tx converts Etherscan ERC-721 tx to normalized transaction
	NormalizeERC721Tx(tx EtherscanTokenTx) (*models.Transaction, error)
	
	// NormalizeERC1155Tx converts Etherscan ERC-1155 tx to normalized transaction
	NormalizeERC1155Tx(tx EtherscanTokenTx) (*models.Transaction, error)
}
