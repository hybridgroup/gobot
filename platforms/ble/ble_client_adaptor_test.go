package ble

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Adaptor = (*ClientAdaptor)(nil)

func TestBLEClientAdaptor(t *testing.T) {
	a := NewClientAdaptor("D7:99:5A:26:EC:38")
	assert.Equal(t, "D7:99:5A:26:EC:38", a.Address())
	assert.True(t, strings.HasPrefix(a.Name(), "BLEClient"))
}

func TestBLEClientAdaptorName(t *testing.T) {
	a := NewClientAdaptor("D7:99:5A:26:EC:38")
	a.SetName("awesome")
	assert.Equal(t, "awesome", a.Name())
}
