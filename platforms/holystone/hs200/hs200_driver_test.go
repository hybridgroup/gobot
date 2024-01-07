package hs200

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*Driver)(nil)

func TestHS200Driver(t *testing.T) {
	d := NewDriver("127.0.0.1:8080", "127.0.0.1:9090")

	assert.Equal(t, "127.0.0.1:8080", d.tcpaddress)
	assert.Equal(t, "127.0.0.1:9090", d.udpaddress)
}
