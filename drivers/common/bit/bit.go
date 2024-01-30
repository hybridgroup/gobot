package bit

// Set is used to set a bit in the given integer at a given position to 1.
func Set(n int, pos uint8) int {
	n |= (1 << pos)
	return n
}

// Clear is used to set a bit in the given integer at a given position to 0.
func Clear(n int, pos uint8) int {
	mask := ^int(1 << pos)
	n &= mask
	return n
}

// IsSet tests if the bit at the given position is set in the given integer.
func IsSet(n int, pos uint8) bool {
	val := n & (1 << uint(pos))
	return (val > 0)
}
