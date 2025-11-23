package providers

import (
	"conintracker-hiring/pkg/models"
	"math"
	"math/big"
	"strconv"
	"time"
)

// EtherscanNormalizer implements the Normalizer interface for Etherscan responses
type EtherscanNormalizer struct{}

// NewEtherscanNormalizer creates a new normalizer instance
func NewEtherscanNormalizer() *EtherscanNormalizer {
	return &EtherscanNormalizer{}
}

// weiToETH converts wei (big.Int) to ETH with proper decimal formatting
func weiToETH(weiStr string) string {
	if weiStr == "" || weiStr == "0" {
		return "0"
	}

	wei := new(big.Int)
	wei.SetString(weiStr, 10)

	// 1 ETH = 10^18 wei
	divisor := big.NewInt(1e18)
	eth := new(big.Rat).SetInt(wei)
	eth.Quo(eth, new(big.Rat).SetInt(divisor))

	f, _ := eth.Float64()
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// parseUint64 safely parses a string to uint64
func parseUint64(s string) uint64 {
	val, _ := strconv.ParseUint(s, 10, 64)
	return val
}

// parseTimestamp converts Unix timestamp string to time.Time
func parseTimestamp(timestampStr string) time.Time {
	ts, _ := strconv.ParseInt(timestampStr, 10, 64)
	return time.Unix(ts, 0)
}

// calculateGasFeeETH calculates gas fee in ETH (gasUsed * gasPrice / 1e18)
func calculateGasFeeETH(gasUsedStr, gasPriceStr string) string {
	gasUsed := new(big.Int)
	gasUsed.SetString(gasUsedStr, 10)

	gasPrice := new(big.Int)
	gasPrice.SetString(gasPriceStr, 10)

	// totalFeeWei = gasUsed * gasPrice
	totalFeeWei := new(big.Int)
	totalFeeWei.Mul(gasUsed, gasPrice)

	// Convert wei to ETH
	divisor := big.NewInt(1e18)
	fee := new(big.Rat).SetInt(totalFeeWei)
	fee.Quo(fee, new(big.Rat).SetInt(divisor))

	f, _ := fee.Float64()
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// adjustForDecimals scales a token value based on its decimal places
func adjustForDecimals(valueStr string, decimals int) string {
	if valueStr == "" || valueStr == "0" {
		return "0"
	}

	val := new(big.Int)
	val.SetString(valueStr, 10)

	// If decimals = 6, we divide by 1e6
	if decimals == 0 {
		return val.String()
	}

	divisor := big.NewInt(int64(math.Pow(10, float64(decimals))))
	result := new(big.Rat).SetInt(val)
	result.Quo(result, new(big.Rat).SetInt(divisor))

	f, _ := result.Float64()
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// NormalizeNormalTx implements Normalizer interface for normal ETH transfers
func (n *EtherscanNormalizer) NormalizeNormalTx(tx EtherscanNormalTx) (*models.Transaction, error) {
	isError := tx.IsError == "1"
	blockNum := parseUint64(tx.BlockNumber)

	return &models.Transaction{
		Hash:      tx.Hash,
		Timestamp: parseTimestamp(tx.TimeStamp),
		From:      tx.From,
		To:        tx.To,
		Type:      models.TypeEthTransfer,
		Amount:    weiToETH(tx.Value),
		GasFeeETH: calculateGasFeeETH(tx.GasUsed, tx.GasPrice),
		BlockNumber: blockNum,
		GasUsed:     parseUint64(tx.GasUsed),
		GasPrice:    tx.GasPrice,
		TransactionFee: tx.GasUsed, // This is calculated later
		IsError:     isError,
		Input:       tx.Input,
		MethodID:    tx.MethodId,
		FunctionName: tx.FunctionName,
	}, nil
}

// NormalizeInternalTx implements Normalizer interface for internal transfers
func (n *EtherscanNormalizer) NormalizeInternalTx(tx EtherscanInternalTx) (*models.Transaction, error) {
	isError := tx.IsError == "1"
	blockNum := parseUint64(tx.BlockNumber)

	return &models.Transaction{
		Hash:      tx.Hash,
		Timestamp: parseTimestamp(tx.TimeStamp),
		From:      tx.From,
		To:        tx.To,
		Type:      models.TypeInternal,
		Amount:    weiToETH(tx.Value),
		BlockNumber: blockNum,
		GasUsed:     parseUint64(tx.GasUsed),
		IsError:     isError,
		Input:       tx.Input,
	}, nil
}

// NormalizeERC20Tx implements Normalizer interface for ERC-20 token transfers
func (n *EtherscanNormalizer) NormalizeERC20Tx(tx EtherscanTokenTx) (*models.Transaction, error) {
	decimals, _ := strconv.Atoi(tx.TokenDecimal)

	return &models.Transaction{
		Hash:                 tx.Hash,
		Timestamp:            parseTimestamp(tx.TimeStamp),
		From:                 tx.From,
		To:                   tx.To,
		Type:                 models.TypeERC20Transfer,
		AssetContractAddress: tx.ContractAddress,
		AssetSymbol:          tx.TokenSymbol,
		Amount:               adjustForDecimals(tx.Value, decimals),
		GasFeeETH:            calculateGasFeeETH(tx.GasUsed, tx.GasPrice),
		BlockNumber:          parseUint64(tx.BlockNumber),
		GasUsed:              parseUint64(tx.GasUsed),
		GasPrice:             tx.GasPrice,
		IsError:              tx.IsError == "1",
		Decimals:             decimals,
	}, nil
}

// NormalizeERC721Tx implements Normalizer interface for ERC-721 NFT transfers
func (n *EtherscanNormalizer) NormalizeERC721Tx(tx EtherscanTokenTx) (*models.Transaction, error) {
	return &models.Transaction{
		Hash:                 tx.Hash,
		Timestamp:            parseTimestamp(tx.TimeStamp),
		From:                 tx.From,
		To:                   tx.To,
		Type:                 models.TypeERC721Transfer,
		AssetContractAddress: tx.ContractAddress,
		AssetSymbol:          tx.TokenSymbol,
		TokenID:              tx.TokenID,
		Amount:               "1", // NFTs are always 1
		GasFeeETH:            calculateGasFeeETH(tx.GasUsed, tx.GasPrice),
		BlockNumber:          parseUint64(tx.BlockNumber),
		GasUsed:              parseUint64(tx.GasUsed),
		GasPrice:             tx.GasPrice,
		IsError:              tx.IsError == "1",
	}, nil
}

// NormalizeERC1155Tx implements Normalizer interface for ERC-1155 multi-token transfers
func (n *EtherscanNormalizer) NormalizeERC1155Tx(tx EtherscanTokenTx) (*models.Transaction, error) {
	// For ERC-1155, use TokenValue if available, otherwise Value
	amount := tx.TokenValue
	if amount == "" {
		amount = tx.Value
	}

	return &models.Transaction{
		Hash:                 tx.Hash,
		Timestamp:            parseTimestamp(tx.TimeStamp),
		From:                 tx.From,
		To:                   tx.To,
		Type:                 models.TypeERC1155Transfer,
		AssetContractAddress: tx.ContractAddress,
		AssetSymbol:          tx.TokenSymbol,
		TokenID:              tx.TokenID,
		Amount:               amount,
		GasFeeETH:            calculateGasFeeETH(tx.GasUsed, tx.GasPrice),
		BlockNumber:          parseUint64(tx.BlockNumber),
		GasUsed:              parseUint64(tx.GasUsed),
		GasPrice:             tx.GasPrice,
		IsError:              tx.IsError == "1",
	}, nil
}
