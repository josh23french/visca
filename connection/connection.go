package connection

import (
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/josh23french/visca"
)

// Connection represents a VISCA interface
type Connection interface {
	Start() error                       // Start the Iface; returns error if failed be started
	Stop()                              // Stop the Iface, can have no error
	Send(*visca.Packet) error           // Send a packet out the Iface, may be queued
	SetReceiveQueue(chan *visca.Packet) // Tell the Iface where to send received packets
}

// Error constants
var (
	ErrIfaceNotStarted          = errors.New("iface not started")
	ErrIncompletePacketSent     = errors.New("incomplete packet sent")
	ErrConnectionPathNotCharDev = errors.New("connection path is not a character device")
	ErrConnectionPathInvalid    = errors.New("connection path invalid")
)

// FromString creates a Connection from a string
func FromString(connString string) (Connection, error) {
	if info, err := os.Stat(connString); err == nil {
		// connString is a path to a file...
		mode := info.Mode()
		if (mode&os.ModeDevice != 0) && (mode&os.ModeCharDevice != 0) {
			// and now it's a path to a character device!
			return NewSerial(connString)
		}

		return nil, ErrConnectionPathNotCharDev
	}
	if strings.HasPrefix(connString, "udp://") {
		hostPort := strings.TrimPrefix(connString, "udp://")
		return NewNetworkConnection("udp", hostPort)
	}
	if strings.HasPrefix(connString, "unix://") {
		path := strings.TrimPrefix(connString, "unix://")
		return NewNetworkConnection("unix", path)
	}
	if strings.HasPrefix(connString, "tcp://") {
		hostPort := strings.TrimPrefix(connString, "tcp://")
		return NewNetworkConnection("tcp", hostPort)
	}

	u, err := url.Parse("tcp://" + connString)
	if err != nil {
		println(connString)
		return nil, ErrConnectionPathInvalid
	}

	return NewNetworkConnection("tcp", u.Host)
}
