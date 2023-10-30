package aio

import "sync"

type aioTestAdaptor struct {
	name            string
	port            string
	mtx             sync.Mutex
	analogReadFunc  func() (val int, err error)
	analogWriteFunc func(val int) (err error)
	written         []int
}

func newAioTestAdaptor() *aioTestAdaptor {
	t := aioTestAdaptor{
		name: "aio_test_adaptor",
		port: "/dev/null",
		analogReadFunc: func() (val int, err error) {
			return 99, nil
		},
		analogWriteFunc: func(val int) (err error) {
			return nil
		},
	}

	return &t
}

// AnalogRead capabilities (interface AnalogReader)
func (t *aioTestAdaptor) AnalogRead(pin string) (val int, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.analogReadFunc()
}

// AnalogWrite capabilities (interface AnalogWriter)
func (t *aioTestAdaptor) AnalogWrite(pin string, val int) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.written = append(t.written, val)
	return t.analogWriteFunc(val)
}

func (t *aioTestAdaptor) Connect() (err error)  { return }
func (t *aioTestAdaptor) Finalize() (err error) { return }
func (t *aioTestAdaptor) Name() string          { return t.name }
func (t *aioTestAdaptor) SetName(n string)      { t.name = n }
func (t *aioTestAdaptor) Port() string          { return t.port }
