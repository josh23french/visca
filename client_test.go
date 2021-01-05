package visca

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestController(t *testing.T) {
	ctrl := NewController()
	assert.Equal(t, ctrl, ctrl)
}
