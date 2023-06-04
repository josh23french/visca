package commands

import (
	"errors"

	"github.com/josh23french/visca"
)

// PanTiltParams are common to several Pan/Tilt commands
type PanTiltParams struct {
	panSpeed  uint8
	tiltSpeed uint8
}

// SetPanSpeed sets the panSpeed
func (p *PanTiltParams) SetPanSpeed(speed int) error {
	if speed < 0x01 || speed > 0x18 {
		return errors.New("invalid speed")
	}
	p.panSpeed = uint8(speed)
	return nil
}

// SetTiltSpeed sets the tiltSpeed
func (p *PanTiltParams) SetTiltSpeed(speed int) error {
	if speed < 0x01 || speed > 0x18 {
		return errors.New("invalid speed")
	}
	p.tiltSpeed = uint8(speed)
	return nil
}

// PanSpeed returns the panSpeed
func (p *PanTiltParams) PanSpeed() int {
	return int(p.panSpeed)
}

// TiltSpeed returns the panSpeed
func (p *PanTiltParams) TiltSpeed() int {
	return int(p.tiltSpeed)
}

// PanTiltUp tilts the camera up
type PanTiltUp struct {
	PanTiltParams
}

// Message returns the command as a Message
func (c *PanTiltUp) Message() visca.Message {
	return []byte{0x01, 0x06, 0x01, c.panSpeed, c.tiltSpeed, 0x03, 0x01}
}

// ParseCompletion does nothing, this is not an inquiry
func (c *PanTiltUp) ParseCompletion(msg visca.Message) {}
