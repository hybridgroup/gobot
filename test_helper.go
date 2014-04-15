package gobot

type null struct{}

func (null) Write(p []byte) (int, error) {
	return len(p), nil
}

type testDriver struct {
	Driver
}

func (me *testDriver) Init() bool  { return true }
func (me *testDriver) Start() bool { return true }
func (me *testDriver) Halt() bool  { return true }

type testAdaptor struct {
	Adaptor
}

func (me *testAdaptor) Finalize() bool   { return true }
func (me *testAdaptor) Connect() bool    { return true }
func (me *testAdaptor) Disconnect() bool { return true }
func (me *testAdaptor) Reconnect() bool  { return true }

func newTestDriver(name string) *testDriver {
	d := new(testDriver)
	d.Name = name
	d.Commands = []string{
		"DriverCommand1",
		"DriverCommand2",
		"DriverCommand3",
	}
	return d
}

func newTestAdaptor(name string) *testAdaptor {
	a := new(testAdaptor)
	a.Name = name
	a.Params = map[string]interface{}{
		"param1": "1",
		"param2": 2,
	}
	return a
}

func newTestRobot(name string) *Robot {
	return &Robot{
		Name:        name,
		Connections: []Connection{newTestAdaptor("Connection 1"), newTestAdaptor("Connection 2"), newTestAdaptor("Connection 3")},
		Devices:     []Device{newTestDriver("Device 1"), newTestDriver("Device 2"), newTestDriver("Device 3")},
		Work:        func() {},
		Commands: map[string]interface{}{
			"Command1": func() {},
			"Command2": func() {},
		},
	}
}
