package btcplex

// Float rounding to precision 8
func FloatToUint(x float64) uint64 {
	return uint64(int64((x * float64(100000000.0)) + float64(0.5)))
}
