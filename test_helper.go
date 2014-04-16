package gobot

import "fmt"

type testStruct struct {
	i int
	f float64
}

func (me *testStruct) Hello(name string, message string) string {
	return fmt.Sprintf("Hello %v! %v", name, message)
}

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
func (me *testDriver) TestDriverCommand(params map[string]interface{}) string {
	name := params["name"].(string)
	return fmt.Sprintf("hello %v", name)
}

type testAdaptor struct {
	Adaptor
}

func (me *testAdaptor) Finalize() bool { return true }
func (me *testAdaptor) Connect() bool  { return true }

func newTestDriver(name string, adaptor *testAdaptor) *testDriver {
	d := new(testDriver)
	d.Name = name
	d.Adaptor = adaptor
	d.Commands = []string{
		"TestDriverCommand",
		"DriverCommand",
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

func robotTestFunction(params map[string]interface{}) string {
	message := params["message"].(string)
	robotname := params["robotname"].(string)
	return fmt.Sprintf("hey %v, %v", robotname, message)
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
			"robotTestFunction": robotTestFunction,
		},
	}
}

func newTestStruct() *testStruct {
	s := new(testStruct)
	s.i = 10
	s.f = 0.2
	return s
}
