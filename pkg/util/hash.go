package util

import (
	"fmt"
	"hash/fnv"
)

// hashes given interfaces into
// an int value. uint32 to be exact.
func Hash(o ...interface{}) uint32 {
	digester := fnv.New32()
	for _, v := range o {
		// put value into hash stream
		_, err := fmt.Fprint(digester, v)
		if err != nil {
			break
		}
	}
	return digester.Sum32()
}
