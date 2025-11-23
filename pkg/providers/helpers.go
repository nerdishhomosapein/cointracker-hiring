package providers

import "time"

// These functions expose internal normalizer helpers for benchmarking and testing purposes

// WeiToETH is the public wrapper for weiToETH
func WeiToETH(weiStr string) string {
	return weiToETH(weiStr)
}

// CalculateGasFeeETH is the public wrapper for calculateGasFeeETH
func CalculateGasFeeETH(gasUsedStr, gasPriceStr string) string {
	return calculateGasFeeETH(gasUsedStr, gasPriceStr)
}

// AdjustForDecimals is the public wrapper for adjustForDecimals
func AdjustForDecimals(valueStr string, decimals int) string {
	return adjustForDecimals(valueStr, decimals)
}

// ParseUint64Public is the public wrapper for parseUint64
func ParseUint64Public(s string) uint64 {
	return parseUint64(s)
}

// ParseTimestampPublic is the public wrapper for parseTimestamp
func ParseTimestampPublic(timestampStr string) time.Time {
	return parseTimestamp(timestampStr)
}
