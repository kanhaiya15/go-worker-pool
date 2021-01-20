package utils

// NextUint64 out next increment
func NextUint64() func() uint64 {
	var i uint64
	return func() uint64 {
		i++
		return i
	}
}
