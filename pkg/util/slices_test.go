package util

import (
	"testing"
)

func TestInsert(t *testing.T) {
	slice := []interface{}{1, 2, 3, 4, 5}
	before := slice[2]

	slice = Insert(slice, 2, 69)

	if slice[2] == before {
		t.Error("inserted value is not equal to expected value")
	}
}
