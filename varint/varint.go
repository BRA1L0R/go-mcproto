package varint

import (
	"bytes"
	"errors"
	"io"
)

var (
	ErrVarIntTooBig = errors.New("mcproto: var int is too big")
)

// EncoeVarInt is an implementation of the varint encoding as specified on wiki.vg
//
// It takes an int32 as an input, and returns the encoded varint in the form of
// a byte slice, and also returns n which is the number of bytes the varint takes
func EncodeVarInt(inputValue int32) (varint []byte, n int) {
	value := uint32(inputValue)

	buffer := new(bytes.Buffer)

	for {
		temp := (byte)(value & 0b01111111)
		value >>= 7

		if value != 0 {
			temp |= 0b10000000
		}

		buffer.WriteByte(temp)
		n++

		if value == 0 {
			break
		}
	}

	return buffer.Bytes(), n
}

func EncodeVarLong(inputValue int64) (varint []byte, n int) {
	value := uint64(inputValue)

	buffer := new(bytes.Buffer)

	for {
		temp := (byte)(value & 0b01111111)
		value >>= 7

		if value != 0 {
			temp |= 0b10000000
		}

		buffer.WriteByte(temp)
		n++

		if value == 0 {
			break
		}
	}

	return buffer.Bytes(), n
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
func DecodeReaderVarInt(reader io.Reader) (result int32, numRead int, err error) {
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

func DecodeReaderVarLong(reader io.Reader) (result int64, numRead int, err error) {
	read := make([]byte, 1)

	for {
		_, err := reader.Read(read)
		if err != nil {
			return result, numRead, err
		}

		readByte := read[0]

		value := int64(readByte & 0b01111111)
		result |= (value << (7 * numRead))

		numRead++
		if numRead > 10 {
			// panic("VarInt is too big >:/")
			return 0, numRead, ErrVarIntTooBig
		}

		if (readByte & 0b10000000) == 0 {
			break
		}
	}

	return result, numRead, nil
}
