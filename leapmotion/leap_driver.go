package gobotLeap

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"github.com/hybridgroup/gobot"
)

type LeapDriver struct {
	gobot.Driver
	LeapAdaptor *LeapAdaptor
}

func NewLeap(adaptor *LeapAdaptor) *LeapDriver {
	d := new(LeapDriver)
	d.Events = make(map[string]chan interface{})
	d.LeapAdaptor = adaptor
	d.Commands = []string{}
	return d
}

func (me *LeapDriver) Start() bool {
	me.Events["Message"] = make(chan interface{})
	enableGestures := map[string]bool{"enableGestures": true}
	b, _ := json.Marshal(enableGestures)
	_, err := me.LeapAdaptor.Leap.Write(b)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			default:
				var msg []byte
				websocket.Message.Receive(me.LeapAdaptor.Leap, &msg)
				gobot.Publish(me.Events["Message"], me.ParseLeapFrame(msg))
			}
		}
	}()

	return true
}
func (me *LeapDriver) Init() bool { return true }
func (me *LeapDriver) Halt() bool { return true }
