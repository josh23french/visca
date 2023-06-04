//  serial.go - serial interface
// 	Copyright (C) 2020  Joshua French
//
// 	This program is free software: you can redistribute it and/or modify
// 	it under the terms of the GNU Affero General Public License as published
// 	by the Free Software Foundation, either version 3 of the License, or
// 	(at your option) any later version.
//
// 	This program is distributed in the hope that it will be useful,
// 	but WITHOUT ANY WARRANTY; without even the implied warranty of
// 	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// 	GNU Affero General Public License for more details.
//
// 	You should have received a copy of the GNU Affero General Public License
// 	along with this program.  If not, see <https://www.gnu.org/licenses/>.

package connection

import (
	"github.com/josh23french/visca"
	"github.com/rs/zerolog/log"
	"go.bug.st/serial"
)

// Serial implements the Iface interface for serial VISCA connections
type Serial struct {
	device       string
	port         *serial.Port
	scanner      *visca.Scanner
	receiveQueue chan *visca.Packet
	quit         chan struct{}
}

// NewSerial creates a new SerialIface
func NewSerial(device string) (*Serial, error) {
	return &Serial{
		device:       device,
		port:         nil,
		scanner:      nil,
		receiveQueue: nil,
		quit:         make(chan struct{}),
	}, nil
}

// Start the interface
func (i *Serial) Start() error {
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

	i.scanner = visca.NewScanner(*i.port)

	go i.scanner.Scan(i.receiveQueue, i.quit)
	log.Info().Msgf("Started read loop from serial interface %v", i.device)

	return nil
}

// Stop the interface
func (i *Serial) Stop() {
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
func (i *Serial) Send(pkt *visca.Packet) error {
	if i.port == nil {
		log.Warn().Msg("not started")
		return ErrIfaceNotStarted
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
func (i *Serial) SetReceiveQueue(q chan *visca.Packet) {
	i.receiveQueue = q
	return
}
