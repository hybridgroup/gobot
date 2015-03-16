package gobot

type NullReadWriteCloser struct{}

func (NullReadWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

func (NullReadWriteCloser) Read(b []byte) (int, error) {
	return len(b), nil
}

func (NullReadWriteCloser) Close() error {
	return nil
}

type testDriver struct {
	name       string
	pin        string
	connection Connection
	Eventer
	Commander
}

var testDriverStart = func() (errs []error) { return }
var testDriverHalt = func() (errs []error) { return }

func (t *testDriver) Start() (errs []error)  { return testDriverStart() }
func (t *testDriver) Halt() (errs []error)   { return testDriverHalt() }
func (t *testDriver) Name() string           { return t.name }
func (t *testDriver) Pin() string            { return t.pin }
func (t *testDriver) Connection() Connection { return t.connection }

func newTestDriver(adaptor *testAdaptor, name string, pin string) *testDriver {
	t := &testDriver{
		name:       name,
		connection: adaptor,
		pin:        pin,
		Eventer:    NewEventer(),
		Commander:  NewCommander(),
	}

	t.AddEvent("DriverCommand")

	t.AddCommand("DriverCommand", func(params map[string]interface{}) interface{} {
		return t.DriverCommand()
	})

	return t
}

func (t *testDriver) DriverCommand() string {
	Publish(t.Event("DriverCommand"), "DriverCommand")
	return "DriverCommand"
}

type testAdaptor struct {
	name string
	port string
}

var testAdaptorConnect = func() (errs []error) { return }
var testAdaptorFinalize = func() (errs []error) { return }

func (t *testAdaptor) Finalize() (errs []error) { return testAdaptorFinalize() }
func (t *testAdaptor) Connect() (errs []error)  { return testAdaptorConnect() }
func (t *testAdaptor) Name() string             { return t.name }
func (t *testAdaptor) Port() string             { return t.port }

func newTestAdaptor(name string, port string) *testAdaptor {
	return &testAdaptor{
		name: name,
		port: port,
	}
}

func newTestRobot(name string) *Robot {
	adaptor1 := newTestAdaptor("Connection1", "/dev/null")
	adaptor2 := newTestAdaptor("Connection2", "/dev/null")
	adaptor3 := newTestAdaptor("", "/dev/null")
	driver1 := newTestDriver(adaptor1, "Device1", "0")
	driver2 := newTestDriver(adaptor2, "Device2", "2")
	driver3 := newTestDriver(adaptor3, "", "1")
	work := func() {}
	r := NewRobot(name,
		[]Connection{adaptor1, adaptor2, adaptor3},
		[]Device{driver1, driver2, driver3},
		work,
	)
	r.AddCommand("RobotCommand", func(params map[string]interface{}) interface{} { return nil })

	return r
}
