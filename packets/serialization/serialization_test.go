package serialization_test

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
	"testing"

	"github.com/BRA1L0R/go-mcproto/packets/serialization"
	"github.com/BRA1L0R/go-mcproto/varint"
	"github.com/Tnze/go-mc/nbt"
)

func TestBasicSerialization(t *testing.T) {
	type TestStruct struct {
		VarInt int32 `mc:"varint"`
	}

	encodedVarint := int32(0x7B)

	testStruct := TestStruct{VarInt: encodedVarint}
	testBuffer := new(bytes.Buffer)

	err := serialization.SerializeFields(reflect.ValueOf(testStruct), testBuffer)
	if err != nil {
		t.Fatal(err)
	}

	if testBuffer.Len() <= 0 {
		t.Fatal("Buffer length is 0 or below")
	}

	read, err := testBuffer.ReadByte()
	if err != nil {
		t.Fatal(err)
	}

	if int32(read) != encodedVarint {
		t.Fatal("Read varint does not match input varint")
	}
}

func TestFullSerialization(t *testing.T) {
	type NbtStruct struct {
		Test1 string `nbt:"test1"`
		Test2 int32  `nbt:"test2"`
	}

	type TestStruct struct {
		VarInt       int32       `mc:"varint"`
		VarLong      int64       `mc:"varlong"`
		VarString    string      `mc:"string"`
		InheritValue int64       `mc:"inherit"`
		Ignore       interface{} `mc:"ignore" len:"5"`
		Byte         byte        `mc:"inherit"`
		Nbt          NbtStruct   `mc:"nbt"`
		Bytes        []byte      `mc:"bytes"`
	}

	testStruct := TestStruct{
		VarInt:       10,
		VarLong:      -2147483648,
		VarString:    "test",
		InheritValue: -32392839992839239,
		Byte:         0xFA,
		Nbt:          NbtStruct{Test1: "nbt_test", Test2: -1234},
		Bytes:        []byte{0x1, 0x2, 0x3, 0x4},
	}

	testBuffer := new(bytes.Buffer)

	err := serialization.SerializeFields(reflect.ValueOf(testStruct), testBuffer)
	if err != nil {
		t.Fatal(err)
	}

	// VarInt test
	varintDecoded, _, err := varint.DecodeReaderVarInt(testBuffer)
	if err != nil {
		t.Fatal(err)
	}

	if varintDecoded != testStruct.VarInt {
		t.Fatal("VarInt mismatch")
	}

	varlongDecoded, _, err := varint.DecodeReaderVarLong(testBuffer)
	if err != nil {
		t.Fatal(err)
	}

	if varlongDecoded != testStruct.VarLong {
		t.Fatal("VarLong mismatch")
	}

	// VarString testing
	varStringLen, _, err := varint.DecodeReaderVarInt(testBuffer)
	if err != nil {
		t.Fatal(err)
	}

	varString := make([]byte, varStringLen)
	read, err := io.ReadFull(testBuffer, varString)
	if err != nil {
		t.Fatal(err)
	}

	if int32(read) != varStringLen {
		t.Fatal("varstring read mismatch")
	}

	if string(varString) != testStruct.VarString {
		t.Fatal("varstring mismatch")
	}

	// Inherit values (integers encoded in standard bigendian) testing
	var inheritValue int64

	// big endian is the standard format for the minecraft protocol(1)
	// (1) https://wiki.vg/Protocol#Data_types
	err = binary.Read(testBuffer, binary.BigEndian, &inheritValue)
	if err != nil {
		t.Fatal(err)
	}

	if inheritValue != testStruct.InheritValue {
		t.Fatal("inherit value mismatch")
	}

	// Ignore testing
	emptySlice := make([]byte, 5)
	ignoredFields := make([]byte, 5)

	_, err = testBuffer.Read(ignoredFields)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(ignoredFields, emptySlice) {
		t.Fatal("ignore mismatch")
	}

	// Byte testing
	byteDecoded, err := testBuffer.ReadByte()
	if err != nil {
		t.Fatal(err)
	}

	if byteDecoded != testStruct.Byte {
		t.Fatal("byte mismatch")
	}

	// Nbt testing
	decodedNbt := NbtStruct{}
	decoder := nbt.NewDecoder(testBuffer)

	err = decoder.Decode(&decodedNbt)
	if err != nil {
		t.Fatal(err)
	}

	if decodedNbt.Test1 != testStruct.Nbt.Test1 || decodedNbt.Test2 != testStruct.Nbt.Test2 {
		t.Fatal("NBT mismatch")
	}

	// Bytes
	for _, v := range testStruct.Bytes {
		byteRead, err := testBuffer.ReadByte()
		if err != nil {
			t.Fatal(err)
		}

		if byteRead != v {
			t.Fatal("bytes read mismatch")
		}
	}

	// Test for any residual data
	if testBuffer.Len() != 0 {
		t.Fatal("data remaining in the buffer, this means extra data has been encoded which was not expected")
	}
}

func TestUnknownInherit(t *testing.T) {
	type UnknownInheritTest struct {
		UnknownInherit string `mc:"inherit"`
	}

	inheritTest := UnknownInheritTest{UnknownInherit: ""}
	testBuffer := new(bytes.Buffer)

	err := serialization.SerializeFields(reflect.ValueOf(inheritTest), testBuffer)
	if err == nil {
		t.Fatal("SerializeFields did not return an error on an unknown inherit type")
	}
}

func TestBadLenField(t *testing.T) {
	type BadLenTest struct {
		FillerWithBadLen interface{} `mc:"ignore" len:"badvalue"`
	}

	testBuffer := new(bytes.Buffer)
	badLenTest := BadLenTest{FillerWithBadLen: 123}

	err := serialization.SerializeFields(reflect.ValueOf(badLenTest), testBuffer)
	if err == nil {
		t.Fatal("SerializeFields did not return an error on a bad len value")
	}
}

func TestDependency(t *testing.T) {
	type PacketWithDependency struct {
		HasVarint bool  `mc:"inherit"`
		VarInt    int32 `mc:"varint" depends_on:"HasVarint"`
	}

	trueBuffer := new(bytes.Buffer)
	trueDependency := PacketWithDependency{HasVarint: true, VarInt: 0x05}

	err := serialization.SerializeFields(reflect.ValueOf(trueDependency), trueBuffer)
	if err != nil {
		t.Fatal(err)
	}

	if trueBuffer.Len() != 2 {
		t.Fatal("buffer length not as expected")
	}

	falseBuffer := new(bytes.Buffer)
	falseDependency := PacketWithDependency{HasVarint: false}

	err = serialization.SerializeFields(reflect.ValueOf(falseDependency), falseBuffer)
	if err != nil {
		t.Fatal(err)
	}

	if falseBuffer.Len() != 1 {
		t.Fatal("buffer length not as expected")
	}
}

func TestStructArr(t *testing.T) {
	type ChildStruct struct {
		VarInt int32  `mc:"varint"`
		String string `mc:"string"`
	}

	type ParentStruct struct {
		StructArr []ChildStruct `mc:"array"`
	}

	testBuffer := new(bytes.Buffer)
	nestedStructArr := ParentStruct{StructArr: []ChildStruct{
		{VarInt: 1, String: "Hello"},
		{VarInt: 5, String: "World"},
	}}

	err := serialization.SerializeFields(reflect.ValueOf(nestedStructArr), testBuffer)
	if err != nil {
		t.Fatal(err)
	}

	if testBuffer.Len() != 14 {
		t.Fatal("buffer length not as expected")
	}
}
