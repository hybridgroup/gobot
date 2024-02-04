package bleclient

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2"
)

var (
	_ gobot.Adaptor      = (*Adaptor)(nil)
	_ gobot.BLEConnector = (*Adaptor)(nil)
)

func TestNewAdaptor(t *testing.T) {
	a := NewAdaptor("D7:99:5A:26:EC:38")
	assert.Equal(t, "D7:99:5A:26:EC:38", a.Address())
	assert.True(t, strings.HasPrefix(a.Name(), "BLEClient"))
}

func TestName(t *testing.T) {
	a := NewAdaptor("D7:99:5A:26:EC:38")
	a.SetName("awesome")
	assert.Equal(t, "awesome", a.Name())
}
