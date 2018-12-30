package util

import (
	"testing"
	"time"
)

// makes sure that random integer values with same seed
// return the same value as well.
func TestRandInt(t *testing.T) {
	seed := []interface{}{time.Now().Unix(), 1337, 42}

	n1 := RandInt(0, 100, seed...)
	n2 := RandInt(0, 100, seed...)

	if n1 != n2 {
		t.Error("random values are not of expected behaviour")
	}
}
