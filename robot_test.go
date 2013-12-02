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
    It("initName should not change name when already set", func() {
      someRobot.Name = "Bumblebee"
      someRobot.initName()
      Expect(someRobot.Name).To(Equal("Bumblebee"))
    })
    It("initName should set random name when not set", func() {
      someRobot.initName()
      Expect(someRobot.Name).NotTo(BeNil())
      Expect(someRobot.Name).NotTo(Equal("Bumblebee"))
    })
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
