package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
)

type generate struct {
	Name        string
	UpperName   string
	FirstLetter string
}

func Generate() cli.Command {
	return cli.Command{
		Name:  "generate",
		Usage: "Generate new Gobot skeleton project",
		Action: func(c *cli.Context) {
			args := c.Args()
			if len(args) == 0 || len(args) > 1 {
				fmt.Println("Please provide a one word package name.")
				return
			}
			pwd, _ := os.Getwd()
			dir := fmt.Sprintf("%s/%s", pwd, "gobot-"+args[0])
			fmt.Println("Creating", dir)
			err := os.MkdirAll(dir, 0700)
			if err != nil {
				fmt.Println(err)
				err = nil
			}

			examplesDir := dir + "/examples"
			fmt.Println("Creating", examplesDir)
			err = os.MkdirAll(examplesDir, 0700)
			if err != nil {
				fmt.Println(err)
				err = nil
			}

			upperName := fmt.Sprintf("%s%s",
				strings.ToUpper(string(args[0][0])),
				string(args[0][1:]))

			name := generate{
				UpperName:   upperName,
				Name:        string(args[0]),
				FirstLetter: string(args[0][0]),
			}

			adaptor, _ := template.New("").Parse(adaptor())
			fileLocation := fmt.Sprintf("%s/%s_adaptor.go", dir, args[0])
			fmt.Println("Creating", fileLocation)
			f, err := os.Create(fileLocation)
			if err != nil {
				fmt.Println(err)
				err = nil
			}
			adaptor.Execute(f, name)
			f.Close()

			driver, _ := template.New("").Parse(driver())
			fileLocation = fmt.Sprintf("%s/%s_driver.go", dir, args[0])
			fmt.Println("Creating", fileLocation)
			f, err = os.Create(fileLocation)
			if err != nil {
				fmt.Println(err)
				err = nil
			}
			driver.Execute(f, name)
			f.Close()

			fileLocation = fmt.Sprintf("%s/LICENSE", dir)
			fmt.Println("Creating", fileLocation)
			f, err = os.Create(fileLocation)
			if err != nil {
				fmt.Println(err)
				err = nil
			}
			f.Close()

			driverTest, _ := template.New("").Parse(driverTest())
			fileLocation = fmt.Sprintf("%s/%s_driver_test.go", dir, args[0])
			fmt.Println("Creating", fileLocation)
			f, err = os.Create(fileLocation)
			if err != nil {
				fmt.Println(err)
				err = nil
			}
			driverTest.Execute(f, name)
			f.Close()

			adaptorTest, _ := template.New("").Parse(adaptorTest())
			fileLocation = fmt.Sprintf("%s/%s_adaptor_test.go", dir, args[0])
			fmt.Println("Creating", fileLocation)
			f, err = os.Create(fileLocation)
			if err != nil {
				fmt.Println(err)
				err = nil
			}
			adaptorTest.Execute(f, name)
			f.Close()

			example, _ := template.New("").Parse(example())
			fileLocation = fmt.Sprintf("%s/main.go", examplesDir)
			fmt.Println("Creating", fileLocation)
			f, err = os.Create(fileLocation)
			if err != nil {
				fmt.Println(err)
				err = nil
			}
			example.Execute(f, name)
			f.Close()

			readme, _ := template.New("").Parse(readme())
			fileLocation = fmt.Sprintf("%s/README.md", dir)
			fmt.Println("Creating", fileLocation)
			f, err = os.Create(fileLocation)
			if err != nil {
				fmt.Println(err)
				err = nil
			}

			ex, _ := ioutil.ReadFile(examplesDir + "/main.go")
			data := struct {
				Name    string
				Example string
			}{
				name.Name,
				string(ex),
			}
			readme.Execute(f, data)
			f.Close()
		},
	}
}

func adaptor() string {
	return `package {{.Name}}

import (
	"github.com/hybridgroup/gobot"
)

var _ gobot.Adaptor = (*{{.UpperName}}Adaptor)(nil)

type {{.UpperName}}Adaptor struct {
	name string
}

func New{{.UpperName}}Adaptor(name string) *{{.UpperName}}Adaptor {
	return &{{.UpperName}}Adaptor{
		name: name,
	}
}

func ({{.FirstLetter}} *{{.UpperName}}Adaptor) Name() string { return {{.FirstLetter}}.name }

func ({{.FirstLetter}} *{{.UpperName}}Adaptor) Connect() []error { return nil }

func ({{.FirstLetter}} *{{.UpperName}}Adaptor) Finalize() []error { return nil }

func ({{.FirstLetter}} *{{.UpperName}}Adaptor) Ping() string { return "pong" }
`
}

func driver() string {
	return `package {{.Name }}

import (
	"time"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*{{.UpperName}}Driver)(nil)

const Hello string = "hello"

type {{.UpperName}}Driver struct {
	name string
	connection gobot.Connection
	interval time.Duration
	halt chan bool
	gobot.Eventer
	gobot.Commander
}

func New{{.UpperName}}Driver(a *{{.UpperName}}Adaptor, name string) *{{.UpperName}}Driver {
	{{.FirstLetter}} := &{{.UpperName}}Driver{
		name: name,
		connection: a,
		interval: 500*time.Millisecond,
		halt: make(chan bool, 0),
    Eventer:    gobot.NewEventer(),
    Commander:  gobot.NewCommander(),
	}

	{{.FirstLetter}}.AddEvent(Hello)

	{{.FirstLetter}}.AddCommand(Hello, func(params map[string]interface{}) interface{} {
		return {{.FirstLetter}}.Hello()
	})

	return {{.FirstLetter}}
}

func ({{.FirstLetter}} *{{.UpperName}}Driver) Name() string { return {{.FirstLetter}}.name }

func ({{.FirstLetter}} *{{.UpperName}}Driver) Connection() gobot.Connection {
	return {{.FirstLetter}}.connection
}

func ({{.FirstLetter}} *{{.UpperName}}Driver) adaptor() *{{.UpperName}}Adaptor {
	return {{.FirstLetter}}.Connection().(*{{.UpperName}}Adaptor)
}

func ({{.FirstLetter}} *{{.UpperName}}Driver) Hello() string {
	return "hello from " + {{.FirstLetter}}.Name() + "!"
}

func ({{.FirstLetter}} *{{.UpperName}}Driver) Ping() string {
	return {{.FirstLetter}}.adaptor().Ping()
}

func ({{.FirstLetter}} *{{.UpperName}}Driver) Start() []error {
	go func() {
		for {
			gobot.Publish({{.FirstLetter}}.Event(Hello), {{.FirstLetter}}.Hello())

			select {
			case <- time.After({{.FirstLetter}}.interval):
			case <- {{.FirstLetter}}.halt:
				return
			}
		}
	}()
	return nil
}

func ({{.FirstLetter}} *{{.UpperName}}Driver) Halt() []error {
	{{.FirstLetter}}.halt <- true
	return nil
}

`
}

func example() string {
	return `
package main

import (
  "../"
  "fmt"
  "time"

  "github.com/hybridgroup/gobot"
)

func main() {
  gbot := gobot.NewGobot()

  conn := {{.Name}}.New{{.UpperName}}Adaptor("conn")
  dev := {{.Name}}.New{{.UpperName}}Driver(conn, "dev")

  work := func() {
    gobot.On(dev.Event({{.Name}}.Hello), func(data interface{}) {
      fmt.Println(data)
    })

    gobot.Every(1200*time.Millisecond, func() {
      fmt.Println(dev.Ping())
    })
  }

  robot := gobot.NewRobot(
    "robot",
    []gobot.Connection{conn},
    []gobot.Device{dev},
    work,
  )

  gbot.AddRobot(robot)
  gbot.Start()
}
`
}

func driverTest() string {
	return `package {{.Name}}

import (
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
)

func Test{{.UpperName}}Driver(t *testing.T) {
	d := New{{.UpperName}}Driver(New{{.UpperName}}Adaptor("conn"), "dev")

	gobot.Assert(t, d.Name(), "dev")
	gobot.Assert(t, d.Connection().Name(), "conn")

	ret := d.Command(Hello)(nil)
	gobot.Assert(t, ret.(string), "hello from dev!")

	gobot.Assert(t, d.Ping(), "pong")

	gobot.Assert(t, len(d.Start()), 0)

	<-time.After(d.interval)

	sem := make(chan bool, 0)

	gobot.On(d.Event(Hello), func(data interface{}) {
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(600 * time.Millisecond):
		t.Errorf("Hello Event was not published")
	}

	gobot.Assert(t, len(d.Halt()), 0)

	gobot.On(d.Event(Hello), func(data interface{}) {
		sem <- true
	})

	select {
	case <-sem:
		t.Errorf("Hello Event should not publish after Halt")
	case <-time.After(600 * time.Millisecond):
	}
}

`
}

func adaptorTest() string {
	return `package {{.Name}}

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func Test{{.UpperName}}Adaptor(t *testing.T) {
	a := New{{.UpperName}}Adaptor("tester")

	gobot.Assert(t, a.Name(), "tester")

	gobot.Assert(t, len(a.Connect()), 0)

	gobot.Assert(t, a.Ping(), "pong")

	gobot.Assert(t, len(a.Connect()), 0)

	gobot.Assert(t, len(a.Finalize()), 0)
}
`
}

func readme() string {
	return `# {{.Name}}

Gobot (http://gobot.io/) is a framework for robotics and physical computing using Go

This repository contains the Gobot adaptor and driver for {{.Name}}.

For more information about Gobot, check out the github repo at
https://github.com/hybridgroup/gobot

## Installing
` + "```bash\ngo get path/to/repo/{{.Name}}\n```" + `

## Using
` + "```go{{.Example}}\n```" + `

## Connecting

Explain how to connect to the device here...

## License

Copyright (c) 2014 Your Name Here. See LICENSE for more details
`
}
