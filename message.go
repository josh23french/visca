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

// MessageType constants
const (
	MsgInvalid       MessageType = iota
	MsgCommand                   // 0x01
	MsgInquiry                   // 0x09
	MsgCancel                    // 0x2y
	MsgNetworkChange             // 0x38
	MsgACK                       // 0x4y
	MsgCompletion                // 0x5y
	MsgError                     // 0x60
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
