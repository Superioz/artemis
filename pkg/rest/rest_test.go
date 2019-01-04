package rest

import "testing"

func TestNew(t *testing.T) {
	server := New("localhost", 2310)
	err := server.Up()
	if err != nil {
		t.Fatal(err)
	}

	// TODO unfinished
}
