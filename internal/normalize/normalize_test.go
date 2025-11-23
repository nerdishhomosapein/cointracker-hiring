package normalize

import (
	"testing"
	"time"

	"conintracker-hiring/internal/etherscan"
)

func TestNormalizeAggregatesAndSorts(t *testing.T) {
	raw := RawData{
		Normal: []etherscan.NormalTx{
			{
				Hash:         "0xhash1",
				BlockNumber:  "1",
				TimeStamp:    "1609459200",
				From:         "0xfrom1",
				To:           "0xto1",
				Value:        "1000000000000000000",
				GasPrice:     "1000000000",
				GasUsed:      "21000",
				Nonce:        "1",
				TransactionIndex: "0",
			},
		},
		Internal: []etherscan.InternalTx{
			{
				Hash:            "0xhash3",
				BlockNumber:     "3",
				TimeStamp:       "1609459210",
				From:            "0xfrom3",
				To:              "0xto3",
				Value:           "5000000000000000",
				ContractAddress: "0xcontract1",
				Gas:             "80000",
				GasUsed:         "60000",
				IsError:         "0",
				Type:            "call",
				TraceID:         "0",
			},
		},
		ERC20: []etherscan.TokenTx{
			{
				Hash:           "0xhash4",
				BlockNumber:    "4",
				TimeStamp:      "1609459220",
				From:           "0xfrom4",
				To:             "0xto4",
				Value:          "1000000",
				TokenName:      "USD Coin",
				TokenSymbol:    "USDC",
				TokenDecimal:   "6",
				ContractAddress: "0xcontractUsdc",
				GasPrice:       "1200000000",
				GasUsed:        "50000",
			},
		},
		ERC721: []etherscan.ERC721Tx{
			{
				Hash:            "0xhash5",
				BlockNumber:     "5",
				TimeStamp:       "1609459230",
				From:            "0xfrom5",
				To:              "0xto5",
				TokenID:         "12345",
				TokenName:       "CoolNFTs",
				TokenSymbol:     "COOL",
				ContractAddress: "0xcontractNft",
				GasPrice:        "1300000000",
				GasUsed:         "55000",
			},
		},
		ERC1155: []etherscan.ERC1155Tx{
			{
				Hash:            "0xhash6",
				BlockNumber:     "6",
				TimeStamp:       "1609459240",
				From:            "0xfrom6",
				To:              "0xto6",
				TokenID:         "777",
				TokenValue:      "3",
				TokenName:       "Items1155",
				TokenSymbol:     "ITM",
				ContractAddress: "0xcontract1155",
				GasPrice:        "1400000000",
				GasUsed:         "60000",
			},
		},
	}

	out, err := Normalize(raw)
	if err != nil {
		t.Fatalf("Normalize returned error: %v", err)
	}
	if len(out) != 5 {
		t.Fatalf("expected 5 normalized entries, got %d", len(out))
	}

	checkOrder(t, out, []string{"0xhash1", "0xhash3", "0xhash4", "0xhash5", "0xhash6"})

	ext := find(t, out, "0xhash1")
	if ext.Type != TypeExternal {
		t.Fatalf("expected external type got %s", ext.Type)
	}
	if ext.Amount != "1.000000000000000000" {
		t.Fatalf("unexpected amount for external: %s", ext.Amount)
	}
	if ext.GasFeeEth != "0.000021000000000000" {
		t.Fatalf("unexpected gas fee: %s", ext.GasFeeEth)
	}
	expectTime(t, ext.Timestamp, 1609459200)

	internal := find(t, out, "0xhash3")
	if internal.Type != TypeInternal || internal.ContractAddress != "0xcontract1" {
		t.Fatalf("unexpected internal tx: %+v", internal)
	}
	if internal.Amount != "0.005000000000000000" {
		t.Fatalf("unexpected internal value: %s", internal.Amount)
	}

	erc20 := find(t, out, "0xhash4")
	if erc20.Type != TypeERC20 || erc20.AssetSymbol != "USDC" || erc20.ContractAddress != "0xcontractUsdc" {
		t.Fatalf("unexpected erc20 tx: %+v", erc20)
	}
	if erc20.Amount != "1.000000" {
		t.Fatalf("unexpected erc20 amount: %s", erc20.Amount)
	}

	erc721 := find(t, out, "0xhash5")
	if erc721.Type != TypeERC721 || erc721.TokenID != "12345" {
		t.Fatalf("unexpected erc721 tx: %+v", erc721)
	}
	if erc721.Amount != "1" {
		t.Fatalf("unexpected erc721 amount: %s", erc721.Amount)
	}

	erc1155 := find(t, out, "0xhash6")
	if erc1155.Type != TypeERC1155 || erc1155.TokenID != "777" {
		t.Fatalf("unexpected erc1155 tx: %+v", erc1155)
	}
	if erc1155.Amount != "3" {
		t.Fatalf("unexpected erc1155 amount: %s", erc1155.Amount)
	}
}

func checkOrder(t *testing.T, txs []NormalizedTx, expected []string) {
	t.Helper()
	if len(txs) != len(expected) {
		t.Fatalf("length mismatch: got %d expected %d", len(txs), len(expected))
	}
	for i, hash := range expected {
		if txs[i].Hash != hash {
			t.Fatalf("order mismatch at %d: got %s expected %s", i, txs[i].Hash, hash)
		}
	}
}

func find(t *testing.T, txs []NormalizedTx, hash string) NormalizedTx {
	t.Helper()
	for _, tx := range txs {
		if tx.Hash == hash {
			return tx
		}
	}
	t.Fatalf("hash %s not found", hash)
	return NormalizedTx{}
}

func expectTime(t *testing.T, ts time.Time, expectedUnix int64) {
	t.Helper()
	if ts.Unix() != expectedUnix {
		t.Fatalf("unexpected timestamp: %v", ts)
	}
}
