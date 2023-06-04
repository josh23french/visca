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
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{
		0x81, 0x00, 0xFF,
		0x90, 0x00, 0xFF,
		0x81, 0xFF, // bad packet shouldn't count; still expecting 2 packets
	})

	scanner := NewScanner(buffer)
	c := make(chan *Packet)
	quit := make(chan struct{})
	defer close(quit)
	go scanner.Scan(c, quit)

	packets := make([]*Packet, 0)
	for packet := range c {
		packets = append(packets, packet)
	}

	assert.Equal(t, 2, len(packets), "should have two packets")

	fmt.Printf("Bytes: %v\n", packets[0])

	pkt := packets[0]
	assert.Equal(t, 0, pkt.Source(), "source should be 0")
	assert.Equal(t, 1, pkt.Destination(), "destination should be 1")
	assert.Equal(t, Message{0x00}, pkt.Message, "message should be 0x00")
}

func TestScannerWithBadPackets(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{
		0x81, 0xFF, // too short
		0x81, 0x00, 0x90, 0x00, 0x00, 0x00, 0x90, 0x00, 0x00, 0x00, 0x90, 0x00, 0x00, 0x00, 0x90, 0x00, 0xFF, // too long
	})

	scanner := NewScanner(buffer)
	c := make(chan *Packet)
	quit := make(chan struct{})
	defer close(quit)
	go scanner.Scan(c, quit)

	packets := make([]*Packet, 0)
	for packet := range c {
		packets = append(packets, packet)
	}

	assert.Equal(t, 0, len(packets), "should have no packets")
}
