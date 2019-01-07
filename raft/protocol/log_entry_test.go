package protocol

import (
	"encoding/json"
	"testing"
)

// makes sure that log entries can be marshalled correctly
func TestLogEntry_MarshalJSON(t *testing.T) {
	entry := LogEntry{
		Index:   1,
		Term:    4,
		Content: []byte("hallo!"),
	}
	j, err := json.Marshal(entry)
	if err != nil {
		t.Fatal(err)
	}
	if j == nil {
		t.Fatal("nil object")
	}

	entry2 := LogEntry{}
	err = json.Unmarshal(j, &entry2)
	if err != nil {
		t.Fatal(err)
	}

	if entry.Index != entry2.Index || entry.Term != entry2.Term || string(entry.Content) != string(entry2.Content) {
		t.Fatal("unexpected unmarshalled object")
	}
}
