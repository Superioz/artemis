package uid

import (
	"testing"
)

// makes sure that shortened uuids are
// still unique
func TestUID_String(t *testing.T) {
	m := make(map[UID]int)

	for i := 0; i < 10000; i++ {
		uid := NewUID()

		s := m[uid]
		m[uid] = s + 1
	}

	for uid, size := range m {
		if size > 1 {
			t.Fatal("unique ids not unique", uid, size)
		}
	}
}

// makes sure that the nil value can be used
// to compare empty uid structs
func TestNil(t *testing.T) {
	uid := UID{}

	if uid != Nil {
		t.Fatal("nil value not nil enough")
	}
}

// makes sure that uids string can be transformed
// into a uid again
func TestFromString(t *testing.T) {
	uid1 := NewUID()
	uid2, err := FromString(uid1.String())
	if err != nil {
		t.Fatal(err)
	}
	if uid1 != uid2 {
		t.Fatal("uid1 does not match uid2", uid1, uid2)
	}
}
