package aio

import "sync"

type aioTestAdaptor struct {
	name            string
	port            string
	mtx             sync.Mutex
	analogReadFunc  func() (val int, err error)
	analogWriteFunc func(val int) error
	written         []int
}

func newAioTestAdaptor() *aioTestAdaptor {
	t := aioTestAdaptor{
		name: "aio_test_adaptor",
		port: "/dev/null",
		analogReadFunc: func() (int, error) {
			return 99, nil
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
	return t.analogReadFunc()
}

// AnalogWrite capabilities (interface AnalogWriter)
func (t *aioTestAdaptor) AnalogWrite(pin string, val int) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.written = append(t.written, val)
	return t.analogWriteFunc(val)
}

func (t *aioTestAdaptor) Connect() error   { return nil }
func (t *aioTestAdaptor) Finalize() error  { return nil }
func (t *aioTestAdaptor) Name() string     { return t.name }
func (t *aioTestAdaptor) SetName(n string) { t.name = n }
func (t *aioTestAdaptor) Port() string     { return t.port }
