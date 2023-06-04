package connection

import "github.com/josh23french/visca"

// NullIface implements the Iface interface without sending packets anywhere
type NullIface struct {
}

// Start the interface
func (i *NullIface) Start() error {
	return nil
}

// Stop the interface
func (i *NullIface) Stop() {
	return
}

// Send a packet
func (i *NullIface) Send(pkt *visca.Packet) error {
	return nil
}

// SetReceiveQueue for received packets
func (i *NullIface) SetReceiveQueue(q chan *visca.Packet) {
	return
}
