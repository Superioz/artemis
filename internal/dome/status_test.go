package dome

import (
	"encoding/json"
	"fmt"
	"github.com/superioz/artemis/pkg/uid"
	"testing"
)

// makes sure that the current status can be marshalled
// and unmarshalled
func TestCurrentStatus_Marshal(t *testing.T) {
	status := Status{
		Id: uid.NewUID(),
	}
	j, err := json.Marshal(status)
	if err != nil {
		t.Fatal(err)
	}
	if j == nil {
		t.Fatal("nil object")
	}
	fmt.Println(string(j))

	status2 := Status{}
	err = json.Unmarshal(j, &status2)
	if err != nil {
		t.Fatal(err)
	}

	if status.Id != status2.Id {
		t.Fatal("unexpected unmarshalled object")
	}
}
