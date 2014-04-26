package gobotJoystick

import (
	"github.com/hybridgroup/go-sdl2/sdl"
	"github.com/hybridgroup/gobot"
)

type JoystickAdaptor struct {
	gobot.Adaptor
	joystick *sdl.Joystick
}

func (me *JoystickAdaptor) Connect() bool {
	sdl.Init(sdl.INIT_JOYSTICK)
	if sdl.NumJoysticks() > 0 {
		me.joystick = sdl.JoystickOpen(0)
	} else {
		panic("No joystick available")
	}
	return true
}

func (me *JoystickAdaptor) Reconnect() bool {
	return true
}

func (me *JoystickAdaptor) Disconnect() bool {
	me.joystick.Close()
	return true
}

func (me *JoystickAdaptor) Finalize() bool {
	return true
}
