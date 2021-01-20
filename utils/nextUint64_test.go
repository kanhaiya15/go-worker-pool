package utils

import "testing"

func TestNextUint64(t *testing.T) {
	nextUint64 := NextUint64()
	i := nextUint64()
	if i != 1 {
		t.Errorf("NextUint64 is incorrect, got: %d, want: %d.", i, 2)
	}
	ii := nextUint64()
	if ii != 2 {
		t.Errorf("NextUint64 is incorrect, got: %d, want: %d.", ii, 2)
	}
}
