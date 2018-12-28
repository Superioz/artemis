package buffer

import (
	"bytes"
	"testing"
)

func TestWriteUint16(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	WriteUint16(buf, 16)

	val, _ := ReadUint16(buf)
	if val != 16 {
		t.Error("read value is not equal to written value")
	}
}
