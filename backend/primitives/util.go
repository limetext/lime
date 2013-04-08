package primitives

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Clamp(_min, _max, v int) int {
	return Max(_min, Min(_max, v))
}

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
