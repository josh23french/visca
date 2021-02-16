//  packet.go - structures and functions for dealing with VISCA packets
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
)

// Terminator is the last byte in a VISCA packet
const Terminator = 0xFF

// Errors
var (
	ErrInvalidVISCAPacket = errors.New("Invalid VISCA packet")
	ErrAddressOutOfBounds = errors.New("VISCA Address out of bounds")
)

// Packet is a high-level representation of a VISCA packet; source/dest are private to ensure packets are always valid
type Packet struct {
	source      int     // Where the packet came from; 0-7
	destination int     // Where the packet is going; 0-7 or 8 for broadcast
	Message     Message // The VISCA command or inquiry (to be decoded at a higher level)
}

// NewPacket constructs a new Packet or returns an error
func NewPacket(source, destination int, message []byte) (*Packet, error) {
	if source > 7 || destination > 8 {
		return nil, ErrAddressOutOfBounds
	}
	return &Packet{
		source:      source,
		destination: destination,
		Message:     message,
	}, nil
}

// PacketFromBytes constructs a Packet from raw bytes or returns an error
func PacketFromBytes(byteArr []byte) (*Packet, error) {
	// Check for packet small/big
	if len(byteArr) < 3 || len(byteArr) > 16 {
		return nil, ErrInvalidVISCAPacket
	}

	// Check for valid packet footer
	if byteArr[len(byteArr)-1] != Terminator {
		return nil, ErrInvalidVISCAPacket
	}

	header := uint8(byteArr[0])

	if (header & 0b1000_0000) != 0b1000_0000 {
		return nil, ErrInvalidVISCAPacket
	}

	// source can only be valid at this point... 0-7 since it's only 3 bits
	source := int((header >> 4) & 0b0000_0111)
	destination := int(header & 0b0000_1111)

	// NewPacket will do some validation for us
	return NewPacket(source, destination, byteArr[1:len(byteArr)-1])
}

// Bytes returns the []byte representation of the packet
func (p *Packet) Bytes() (bytes []byte) {
	// Create the header
	leftNibble := uint8(8) + uint8(p.source)
	rightNibble := uint8(p.destination)
	header := leftNibble<<4 + rightNibble

	// Construct the packet
	bytes = make([]byte, 0)
	bytes = append(bytes, header)
	bytes = append(bytes, p.Message...)
	bytes = append(bytes, Terminator)
	return
}

// IsBroadcast returns true if the packet is a broadcast packet, false otherwise
func (p *Packet) IsBroadcast() bool {
	return p.destination == 8
}

// Source returns the source address
func (p *Packet) Source() int {
	return p.source
}

// Destination returns the destination address
func (p *Packet) Destination() int {
	return p.destination
}
