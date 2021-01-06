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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPacketFromBytes(t *testing.T) {
	pkt, err := PacketFromBytes([]byte{0x81, 0xF4, 0x66, 0xFF})
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, 0, pkt.Source(), "source should be zero")
	assert.Equal(t, 1, pkt.Destination(), "destination should be one")
	assert.Equal(t, []byte{0xF4, 0x66}, pkt.Message, "message should match")

	pkt, err = PacketFromBytes([]byte{0x82, 0x00, 0xFF})
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, 0, pkt.Source(), "source should be zero")
	assert.Equal(t, 2, pkt.Destination(), "destination should be one")
	assert.Equal(t, []byte{0x00}, pkt.Message, "message should match")

	pkt, err = PacketFromBytes([]byte{0xA0, 0x00, 0xFF})
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, 2, pkt.Source(), "source should be two")
	assert.Equal(t, 0, pkt.Destination(), "destination should be zero")
	assert.Equal(t, []byte{0x00}, pkt.Message, "message should match")

	pkt, err = PacketFromBytes([]byte{0xA0, 0x00, 0xF0})
	assert.Nil(t, pkt, "packet should be nil")
	assert.NotNil(t, err, "error should not be nil")
	assert.Equal(t, ErrInvalidVISCAPacket, err, "should get invalid packet error")

	pkt, err = PacketFromBytes([]byte{0x70, 0x00, 0xFF})
	assert.Nil(t, pkt, "packet should be nil")
	assert.NotNil(t, err, "error should not be nil")
	assert.Equal(t, ErrInvalidVISCAPacket, err, "should get invalid packet error")

	pkt, err = PacketFromBytes([]byte{0xAF, 0x00, 0xFF})
	assert.Nil(t, pkt, "packet should be nil")
	assert.NotNil(t, err, "error should not be nil")
	assert.Equal(t, ErrAddressOutOfBounds, err, "should get address out of bounds error")
}

func TestBytes(t *testing.T) {
	pkt, err := NewPacket(0, 1, []byte{0x00})
	assert.Nil(t, err)
	bytes := pkt.Bytes()
	assert.Equal(t, []byte{0x81, 0x00, 0xFF}, bytes)

	pkt, err = NewPacket(1, 0, []byte{0xCC, 0xEF})
	assert.Nil(t, err)
	bytes = pkt.Bytes()
	assert.Equal(t, []byte{0x90, 0xCC, 0xEF, 0xFF}, bytes)
}

func TestIsBroadcast(t *testing.T) {
	pkt, err := NewPacket(0, 1, []byte{0x00})
	assert.Nil(t, err)
	assert.Equal(t, false, pkt.IsBroadcast())

	pkt.destination = 8
	assert.Equal(t, true, pkt.IsBroadcast())
}
