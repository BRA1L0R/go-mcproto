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
		VarInt int32 `type:"varint"`
	}

	encodedVarint := int32(0x7B)

	testStruct := TestStruct{VarInt: encodedVarint}
	testBuffer := new(bytes.Buffer)

	err := serialization.SerializeFields(reflect.ValueOf(testStruct), testBuffer)
	if err != nil {
		t.Error(err)
	}

	if testBuffer.Len() <= 0 {
		t.Error("Buffer length is 0 or below")
	}

	read, err := testBuffer.ReadByte()
	if err != nil {
		t.Error(err)
	}

	if int32(read) != encodedVarint {
		t.Error("Read varint does not match input varint")
	}
}

func TestFullSerialization(t *testing.T) {
	type NbtStruct struct {
		Test1 string `nbt:"test1"`
		Test2 int32  `nbt:"test2"`
	}

	type TestStruct struct {
		VarInt       int32       `type:"varint"`
		VarString    string      `type:"string"`
		InheritValue int64       `type:"inherit"`
		Ignore       interface{} `type:"ignore" len:"5"`
		Byte         byte        `type:"inherit"`
		Nbt          NbtStruct   `type:"nbt"`
		VarIntArr    []int32     `type:"varint"`
	}

	testStruct := TestStruct{
		VarInt:       10,
		VarString:    "test",
		InheritValue: -32392839992839239,
		Byte:         0xFA,
		Nbt:          NbtStruct{Test1: "nbt_test", Test2: -1234},
		VarIntArr:    []int32{34, 12, 10, 56},
	}

	testBuffer := new(bytes.Buffer)

	err := serialization.SerializeFields(reflect.ValueOf(testStruct), testBuffer)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// VarInt test
	varintDecoded, _, err := varint.DecodeReaderVarInt(testBuffer)
	if err != nil {
		t.Error(err)
	}

	if varintDecoded != testStruct.VarInt {
		t.Error("VarInt mismatch")
	}

	// VarString testing
	varStringLen, _, err := varint.DecodeReaderVarInt(testBuffer)
	if err != nil {
		t.Error(err)
	}

	varString := make([]byte, varStringLen)
	read, err := io.ReadFull(testBuffer, varString)
	if err != nil {
		t.Error(err)
	}

	if int32(read) != varStringLen {
		t.Error("varstring read mismatch")
	}

	if string(varString) != testStruct.VarString {
		t.Error("varstring mismatch")
	}

	// Inherit values (integers encoded in standard bigendian) testing
	var inheritValue int64

	// big endian is the standard format for the minecraft protocol(1)
	// (1) https://wiki.vg/Protocol#Data_types
	err = binary.Read(testBuffer, binary.BigEndian, &inheritValue)
	if err != nil {
		t.Error(err)
	}

	if inheritValue != testStruct.InheritValue {
		t.Error("inherit value mismatch")
	}

	// Ignore testing
	emptySlice := make([]byte, 5)
	ignoredFields := make([]byte, 5)

	_, err = testBuffer.Read(ignoredFields)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(ignoredFields, emptySlice) {
		t.Error("ignore mismatch")
	}

	// Byte testing
	byteDecoded, err := testBuffer.ReadByte()
	if err != nil {
		t.Error(err)
	}

	if byteDecoded != testStruct.Byte {
		t.Error("byte mismatch")
	}

	// Nbt testing
	decodedNbt := NbtStruct{}
	decoder := nbt.NewDecoder(testBuffer)

	err = decoder.Decode(&decodedNbt)
	if err != nil {
		t.Error(err)
	}

	if decodedNbt.Test1 != testStruct.Nbt.Test1 || decodedNbt.Test2 != testStruct.Nbt.Test2 {
		t.Error("NBT mismatch")
	}

	// VarIntArr testing
	for _, v := range testStruct.VarIntArr {
		decodedVarint, _, err := varint.DecodeReaderVarInt(testBuffer)
		if err != nil {
			t.Error(err)
		}

		if v != decodedVarint {
			t.Error("VarInt array element mismatch")
		}
	}

	// Test for any residual data
	if testBuffer.Len() != 0 {
		t.Error("data remaining in the buffer, this means extra data has been encoded which was not expected")
	}
}

func TestUnknownInherit(t *testing.T) {
	type UnknownInheritTest struct {
		UnknownInherit string `type:"inherit"`
	}

	inheritTest := UnknownInheritTest{UnknownInherit: ""}
	testBuffer := new(bytes.Buffer)

	err := serialization.SerializeFields(reflect.ValueOf(inheritTest), testBuffer)
	if err == nil {
		t.Error("SerializeFields did not return an error on an unknown inherit type")
	}
}

func TestBadLenField(t *testing.T) {
	type BadLenTest struct {
		FillerWithBadLen interface{} `type:"ignore" len:"badvalue"`
	}

	testBuffer := new(bytes.Buffer)
	badLenTest := BadLenTest{FillerWithBadLen: 123}

	err := serialization.SerializeFields(reflect.ValueOf(badLenTest), testBuffer)
	if err == nil {
		t.Error("SerializeFields did not return an error on a bad len value")
	}
}

func TestDependency(t *testing.T) {
	type PacketWithDependency struct {
		HasVarint bool  `type:"inherit"`
		VarInt    int32 `type:"varint" depends_on:"HasVarint"`
	}

	trueBuffer := new(bytes.Buffer)
	trueDependency := PacketWithDependency{HasVarint: true, VarInt: 0x05}

	err := serialization.SerializeFields(reflect.ValueOf(trueDependency), trueBuffer)
	if err != nil {
		t.Error(err)
	}

	if trueBuffer.Len() != 2 {
		t.Error("buffer length not as expected")
	}

	falseBuffer := new(bytes.Buffer)
	falseDependency := PacketWithDependency{HasVarint: false}

	err = serialization.SerializeFields(reflect.ValueOf(falseDependency), falseBuffer)
	if err != nil {
		t.Error(err)
	}

	if falseBuffer.Len() != 1 {
		t.Error("buffer length not as expected")
	}
}

func TestStructArr(t *testing.T) {
	type ChildStruct struct {
		VarInt int32  `type:"varint"`
		String string `type:"string"`
	}

	type ParentStruct struct {
		StructArr []ChildStruct `type:"array"`
	}

	testBuffer := new(bytes.Buffer)
	nestedStructArr := ParentStruct{StructArr: []ChildStruct{
		{VarInt: 1, String: "Hello"},
		{VarInt: 5, String: "World"},
	}}

	err := serialization.SerializeFields(reflect.ValueOf(nestedStructArr), testBuffer)
	if err != nil {
		t.Error(err)
	}

	if testBuffer.Len() != 14 {
		t.Error("buffer length not as expected")
	}
}
