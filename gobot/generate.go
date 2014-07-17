package main

import (
	"fmt"
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
			dir := fmt.Sprintf("%s/%s", pwd, args[0])
			fmt.Println("Creating", dir)
			err := os.MkdirAll(dir, 0700)
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

			readme, _ := template.New("").Parse(readme())
			fileLocation = fmt.Sprintf("%s/README.md", dir)
			fmt.Println("Creating", fileLocation)
			f, err = os.Create(fileLocation)
			if err != nil {
				fmt.Println(err)
				err = nil
			}
			readme.Execute(f, name)
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
		},
	}
}

func adaptor() string {
	return `package {{ .Name }}

import (
  "github.com/hybridgroup/gobot"
)

type {{ .UpperName }}Adaptor struct {
  gobot.Adaptor
}

func New{{.UpperName}}Adaptor(name string) *{{.UpperName}}Adaptor {
  return &{{.UpperName}}Adaptor{
    Adaptor: *gobot.NewAdaptor(
      name,
      "{{.Name}}.{{.UpperName}}Adaptor",
    ),
  }
}

func ({{.FirstLetter}} *{{ .UpperName }}Adaptor) Connect() bool {
  return true
}

func ({{.FirstLetter}} *{{ .UpperName }}Adaptor) Finalize() bool {
  return true
}
`
}

func driver() string {
	return `package {{ .Name }}

import (
  "github.com/hybridgroup/gobot"
)

type {{ .UpperName }}Driver struct {
  gobot.Driver
}

type {{ .UpperName }}Interface interface {
}

func New{{.UpperName}}Driver(a *{{.UpperName}}Adaptor, name string) *{{.UpperName}}Driver {
  return &{{.UpperName}}Driver{
    Driver: *gobot.NewDriver(
      name,
      "{{.Name}}.{{.UpperName}}Driver",
      a,
    ),
  }
}

func ({{.FirstLetter}} *{{ .UpperName }}Driver) adaptor() *{{ .UpperName }}Adaptor {
  return {{ .FirstLetter }}.Driver.Adaptor().(*{{ .UpperName }}Adaptor)
}

func ({{.FirstLetter}} *{{ .UpperName }}Driver) Start() bool { return true }
func ({{.FirstLetter}} *{{ .UpperName }}Driver) Halt() bool { return true }
`
}

func driverTest() string {
	return `package {{ .Name }}

import (
  "github.com/hybridgroup/gobot"
  "testing"
)

func initTest{{ .UpperName }}Driver() *{{ .UpperName }}Driver {
  return New{{ .UpperName}}Driver(New{{ .UpperName }}Adaptor("myAdaptor"), "myDriver")
}

func Test{{ .UpperName }}DriverStart(t *testing.T) {
  d := initTest{{.UpperName }}Driver()
  gobot.Assert(t, d.Start(), true)
}

func Test{{ .UpperName }}DriverHalt(t *testing.T) {
  d := initTest{{.UpperName }}Driver()
  gobot.Assert(t, d.Halt(), true)
}
`
}

func adaptorTest() string {
	return `package {{ .Name }}

import (
  "github.com/hybridgroup/gobot"
  "testing"
)

func initTest{{ .UpperName }}Adaptor() *{{ .UpperName }}Adaptor {
  return New{{ .UpperName }}Adaptor("myAdaptor")
}

func Test{{ .UpperName }}AdaptorConnect(t *testing.T) {
  a := initTest{{.UpperName }}Adaptor()
  gobot.Assert(t, a.Connect(), true)
}

func Test{{ .UpperName }}AdaptorFinalize(t *testing.T) {
  a := initTest{{.UpperName }}Adaptor()
  gobot.Assert(t, a.Finalize(), true)
}
`
}

func readme() string {
	return `# {{ .Name }}

Gobot (http://gobot.io/) is a framework for robotics and physical computing using Go

This repository contains the Gobot adaptor and driver for {{ .Name }}.

For more information about Gobot, check out the github repo at
https://github.com/hybridgroup/gobot

## Installing

    go get path/to/repo/{{ .Name }}

## Using

    your example code here...

## Connecting

Explain how to connect from the computer to the device here...

## License

Copyright (c) 2014 Your Name Here. See LICENSE for more details
`
}
