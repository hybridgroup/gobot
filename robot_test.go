package gobot

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Robot", func() {

	var (
		someRobot *Robot
	)

	Context("when valid", func() {
		BeforeEach(func() {
			someRobot = newTestRobot("")
			someRobot.Start()
		})

		It("should set random name when not set", func() {
			Expect(someRobot.Name).NotTo(BeNil())
		})
		It("GetDevice should return nil if device doesn't exist", func() {
			Expect(someRobot.GetDevice("Device 4")).To(BeNil())
		})
		It("GetDevice should return device", func() {
			Expect(someRobot.GetDevice("Device 1").Name).To(Equal("Device 1"))
		})
		It("GetDevices should return devices", func() {
			Expect(len(someRobot.GetDevices())).To(Equal(3))
		})
		It("GetConnection should return nil if connection doesn't exist", func() {
			Expect(someRobot.GetConnection("Connection 4")).To(BeNil())
		})
		It("GetConnection should return connection", func() {
			Expect(someRobot.GetConnection("Connection 1").Name).To(Equal("Connection 1"))
		})
		It("GetConnections should return connections", func() {
			Expect(len(someRobot.GetConnections())).To(Equal(3))
		})
	})
})
