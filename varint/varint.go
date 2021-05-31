package varint

import (
	"errors"
	"io"
)

var (
	ErrVarIntTooBig = errors.New("var int is too big")
)

func EncodeVarInt(inputValue int) ([]byte, int) {
	value := uint32(inputValue)

	buffer := []byte{}
	written := 0

	for {
		temp := (byte)(value & 0b01111111)
		value >>= 7

		if value != 0 {
			temp |= 0b10000000
		}

		buffer = append(buffer, temp)
		written++

		if value == 0 {
			break
		}
	}

	return buffer, written
}

func DecodeReaderVarInt(reader io.Reader) (int, int, error) {
	numRead := 0
	result := 0

	read := make([]byte, 1)

	for {
		reader.Read(read)
		readByte := read[0]

		value := int(readByte & 0b01111111)
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
