package varint

import (
	"bytes"
	"errors"
	"io"
)

var (
	ErrVarIntTooBig = errors.New("mcproto: var int is too big")
)

func EncodeVarInt(inputValue int32) ([]byte, int32) {
	value := uint32(inputValue)

	buffer := new(bytes.Buffer)
	written := int32(0)

	for {
		temp := (byte)(value & 0b01111111)
		value >>= 7

		if value != 0 {
			temp |= 0b10000000
		}

		// buffer = append(buffer, temp)
		buffer.WriteByte(temp)
		written++

		if value == 0 {
			break
		}
	}

	return buffer.Bytes(), written
}

// DecodeReaderVarInt takes an io.Reader as a parameter and returns in order:
//
// - the result varint
//
// - the number of bytes read
//
// - an eventual error
//
// NOTE: It returns the number of bytes read even in occurance of an error
func DecodeReaderVarInt(reader io.Reader) (result int32, numRead int32, err error) {
	read := make([]byte, 1)

	for {
		_, err := reader.Read(read)
		if err != nil {
			return result, numRead, err
		}

		readByte := read[0]

		value := int32(readByte & 0b01111111)
		result |= (value << (7 * numRead))

		numRead++
		if numRead > 5 {
			// panic("VarInt is too big >:/")
			return 0, numRead, ErrVarIntTooBig
		}

		if (readByte & 0b10000000) == 0 {
			break
		}
	}

	return result, numRead, nil
}
