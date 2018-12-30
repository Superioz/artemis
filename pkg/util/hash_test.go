package util

import (
	"testing"
)

// makes sure that hash values are consistent
func TestHash(t *testing.T) {
	objs := []interface{}{"hello", 1, 42}
	objs2 := []interface{}{0, "k", 1337}

	hash1 := Hash(objs...)
	hash2 := Hash(objs...)
	hash3 := Hash(objs2...)

	if hash1 != hash2 || hash3 == hash1 {
		t.Error("hash function didn't return expected values")
	}
}
