package ble

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func initTestBLESerialPort() *SerialPort {
	return NewSerialPort("TEST123", "123", "456")
}

func TestBLESerialPort(t *testing.T) {
	d := initTestBLESerialPort()
	assert.Equal(t, "TEST123", d.Address())
}
