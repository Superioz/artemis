package buffer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// Writes an unsigned int with byte length 2 to a byte buffer
// Uses byte order of BigEndian
func WriteUint16(buf *bytes.Buffer, i uint16) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)

	err := binary.Write(buf, binary.BigEndian, b)
	if err != nil {
		fmt.Println(err)
	}
}

// Reads an unsigned int with byte length 2 from byte buffer
// Uses byte order of BigEndian
func ReadUint16(buf *bytes.Buffer) (uint16, error) {
	if buf.Len() < 2 {
		return 0, errors.New("buffer size lower than 2")
	}
	return binary.BigEndian.Uint16(buf.Next(2)), nil
}
