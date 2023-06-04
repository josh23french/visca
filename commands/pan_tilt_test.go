package commands

import (
	"testing"

	"github.com/josh23french/visca"
	"github.com/stretchr/testify/assert"
)

func TestPanTiltParams(t *testing.T) {
	cmd := PanTiltUp{}

	err := cmd.SetPanSpeed(0x13)
	assert.Nil(t, err)

	err = cmd.SetTiltSpeed(0x15)
	assert.Nil(t, err)

	assert.Equal(t, 0x13, cmd.PanSpeed())
	assert.Equal(t, 0x15, cmd.TiltSpeed())

	err = cmd.SetPanSpeed(200)
	assert.NotNil(t, err)

	err = cmd.SetTiltSpeed(0)
	assert.NotNil(t, err)

	assert.Equal(t, visca.Message([]byte{0x01, 0x06, 0x01, 0x13, 0x15, 0x03, 0x01}), cmd.Message())
}
