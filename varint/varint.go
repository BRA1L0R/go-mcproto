package varint

import (
	"bytes"
	"errors"
	"io"
)

var (
	ErrVarIntTooBig = errors.New("var int is too big")
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
// func DecodeReaderVarInt(reader io.Reader) (int, int, error) {
// 	numRead := 0
// 	result := 0

// 	read := make([]byte, 1)

// 	for {
// 		_, err := reader.Read(read)
// 		if err != nil {
// 			return result, numRead, err
// 		}

// 		readByte := read[0]

// 		value := int(readByte & 0b01111111)
// 		result |= (value << (7 * numRead))

// 		numRead++
// 		if numRead > 5 {
// 			// panic("VarInt is too big >:/")
// 			return 0, numRead, ErrVarIntTooBig
// 		}

// 		if (readByte & 0b10000000) == 0 {
// 			break
// 		}
// 	}

// 	return result, numRead, nil
// }

func DecodeReaderVarInt(r io.Reader) (result int32, read int32, err error) {
	var V uint32
	readbuf := make([]byte, 1)

	for sec := byte(0x80); sec&0x80 != 0; read++ {
		if read > 5 {
			return 0, read, errors.New("VarInt is too big")
		}

		_, err := r.Read(readbuf)
		if err != nil {
			return 0, read, err
		}

		sec = readbuf[0]

		V |= uint32(sec&0x7F) << uint32(7*read)
	}

	return int32(V), read, nil
}
