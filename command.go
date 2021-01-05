package visca

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

// Message is any message that can be sent/received
type Message interface {
	ToBytes() []byte
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
