package packets

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
	"reflect"

	"github.com/BRA1L0R/go-mcproto/packets/serialization"
	"github.com/BRA1L0R/go-mcproto/varint"
)

func getReflectValue(inter interface{}) (val reflect.Value, err error) {
	v := reflect.ValueOf(inter)

	if v.Kind() != reflect.Ptr {
		err = errors.New("mcproto: inter should be a pointer to an interface")
		return
	}

	val = v.Elem()
	return
}

// serialization with the uncompressed format
func (p *MinecraftPacket) serializeUncompressedPacket() ([]byte, error) {
	encodedPacketId, pIdLength := varint.EncodeVarInt(p.PacketID)

	length := pIdLength + p.Data.Len()
	encodedLength, _ := varint.EncodeVarInt(int32(length))

	dataBuffer := new(bytes.Buffer)
	dataBuffer.Write(encodedLength)
	dataBuffer.Write(encodedPacketId)

	_, err := dataBuffer.ReadFrom(p.Data)
	if err != nil {
		return nil, err
	}

	return dataBuffer.Bytes(), nil
}

// serialization with the compressed format.
//
// Returns a slice with the packet serialized in bytes
func (p *MinecraftPacket) serializeCompressedPacket(compressionTreshold int32) ([]byte, error) {
	packetBuffer := new(bytes.Buffer)

	dataSize := int32(p.Data.Len())
	if dataSize >= compressionTreshold {
		err := p.serializeWithCompression(dataSize, packetBuffer)
		if err != nil {
			return nil, err
		}
	} else {
		err := p.serializeWithoutCompression(dataSize, packetBuffer)
		if err != nil {
			return nil, err
		}
	}

	return packetBuffer.Bytes(), nil
}

// serialization with the compressed format but without compression
func (p *MinecraftPacket) serializeWithoutCompression(dataSize int32, packetBuffer *bytes.Buffer) error {
	packetID, packetIDSize := varint.EncodeVarInt(p.PacketID)
	dataLength, dataLengthSize := varint.EncodeVarInt(0) // data length must be 0 in uncompressed packets
	packetLength, _ := varint.EncodeVarInt(int32(packetIDSize) + int32(dataLengthSize) + dataSize)

	packetBuffer.Write(packetLength)
	packetBuffer.Write(dataLength)
	packetBuffer.Write(packetID)

	_, err := packetBuffer.ReadFrom(p.Data)
	if err != nil {
		return err
	}

	return nil
}

// serialization with the compressed format but compressed
func (p *MinecraftPacket) serializeWithCompression(dataSize int32, packetBuffer *bytes.Buffer) error {
	packetID, packetIdSize := varint.EncodeVarInt(p.PacketID)
	dataLength, dataLengthSize := varint.EncodeVarInt(int32(packetIdSize) + dataSize)

	compressedBuf := new(bytes.Buffer)
	writer := zlib.NewWriter(compressedBuf)

	_, err := writer.Write(packetID)
	if err != nil {
		return err
	}

	_, err = writer.Write(p.Data.Bytes())
	if err != nil {
		return err
	}

	// Flush the writer to the buffer
	// and close it
	err = writer.Flush()
	if err != nil {
		return err
	}

	writer.Close()

	// Packet length, in compressed packets, is the data remaining
	// to the end of the packet, not the length of the uncompressed data itself
	packetLength, _ := varint.EncodeVarInt(int32(dataLengthSize + compressedBuf.Len()))

	packetBuffer.Write(packetLength)
	packetBuffer.Write(dataLength)

	_, err = packetBuffer.ReadFrom(compressedBuf)
	if err != nil {
		return err
	}

	return nil
}

// SerializeData serializes all the fields in the struct with a type struct tag
// into the data buffer
func (p *MinecraftPacket) SerializeData(inter interface{}) error {
	// ensure Data is a clean buffer before starting
	p.Data = new(bytes.Buffer)

	val, err := getReflectValue(inter)
	if err != nil {
		return err
	}

	return serialization.SerializeFields(val, p.Data)
}

func (p *MinecraftPacket) UncompressData() error {
	if !p.compressed {
		return errors.New("mcproto: packet data has already been uncompressed")
	}

	uncompressed, err := io.ReadAll(p.compressedData)
	if err != nil {
		return err
	}

	p.Data = bytes.NewBuffer(uncompressed)
	p.compressed = false

	return nil
}

// DeserializeData deserializes the data buffer into the struct fields which
// have a type struct tag
func (p *MinecraftPacket) DeserializeData(inter interface{}) error {
	// if data is compressed proceed to uncompress it
	// fmt.Println(p.compressedBuffer)
	if p.compressed {
		err := p.UncompressData()
		if err != nil {
			return err
		}
	}

	val, err := getReflectValue(inter)
	if err != nil {
		return err
	}

	return serialization.DeserializeFields(val, p.Data)
}
