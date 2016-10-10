// Based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)

package audio

import (
	"time"

	"github.com/hybridgroup/gobot"
)

type Driver struct {
	name       string
	connection gobot.Connection
	interval   time.Duration
	halt       chan bool
	gobot.Eventer
	gobot.Commander
	filename string
}

func NewDriver(a *Adaptor, filename string) *Driver {
	d := &Driver{
		name:       "Audio",
		connection: a,
		interval:   500 * time.Millisecond,
		filename:   filename,
		halt:       make(chan bool, 0),
		Eventer:    gobot.NewEventer(),
		Commander:  gobot.NewCommander(),
	}
	return d
}

func (d *Driver) Name() string { return d.name }

func (d *Driver) SetName(n string) { d.name = n }

func (d *Driver) Filename() string { return d.filename }

func (d *Driver) Connection() gobot.Connection {
	return d.connection
}

func (d *Driver) Sound(fileName string) []error {
	return d.Connection().(*Adaptor).Sound(fileName)
}

func (d *Driver) Play() []error {
	return d.Sound(d.Filename())
}

func (d *Driver) adaptor() *Adaptor {
	return d.Connection().(*Adaptor)
}

func (d *Driver) Start() (err []error) {
	return
}

func (d *Driver) Halt() (err []error) {
	return
}
