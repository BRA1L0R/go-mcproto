package packets

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"

	"github.com/BRA1L0R/go-mcproto/varint"
)

var (
	maxPacketLength    = 2097151 // https://wiki.vg/Protocol#Packet_format
	ErrMaxPacketLength = errors.New("mcproto: counterpart sent a packet which was too big")
)

type MinecraftRawPacket struct {
	packetLength int
	dataLength   int
	data         []byte
}

// NewReader opens a bytes.Buffer or zlib readcloser depending if the packet
// is compressed or not
func (rp *MinecraftRawPacket) NewReader() (io.ReadCloser, error) {
	var reader io.ReadCloser

	// create a default reader from the data buffer
	reader = io.NopCloser(bytes.NewBuffer(rp.data))

	// if the data length is >0 it means that the data is compressed
	// using zlib compression
	if rp.dataLength > 0 {
		zlibReader, err := zlib.NewReader(reader)
		if err != nil {
			return nil, err
		}

		reader = zlibReader
	}

	return reader, nil
}

// ReadPacketId opens a reader, reads the packet id but doesn't go further
func (rp *MinecraftRawPacket) ReadPacketId() (int32, error) {
	reader, err := rp.NewReader()
	if err != nil {
		return 0, err
	}

	defer reader.Close()

	// packetid is always the first ever value to appear in data, in both formats
	packetId, _, err := varint.DecodeReaderVarInt(reader)
	return packetId, err
}

// ReadAll reads all the available content in the packet,
// returning data as an uncompressed byte slice
func (rp *MinecraftRawPacket) ReadAll() (packetId int32, data []byte, err error) {
	reader, err := rp.NewReader()
	if err != nil {
		return
	}

	defer reader.Close()

	packetId, _, err = varint.DecodeReaderVarInt(reader)
	if err != nil {
		return
	}

	data, err = io.ReadAll(reader)
	return
}

// Write compress calculates packet length and writes
// the complete packet to writer using the uncompressed format
func (rp *MinecraftRawPacket) WriteUncompressed(writer io.Writer) error {
	rp.packetLength = len(rp.data)

	encPlen, _ := varint.EncodeVarInt(int32(rp.packetLength))

	// Write packet length
	if _, err := writer.Write(encPlen); err != nil {
		return err
	}

	// Write packet data (packetid + data)
	if _, err := writer.Write(rp.data); err != nil {
		return err
	}

	return nil
}

// WriteCompressed calculates the packetlength based on the dataLength and the
// dimension of the compressed/uncompressed data.
func (rp *MinecraftRawPacket) WriteCompressed(writer io.Writer) error {
	encDlen, encDlenLen := varint.EncodeVarInt(int32(rp.dataLength))
	rp.packetLength = encDlenLen + len(rp.data)

	encPlen, _ := varint.EncodeVarInt(int32(rp.packetLength))

	// Write packet length
	if _, err := writer.Write(encPlen); err != nil {
		return err
	}

	// Write dataLength
	if _, err := writer.Write(encDlen); err != nil {
		return err
	}

	// Write remaining packet data
	if _, err := writer.Write(rp.data); err != nil {
		return err
	}

	return nil
}

func readLength(reader io.Reader) (res int, n int, err error) {
	decoded, n, err := varint.DecodeReaderVarInt(reader)
	if decoded > int32(maxPacketLength) {
		return 0, 0, ErrMaxPacketLength
	}

	res = int(decoded)
	return
}

// FromUncompressedReader reads a packet in its uncompressed format (without dataLength)
// and returns a new MinecraftRawPacket
func FromUncompressedReader(reader io.Reader) (*MinecraftRawPacket, error) {
	packet := new(MinecraftRawPacket)

	length, _, err := readLength(reader)
	if err != nil {
		return nil, err
	}

	packet.packetLength = length
	packet.dataLength = -1

	packet.data = make([]byte, length)
	_, err = io.ReadFull(reader, packet.data)

	return packet, err
}

// FromCompressedReader reads a packet in its compressed format, but it does not necessarily
// mean the packet must be compressed. Can be uncompressed later with NewReader
func FromCompressedReader(reader io.Reader) (*MinecraftRawPacket, error) {
	packet := new(MinecraftRawPacket)

	length, _, err := readLength(reader)
	if err != nil {
		return nil, err
	}

	dataLength, dLenLen, err := readLength(reader)
	if err != nil {
		return nil, err
	}

	packet.packetLength = length
	packet.dataLength = dataLength

	packet.data = make([]byte, length-dLenLen)
	_, err = io.ReadFull(reader, packet.data)

	return packet, err
}
