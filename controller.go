//  controller.go
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

	"github.com/rs/zerolog/log"
)

// Error constants
var (
	ErrNoCameraConnection  = errors.New("no camera connection")
	ErrInvalidCameraNumber = errors.New("invalid camera number")
)

// Controller represents a high-level VISCA PTZ controller
//
// Example
//
//  ctrl := NewController()
//  err := ctrl.Start()
//  if err != nil {
//    // Do something
//  }
type Controller struct {
	connections  []Connection
	camera       int
	sendQueue    []Packet
	receiveQueue chan *Packet
	quit         chan struct{}
}

// NewController creates a new controller with no cameras
func NewController() *Controller {
	return &Controller{
		connections:  make([]Connection, 8), // 7 connections total; 0 is not used
		camera:       1,                     // starts with camera 1 selected
		sendQueue:    make([]Packet, 0),     // stores commands/inquiries/cancels to be sent
		receiveQueue: make(chan *Packet),    // channel of incoming packets
		quit:         make(chan struct{}),   // used to stop the processReceiveQueue goroutine
	}
}

// Start the Controller
func (c *Controller) Start() error {
	go c.processReceiveQueue()
	return nil
}

// Stop the Controller
func (c *Controller) Stop() {
	close(c.quit)
}

// processReceiveQueue processes packets from the receiveQueue
func (c *Controller) processReceiveQueue() {
loop:
	for {
		select {
		case <-c.quit:
			break loop
		default:
			pkt := <-c.receiveQueue
			// we're only interested in packets for 0 (controller) or 8 (broadcast, which includes the controller)
			if pkt.destination != 0 && pkt.destination != 8 {
				log.Debug().Msg("ignoring packet not for us")
				continue
			}
			switch pkt.Message.Type() {
			case MsgACK:
				// peek the most recent command from the send queue and send something on the ack channel
				// socket := pkt.Message.Socket()
				// c.sendQueue[socket].pop()
				continue
			case MsgCompletion, MsgError:
				// pop the most recent command from the correct buffer and send the packet on the completion channel
				continue
			default:
				// shouldn't get any other message types to the controller...
				// but there's nothing we can do with them but log them
				log.Warn().Msgf("got %v message from %v", pkt.Message.Type(), pkt.Source())
				continue
			}
		}
	}
	if c.receiveQueue != nil {
		close(c.receiveQueue)
	}
}

// AddCamera adds a camera to the controller
func (c *Controller) AddCamera(num int, camera Connection) error {
	if num > 7 || num <= 0 {
		return ErrInvalidCameraNumber
	}
	c.connections[num] = camera
	camera.SetReceiveQueue(c.receiveQueue)
	camera.Start()
	return nil
}

// RemoveCamera removes a camera from the controller
func (c *Controller) RemoveCamera(num int) error {
	if num > 7 || num <= 0 {
		return ErrInvalidCameraNumber
	}
	if c.connections[num] == nil {
		return ErrNoCameraConnection
	}
	c.connections[num].Stop()
	c.connections[num] = nil
	return nil
}

// SetCamera selects the camera the controller is currently working on
func (c *Controller) SetCamera(num int) {
	c.camera = num
}

// sendMessage crafts a packet from the given Message and sends it to the current Camera
func (c *Controller) sendMessage(msg Message) error {
	conn := c.connections[c.camera]
	if conn == nil {
		return ErrNoCameraConnection
	}
	pkt, err := NewPacket(0, c.camera, msg)
	if err != nil {
		return err
	}
	return conn.Send(pkt)
}

//
// Command Set: PRESET
//

// PresetReset resets the preset on the current camera
func (c *Controller) PresetReset(num uint8) {
	c.sendMessage([]byte{0x01, 0x04, 0x3F, 0x00, num})
}

// PresetSet sets the preset on the current camera
func (c *Controller) PresetSet(num uint8) {
	c.sendMessage([]byte{0x01, 0x04, 0x3F, 0x01, num})
}

// PresetRecall recalls the preset on the current camera
func (c *Controller) PresetRecall(num uint8) {
	c.sendMessage([]byte{0x01, 0x04, 0x3F, 0x02, num})
}

//
// Command Set: Pan Tilt Drive
//

// PanTilt is the high-level PT control
func (c *Controller) PanTilt(pan int, tilt int) {
	if pan == 0 && tilt == 0 {
		c.PanTiltStop()
		return
	}
	c.sendMessage([]byte{0x01, 0x04, 0x07, 0x00})
}

// PanTiltStop stops all PT movement
func (c *Controller) PanTiltStop() {
	c.sendMessage([]byte{0x01, 0x04, 0x07, 0x00})
}

//
// Command Set: Zoom
//

// ZoomStop stops zoom movement
func (c *Controller) ZoomStop() {
	c.sendMessage([]byte{0x01, 0x04, 0x07, 0x00})
}

// ZoomIn (re)starts a zoom in
func (c *Controller) ZoomIn() {
	c.sendMessage([]byte{0x01, 0x04, 0x07, 0x02})
}

// ZoomOut (re)starts a zoom out
func (c *Controller) ZoomOut() {
	c.sendMessage([]byte{0x01, 0x04, 0x07, 0x03})
}

// ZoomTo changes zoom to a specific value
func (c *Controller) ZoomTo(value int) {
	c.sendMessage([]byte{0x01, 0x04, 0x47})
}
