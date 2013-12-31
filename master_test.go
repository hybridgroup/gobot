package gobot

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Master", func() {

	var (
		myMaster Master
	)

	BeforeEach(func() {
		myMaster = Master{
			Robots: []Robot{
				Robot{
					Name:        "Robot 1",
					Connections: []Connection{newTestAdaptor("Connection 1")},
					Devices:     []Device{newTestDriver("Device 1")},
				},
				Robot{
					Name:        "Robot 2",
					Connections: []Connection{newTestAdaptor("Connection 2")},
					Devices:     []Device{newTestDriver("Device 2")},
				},
				Robot{
					Name:        "Robot 3",
					Connections: []Connection{newTestAdaptor("Connection 3")},
					Devices:     []Device{newTestDriver("Device 3")},
				},
			},
		}
		myMaster.Robots[0].initDevices()
		myMaster.Robots[0].initConnections()
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
