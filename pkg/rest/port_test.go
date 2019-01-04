package rest

import (
	"fmt"
	"testing"
)

// makes sure that at least one free port can be found
func TestGetFreePort(t *testing.T) {
	sl := []int{2310, 2311, 2312, 2313, 2314}
	port, err := GetFreePort("localhost", sl...)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("found free port at %d\n", port)
}

// makes sure that at least one free port can be found
func TestGetFreePortInRange(t *testing.T) {
	port, err := GetFreePortInRange("localhost", 2310, 2314)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("found free port at %d\n", port)
}
