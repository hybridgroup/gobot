package gobot

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Master", func() {
	var (
		myMaster *Master
	)

	BeforeEach(func() {
		myMaster = GobotMaster()
		myMaster.Robots = []*Robot{
			newTestRobot("Robot 1"),
			newTestRobot("Robot 2"),
			newTestRobot("Robot 3"),
		}
		trap = func(c chan os.Signal) {
			c <- os.Interrupt
		}
		myMaster.Start()
	})

	Context("when valid", func() {
		It("should Find the specific robot", func() {
			Expect(myMaster.FindRobot("Robot 1").Name).To(Equal("Robot 1"))
		})
		It("should return nil if Robot doesn't exist", func() {
			Expect(myMaster.FindRobot("Robot 4")).To(BeNil())
		})
		It("should Find the specific robot device", func() {
			Expect(myMaster.FindRobotDevice("Robot 2", "Device 2").Name).To(Equal("Device 2"))
		})
		It("should return nil if the robot device doesn't exist", func() {
			Expect(myMaster.FindRobotDevice("Robot 4", "Device 2")).To(BeNil())
		})
		It("should Find the specific robot connection", func() {
			Expect(myMaster.FindRobotConnection("Robot 3", "Connection 1").Name).To(Equal("Connection 1"))
		})
		It("should return nil if the robot connection doesn't exist", func() {
			Expect(myMaster.FindRobotConnection("Robot 4", "Connection 1")).To(BeNil())
		})
		It("Commands should return device commands", func() {
			Expect(myMaster.FindRobotDevice("Robot 2", "Device 1").Commands()).To(Equal([]string{"DriverCommand1", "DriverCommand2", "DriverCommand3"}))
		})
	})
})
