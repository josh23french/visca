package visca

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestControllerInvalidCamera(t *testing.T) {
	ctrl := NewController()

	conn1 := &MockConnection{}
	conn1.On("SetReceiveQueue", mock.Anything).Return()
	conn1.On("Start").Return(nil)
	conn1.On("Stop").Return()

	// AddCamera
	// Too high
	err := ctrl.AddCamera(9, conn1)
	assert.NotNil(t, err)

	// Too low
	err = ctrl.AddCamera(0, conn1)
	assert.NotNil(t, err)

	// Just right
	err = ctrl.AddCamera(4, conn1)
	assert.Nil(t, err)

	// RemoveCamera
	// Too high
	err = ctrl.RemoveCamera(10)
	assert.Equal(t, ErrInvalidCameraNumber, err)

	// Too low
	err = ctrl.RemoveCamera(0)
	assert.Equal(t, ErrInvalidCameraNumber, err)

	// Just right
	err = ctrl.RemoveCamera(4)
	assert.Nil(t, err)

	// Just right (but doesn't exist)
	err = ctrl.RemoveCamera(5)
	assert.Equal(t, ErrNoCameraConnection, err)
}

func TestController(t *testing.T) {
	conn1 := &MockConnection{}
	conn2 := &MockConnection{}

	ctrl := NewController()

	conn1.On("SetReceiveQueue", mock.Anything).Return()
	conn2.On("SetReceiveQueue", mock.Anything).Return()

	conn1.On("Start").Return(nil)
	conn2.On("Start").Return(nil)

	ctrl.AddCamera(1, conn1)
	ctrl.AddCamera(2, conn2)
	ctrl.AddCamera(3, conn1)

	ctrl.AddCamera(4, conn1)
	conn1.On("Stop").Once().Return(nil)
	err := ctrl.RemoveCamera(4)
	assert.Nil(t, err)
	conn1.AssertExpectations(t)

	ctrl.SetCamera(1)

	// Check PresetRecall(1)
	conn1.On("Send", &Packet{
		source:      0,
		destination: 1,
		Message: []byte{
			0x01, 0x04, 0x3f, 0x02, 0x01,
		},
	}).Return(nil).Once()
	ctrl.PresetRecall(1)
	conn1.AssertExpectations(t)

	// Check PresetRecall(2)
	conn1.On("Send", &Packet{
		source:      0,
		destination: 1,
		Message: []byte{
			0x01, 0x04, 0x3f, 0x02, 0x02,
		},
	}).Return(nil).Once()
	ctrl.PresetRecall(2)
	conn1.AssertExpectations(t)

	ctrl.SetCamera(2)
	// Check PresetRecall(2) after changing camera to 2
	conn2.On("Send", &Packet{
		source:      0,
		destination: 2,
		Message: []byte{
			0x01, 0x04, 0x07, 0x00,
		},
	}).Return(nil).Once()
	ctrl.ZoomStop()
	conn2.AssertExpectations(t)

	ctrl.SetCamera(3)
	// Check PresetRecall(12) after changing camera to 3
	conn1.On("Send", &Packet{
		source:      0,
		destination: 3,
		Message: []byte{
			0x01, 0x04, 0x3f, 0x02, 0x0c,
		},
	}).Return(nil).Once()
	ctrl.PresetRecall(12)
	conn1.AssertExpectations(t)

	assert.Equal(t, ctrl, ctrl)
}

func TestSendMessage(t *testing.T) {
	ctrl := NewController()

	conn := &MockConnection{}
	conn.On("SetReceiveQueue", mock.Anything).Return(nil)
	conn.On("Start").Return(nil)

	ctrl.AddCamera(1, conn)

	// Invalid packet
	err := ctrl.sendMessage([]byte{})
	assert.Equal(t, ErrInvalidVISCAPacket, err)

	// Valid packet
	conn.On("Send", &Packet{
		source:      0,
		destination: 1,
		Message:     []byte{0x50},
	}).Return(nil).Once()
	err = ctrl.sendMessage([]byte{0x50})
	assert.Nil(t, err)
	conn.AssertExpectations(t)
}
