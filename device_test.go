package gobot

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Device", func() {

	var (
		someRobot Robot
	)

	BeforeEach(func() {
		someRobot = Robot{
			Connections: []Connection{newTestAdaptor("Connection 1")},
			Devices:     []Device{newTestDriver("Device 1")},
		}
	})

	Context("when valid", func() {
		It("Commands should return device commands", func() {
			someRobot.initDevices()
			Expect(someRobot.devices[0].Commands()).To(Equal([]string{"DriverCommand1", "DriverCommand2", "DriverCommand3"}))
		})
	})
})
