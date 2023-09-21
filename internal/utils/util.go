package utils

import "github.com/cockroachdb/pebble"

func KeyUpperBound(prefix []byte) []byte {
	ub := make([]byte, len(prefix))
	copy(ub, prefix)
	for i := len(ub) - 1; i >= 0; i-- {
		ub[i] = ub[i] + 1
		if ub[i] != 0 {
			return ub[:i+1]
		}
	}
	return nil
}

func PrefixIterOptions(prefix string) *pebble.IterOptions {
	bytes := []byte(prefix)
	return &pebble.IterOptions{LowerBound: bytes, UpperBound: KeyUpperBound(bytes)}
}
