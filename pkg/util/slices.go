package util

// Inserts given `val` into given `slice` at `index`.
// First, it increases the slice size, to be sure, to not
// have an index out of bounds.
// Afterwars it moves the slice content around, to make the appended
// value to index = `index`.
// Then we can replace this value.
func Insert(slice []interface{}, index int, val interface{}) []interface{} {
	slice = append(slice, 0)
	copy(slice[index+1:], slice[index:])
	slice[index] = val

	// return the finished slice, which should be a copy of the original one
	return slice
}
