package sysfs

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func TestDigitalPin(t *testing.T) {
	lastPath := ""
	lastData := []byte{}

	writeFile = func(path string, data []byte) (i int, err error) {
		lastPath = path
		lastData = data
		return
	}

	readFile = func(path string) (b []byte, err error) {
		lastPath = path
		return []byte("0"), nil
	}

	pin := NewDigitalPin(10, "custom")
	gobot.Assert(t, pin.pin, "10")
	gobot.Assert(t, pin.label, "custom")

	pin = NewDigitalPin(10)
	gobot.Assert(t, pin.label, "gpio10")

	pin.Unexport()
	gobot.Assert(t, lastPath, "/sys/class/gpio/unexport")
	gobot.Assert(t, string(lastData), "10")

	pin.Export()
	gobot.Assert(t, lastPath, "/sys/class/gpio/export")
	gobot.Assert(t, string(lastData), "10")

	pin.Write(1)
	gobot.Assert(t, lastPath, "/sys/class/gpio/gpio10/value")
	gobot.Assert(t, string(lastData), "1")

	pin.SetDirection(IN)
	gobot.Assert(t, lastPath, "/sys/class/gpio/gpio10/direction")
	gobot.Assert(t, string(lastData), "in")

	pin.Direction()
	gobot.Assert(t, pin.direction, "in")
	pin.SetDirection(OUT)
	gobot.Assert(t, pin.direction, "out")

	data, _ := pin.Read()
	gobot.Assert(t, data, 0)
	gobot.Assert(t, lastPath, "/sys/class/gpio/gpio10/value")
}
