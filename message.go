//  message.go
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

// SocketMask AND-ed with the message type reveals the message's Socket
const SocketMask = 0b0000_0011

// MessageType indicates the type of packet
type MessageType int

const (
	// MsgInvalid is NOT a valid message (0x00)
	MsgInvalid MessageType = iota
	// MsgCommand sends operational Commands to the Camera. (0x01)
	MsgCommand
	// MsgInquiry is used for inquiring about the current state of the Camera. (0x09)
	MsgInquiry
	// MsgCancel cancels the command in the given Socket of the Camera (0x2y)
	MsgCancel
	// MsgAddressSet sets the address of a Camera. (0x30)
	//
	// Use when initializing the network. Always broadcasted.
	// It is returned by the last Camera in the series to the controller (still broadcast) with the number of cameras + 1.
	MsgAddressSet
	// MsgNetworkChange is sent by a Camera when a device is removed from or added to the network. (0x38)
	MsgNetworkChange
	// MsgACK is returned by a Camera when it receives a Command. No ACK message is returned for an inquiry. (0x4y)
	MsgACK
	// MsgCompletion is returned by a Camera when execution of Commands or Inquiries is completed. (0x5y)
	MsgCompletion
	// MsgError is returned instead of a Completion Message when a Command or Inquiry could not be executed or failed. (0x60)
	MsgError
)

//go:generate enumer -json -text -yaml -type=MessageType
//go:generate enumer -json -text -yaml -type=CategoryCode

// CategoryCode indicates the category of the packet
type CategoryCode int

// CategoryCode constants
const (
	CatInvalid   CategoryCode = -1
	CatInterface CategoryCode = 0x00
	CatCamera1   CategoryCode = 0x04
	CatCamera2   CategoryCode = 0x05
	CatPanTilter CategoryCode = 0x06
	CatDisplay   CategoryCode = 0x7E
)

// Message is any message that can be sent/received
type Message []byte

// Type returns the type of the packet
func (m Message) Type() MessageType {
	switch uint(m[0]) {
	case 0x01:
		return MsgCommand
	case 0x09:
		return MsgInquiry
	case 0x20, 0x21, 0x22:
		return MsgCancel
	case 0x30:
		return MsgAddressSet
	case 0x38:
		return MsgNetworkChange
	case 0x40, 0x41, 0x42:
		return MsgACK
	case 0x50, 0x51, 0x52:
		return MsgCompletion
	case 0x60, 0x61, 0x62:
		return MsgError
	default:
		return MsgInvalid
	}
}

// Socket returns the socket for the packet, if applicable
func (m Message) Socket() uint8 {
	switch m.Type() {
	case MsgACK:
	case MsgCompletion:
	case MsgError:
		return uint8(m[0]) & 0b0000_0011
	default:
	}
	return 0
}

// SetSocket sets the socket or returns an error
func (m Message) SetSocket(s uint8) error {
	if s > 2 {
		return errors.New("invalid socket number")
	}

	m[0] = m[0]&0b1111_1100 + s&0b0000_0011

	return nil
}

// Category returns the message's category code
func (m Message) Category() CategoryCode {
	if len(m) >= 2 {
		code := m[1]
		cc := CategoryCode(uint(code))
		if cc.IsACategoryCode() {
			return cc
		}
	}
	return CatInvalid
}
