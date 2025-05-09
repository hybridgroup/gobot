package tello

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*Driver)(nil)

type WriteCloserDoNothing struct{}

func (w *WriteCloserDoNothing) Write(p []byte) (int, error) {
	return 0, nil
}

func (w *WriteCloserDoNothing) Close() error {
	return nil
}

func TestNewDriver(t *testing.T) {
	d := NewDriver("8888")

	assert.Equal(t, "8888", d.respPort)
}

func Test_handleResponse(t *testing.T) {
	tests := map[string]struct {
		msg       []byte
		wantEvent string
		wantData  (interface{})
		err       error
	}{
		"[empty message]": {
			msg: nil,
			err: io.EOF,
		},
		"wifiMessage": {
			msg:       statusMessage(wifiMessage, 0x07, 0x08, 0xA3, 0x0A),
			wantEvent: WifiDataEvent,
			wantData:  &WifiData{Strength: -93, Disturb: 10},
		},
		"lightMessage": {
			msg:       statusMessage(lightMessage, 0x17, 0x18, 0xFF),
			wantEvent: LightStrengthEvent,
			wantData:  int8(-1),
		},
		"logMessage": {
			msg:       statusMessage(logMessage),
			wantEvent: LogEvent,
			wantData:  make([]byte, 2048-9),
		},
		"timeCommand": {
			msg:       statusMessage(timeCommand, 0x27),
			wantEvent: TimeEvent,
			wantData:  []uint8{0x27},
		},
		"bounceCommand": {
			msg:       statusMessage(bounceCommand, 0x37),
			wantEvent: BounceEvent,
			wantData:  []uint8{0x37},
		},
		"takeoffCommand": {
			msg:       statusMessage(takeoffCommand, 0x47),
			wantEvent: TakeoffEvent,
			wantData:  []uint8{0x47},
		},
		"landCommand": {
			msg:       statusMessage(landCommand, 0x57),
			wantEvent: LandingEvent,
			wantData:  []uint8{0x57},
		},
		"palmLandCommand": {
			msg:       statusMessage(palmLandCommand, 0x67),
			wantEvent: PalmLandingEvent,
			wantData:  []uint8{0x67},
		},
		"flipCommand": {
			msg:       statusMessage(flipCommand, 0x77),
			wantEvent: FlipEvent,
			wantData:  []uint8{0x77},
		},
		"flightMessage": {
			msg:       statusMessage(flightMessage, 0x87, 0x88, 0x60, 0xA4),
			wantEvent: FlightDataEvent,
			wantData:  &FlightData{Height: -23456},
		},
		"exposureCommand": {
			msg:       statusMessage(exposureCommand, 0x97),
			wantEvent: SetExposureEvent,
			wantData:  []uint8{0x97},
		},
		"videoEncoderRateCommand": {
			msg:       statusMessage(videoEncoderRateCommand, 0xa7),
			wantEvent: SetVideoEncoderRateEvent,
			wantData:  []uint8{0xA7},
		},
		"ConnectedEvent": {
			msg:       []byte{0x63, 0x6f, 0x6e},
			wantEvent: ConnectedEvent,
			wantData:  nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d := NewDriver("8888")
			events := d.Subscribe()
			err := d.handleResponse(bytes.NewReader(tc.msg))
			require.ErrorIs(t, tc.err, err)
			if tc.wantEvent != "" {
				select {
				case ev, ok := <-events:
					if !ok {
						t.Error("subscription channel is closed")
					}
					if ev.Name != tc.wantEvent {
						require.Fail(t, "\ngot: %s\nwant: %s\n", ev.Name, tc.wantEvent)
					}
					got := fmt.Sprintf("%T %+[1]v", ev.Data)
					want := fmt.Sprintf("%T %+[1]v", tc.wantData)
					if got != want {
						require.Fail(t, "\ngot: %s\nwant: %s\n", got, want)
					}
				case <-time.After(5 * time.Millisecond):
					t.Error("subscription channel seems empty")
				}
			}
		})
	}
}

func TestHaltShouldTerminateAllTheRelatedGoroutines(t *testing.T) {
	d := NewDriver("8888")
	d.cmdConn = &WriteCloserDoNothing{}

	var wg sync.WaitGroup
	wg.Add(3)

	d.addDoneChReaderCount(1)
	go func() {
		<-d.doneCh
		d.addDoneChReaderCount(-1)
		wg.Done()
		fmt.Println("Done routine 1.")
	}()

	d.addDoneChReaderCount(1)
	go func() {
		<-d.doneCh
		d.addDoneChReaderCount(-1)
		wg.Done()
		fmt.Println("Done routine 2.")
	}()

	d.addDoneChReaderCount(1)
	go func() {
		<-d.doneCh
		d.addDoneChReaderCount(-1)
		wg.Done()
		fmt.Println("Done routine 3.")
	}()

	_ = d.Halt()
	wg.Wait()

	assert.Equal(t, int32(0), d.doneChReaderCount)
}

func TestHaltNotWaitForeverWhenCalledMultipleTimes(t *testing.T) {
	d := NewDriver("8888")
	d.cmdConn = &WriteCloserDoNothing{}

	_ = d.Halt()
	_ = d.Halt()
	_ = d.Halt()
}

func statusMessage(msgType uint16, msgAfter7 ...byte) []byte {
	msg := make([]byte, 7, len(msgAfter7)+7)
	msg[0] = messageStart
	binary.LittleEndian.PutUint16(msg[5:7], msgType)
	msg = append(msg, msgAfter7...)
	return msg
}
