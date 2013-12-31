package gobot

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Master", func() {

	var (
		myMaster *Master
	)

	BeforeEach(func() {
		myMaster = GobotMaster()
		myMaster.Robots = []Robot{
			newTestRobot("Robot 1"),
			newTestRobot("Robot 2"),
			newTestRobot("Robot 3"),
		}
		for s := range myMaster.Robots {
			myMaster.Robots[s].startRobot()
		}
	})

	Context("when valid", func() {
		It("should Find the specific robot", func() {
			Expect(myMaster.FindRobot("Robot 1").Name).To(Equal("Robot 1"))
		})
		It("should Find the specific robot device", func() {
			Expect(myMaster.FindRobotDevice("Robot 1", "Device 1").Name).To(Equal("Device 1"))
		})
		It("should Find the specific robot connection", func() {
			Expect(myMaster.FindRobotConnection("Robot 1", "Connection 1").Name).To(Equal("Connection 1"))
		})
	})
})
