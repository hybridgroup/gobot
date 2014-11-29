package gpio

import (
	"time"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*MakeyButtonDriver)(nil)

// Represents a Makey Button
type MakeyButtonDriver struct {
	name       string
	pin        string
	connection gobot.Connection
	Active     bool
	data       []int
	interval   time.Duration
	gobot.Eventer
}

// NewMakeyButtonDriver returns a new MakeyButtonDriver given a DigitalRead, name and pin.
func NewMakeyButtonDriver(a DigitalReader, name string, pin string, v ...time.Duration) *MakeyButtonDriver {
	m := &MakeyButtonDriver{
		name:       name,
		connection: a.(gobot.Connection),
		pin:        pin,
		Active:     false,
		Eventer:    gobot.NewEventer(),
		interval:   10 * time.Millisecond,
	}

	if len(v) > 0 {
		m.interval = v[0]
	}

	m.AddEvent("error")
	m.AddEvent("push")
	m.AddEvent("release")

	return m
}

func (b *MakeyButtonDriver) Name() string                 { return b.name }
func (b *MakeyButtonDriver) Pin() string                  { return b.pin }
func (b *MakeyButtonDriver) Connection() gobot.Connection { return b.connection }

func (b *MakeyButtonDriver) adaptor() DigitalReader {
	return b.Connection().(DigitalReader)
}

// Starts the MakeyButtonDriver and reads the state of the button at the given Driver.Interval().
// Returns true on successful start of the driver.
//
// Emits the Events:
// 	"push"    int - On button push
//	"release" int - On button release
func (m *MakeyButtonDriver) Start() (errs []error) {
	state := 0
	go func() {
		for {
			newValue, err := m.readState()
			if err != nil {
				gobot.Publish(m.Event("error"), err)
			} else if newValue != state && newValue != -1 {
				state = newValue
				if newValue == 0 {
					m.Active = true
					gobot.Publish(m.Event("push"), newValue)
				} else {
					m.Active = false
					gobot.Publish(m.Event("release"), newValue)
				}
			}
		}
		<-time.After(m.interval)
	}()
	return
}

// Halt returns true on a successful halt of the driver
func (m *MakeyButtonDriver) Halt() (errs []error) { return }

func (m *MakeyButtonDriver) readState() (val int, err error) {
	return m.adaptor().DigitalRead(m.Pin())
}
