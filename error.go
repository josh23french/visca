//  error.go - VISCA error messages
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

import "errors"

// Error is an "enum"
type Error byte

// VISCAError constants
const (
	MessageLengthError   Error = 0x01
	SyntaxError          Error = 0x02
	CommandBufferFull    Error = 0x03
	CommandCanceled      Error = 0x04
	NoSocket             Error = 0x05
	CommandNotExecutable Error = 0x41
)

// NewErrorMessage creates an error message
func NewErrorMessage(socket uint8, err Error) (Message, error) {
	if socket > 2 {
		return nil, errors.New("invalid socket number")
	}

	return []byte{
		0x60 + 0b0000_0011&socket,
		byte(err),
	}, nil
}
