package tello

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

func TestTelloDriver(t *testing.T) {
	d := NewDriver("8888")

	gobottest.Assert(t, d.respPort, "8888")
}

func statusMessage(msgType uint16, msgAfter7 ...byte) []byte {
	msg := make([]byte, 7, len(msgAfter7)+7)
	msg[0] = messageStart
	binary.LittleEndian.PutUint16(msg[5:7], msgType)
	msg = append(msg, msgAfter7...)
	return msg
}

func TestHandleResponse(t *testing.T) {
	cc := []struct {
		name   string
		msg    io.Reader
		events []gobot.Event
		err    error
	}{
		{
			name: "[empty messsage]",
			msg:  bytes.NewReader(nil),
			err:  io.EOF,
		},
		{
			name:   "wifiMessage",
			msg:    bytes.NewReader(statusMessage(wifiMessage)),
			events: []gobot.Event{{Name: WifiDataEvent}},
		},
		{
			name:   "lightMessage",
			msg:    bytes.NewReader(statusMessage(lightMessage)),
			events: []gobot.Event{{Name: LightStrengthEvent}},
		},
		{
			name:   "logMessage",
			msg:    bytes.NewReader(statusMessage(logMessage)),
			events: []gobot.Event{{Name: LogEvent}},
		},
		{
			name:   "timeCommand",
			msg:    bytes.NewReader(statusMessage(timeCommand)),
			events: []gobot.Event{{Name: TimeEvent}},
		},
		{
			name:   "bounceCommand",
			msg:    bytes.NewReader(statusMessage(bounceCommand)),
			events: []gobot.Event{{Name: BounceEvent}},
		},
		{
			name:   "takeoffCommand",
			msg:    bytes.NewReader(statusMessage(takeoffCommand)),
			events: []gobot.Event{{Name: TakeoffEvent}},
		},
		{
			name:   "landCommand",
			msg:    bytes.NewReader(statusMessage(landCommand)),
			events: []gobot.Event{{Name: LandingEvent}},
		},
		{
			name:   "palmLandCommand",
			msg:    bytes.NewReader(statusMessage(palmLandCommand)),
			events: []gobot.Event{{Name: PalmLandingEvent}},
		},
		{
			name:   "flipCommand",
			msg:    bytes.NewReader(statusMessage(flipCommand)),
			events: []gobot.Event{{Name: FlipEvent}},
		},
		{
			name:   "flightMessage",
			msg:    bytes.NewReader(statusMessage(flightMessage)),
			events: []gobot.Event{{Name: FlightDataEvent}},
		},
		{
			name:   "exposureCommand",
			msg:    bytes.NewReader(statusMessage(exposureCommand)),
			events: []gobot.Event{{Name: SetExposureEvent}},
		},
		{
			name:   "videoEncoderRateCommand",
			msg:    bytes.NewReader(statusMessage(videoEncoderRateCommand)),
			events: []gobot.Event{{Name: SetVideoEncoderRateEvent}},
		},
		{
			name:   "ConnectedEvent",
			msg:    bytes.NewReader([]byte{0x63, 0x6f, 0x6e}),
			events: []gobot.Event{{Name: ConnectedEvent}},
		},
	}

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			d := NewDriver("8888")
			events := d.Subscribe()
			err := d.handleResponse(c.msg)
			if c.err != err {
				t.Errorf("expected '%v' error, got: %v", c.err, err)
			}
			for i, cev := range c.events {
				t.Run(fmt.Sprintf("event %d", i), func(t *testing.T) {
					t.Logf("expect: %#v", cev)
					select {
					case ev, ok := <-events:
						if !ok {
							t.Error("subscription channel is closed")
						}
						if ev.Name != cev.Name {
							t.Errorf("got: %s", ev.Name)
						}
					case <-time.After(time.Millisecond):
						t.Error("subscription channel seems empty")
					}
				})
			}
		})
	}
}
