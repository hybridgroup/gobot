package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"strings"
	"text/template"
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

			upperName := fmt.Sprintf("%s%s", strings.ToUpper(string(args[0][0])), string(args[0][1:]))

			name := generate{UpperName: upperName, Name: string(args[0]), FirstLetter: string(args[0][0])}

			adaptor, _ := template.New("").Parse(adaptor())
			file_location := fmt.Sprintf("%s/%s_adaptor.go", dir, args[0])
			fmt.Println("Creating", file_location)
			f, err := os.Create(file_location)
			if err != nil {
				fmt.Println(err)
				err = nil
			}
			adaptor.Execute(f, name)
			f.Close()

			driver, _ := template.New("").Parse(driver())
			file_location = fmt.Sprintf("%s/%s_driver.go", dir, args[0])
			fmt.Println("Creating", file_location)
			f, err = os.Create(file_location)
			if err != nil {
				fmt.Println(err)
				err = nil
			}
			driver.Execute(f, name)
			f.Close()

			readme, _ := template.New("").Parse(readme())
			file_location = fmt.Sprintf("%s/README.md", dir)
			fmt.Println("Creating", file_location)
			f, err = os.Create(file_location)
			if err != nil {
				fmt.Println(err)
				err = nil
			}
			readme.Execute(f, name)
			f.Close()

			file_location = fmt.Sprintf("%s/LICENSE", dir)
			fmt.Println("Creating", file_location)
			f, err = os.Create(file_location)
			if err != nil {
				fmt.Println(err)
				err = nil
			}
			f.Close()

			driverTest, _ := template.New("").Parse(driverTest())
			file_location = fmt.Sprintf("%s/%s_driver_test.go", dir, args[0])
			fmt.Println("Creating", file_location)
			f, err = os.Create(file_location)
			if err != nil {
				fmt.Println(err)
				err = nil
			}
			driverTest.Execute(f, name)
			f.Close()

			adaptorTest, _ := template.New("").Parse(adaptorTest())
			file_location = fmt.Sprintf("%s/%s_adaptor_test.go", dir, args[0])
			fmt.Println("Creating", file_location)
			f, err = os.Create(file_location)
			if err != nil {
				fmt.Println(err)
				err = nil
			}
			adaptorTest.Execute(f, name)
			f.Close()

			testSuite, _ := template.New("").Parse(testSuite())
			file_location = fmt.Sprintf("%s/%s_suite_test.go", dir, args[0])
			fmt.Println("Creating", file_location)
			f, err = os.Create(file_location)
			if err != nil {
				fmt.Println(err)
				err = nil
			}
			testSuite.Execute(f, name)
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
    Adaptor: gobot.Adaptor{
      Name: name,
    },
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
  Adaptor *{{ .UpperName }}Adaptor
}

type {{ .UpperName }}Interface interface {
}

func New{{.UpperName}}Driver(a *{{.UpperName}}Adaptor, name string) *{{.UpperName}}Driver {
  return &{{.UpperName}}Driver{
    Driver: gobot.Driver{
      Name: name,
      Events: make(map[string]chan interface{}),
      Commands: []string{},
    },
    Adaptor: a,
  }
}

func ({{.FirstLetter}} *{{ .UpperName }}Driver) Start() bool { return true }
func ({{.FirstLetter}} *{{ .UpperName }}Driver) Halt() bool { return true }
`
}

func driverTest() string {
	return `package {{ .Name }}

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("{{ .UpperName }}Driver", func() {
  var (
    driver *{{ .UpperName }}Driver
  )

  BeforeEach(func() {
    driver = New{{ .UpperName }}Driver(New{{ .UpperName }}Adaptor("adaptor"), "driver")
  })

  It("Must be able to Start", func() {
    Expect(driver.Start()).To(Equal(true))
  })
  It("Must be able to Halt", func() {
    Expect(driver.Halt()).To(Equal(true))
  })
})
`
}

func adaptorTest() string {
	return `package {{ .Name }}

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("{{ .UpperName }}Adaptor", func() {
  var (
    adaptor *{{ .UpperName }}Adaptor
  )

  BeforeEach(func() {
    adaptor = New{{ .UpperName }}Adaptor("adaptor")
  })

  It("Must be able to Finalize", func() {
    Expect(adaptor.Finalize()).To(Equal(true))
  })
  It("Must be able to Connect", func() {
    Expect(adaptor.Connect()).To(Equal(true))
  })
})
`
}

func testSuite() string {
	return `package {{ .Name }}

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "testing"
)

func TestGobot{{ .UpperName }}(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "{{ .UpperName }} Suite")
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
