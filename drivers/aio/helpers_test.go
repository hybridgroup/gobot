package aio

import (
	"fmt"
	"sync"
)

const analogReadReturnValue = 99

type aioTestBareAdaptor struct{}

func (t *aioTestBareAdaptor) Connect() error   { return nil }
func (t *aioTestBareAdaptor) Finalize() error  { return nil }
func (t *aioTestBareAdaptor) Name() string     { return "bare" }
func (t *aioTestBareAdaptor) SetName(n string) {}

type aioTestWritten struct {
	pin string
	val int
}

type aioTestAdaptor struct {
	name               string
	written            []aioTestWritten
	simulateWriteError bool
	simulateReadError  bool
	port               string
	mtx                sync.Mutex
	analogReadFunc     func() (val int, err error)
	analogWriteFunc    func(val int) error
}

func newAioTestAdaptor() *aioTestAdaptor {
	t := aioTestAdaptor{
		name: "aio_test_adaptor",
		port: "/dev/null",
		analogReadFunc: func() (int, error) {
			return analogReadReturnValue, nil
		},
		analogWriteFunc: func(val int) error {
			return nil
		},
	}

	return &t
}

// AnalogRead capabilities (interface AnalogReader)
func (t *aioTestAdaptor) AnalogRead(pin string) (int, error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	if t.simulateReadError {
		return 0, fmt.Errorf("read error")
	}

	return t.analogReadFunc()
}

// AnalogWrite capabilities (interface AnalogWriter)
func (t *aioTestAdaptor) AnalogWrite(pin string, val int) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	if t.simulateWriteError {
		return fmt.Errorf("write error")
	}

	w := aioTestWritten{pin: pin, val: val}
	t.written = append(t.written, w)
	return t.analogWriteFunc(val)
}

func (t *aioTestAdaptor) Connect() error   { return nil }
func (t *aioTestAdaptor) Finalize() error  { return nil }
func (t *aioTestAdaptor) Name() string     { return t.name }
func (t *aioTestAdaptor) SetName(n string) { t.name = n }
func (t *aioTestAdaptor) Port() string     { return t.port }
