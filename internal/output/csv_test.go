package output

import (
	"strings"
	"testing"
	"time"

	"conintracker-hiring/internal/normalize"
)

func TestCSVWriterProducesHeadersAndRows(t *testing.T) {
	rows := []normalize.NormalizedTx{
		{
			Hash:            "0xhash1",
			Timestamp:       time.Unix(1609459200, 0).UTC(),
			From:            "0xfrom1",
			To:              "0xto1",
			Type:            normalize.TypeExternal,
			ContractAddress: "",
			AssetSymbol:     "ETH",
			TokenID:         "",
			Amount:          "1.000000000000000000",
			GasFeeEth:       "0.000021000000000000",
		},
		{
			Hash:            "0xhash4",
			Timestamp:       time.Unix(1609459220, 0).UTC(),
			From:            "0xfrom4",
			To:              "0xto4",
			Type:            normalize.TypeERC20,
			ContractAddress: "0xcontractUsdc",
			AssetSymbol:     "USDC",
			TokenID:         "",
			Amount:          "1.000000",
			GasFeeEth:       "0.000060000000000000",
		},
	}

	var buf strings.Builder
	writer, err := NewWriter("csv", &buf)
	if err != nil {
		t.Fatalf("NewWriter error: %v", err)
	}
	if err := writer.Write(rows); err != nil {
		t.Fatalf("write error: %v", err)
	}

	got := buf.String()
	if !strings.HasPrefix(got, "Transaction Hash,Date & Time,From Address,To Address,Transaction Type,Asset Contract Address,Asset Symbol / Name,Token ID,Value / Amount,Gas Fee (ETH)\n") {
		t.Fatalf("missing or incorrect header: %s", got)
	}
	if !strings.Contains(got, "0xhash1,2021-01-01T00:00:00Z,0xfrom1,0xto1,eth_transfer,,ETH,,1.000000000000000000,0.000021000000000000") {
		t.Fatalf("missing row for hash1: %s", got)
	}
	if !strings.Contains(got, "0xhash4,2021-01-01T00:00:20Z,0xfrom4,0xto4,erc20,0xcontractUsdc,USDC,,1.000000,0.000060000000000000") {
		t.Fatalf("missing row for hash4: %s", got)
	}
}

func TestUnknownFormatErrors(t *testing.T) {
	_, err := NewWriter("pdf", &strings.Builder{})
	if err == nil {
		t.Fatalf("expected error for unknown format")
	}
}
