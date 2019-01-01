package util

import (
	"fmt"
	"reflect"
)

// Inserts given `val` into given `slice` at `index`.
// First, it increases the slice size, to be sure, to not
// have an index out of bounds.
// Afterwars it moves the slice content around, to make the appended
// value to index = `index`.
// Then we can replace this value.
func Insert(slice interface{}, index int, val interface{}) (interface{}, error) {
	inValue := reflect.ValueOf(slice)
	inKind := inValue.Type().Kind()
	if inKind != reflect.Slice && inKind != reflect.Array {
		return nil, fmt.Errorf("can only insert to slices and arrays")
	}

	rs := slice.([]interface{})
	rs = append(rs, 0)
	copy(rs[index+1:], rs[index:])
	rs[index] = val

	// return the finished slice, which should be a copy of the original one
	return rs, nil
}

// deletes items from given slice in range [indexf; indext]
func Delete(slice interface{}, indexf int, indext int) (interface{}, error) {
	inValue := reflect.ValueOf(slice)
	inKind := inValue.Type().Kind()
	if inKind != reflect.Slice && inKind != reflect.Array {
		return nil, fmt.Errorf("can only delete from slices and arrays")
	}

	rs := slice.([]interface{})
	rs = append(rs[:indexf], rs[indext+1:]...)
	return rs, nil
}
