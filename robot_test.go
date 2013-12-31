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
			Connections: []Connection{newTestAdaptor("Connection 1"), newTestAdaptor("Connection 2"), newTestAdaptor("Connection 3")},
			Devices:     []Device{newTestDriver("Device 1"), newTestDriver("Device 2"), newTestDriver("Device 3")},
			Commands: map[string]interface{}{
				"Command1": func() {},
				"Command2": func() {},
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
		It("initCommands should set RobotCommands equal to Commands Key", func() {
			someRobot.initCommands()
			Expect(someRobot.RobotCommands).To(Equal([]string{"Command1", "Command2"}))
		})
		It("GetDevices should return robot devices", func() {
			someRobot.initDevices()
			Expect(someRobot.GetDevices).NotTo(BeNil())
		})
		It("GetDevice should return a robot device", func() {
			someRobot.initDevices()
			Expect(someRobot.GetDevice("Device 1").Name).To(Equal("Device 1"))
		})
		It("initConnections should initialize connections", func() {
			someRobot.initConnections()
			Expect(len(someRobot.connections)).To(Equal(3))
		})
		It("initDevices should initialize devices", func() {
			someRobot.initDevices()
			Expect(len(someRobot.devices)).To(Equal(3))
		})
		It("startConnections should connect all connections", func() {
			someRobot.initConnections()
			Expect(someRobot.startConnections()).To(Equal(true))
		})
		It("startDevices should start all devices", func() {
			someRobot.initDevices()
			Expect(someRobot.startDevices()).To(Equal(true))
		})
	})
})
