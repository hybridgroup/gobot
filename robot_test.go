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
		someRobot = newTestRobot("")
		start = func(r *Robot) {
			r.startRobot()
		}
		someRobot.Start()
	})

	Context("when valid", func() {
		It("initName should not change name when already set", func() {
			someRobot.Name = "Bumblebee"
			Expect(someRobot.Name).To(Equal("Bumblebee"))
		})
		It("initName should set random name when not set", func() {
			Expect(someRobot.Name).NotTo(BeNil())
			Expect(someRobot.Name).NotTo(Equal("Bumblebee"))
		})
		It("initCommands should set RobotCommands equal to Commands Key", func() {
			Expect(someRobot.RobotCommands).To(Equal([]string{"Command1", "Command2"}))
		})
		It("GetDevices should return robot devices", func() {
			Expect(someRobot.GetDevices).NotTo(BeNil())
		})
		It("GetDevice should return a robot device", func() {
			Expect(someRobot.GetDevice("Device 1").Name).To(Equal("Device 1"))
		})
		It("initConnections should initialize connections", func() {
			Expect(len(someRobot.connections)).To(Equal(3))
		})
		It("initDevices should initialize devices", func() {
			Expect(len(someRobot.devices)).To(Equal(3))
		})
	})
})
