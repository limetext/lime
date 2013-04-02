package primitives

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clamp(_min, _max, v int) int {
	return max(_min, min(_max, v))
}
