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
		someRobot = newTestRobot("")
		start = func(r *Robot) {
			r.startRobot()
		}
		someRobot.Start()
	})

	Context("when valid", func() {
		It("Commands should return device commands", func() {
			Expect(someRobot.devices[0].Commands()).To(Equal([]string{"DriverCommand1", "DriverCommand2", "DriverCommand3"}))
		})
		It("Start should call driver start", func() {
			Expect(someRobot.Devices[0].Start()).To(Equal(true))
		})
	})
})
