package gobot

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("Robot", func() {

  var (
    someRobot Robot
  )

  BeforeEach(func() {
    someRobot = Robot{
      Work: func() {
      },
    }
  })

  Context("when valid", func() {
    PIt("should Start", func() {
      Expect(true)
    })
    PIt("should initConnections", func() {
      Expect(true)
    })
    PIt("should initDevices", func() {
      Expect(true)
    })
    PIt("should startConnections", func() {
      Expect(true)
    })
    PIt("should startDevices", func() {
      Expect(true)
    })
    PIt("should GetDevices", func() {
      Expect(true)
    })
    PIt("should GetDevice", func() {
      Expect(true)
    })
  })
})
