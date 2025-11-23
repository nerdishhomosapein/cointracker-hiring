package output

import (
	"bytes"
	"conintracker-hiring/pkg/models"
	"context"
	"sync"
	"testing"
	"time"
)

// BenchmarkStreamingCSVWriter benchmarks writing transactions via streaming
func BenchmarkStreamingCSVWriter(b *testing.B) {
	// Generate test transactions
	generateTransactions := func(count int) chan *models.Transaction {
		txChan := make(chan *models.Transaction, 100)
		go func() {
			for i := 0; i < count; i++ {
				txChan <- &models.Transaction{
					Hash:      "0x" + string(rune(48+(i%10))),
					Timestamp: time.Now(),
					From:      "0x1111111111111111111111111111111111111111",
					To:        "0x2222222222222222222222222222222222222222",
					Type:      models.TypeEthTransfer,
					Amount:    "1.5",
					GasFeeETH: "0.001",
				}
			}
			close(txChan)
		}()
		return txChan
	}

	b.Run("1000Transactions", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf := &bytes.Buffer{}
			writer := NewStreamingCSVWriter(buf)
			ctx := context.Background()
			txChan := generateTransactions(1000)
			writer.WriteStream(ctx, txChan, nil)
		}
	})

	b.Run("10000Transactions", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf := &bytes.Buffer{}
			writer := NewStreamingCSVWriter(buf)
			ctx := context.Background()
			txChan := generateTransactions(10000)
			writer.WriteStream(ctx, txChan, nil)
		}
	})
}

// BenchmarkStreamingWithProgress benchmarks streaming with progress callback
func BenchmarkStreamingWithProgress(b *testing.B) {
	generateTransactions := func(count int) chan *models.Transaction {
		txChan := make(chan *models.Transaction, 100)
		go func() {
			for i := 0; i < count; i++ {
				txChan <- &models.Transaction{
					Hash:      "0x" + string(rune(48+(i%10))),
					Timestamp: time.Now(),
					From:      "0x1111111111111111111111111111111111111111",
					To:        "0x2222222222222222222222222222222222222222",
					Type:      models.TypeEthTransfer,
					Amount:    "1.5",
					GasFeeETH: "0.001",
				}
			}
			close(txChan)
		}()
		return txChan
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := &bytes.Buffer{}
		writer := NewStreamingCSVWriter(buf)
		ctx := context.Background()
		txChan := generateTransactions(1000)

		progressCount := 0
		writer.WriteStream(ctx, txChan, func(count int) {
			progressCount = count
		})
		_ = progressCount
	}
}

// BenchmarkMetricsCollector benchmarks the metrics collection overhead
func BenchmarkMetricsCollector(b *testing.B) {
	collector := NewMetricsCollector()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate recording writes and reads
		collector.RecordWrite(1, 100)
		collector.RecordWrite(10, 1000)
		collector.RecordError()
		metrics := collector.GetMetrics()
		_ = metrics
	}
}

// TestStreamingCSVWriter tests the streaming writer functionality
func TestStreamingCSVWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	writer := NewStreamingCSVWriter(buf)
	writer.SetBatchSize(2) // Small batch for testing

	// Create transaction channel
	txChan := make(chan *models.Transaction)

	var wg sync.WaitGroup
	wg.Add(1)
	var err error
	go func() {
		defer wg.Done()
		ctx := context.Background()
		err = writer.WriteStream(ctx, txChan, nil)
	}()

	// Send some transactions
	for i := 0; i < 5; i++ {
		txChan <- &models.Transaction{
			Hash:      "0xabc" + string(rune(48+i)),
			Timestamp: time.Now(),
			From:      "0x1111111111111111111111111111111111111111",
			To:        "0x2222222222222222222222222222222222222222",
			Type:      models.TypeEthTransfer,
			Amount:    "1.0",
			GasFeeETH: "0.001",
		}
	}
	close(txChan)

	wg.Wait()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that we have output
	output := buf.String()
	if len(output) == 0 {
		t.Fatal("expected CSV output but got empty string")
	}

	// Check for header
	if !bytes.Contains(buf.Bytes(), []byte("Transaction Hash")) {
		t.Fatal("CSV header not found in output")
	}
}

// TestMetricsCollector tests metrics collection
func TestMetricsCollector(t *testing.T) {
	collector := NewMetricsCollector()

	collector.RecordWrite(100, 1000)
	collector.RecordWrite(50, 500)
	collector.RecordError()
	collector.RecordError()

	metrics := collector.Finalize()

	if metrics.TotalWritten != 150 {
		t.Errorf("expected 150 total written, got %d", metrics.TotalWritten)
	}

	if metrics.TotalErrors != 2 {
		t.Errorf("expected 2 errors, got %d", metrics.TotalErrors)
	}

	if metrics.BytesWritten != 1500 {
		t.Errorf("expected 1500 bytes written, got %d", metrics.BytesWritten)
	}
}
