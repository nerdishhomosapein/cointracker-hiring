package output

import (
	"encoding/csv"
	"fmt"
	"io"

	"conintracker-hiring/internal/normalize"
)

// Writer represents a transaction output writer
type Writer interface {
	Write([]normalize.NormalizedTx) error
}

// NewWriter creates a new writer for the specified format
func NewWriter(format string, w io.Writer) (Writer, error) {
	switch format {
	case "csv":
		return &CSVWriter{writer: csv.NewWriter(w)}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// CSVWriter writes normalized transactions to CSV format
type CSVWriter struct {
	writer *csv.Writer
}

// Write writes the transactions to CSV format
func (w *CSVWriter) Write(txs []normalize.NormalizedTx) error {
	// Write header
	header := []string{
		"Transaction Hash",
		"Date & Time",
		"From Address",
		"To Address",
		"Transaction Type",
		"Asset Contract Address",
		"Asset Symbol / Name",
		"Token ID",
		"Value / Amount",
		"Gas Fee (ETH)",
	}
	
	if err := w.writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write transaction rows
	for _, tx := range txs {
		row := []string{
			tx.Hash,
			tx.Timestamp.Format("2006-01-02T15:04:05Z"),
			tx.From,
			tx.To,
			formatTxType(tx.Type),
			tx.ContractAddress,
			tx.AssetSymbol,
			tx.TokenID,
			tx.Amount,
			tx.GasFeeEth,
		}
		
		if err := w.writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	// Flush the writer
	w.writer.Flush()
	return w.writer.Error()
}

// formatTxType converts the transaction type to the expected string format
func formatTxType(txType normalize.TxType) string {
	switch txType {
	case normalize.TypeExternal:
		return "eth_transfer"
	case normalize.TypeInternal:
		return "internal"
	case normalize.TypeERC20:
		return "erc20"
	case normalize.TypeERC721:
		return "erc721"
	case normalize.TypeERC1155:
		return "erc1155"
	default:
		return string(txType)
	}
}