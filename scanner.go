//  scanner.go - VISCA packet scanner
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
	"bufio"
	"io"

	"github.com/rs/zerolog/log"
)

// Scanner represents a
type Scanner struct {
	scanner *bufio.Scanner
	buffer  io.Reader
}

func splitPackets(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		// If we found a terminator byte, we have a packet
		if data[i] == 0xFF {
			// Include that FF in the packet!
			return i + 1, data[:i+1], nil
		}
	}
	if !atEOF {
		return 0, nil, nil
	}
	return 0, data, bufio.ErrFinalToken
}

// NewScanner constructs a Scanner
func NewScanner(buffer io.Reader) *Scanner {
	scanner := bufio.NewScanner(buffer)
	scanner.Buffer([]byte{}, 32)
	scanner.Split(splitPackets)

	return &Scanner{
		scanner,
		buffer,
	}
}

// Scan sends packets to the given channel
func (s *Scanner) Scan(c chan *Packet, quit chan struct{}) {
loop:
	for {
		select {
		case <-quit:
			break loop
		default:
			ok := s.scanner.Scan()
			if !ok {
				log.Err(s.scanner.Err()).Msg("Scan stopped")
				break loop
			}
			packetBytes := s.scanner.Bytes()
			if len(packetBytes) == 0 {
				log.Warn().Msg("Packet is empty!!!")
				break loop
			}
			if len(packetBytes) <= 2 {
				log.Warn().Msg("Packet is smaller than 3 bytes!")
				// Skip runts (scanner returns a zero-length slice last in the common case)
				continue
			}
			packet, err := PacketFromBytes(packetBytes)
			if err != nil {
				log.Err(err).Msg("error creating packet from bytes")
				continue
			}
			c <- packet
		}
	}
	close(c)
}
