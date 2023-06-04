//  serial.go - serial connection
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
	"github.com/rs/zerolog/log"
	"go.bug.st/serial"
)

// SerialConnection implements the Iface interface for serial VISCA connections
type SerialConnection struct {
	device       string
	port         *serial.Port
	scanner      *Scanner
	receiveQueue chan *Packet
	quit         chan struct{}
}

// NewSerialConnection creates a new SerialIface
func NewSerialConnection(device string) (*SerialConnection, error) {
	return &SerialConnection{
		device:       device,
		port:         nil,
		scanner:      nil,
		receiveQueue: nil,
		quit:         make(chan struct{}),
	}, nil
}

// Start the interface
func (i *SerialConnection) Start() error {
	log.Info().Msgf("Opening serial interface %v...", i.device)
	mode := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: 1,
	}
	port, err := serial.Open(i.device, mode)
	if err != nil {
		return err
	}
	log.Info().Msgf("Opened serial interface %v", i.device)
	i.port = &port

	i.quit = make(chan struct{})

	i.scanner = NewScanner(*i.port)

	go i.scanner.Scan(i.receiveQueue, i.quit)
	log.Info().Msgf("Started read loop from serial interface %v", i.device)

	return nil
}

// Stop the connection
func (i *SerialConnection) Stop() {
	if i.port == nil {
		log.Warn().Msg("Never Started")
		return
	}

	// Stop the receive goroutine first
	close(i.quit)

	// Then close the port
	port := *i.port
	err := port.Close()
	if err != nil {
		log.Warn().Err(err).Msgf("Error stopping serial interface %v", i.device)
	}
}

// Send a packet
func (i *SerialConnection) Send(pkt *Packet) error {
	if i.port == nil {
		log.Warn().Msg("not started")
		return ErrNotStarted
	}
	port := *i.port
	log.Debug().Msgf("Sending packet %v via %v", pkt, i.device)
	written, err := port.Write(pkt.Bytes())
	if err != nil {
		return err
	}

	if written != len(pkt.Bytes()) {
		return ErrIncompletePacketSent
	}

	return nil
}

// SetReceiveQueue for received packets
func (i *SerialConnection) SetReceiveQueue(q chan *Packet) {
	i.receiveQueue = q
	return
}
