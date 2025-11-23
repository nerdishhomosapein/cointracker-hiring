package providers

// BaselineBenchmarkResults stores the baseline performance metrics
// These are populated when baseline benchmarks are run
type BaselineBenchmarkResults struct {
	// Normalization benchmarks (per operation in nanoseconds)
	WeiToETHNs              int64 // ns/op
	CalculateGasFeeETHNs    int64 // ns/op
	AdjustForDecimalsNs     int64 // ns/op
	ParseUint64Ns           int64 // ns/op
	ParseTimestampNs        int64 // ns/op

	// Transaction normalization (per transaction in nanoseconds)
	NormalizeNormalTxNs     int64 // ns/op
	NormalizeInternalTxNs   int64 // ns/op
	NormalizeERC20TxNs      int64 // ns/op
	NormalizeERC721TxNs     int64 // ns/op
	NormalizeERC1155TxNs    int64 // ns/op

	// Full pipeline
	NormalizationPipelineNs int64 // ns/op for processing all 5 types

	// Fetch orchestration
	FetchAllTransactionsNs  int64 // ns/op

	// Memory allocations
	// These will be populated by benchstat post-processing
	TxNormalizationAllocsPerOp int64
	FetchAllTransactionsAllocsPerOp int64
}

// GetExpectedBaseline returns conservative baseline expectations based on the platform
// These are used for regression testing
func GetExpectedBaseline() *BaselineBenchmarkResults {
	return &BaselineBenchmarkResults{
		// Conservative estimates - actual values will be measured
		WeiToETHNs:              2000,   // ~2µs per wei to ETH conversion (big.Int operations)
		CalculateGasFeeETHNs:    3000,   // ~3µs per gas fee calculation
		AdjustForDecimalsNs:     2500,   // ~2.5µs per decimal adjustment
		ParseUint64Ns:           200,    // ~0.2µs per uint64 parse
		ParseTimestampNs:        300,    // ~0.3µs per timestamp parse

		NormalizeNormalTxNs:     10000,  // ~10µs per normal tx (calls several helpers)
		NormalizeInternalTxNs:   8000,   // ~8µs per internal tx
		NormalizeERC20TxNs:      12000,  // ~12µs per ERC20 tx (includes decimal parsing)
		NormalizeERC721TxNs:     11000,  // ~11µs per ERC721 tx
		NormalizeERC1155TxNs:    12000,  // ~12µs per ERC1155 tx

		NormalizationPipelineNs: 15000000, // ~15ms for 1000 transactions total (all 5 types)

		FetchAllTransactionsNs: 20000000, // ~20ms for orchestration with 1000 txs
	}
}

// RegressionThreshold defines acceptable deviation from baseline
type RegressionThreshold struct {
	// PercentageIncrease is the acceptable % increase from baseline (default 10%)
	PercentageIncrease float64
	// AbsoluteNsIncrease is additional absolute nanosecond tolerance
	AbsoluteNsIncrease int64
}

// GetDefaultRegressionThreshold returns sensible defaults for performance regression detection
func GetDefaultRegressionThreshold() *RegressionThreshold {
	return &RegressionThreshold{
		PercentageIncrease: 10.0,  // 10% degradation allowed
		AbsoluteNsIncrease: 5000,  // plus 5µs absolute tolerance
	}
}
