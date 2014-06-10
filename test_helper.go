package gobot

import (
	"fmt"
)

type testStruct struct {
	i int
	f float64
}

func (t *testStruct) Hello(name string, message string) string {
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

func (t *testDriver) Init() bool  { return true }
func (t *testDriver) Start() bool { return true }
func (t *testDriver) Halt() bool  { return true }
func (t *testDriver) TestDriverCommand(params map[string]interface{}) string {
	name := params["name"].(string)
	return fmt.Sprintf("hello %v", name)
}

type testAdaptor struct {
	Adaptor
}

func (t *testAdaptor) Finalize() bool { return true }
func (t *testAdaptor) Connect() bool  { return true }

func newTestDriver(name string, adaptor *testAdaptor) *testDriver {
	return &testDriver{
		Driver: Driver{
			Commands: []string{
				"TestDriverCommand",
				"DriverCommand",
			},
			Name: name,
		},
		Adaptor: adaptor,
	}
}

func newTestAdaptor(name string) *testAdaptor {
	return &testAdaptor{
		Adaptor: Adaptor{
			Name: name,
			Params: map[string]interface{}{
				"param1": "1",
				"param2": 2,
			},
		},
	}
}

func robotTestFunction(params map[string]interface{}) string {
	message := params["message"].(string)
	robotname := params["robotname"].(string)
	return fmt.Sprintf("hey %v, %v", robotname, message)
}

func NewTestRobot(name string) *Robot {
	return newTestRobot(name)
}
func newTestRobot(name string) *Robot {
	adaptor1 := newTestAdaptor("Connection 1")
	adaptor2 := newTestAdaptor("Connection 2")
	adaptor3 := newTestAdaptor("Connection 3")
	driver1 := newTestDriver("Device 1", adaptor1)
	driver2 := newTestDriver("Device 2", adaptor2)
	driver3 := newTestDriver("Device 3", adaptor3)
	work := func() {}
	//Commands := map[string]interface{}{
	//	"robotTestFunction": robotTestFunction,
	//}
	return NewRobot(name, []Connection{adaptor1, adaptor2, adaptor3}, []Device{driver1, driver2, driver3}, work)
}

func newTestStruct() *testStruct {
	s := new(testStruct)
	s.i = 10
	s.f = 0.2
	return s
}
