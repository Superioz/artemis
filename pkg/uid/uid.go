package uid

import (
	"fmt"
	"github.com/satori/go.uuid"
	"regexp"
)

const shortenedSize = 8

var (
	Nil   = UID{}
	valid = regexp.MustCompile("^[a-zA-Z0-9]{8}$")
)

type UID struct {
	repr string
}

func NewUID() UID {
	return UID{
		repr: uuid.NewV4().String()[:shortenedSize],
	}
}

func FromString(str string) (UID, error) {
	if !valid.MatchString(str) {
		return Nil, fmt.Errorf("string does not match format %s", valid)
	}
	return UID{repr: str}, nil
}

func (id UID) String() string {
	return fmt.Sprintf("%s", id.repr)
}
