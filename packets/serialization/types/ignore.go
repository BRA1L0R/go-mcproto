package types

import (
	"bytes"
)

func SerializeIgnore(ignoreLen int, databuf *bytes.Buffer) error {
	if ignoreLen < 0 {
		return ErrIgnoreLenUnknown
	}

	ignoreBuf := make([]byte, ignoreLen)

	_, err := databuf.Write(ignoreBuf)
	return err
}

func DeserializeIgnore(ignoreLen int, databuf *bytes.Buffer) error {
	if ignoreLen < 0 {
		return ErrIgnoreLenUnknown
	}

	ignoreBuf := make([]byte, ignoreLen)

	_, err := databuf.Read(ignoreBuf)
	return err
}
