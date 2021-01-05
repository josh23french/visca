package visca

import (
	"errors"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

// ConnectionType is an enum of connection types
//go:generate stringer -type=ConnectionType
type ConnectionType int

// ConnectionType enum variants
const (
	unknown ConnectionType = iota
	TCP
	UDP
	CharDev
	Unix
	sentinel
)

func connectionFromString(connString string, timeout time.Duration) (io.ReadWriteCloser, error) {
	if info, err := os.Stat(connString); err == nil {
		// connString is a path to a file...
		mode := info.Mode()
		if (mode&os.ModeDevice != 0) && (mode&os.ModeCharDevice != 0) {
			// and now it's a path to a character device!
			return os.Open(connString)
		}

		return nil, errors.New("connection string is a file, but is not a character device")
	}
	if strings.HasPrefix(connString, "udp://") {
		hostPort := strings.TrimPrefix(connString, "udp://")
		return net.DialTimeout("udp", hostPort, timeout)
	}
	if strings.HasPrefix(connString, "unix://") {
		path := strings.TrimPrefix(connString, "unix://")
		return net.DialTimeout("unix", path, timeout)
	}

	// Default is a TCP conn
	hostPort := strings.TrimPrefix(connString, "tcp://")
	return net.DialTimeout("tcp", hostPort, timeout)
}

func (t ConnectionType) isValid() bool {
	return t > unknown && t < sentinel
}

// Camera represents a camera
type Camera struct {
	Name           string
	ConnectionType ConnectionType
	conn           io.ReadWriteCloser
}

// NewCamera creates a new Camera
func NewCamera(name, connString string) (*Camera, error) {
	conn, err := connectionFromString(connString, 2*time.Second)
	if err != nil {
		return nil, err
	}
	return &Camera{
		Name: name,
		conn: conn,
	}, nil
}
