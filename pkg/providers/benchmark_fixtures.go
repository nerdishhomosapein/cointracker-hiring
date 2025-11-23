package providers

import (
	"conintracker-hiring/pkg/models"
)

// BenchmarkFixtures contains reusable test data for benchmarks
type BenchmarkFixtures struct {
	NormalTxs      []EtherscanNormalTx
	InternalTxs    []EtherscanInternalTx
	TokenTxs       []EtherscanTokenTx
	NFTTxs         []EtherscanTokenTx
	ERC1155Txs     []EtherscanTokenTx
	NormalizedTxs  []*models.Transaction
}

// NewBenchmarkFixtures creates a set of benchmark fixtures with realistic data
func NewBenchmarkFixtures(size int) *BenchmarkFixtures {
	fixtures := &BenchmarkFixtures{
		NormalTxs:     make([]EtherscanNormalTx, size),
		InternalTxs:   make([]EtherscanInternalTx, size),
		TokenTxs:      make([]EtherscanTokenTx, size),
		NFTTxs:        make([]EtherscanTokenTx, size),
		ERC1155Txs:    make([]EtherscanTokenTx, size),
		NormalizedTxs: make([]*models.Transaction, 0, size*5),
	}

	// Generate normal transactions
	for i := 0; i < size; i++ {
		fixtures.NormalTxs[i] = EtherscanNormalTx{
			BlockNumber: "19000000",
			TimeStamp:   "1700000000",
			Hash:        "0x" + padHex(i, 64),
			From:        "0x" + padHex(i%10, 40),
			To:          "0x" + padHex(i%20, 40),
			Value:       "1000000000000000000", // 1 ETH
			GasUsed:     "21000",
			GasPrice:    "20000000000",
			IsError:     "0",
			Input:       "0x",
			MethodId:    "0x",
			FunctionName: "",
		}
	}

	// Generate internal transactions
	for i := 0; i < size; i++ {
		fixtures.InternalTxs[i] = EtherscanInternalTx{
			BlockNumber: "19000000",
			TimeStamp:   "1700000000",
			Hash:        "0x" + padHex(i, 64),
			From:        "0x" + padHex(i%10, 40),
			To:          "0x" + padHex(i%20, 40),
			Value:       "500000000000000000", // 0.5 ETH
			GasUsed:     "10000",
			IsError:     "0",
			Input:       "0x",
		}
	}

	// Generate token transfers
	for i := 0; i < size; i++ {
		fixtures.TokenTxs[i] = EtherscanTokenTx{
			BlockNumber:     "19000000",
			TimeStamp:       "1700000000",
			Hash:            "0x" + padHex(i, 64),
			From:            "0x" + padHex(i%10, 40),
			To:              "0x" + padHex(i%20, 40),
			Value:           "1000000000000000000",
			TokenDecimal:    "18",
			TokenSymbol:     "USDC",
			ContractAddress: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			GasUsed:         "65000",
			GasPrice:        "20000000000",
			IsError:         "0",
		}
	}

	// Generate NFT transfers (ERC-721)
	for i := 0; i < size; i++ {
		fixtures.NFTTxs[i] = EtherscanTokenTx{
			BlockNumber:     "19000000",
			TimeStamp:       "1700000000",
			Hash:            "0x" + padHex(i, 64),
			From:            "0x" + padHex(i%10, 40),
			To:              "0x" + padHex(i%20, 40),
			Value:           "1",
			TokenID:         padDecimal(i),
			TokenSymbol:     "BLUR",
			ContractAddress: "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
			GasUsed:         "70000",
			GasPrice:        "20000000000",
			IsError:         "0",
		}
	}

	// Generate ERC-1155 transfers
	for i := 0; i < size; i++ {
		fixtures.ERC1155Txs[i] = EtherscanTokenTx{
			BlockNumber:     "19000000",
			TimeStamp:       "1700000000",
			Hash:            "0x" + padHex(i, 64),
			From:            "0x" + padHex(i%10, 40),
			To:              "0x" + padHex(i%20, 40),
			Value:           "100",
			TokenValue:      "50",
			TokenID:         padDecimal(i),
			TokenSymbol:     "ERC1155",
			ContractAddress: "0x0000000000000000000000000000000000000000",
			GasUsed:         "85000",
			GasPrice:        "20000000000",
			IsError:         "0",
		}
	}

	return fixtures
}

// padHex pads an integer to a hex string of specified length with leading zeros
func padHex(i int, length int) string {
	hexStr := ""
	num := i
	if num == 0 {
		hexStr = "0"
	} else {
		for num > 0 {
			digit := num % 16
			if digit < 10 {
				hexStr = string(rune(48+digit)) + hexStr
			} else {
				hexStr = string(rune(97+digit-10)) + hexStr
			}
			num /= 16
		}
	}
	
	// Pad with leading zeros if needed
	for len(hexStr) < length {
		hexStr = "0" + hexStr
	}
	
	// Truncate if too long
	if len(hexStr) > length {
		hexStr = hexStr[len(hexStr)-length:]
	}
	
	return hexStr
}

// padDecimal pads an integer to a decimal string
func padDecimal(i int) string {
	hexStr := ""
	num := i
	if num == 0 {
		return "0"
	}
	for num > 0 {
		hexStr = string(rune(48+(num%10))) + hexStr
		num /= 10
	}
	return hexStr
}

// GetLargeFixture returns a fixture set suitable for large-scale benchmarking
func GetLargeFixture() *BenchmarkFixtures {
	return NewBenchmarkFixtures(10000)
}

// GetMediumFixture returns a fixture set suitable for medium-scale benchmarking
func GetMediumFixture() *BenchmarkFixtures {
	return NewBenchmarkFixtures(1000)
}

// GetSmallFixture returns a fixture set suitable for small-scale benchmarking
func GetSmallFixture() *BenchmarkFixtures {
	return NewBenchmarkFixtures(100)
}
