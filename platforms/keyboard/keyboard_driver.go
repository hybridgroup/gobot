package keyboard

import (
	"log"
	"os"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*KeyboardDriver)(nil)

type KeyboardDriver struct {
	name    string
	connect func(*KeyboardDriver) (err error)
	listen  func(*KeyboardDriver)
	stdin   *os.File
	gobot.Eventer
}

func NewKeyboardDriver(name string) *KeyboardDriver {
	k := &KeyboardDriver{
		name: name,
		connect: func(k *KeyboardDriver) (err error) {
			if err := configure(); err != nil {
				return err
			}

			k.stdin = os.Stdin
			return
		},
		listen: func(k *KeyboardDriver) {
			ctrlc := bytes{3}

			for {
				var keybuf bytes
				k.stdin.Read(keybuf[0:3])

				if keybuf == ctrlc {
					proc, err := os.FindProcess(os.Getpid())
					if err != nil {
						log.Fatal(err)
					}

					proc.Signal(os.Interrupt)
					break
				}

				gobot.Publish(k.Event("key"), Parse(keybuf))

			}
		},
		Eventer: gobot.NewEventer(),
	}

	k.AddEvent("key")

	return k
}

func (k *KeyboardDriver) Name() string                 { return k.name }
func (k *KeyboardDriver) Connection() gobot.Connection { return nil }

// Start initializes keyboard by grabbing key events as they come in and
// publishing a key event
func (k *KeyboardDriver) Start() (errs []error) {
	if err := k.connect(k); err != nil {
		return []error{err}
	}

	go k.listen(k)

	return
}

// Halt stops camera driver
func (k *KeyboardDriver) Halt() (errs []error) {
	if originalState != "" {
		return restore()
	}
	return
}
