// Based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)

package audio

import (
	"github.com/hybridgroup/gobot"
	"time"
)

type AudioDriver struct {
	name       string
	connection gobot.Connection
	interval   time.Duration
	halt       chan bool
	gobot.Eventer
	gobot.Commander
	filename string
}

func NewAudioDriver(a *AudioAdaptor, name string, filename string) *AudioDriver {
	d := &AudioDriver{
		name:       name,
		connection: a,
		interval:   500 * time.Millisecond,
		filename:   filename,
		halt:       make(chan bool, 0),
		Eventer:    gobot.NewEventer(),
		Commander:  gobot.NewCommander(),
	}
	return d
}

func (d *AudioDriver) Name() string { return d.name }

func (d *AudioDriver) Filename() string { return d.filename }

func (d *AudioDriver) Connection() gobot.Connection {
	return d.connection
}

func (d *AudioDriver) Sound(fileName string) []error {
	return d.Connection().(*AudioAdaptor).Sound(fileName)
}

func (d *AudioDriver) Play() []error {
	return d.Sound(d.Filename())
}

func (d *AudioDriver) adaptor() *AudioAdaptor {
	return d.Connection().(*AudioAdaptor)
}

func (d *AudioDriver) Start() (err []error) {
	return
}

func (d *AudioDriver) Halt() (err []error) {
	return
}
