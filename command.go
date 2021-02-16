package visca

import (
	"errors"
	"math"
)

// Header represents the Visca header
type Header struct {
	From uint8
	To   uint8
}

// ToByte returns the wire representation of the Header
func (h *Header) ToByte() byte {
	var header uint8
	header = 1 << 7
	header += (h.From << 4)
	header += h.To
	return header
}

// Command represents a command to a camera
type Command struct {
	header Header
}

// NewCommand creates a new command
func NewCommand(from, to uint8) *Command {
	return &Command{
		header: Header{
			From: from,
			To:   to,
		},
	}
}

// ToBytes satisfies the Message interface
func (c *Command) ToBytes() []byte {
	buf := make([]byte, 0)
	buf = append(buf, c.header.ToByte())
	return buf
}

func parseIntFromNibbles(b []byte) (int64, error) {
	if len(b) > 16 {
		return 0, errors.New("value does not fit in a int64")
	}

	var n int64
	for i, v := range b {
		shift := uint((len(b) - i - 1) * 4)
		if i == 0 && v&0x08 != 0 {
			n -= 0x08 << shift
			v &= 0x07
		}
		n += int64(v) << shift
	}
	return n, nil
}

func decodePosition(nibbles []byte) (float64, error) {
	divisor := 14.4
	if len(nibbles) == 5 {
		divisor = 235.9
	} else if len(nibbles) != 4 {
		return 0, errors.New("invalid length")
	}
	d, err := parseIntFromNibbles(nibbles)
	return float64(d) / divisor, err
}

func decodeTilt(nibbles []byte) (float64, error) {
	divisor := 235.9
	if len(nibbles) != 4 {
		return 0, errors.New("invalid length")
	}
	d, err := parseIntFromNibbles(nibbles)
	return float64(d) / divisor, err
}

func decodeZoom(nibbles []byte) (float64, error) {
	// divisor := 6144.0
	if len(nibbles) != 4 {
		return 0, errors.New("invalid length")
	}
	d, err := parseIntFromNibbles(nibbles)
	return 6105.543 + 4295.1494*math.Log1p(float64(d)), err
}
