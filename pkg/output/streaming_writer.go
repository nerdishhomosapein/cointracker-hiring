package output

import (
	"conintracker-hiring/pkg/models"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"sync"
	"time"
)

// StreamingCSVWriter writes transactions to CSV as they arrive via a channel
type StreamingCSVWriter struct {
	writer        *csv.Writer
	file          io.Writer
	batchSize     int
	flushInterval time.Duration
	headerWritten bool
	mu            sync.Mutex
}

// NewStreamingCSVWriter creates a new streaming CSV writer
func NewStreamingCSVWriter(w io.Writer) *StreamingCSVWriter {
	return &StreamingCSVWriter{
		writer:        csv.NewWriter(w),
		file:          w,
		batchSize:     100,
		flushInterval: 5 * time.Second,
		headerWritten: false,
	}
}

// SetBatchSize sets the number of transactions to batch before flushing
func (scw *StreamingCSVWriter) SetBatchSize(size int) {
	if size > 0 && size <= 10000 {
		scw.batchSize = size
	}
}

// SetFlushInterval sets the maximum time between flushes
func (scw *StreamingCSVWriter) SetFlushInterval(interval time.Duration) {
	if interval > 0 {
		scw.flushInterval = interval
	}
}

// WriteStream reads transactions from a channel and writes them to CSV
// Returns error if writing fails; returns ctx.Err() on context cancellation
func (scw *StreamingCSVWriter) WriteStream(
	ctx context.Context,
	txChan <-chan *models.Transaction,
	onProgress func(count int),
) error {
	// Write header once
	scw.mu.Lock()
	if !scw.headerWritten {
		if err := scw.writeHeader(); err != nil {
			scw.mu.Unlock()
			return fmt.Errorf("failed to write CSV header: %w", err)
		}
		scw.headerWritten = true
	}
	scw.mu.Unlock()

	batch := make([]*models.Transaction, 0, scw.batchSize)
	ticker := time.NewTicker(scw.flushInterval)
	defer ticker.Stop()

	count := 0

	for {
		select {
		case <-ctx.Done():
			// Flush remaining batch before exiting
			if len(batch) > 0 {
				scw.mu.Lock()
				if err := scw.writeBatch(batch); err != nil {
					scw.mu.Unlock()
					return fmt.Errorf("failed to write final batch: %w", err)
				}
				scw.mu.Unlock()
				if onProgress != nil {
					onProgress(count)
				}
			}
			return ctx.Err()

		case tx, ok := <-txChan:
			if !ok {
				// Channel closed, flush remaining batch
				if len(batch) > 0 {
					scw.mu.Lock()
					if err := scw.writeBatch(batch); err != nil {
						scw.mu.Unlock()
						return fmt.Errorf("failed to write final batch: %w", err)
					}
					scw.mu.Unlock()
					if onProgress != nil {
						onProgress(count)
					}
				}
				// Flush all remaining data
				scw.mu.Lock()
				scw.writer.Flush()
				scw.mu.Unlock()
				return nil
			}

			batch = append(batch, tx)
			count++

			// Flush batch if it reaches the batch size
			if len(batch) >= scw.batchSize {
				scw.mu.Lock()
				if err := scw.writeBatch(batch); err != nil {
					scw.mu.Unlock()
					return fmt.Errorf("failed to write batch: %w", err)
				}
				scw.mu.Unlock()
				batch = batch[:0] // Reset batch
				if onProgress != nil {
					onProgress(count)
				}
			}

		case <-ticker.C:
			// Periodic flush even if batch isn't full
			if len(batch) > 0 {
				scw.mu.Lock()
				if err := scw.writeBatch(batch); err != nil {
					scw.mu.Unlock()
					return fmt.Errorf("failed to write batch: %w", err)
				}
				scw.mu.Unlock()
				batch = batch[:0]
				if onProgress != nil {
					onProgress(count)
				}
			}
		}
	}
}

// writeBatch writes a batch of transactions (must be called with mutex held)
func (scw *StreamingCSVWriter) writeBatch(txs []*models.Transaction) error {
	for _, tx := range txs {
		record := []string{
			tx.Hash,
			tx.Timestamp.Format("2006-01-02 15:04:05 MST"),
			tx.From,
			tx.To,
			string(tx.Type),
			tx.AssetContractAddress,
			tx.AssetSymbol,
			tx.TokenID,
			tx.Amount,
			tx.GasFeeETH,
		}
		if err := scw.writer.Write(record); err != nil {
			return err
		}
	}
	scw.writer.Flush()
	return scw.writer.Error()
}

// writeHeader writes the CSV header row (must be called with mutex held)
func (scw *StreamingCSVWriter) writeHeader() error {
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
	if err := scw.writer.Write(header); err != nil {
		return err
	}
	scw.writer.Flush()
	return scw.writer.Error()
}

// StreamingOutputMetrics tracks performance metrics during streaming output
type StreamingOutputMetrics struct {
	TotalWritten     int64
	TotalErrors      int64
	StartTime        time.Time
	EndTime          time.Time
	BytesWritten     int64
	TransactionsPerSecond float64
}

// MetricsCollector collects streaming metrics
type MetricsCollector struct {
	mu      sync.RWMutex
	metrics StreamingOutputMetrics
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: StreamingOutputMetrics{
			StartTime: time.Now(),
		},
	}
}

// RecordWrite records a successful write
func (mc *MetricsCollector) RecordWrite(count int64, bytes int64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.metrics.TotalWritten += count
	mc.metrics.BytesWritten += bytes
}

// RecordError records a write error
func (mc *MetricsCollector) RecordError() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.metrics.TotalErrors++
}

// Finalize computes final metrics
func (mc *MetricsCollector) Finalize() StreamingOutputMetrics {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.metrics.EndTime = time.Now()
	duration := mc.metrics.EndTime.Sub(mc.metrics.StartTime).Seconds()
	if duration > 0 {
		mc.metrics.TransactionsPerSecond = float64(mc.metrics.TotalWritten) / duration
	}
	return mc.metrics
}

// GetMetrics returns a copy of current metrics
func (mc *MetricsCollector) GetMetrics() StreamingOutputMetrics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	metrics := mc.metrics
	if metrics.EndTime.IsZero() {
		metrics.EndTime = time.Now()
	}
	return metrics
}
