//  connection.go - VISCA connections
// 	Copyright (C) 2021  Joshua French
//
// 	This program is free software: you can redistribute it and/or modify
// 	it under the terms of the GNU Lesser General Public License as published
// 	by the Free Software Foundation, either version 3 of the License, or
// 	(at your option) any later version.
//
// 	This program is distributed in the hope that it will be useful,
// 	but WITHOUT ANY WARRANTY; without even the implied warranty of
// 	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// 	GNU Lesser General Public License for more details.
//
// 	You should have received a copy of the GNU Lesser General Public License
// 	along with this program.  If not, see <https://www.gnu.org/licenses/>.

package visca

import (
	"errors"
	"os"
	"strings"
)

// Connection implements a VISCA connection
type Connection interface {
	Start() error
	Stop()
	Send(*Packet) error
	SetReceiveQueue(chan *Packet)
}

// Error constants
var (
	ErrNotStarted           = errors.New("connection not started")
	ErrIncompletePacketSent = errors.New("incomplete packet sent")
)

// NewConnectionFromString creates a Connection from a string
func NewConnectionFromString(connString string) (Connection, error) {
	if info, err := os.Stat(connString); err == nil {
		// connString is a path to a file...
		mode := info.Mode()
		if (mode&os.ModeDevice != 0) && (mode&os.ModeCharDevice != 0) {
			// and now it's a path to a character device!
			return NewSerialConnection(connString)
		}

		return nil, errors.New("connection string is a file, but is not a character device")
	}
	if strings.HasPrefix(connString, "udp://") {
		hostPort := strings.TrimPrefix(connString, "udp://")
		return NewNetworkConnection("udp", hostPort)
	}
	if strings.HasPrefix(connString, "unix://") {
		path := strings.TrimPrefix(connString, "unix://")
		return NewNetworkConnection("unix", path)
	}

	// Default is a TCP conn
	hostPort := strings.TrimPrefix(connString, "tcp://")
	return NewNetworkConnection("tcp", hostPort)
}
