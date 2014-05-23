package gobot

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Gobot", func() {
	var (
		g *Gobot
	)

	BeforeEach(func() {
		g = NewGobot()
		g.trap = func(c chan os.Signal) {
			c <- os.Interrupt
		}
		g.Robots = []*Robot{
			newTestRobot("Robot 1"),
			newTestRobot("Robot 2"),
			newTestRobot("Robot 3"),
		}
		g.Start()
	})

	Context("when valid", func() {
		It("should Find the specific robot", func() {
			Expect(g.Robot("Robot 1").Name).To(Equal("Robot 1"))
		})
		It("should return nil if Robot doesn't exist", func() {
			Expect(g.Robot("Robot 4")).To(BeNil())
		})
		It("Device should return nil if device doesn't exist", func() {
			Expect(g.Robot("Robot 1").Device("Device 4")).To(BeNil())
		})
		It("Device should return device", func() {
			Expect(g.Robot("Robot 1").Device("Device 1").Name).To(Equal("Device 1"))
		})
		It("Devices should return devices", func() {
			Expect(len(g.Robot("Robot 1").Devices())).To(Equal(3))
		})
		It("Connection should return nil if connection doesn't exist", func() {
			Expect(g.Robot("Robot 1").Connection("Connection 4")).To(BeNil())
		})
		It("Connection should return connection", func() {
			Expect(g.Robot("Robot 1").Connection("Connection 1").Name).To(Equal("Connection 1"))
		})
		It("Connections should return connections", func() {
			Expect(len(g.Robot("Robot 1").Connections())).To(Equal(3))
		})

	})
})
