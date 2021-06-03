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
			t.Error(err)
			t.FailNow()
		}

		if result != w {
			t.Error("VarInt mismatch")
			t.Fail()
		}
	}

	for i := int32(-1000); i < 1000; i++ {
		test, bytesWritten := varint.EncodeVarInt(i)

		buf := bytes.NewBuffer(test)
		varintDecoded, bytesRead, err := varint.DecodeReaderVarInt(buf)

		if err != nil {
			t.Error(err)
		}

		if varintDecoded != i {
			t.Error("varint mismatch, encoded:", i, "decoded:", varintDecoded)
			t.FailNow()
		}

		if bytesRead != bytesWritten {
			t.Error("bytes count mismatch")
			t.FailNow()
		}
	}
}
