package models

import (
	"time"
)

// TransactionType represents the category of transaction
type TransactionType string

const (
	TypeEthTransfer    TransactionType = "ETH"
	TypeERC20Transfer  TransactionType = "ERC-20"
	TypeERC721Transfer TransactionType = "ERC-721"
	TypeERC1155Transfer TransactionType = "ERC-1155"
	TypeInternal       TransactionType = "Internal"
	TypeContractCreate TransactionType = "Contract Creation"
)

// Transaction represents a normalized transaction record
type Transaction struct {
	// Core transaction info
	Hash      string `csv:"Transaction Hash"`
	Timestamp time.Time `csv:"Date & Time"`
	From      string `csv:"From Address"`
	To        string `csv:"To Address"`
	
	// Transaction categorization
	Type TransactionType `csv:"Transaction Type"`
	
	// Asset info
	AssetContractAddress string `csv:"Asset Contract Address"`
	AssetSymbol          string `csv:"Asset Symbol / Name"`
	TokenID              string `csv:"Token ID"` // For NFTs (ERC-721, ERC-1155)
	
	// Values
	Amount  string `csv:"Value / Amount"` // Quantity transferred
	GasFeeETH string `csv:"Gas Fee (ETH)"` // Total gas cost in ETH
	
	// Additional metadata (not in CSV but useful for processing)
	BlockNumber     uint64 `csv:"-"`
	GasUsed         uint64 `csv:"-"`
	GasPrice        string `csv:"-"` // in Wei
	TransactionFee  string `csv:"-"` // in Wei
	Nonce           uint64 `csv:"-"`
	IsError         bool   `csv:"-"`
	Input           string `csv:"-"`
	MethodID        string `csv:"-"`
	FunctionName    string `csv:"-"`
	Decimals        int    `csv:"-"` // For token transfers
}

// TransactionList is a sortable slice of transactions
type TransactionList []*Transaction

// Len implements sort.Interface
func (tl TransactionList) Len() int {
	return len(tl)
}

// Less implements sort.Interface (sort by block number first, then timestamp)
func (tl TransactionList) Less(i, j int) bool {
	if tl[i].BlockNumber != tl[j].BlockNumber {
		return tl[i].BlockNumber < tl[j].BlockNumber
	}
	return tl[i].Timestamp.Before(tl[j].Timestamp)
}

// Swap implements sort.Interface
func (tl TransactionList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}
