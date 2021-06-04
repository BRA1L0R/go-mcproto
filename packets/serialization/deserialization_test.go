package serialization_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/BRA1L0R/go-mcproto/packets/serialization"
	"github.com/Tnze/go-mc/nbt"
)

func TestArrDeserialization(t *testing.T) {
	type TestChildStruct struct {
		VarInt1 int32  `mc:"varint"`
		Text    string `mc:"string"`
	}

	type TestStruct struct {
		Nested []TestChildStruct `mc:"array" len:"2"`
	}

	testStruct := new(TestStruct)
	testBuffer := new(bytes.Buffer)

	testBuffer.Write([]byte{0x01, 0x02, 0x59, 0x68, 0x02, 0x00})

	err := serialization.DeserializeFields(reflect.ValueOf(testStruct).Elem(), testBuffer)
	if err != nil {
		t.Error(err)
	}
}

func TestDeserialization(t *testing.T) {
	type NbtStruct struct {
		String1 string `nbt:"stringone"`
		String2 string `nbt:"stringtwo"`
	}

	type TestStruct struct {
		VarInt  int32       `mc:"varint"`
		String  string      `mc:"string"`
		Inherit uint32      `mc:"inherit"`
		Ignore  interface{} `mc:"ignore" len:"6"`
		Bytes   []byte      `mc:"bytes" len:"3"`
		Nbt     NbtStruct   `mc:"nbt"`
		NbtArr  []NbtStruct `mc:"nbt" len:"3"`
	}

	testBuffer := new(bytes.Buffer)
	testStruct := new(TestStruct)

	testBuffer.Write([]byte{
		0x80, 0x01, // varint
		0x04, 0x59, 0x59, 0x59, 0x59, // string
		0xff, 0x00, 0xff, 0x00, // inherit
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // ignored
		0x01, 0x02, 0x03, // byte
	})

	err := nbt.Marshal(testBuffer, NbtStruct{String1: "Hello", String2: "World"})
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		err := nbt.Marshal(testBuffer, NbtStruct{String1: "ArrayTest", String2: "ArrayTest2"})
		if err != nil {
			t.Fatal(err)
		}
	}

	err = serialization.DeserializeFields(reflect.ValueOf(testStruct).Elem(), testBuffer)
	if err != nil {
		t.Fatal(err)
	}

	// VarInt
	if testStruct.VarInt != 128 {
		t.Fatal("deserialized varint mismatch")
	}

	// String
	if testStruct.String != "YYYY" {
		t.Fatal("deserialized string mismatch")
	}

	// Inherit uint32
	if testStruct.Inherit != 4278255360 {
		t.Fatal("deserialized inherit (uint32) mismatch")
	}

	// Bytes
	if !bytes.Equal(testStruct.Bytes, []byte{0x01, 0x02, 0x03}) {
		t.Fatal("deserialized bytes mismatch")
	}

	if testStruct.Nbt.String1 != "Hello" || testStruct.Nbt.String2 != "World" {
		t.Fatal("deserialized nbt mismatch")
	}

	if len(testStruct.NbtArr) != 3 {
		t.Fatal("nbt array not the desired length")
	}

	for _, v := range testStruct.NbtArr {
		if v.String1 != "ArrayTest" || v.String2 != "ArrayTest2" {
			t.Fatal("deserialized array elements not as encoded")
		}
	}
}
