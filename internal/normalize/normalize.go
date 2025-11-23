package normalize

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"time"

	"conintracker-hiring/internal/etherscan"
)

// TxType represents the type of transaction
type TxType string

const (
	TypeExternal TxType = "External"
	TypeInternal TxType = "Internal"
	TypeERC20    TxType = "ERC-20"
	TypeERC721   TxType = "ERC-721"
	TypeERC1155  TxType = "ERC-1155"
)

// NormalizedTx represents a normalized transaction
type NormalizedTx struct {
	Hash            string
	Timestamp       time.Time
	From            string
	To              string
	Type            TxType
	ContractAddress string
	AssetSymbol     string
	TokenID         string
	Amount          string
	GasFeeEth       string
}

// RawData holds all types of raw transaction data
type RawData struct {
	Normal  []etherscan.NormalTx
	Internal []etherscan.InternalTx
	ERC20   []etherscan.TokenTx
	ERC721  []etherscan.ERC721Tx
	ERC1155 []etherscan.ERC1155Tx
}

// Normalize processes raw transaction data and returns normalized transactions
func Normalize(raw RawData) ([]NormalizedTx, error) {
	var result []NormalizedTx

	// Process normal transactions
	for _, tx := range raw.Normal {
		normalized, err := normalizeNormalTx(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to normalize normal tx %s: %w", tx.Hash, err)
		}
		result = append(result, normalized)
	}

	// Process internal transactions  
	for _, tx := range raw.Internal {
		normalized, err := normalizeInternalTx(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to normalize internal tx %s: %w", tx.Hash, err)
		}
		result = append(result, normalized)
	}

	// Process ERC-20 transactions
	for _, tx := range raw.ERC20 {
		normalized, err := normalizeERC20Tx(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to normalize ERC-20 tx %s: %w", tx.Hash, err)
		}
		result = append(result, normalized)
	}

	// Process ERC-721 transactions
	for _, tx := range raw.ERC721 {
		normalized, err := normalizeERC721Tx(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to normalize ERC-721 tx %s: %w", tx.Hash, err)
		}
		result = append(result, normalized)
	}

	// Process ERC-1155 transactions
	for _, tx := range raw.ERC1155 {
		normalized, err := normalizeERC1155Tx(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to normalize ERC-1155 tx %s: %w", tx.Hash, err)
		}
		result = append(result, normalized)
	}

	// Sort by timestamp
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})

	return result, nil
}

// normalizeNormalTx normalizes a normal Ethereum transaction
func normalizeNormalTx(tx etherscan.NormalTx) (NormalizedTx, error) {
	timestamp, err := parseTimestamp(tx.TimeStamp)
	if err != nil {
		return NormalizedTx{}, fmt.Errorf("invalid timestamp: %w", err)
	}

	amount := weiToETH(tx.Value)
	gasFee := calculateGasFeeETH(tx.GasUsed, tx.GasPrice)

	return NormalizedTx{
		Hash:            tx.Hash,
		Timestamp:       timestamp,
		From:            tx.From,
		To:              tx.To,
		Type:            TypeExternal,
		ContractAddress: tx.ContractAddress,
		AssetSymbol:     "ETH",
		TokenID:         "",
		Amount:          amount,
		GasFeeEth:       gasFee,
	}, nil
}

// normalizeInternalTx normalizes an internal transaction
func normalizeInternalTx(tx etherscan.InternalTx) (NormalizedTx, error) {
	timestamp, err := parseTimestamp(tx.TimeStamp)
	if err != nil {
		return NormalizedTx{}, fmt.Errorf("invalid timestamp: %w", err)
	}

	amount := weiToETH(tx.Value)

	return NormalizedTx{
		Hash:            tx.Hash,
		Timestamp:       timestamp,
		From:            tx.From,
		To:              tx.To,
		Type:            TypeInternal,
		ContractAddress: tx.ContractAddress,
		AssetSymbol:     "ETH",
		TokenID:         "",
		Amount:          amount,
		GasFeeEth:       "0", // Internal txs don't have gas fees
	}, nil
}

// normalizeERC20Tx normalizes an ERC-20 token transaction
func normalizeERC20Tx(tx etherscan.TokenTx) (NormalizedTx, error) {
	timestamp, err := parseTimestamp(tx.TimeStamp)
	if err != nil {
		return NormalizedTx{}, fmt.Errorf("invalid timestamp: %w", err)
	}

	decimals, err := strconv.Atoi(tx.TokenDecimal)
	if err != nil {
		return NormalizedTx{}, fmt.Errorf("invalid token decimals: %w", err)
	}

	amount := adjustForDecimals(tx.Value, decimals)
	gasFee := calculateGasFeeETH(tx.GasUsed, tx.GasPrice)

	return NormalizedTx{
		Hash:            tx.Hash,
		Timestamp:       timestamp,
		From:            tx.From,
		To:              tx.To,
		Type:            TypeERC20,
		ContractAddress: tx.ContractAddress,
		AssetSymbol:     tx.TokenSymbol,
		TokenID:         "",
		Amount:          amount,
		GasFeeEth:       gasFee,
	}, nil
}

// normalizeERC721Tx normalizes an ERC-721 NFT transaction
func normalizeERC721Tx(tx etherscan.ERC721Tx) (NormalizedTx, error) {
	timestamp, err := parseTimestamp(tx.TimeStamp)
	if err != nil {
		return NormalizedTx{}, fmt.Errorf("invalid timestamp: %w", err)
	}

	gasFee := calculateGasFeeETH(tx.GasUsed, tx.GasPrice)

	return NormalizedTx{
		Hash:            tx.Hash,
		Timestamp:       timestamp,
		From:            tx.From,
		To:              tx.To,
		Type:            TypeERC721,
		ContractAddress: tx.ContractAddress,
		AssetSymbol:     tx.TokenSymbol,
		TokenID:         tx.TokenID,
		Amount:          "1", // NFTs are always quantity 1
		GasFeeEth:       gasFee,
	}, nil
}

// normalizeERC1155Tx normalizes an ERC-1155 token transaction
func normalizeERC1155Tx(tx etherscan.ERC1155Tx) (NormalizedTx, error) {
	timestamp, err := parseTimestamp(tx.TimeStamp)
	if err != nil {
		return NormalizedTx{}, fmt.Errorf("invalid timestamp: %w", err)
	}

	gasFee := calculateGasFeeETH(tx.GasUsed, tx.GasPrice)

	return NormalizedTx{
		Hash:            tx.Hash,
		Timestamp:       timestamp,
		From:            tx.From,
		To:              tx.To,
		Type:            TypeERC1155,
		ContractAddress: tx.ContractAddress,
		AssetSymbol:     tx.TokenSymbol,
		TokenID:         tx.TokenID,
		Amount:          tx.TokenValue,
		GasFeeEth:       gasFee,
	}, nil
}

// Helper functions

// parseTimestamp converts a Unix timestamp string to time.Time
func parseTimestamp(timestampStr string) (time.Time, error) {
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(timestamp, 0), nil
}

// weiToETH converts wei (string) to ETH with proper decimal places
func weiToETH(weiStr string) string {
	if weiStr == "" || weiStr == "0" {
		return "0"
	}

	// Parse the wei value
	wei, ok := new(big.Int).SetString(weiStr, 10)
	if !ok {
		return "0"
	}

	// Convert to ETH (divide by 10^18)
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	eth := new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetInt(divisor))

	// Format with 18 decimal places
	return eth.Text('f', 18)
}

// calculateGasFeeETH calculates gas fee in ETH
func calculateGasFeeETH(gasUsedStr, gasPriceStr string) string {
	gasUsed, err1 := strconv.ParseUint(gasUsedStr, 10, 64)
	gasPrice, err2 := strconv.ParseUint(gasPriceStr, 10, 64)
	
	if err1 != nil || err2 != nil {
		return "0"
	}

	// Calculate total gas cost in wei
	totalGasCost := new(big.Int).Mul(
		new(big.Int).SetUint64(gasUsed),
		new(big.Int).SetUint64(gasPrice),
	)

	return weiToETH(totalGasCost.String())
}

// adjustForDecimals adjusts a token amount for its decimal places
func adjustForDecimals(valueStr string, decimals int) string {
	if valueStr == "" || valueStr == "0" {
		return "0"
	}

	// Parse the value
	value, ok := new(big.Int).SetString(valueStr, 10)
	if !ok {
		return "0"
	}

	if decimals <= 0 {
		return value.String()
	}

	// Convert to proper decimal representation
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	result := new(big.Float).Quo(new(big.Float).SetInt(value), new(big.Float).SetInt(divisor))

	// Format with the specified decimal places (don't trim for consistency)
	return result.Text('f', decimals)
}