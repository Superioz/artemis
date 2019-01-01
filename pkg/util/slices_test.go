package util

import (
	"testing"
)

func TestInsert(t *testing.T) {
	slice := []interface{}{1, 2, 3, 4, 5}
	before := slice[2]

	sl, err := Insert(slice, 2, 69)
	slice = sl.([]interface{})

	if err != nil {
		t.Error(err)
	}
	if slice[2] == before {
		t.Error("inserted value is not equal to expected value")
	}
}

func TestDelete(t *testing.T) {
	slice := []interface{}{1, 2, 3, 4, 5}
	sl, err := Delete(slice, 1, len(slice)-1)
	if err != nil {
		t.Error(err)
	}
	slice = sl.([]interface{})

	if len(slice) != 1 || slice[0] != 1 {
		t.Error("unexpected slice length or content")
	}
}
