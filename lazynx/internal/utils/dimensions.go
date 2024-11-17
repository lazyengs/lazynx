package utils

func SafeHalves(n int) (int, int) {
	half := n / 2

	if n%2 == 0 {
		return half, half
	}

	return half, n - half
}
