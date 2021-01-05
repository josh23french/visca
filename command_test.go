package visca

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderToByteWorks(t *testing.T) {
	h := Header{
		From: 0,
		To:   1,
	}
	assert.Equal(t, byte(0x81), h.ToByte(), "it outputs the expected byte")
}

func TestCommandToBytesWorks(t *testing.T) {
	cmd := NewCommand(0, 1)
	assert.Equal(t, []byte{0x81}, cmd.ToBytes(), "it outputs the expected byte sequence")
}
