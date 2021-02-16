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

func TestErrorMessage(t *testing.T) {
	var tests = []struct {
		bytes       []byte
		socket      uint8
		messageType MessageType
		err         Error
	}{
		{[]byte{0x60, 0x01}, 0, MsgError, MessageLengthError},
		{[]byte{0x61, 0x02}, 1, MsgError, SyntaxError},
		{[]byte{0x62, 0x03}, 2, MsgError, CommandBufferFull},
		{[]byte{0x60, 0x04}, 0, MsgError, CommandCanceled},
		{[]byte{0x61, 0x05}, 1, MsgError, NoSocket},
		{[]byte{0x62, 0x41}, 2, MsgError, CommandNotExecutable},
	}

	for _, tt := range tests {
		m := Message(tt.bytes)
		assert.Equal(t, tt.messageType, m.Type())
		assert.Equal(t, tt.socket, m.Socket())
		// assert.Equal(t, tt.err, m.Error())
	}
}
