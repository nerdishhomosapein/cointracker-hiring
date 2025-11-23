package output

import (
	"conintracker-hiring/pkg/models"
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

// CSVWriter writes transactions to a CSV file
type CSVWriter struct {
	writer *csv.Writer
	file   io.WriteCloser
}

// CSVConfig holds configuration for CSV writing
type CSVConfig struct {
	Writer io.WriteCloser
}

// NewCSVWriter creates a new CSV writer
func NewCSVWriter(config CSVConfig) (*CSVWriter, error) {
	cw := &CSVWriter{
		writer: csv.NewWriter(config.Writer),
		file:   config.Writer,
	}

	// Write header
	headers := []string{
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

	if err := cw.writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	return cw, nil
}

// WriteTransaction writes a single transaction to CSV
func (cw *CSVWriter) WriteTransaction(tx *models.Transaction) error {
	// Format timestamp as RFC3339 (ISO 8601)
	timestamp := tx.Timestamp.Format(time.RFC3339)

	record := []string{
		tx.Hash,
		timestamp,
		tx.From,
		tx.To,
		string(tx.Type),
		tx.AssetContractAddress,
		tx.AssetSymbol,
		tx.TokenID,
		tx.Amount,
		tx.GasFeeETH,
	}

	if err := cw.writer.Write(record); err != nil {
		return fmt.Errorf("failed to write CSV record: %w", err)
	}

	return nil
}

// WriteTransactions writes multiple transactions to CSV
func (cw *CSVWriter) WriteTransactions(txs []*models.Transaction) error {
	for _, tx := range txs {
		if err := cw.WriteTransaction(tx); err != nil {
			return err
		}
	}
	return nil
}

// Close flushes the writer and closes the file
func (cw *CSVWriter) Close() error {
	cw.writer.Flush()
	if err := cw.writer.Error(); err != nil {
		return fmt.Errorf("CSV writer error: %w", err)
	}
	return cw.file.Close()
}

// Exporter interface for different output formats
type Exporter interface {
	WriteTransaction(tx *models.Transaction) error
	WriteTransactions(txs []*models.Transaction) error
	Close() error
}

// Writer interface for extensibility
type Writer interface {
	Write(txs []*models.Transaction) error
}

// CSVExporter is the CSV implementation of Exporter
var _ Exporter = (*CSVWriter)(nil)
