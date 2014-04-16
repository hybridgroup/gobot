package gobot

import "fmt"

type null struct{}

func (null) Write(p []byte) (int, error) {
	return len(p), nil
}

type testDriver struct {
	Driver
	Adaptor *testAdaptor
}

func (me *testDriver) Init() bool  { return true }
func (me *testDriver) Start() bool { return true }
func (me *testDriver) Halt() bool  { return true }
func (me *testDriver) DriverCommand1(params map[string]interface{}) string {
	name := params["name"].(string)
	return fmt.Sprintf("hello %v", name)
}

type testAdaptor struct {
	Adaptor
}

func (me *testAdaptor) Finalize() bool   { return true }
func (me *testAdaptor) Connect() bool    { return true }
func (me *testAdaptor) Disconnect() bool { return true }
func (me *testAdaptor) Reconnect() bool  { return true }

func newTestDriver(name string, adaptor *testAdaptor) *testDriver {
	d := new(testDriver)
	d.Name = name
	d.Adaptor = adaptor
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
	adaptor1 := newTestAdaptor("Connection 1")
	adaptor2 := newTestAdaptor("Connection 2")
	adaptor3 := newTestAdaptor("Connection 3")
	driver1 := newTestDriver("Device 1", adaptor1)
	driver2 := newTestDriver("Device 2", adaptor2)
	driver3 := newTestDriver("Device 3", adaptor3)
	return &Robot{
		Name:        name,
		Connections: []Connection{adaptor1, adaptor2, adaptor3},
		Devices:     []Device{driver1, driver2, driver3},
		Work:        func() {},
		Commands: map[string]interface{}{
			"Command1": func() { fmt.Println("hi") },
			"Command2": func() {},
		},
	}
}
