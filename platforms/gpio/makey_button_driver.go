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
	connection DigitalReader
	Active     bool
	data       []int
	interval   time.Duration
	gobot.Eventer
}

// NewMakeyButtonDriver returns a new MakeyButtonDriver given a DigitalRead, name and pin.
func NewMakeyButtonDriver(a DigitalReader, name string, pin string, v ...time.Duration) *MakeyButtonDriver {
	m := &MakeyButtonDriver{
		name:       name,
		connection: a,
		pin:        pin,
		Active:     false,
		Eventer:    gobot.NewEventer(),
		interval:   10 * time.Millisecond,
	}

	if len(v) > 0 {
		m.interval = v[0]
	}

	m.AddEvent(Error)
	m.AddEvent(Push)
	m.AddEvent(Release)

	return m
}

func (b *MakeyButtonDriver) Name() string                 { return b.name }
func (b *MakeyButtonDriver) Pin() string                  { return b.pin }
func (b *MakeyButtonDriver) Connection() gobot.Connection { return b.connection.(gobot.Connection) }

// Starts the MakeyButtonDriver and reads the state of the button at the given Driver.Interval().
// Returns true on successful start of the driver.
//
// Emits the Events:
// 	"push"    int - On button push
//	"release" int - On button release
func (m *MakeyButtonDriver) Start() (errs []error) {
	state := 1
	go func() {
		for {
			newValue, err := m.connection.DigitalRead(m.Pin())
			if err != nil {
				gobot.Publish(m.Event(Error), err)
			} else if newValue != state && newValue != -1 {
				state = newValue
				if newValue == 0 {
					m.Active = true
					gobot.Publish(m.Event(Push), newValue)
				} else {
					m.Active = false
					gobot.Publish(m.Event(Release), newValue)
				}
			}
			<-time.After(m.interval)
		}
	}()
	return
}

// Halt returns true on a successful halt of the driver
func (m *MakeyButtonDriver) Halt() (errs []error) { return }
