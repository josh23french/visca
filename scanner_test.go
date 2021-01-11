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
	buffer := bytes.NewBuffer([]byte{0x81, 0x00, 0xFF, 0x90, 0x00, 0xFF, 0x88})

	scanner := NewScanner(buffer)
	packets := scanner.Scan()
	assert.NotNil(t, scanner, "shouldn't be nil")
	assert.Equal(t, 2, len(packets), "should have two packets")

	fmt.Printf("Bytes: %v\n", packets[0])

	pkt, err := PacketFromBytes(packets[0])
	assert.Nil(t, err, "should have no error")
	assert.Equal(t, pkt.Source(), 0)
	assert.Equal(t, pkt.Destination(), 1)
	assert.Equal(t, pkt.Message, []byte{0x00})
}
