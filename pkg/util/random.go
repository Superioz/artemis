package util

import (
	"math/rand"
)

// generates a random integer between `min` and `max`
// with given interfaces as seed
func RandInt(min int, max int, seed ...interface{}) int {
	var r *rand.Rand
	if len(seed) != 0 {
		r = rand.New(rand.NewSource(int64(Hash(seed))))
	}
	return r.Intn(max-min) + min
}
