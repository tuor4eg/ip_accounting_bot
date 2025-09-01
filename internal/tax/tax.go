package tax

// ExtraOverThreshold returns the extra contribution for income above a threshold.
// yearIncome and threshold are in kopecks; rateBP is in basis points (1% = 100 bp).
// The calculation is integer-only and floors toward zero.
func ExtraOverThreshold(yearIncome, threshold, rateBP int64) int64 {
	if rateBP <= 0 || yearIncome <= threshold {
		return 0
	}

	const bpDen = int64(10_000) // basis points denominator

	excess := yearIncome - threshold

	return excess * rateBP / bpDen
}
