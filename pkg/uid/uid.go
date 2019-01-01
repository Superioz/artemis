package uid

import (
	"fmt"
	"github.com/satori/go.uuid"
	"regexp"
)

const shortenedSize = 8

var (
	// the nil value to check if a value is the
	// default struct value
	Nil = UID{}

	// the validation for the shortened uuid
	valid = regexp.MustCompile("^[a-zA-Z0-9]{8}$")
)

// represents a shortened form of a uuid
// very unlikely to not be unique.
type UID struct {
	repr string
}

// creates a new uid instance with a random
// uuid as backup value.
// We use this backup value to create the shortened
// string from
func NewUID() UID {
	return UID{
		repr: uuid.NewV4().String()[:shortenedSize],
	}
}

// returns a uid which contains given str as representation
func FromString(str string) (UID, error) {
	if !valid.MatchString(str) {
		return Nil, fmt.Errorf("string does not match format %s", valid)
	}
	return UID{repr: str}, nil
}

// String() method to overwrite.
func (id UID) String() string {
	return fmt.Sprintf("%s", id.repr)
}
