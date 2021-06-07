package varint_test

import (
	"bytes"
	"testing"

	"github.com/BRA1L0R/go-mcproto/varint"
)

func TestVarInt(t *testing.T) {

	readTest := bytes.NewBuffer([]byte{
		0x00,
		0x01,
		0x02,
		0x03,
		0x80, 0x01,
		0xff, 0x01,
		0xff, 0xff, 0x7f,
		0xff, 0xff, 0xff, 0xff, 0x07,
		0xff, 0xff, 0xff, 0xff, 0x0f,
		0x80, 0x80, 0x80, 0x80, 0x08,
	})

	writeTest := []int32{
		0,
		1,
		2,
		3,
		128,
		255,
		2097151,
		2147483647,
		-1,
		-2147483648,
	}

	for _, w := range writeTest {
		result, _, err := varint.DecodeReaderVarInt(readTest)
		if err != nil {
			t.Fatal(err)
		}

		if result != w {
			t.Fatal("VarInt mismatch")
		}
	}

	for i := int32(-1000); i < 1000; i++ {
		test, bytesWritten := varint.EncodeVarInt(i)

		buf := bytes.NewBuffer(test)
		varintDecoded, bytesRead, err := varint.DecodeReaderVarInt(buf)

		if err != nil {
			t.Fatal(err)
		}

		if varintDecoded != i {
			t.Fatal("varint mismatch, encoded:", i, "decoded:", varintDecoded)
		}

		if bytesRead != bytesWritten {
			t.Fatal("bytes count mismatch")
		}
	}
}

func TestVarLong(t *testing.T) {
	for i := int64(-70263); i < 78912; i++ {
		test, bytesWritten := varint.EncodeVarLong(i)

		buf := bytes.NewBuffer(test)
		varlongDecoded, bytesRead, err := varint.DecodeReaderVarLong(buf)

		if err != nil {
			t.Fatal(err)
		}

		if varlongDecoded != i {
			t.Fatal("varlong mismatch, encoded:", i, "decoded:", varlongDecoded)
		}

		if bytesWritten != bytesRead {
			t.Fatal("bytes count mismatch")
		}
	}
}
