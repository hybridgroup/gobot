package aio

import "sync"

type aioTestBareAdaptor struct{}

func (t *aioTestBareAdaptor) Connect() (err error)  { return }
func (t *aioTestBareAdaptor) Finalize() (err error) { return }
func (t *aioTestBareAdaptor) Name() string          { return "" }
func (t *aioTestBareAdaptor) SetName(n string)      {}

type aioTestAdaptor struct {
	name                   string
	port                   string
	mtx                    sync.Mutex
	testAdaptorAnalogRead  func() (val int, err error)
	testAdaptorAnalogWrite func(val int) (err error)
	written                []int
}

func (t *aioTestAdaptor) TestAdaptorAnalogRead(f func() (val int, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorAnalogRead = f
}

func (t *aioTestAdaptor) TestAdaptorAnalogWrite(f func(val int) (err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorAnalogWrite = f
}

func (t *aioTestAdaptor) AnalogRead(pin string) (val int, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorAnalogRead()
}

func (t *aioTestAdaptor) AnalogWrite(pin string, val int) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.written = append(t.written, val)
	return t.testAdaptorAnalogWrite(val)
}

func (t *aioTestAdaptor) Connect() (err error)  { return }
func (t *aioTestAdaptor) Finalize() (err error) { return }
func (t *aioTestAdaptor) Name() string          { return t.name }
func (t *aioTestAdaptor) SetName(n string)      { t.name = n }
func (t *aioTestAdaptor) Port() string          { return t.port }

func newAioTestAdaptor() *aioTestAdaptor {
	return &aioTestAdaptor{
		port: "/dev/null",
		testAdaptorAnalogRead: func() (val int, err error) {
			return 99, nil
		},
		testAdaptorAnalogWrite: func(val int) (err error) {
			return nil
		},
	}
}
