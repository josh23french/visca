package visca

// NullConnection implements the Connection interface without sending packets anywhere
type NullConnection struct {
	started bool
}

// Start the interface
func (c *NullConnection) Start() error {
	c.started = true
	return nil
}

// Stop the connection
func (c *NullConnection) Stop() {
	if c.started {
		c.started = false
	}
	return
}

// Send a packet
func (c *NullConnection) Send(pkt *Packet) error {
	if !c.started {
		return ErrNotStarted
	}
	return nil
}

// SetReceiveQueue for received packets
func (c *NullConnection) SetReceiveQueue(q chan *Packet) {
	return
}
