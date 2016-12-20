package aio

type aioTestBareAdaptor struct{}

func (t *aioTestBareAdaptor) Connect() (err error)  { return }
func (t *aioTestBareAdaptor) Finalize() (err error) { return }
func (t *aioTestBareAdaptor) Name() string          { return "" }
func (t *aioTestBareAdaptor) SetName(n string)      {}

type aioTestAdaptor struct {
	name string
	port string
}

var testAdaptorAnalogRead = func() (val int, err error) {
	return 99, nil
}

func (t *aioTestAdaptor) AnalogRead(string) (val int, err error) {
	return testAdaptorAnalogRead()
}
func (t *aioTestAdaptor) Connect() (err error)  { return }
func (t *aioTestAdaptor) Finalize() (err error) { return }
func (t *aioTestAdaptor) Name() string          { return t.name }
func (t *aioTestAdaptor) SetName(n string)      { t.name = n }
func (t *aioTestAdaptor) Port() string          { return t.port }

func newAioTestAdaptor() *aioTestAdaptor {
	return &aioTestAdaptor{
		port: "/dev/null",
	}
}
