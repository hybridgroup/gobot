package gobot

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Master", func() {
	var (
		m *Master
	)

	BeforeEach(func() {
		m = NewMaster()
		m.trap = func(c chan os.Signal) {
			c <- os.Interrupt
		}
		m.Robots = []*Robot{
			newTestRobot("Robot 1"),
			newTestRobot("Robot 2"),
			newTestRobot("Robot 3"),
		}
		m.Api = NewApi()
		m.Api.startFunc = func(m *api) {}
		m.Start()
	})

	Context("when valid", func() {
		It("should Find the specific robot", func() {
			Expect(m.FindRobot("Robot 1").Name).To(Equal("Robot 1"))
		})
		It("should return nil if Robot doesn't exist", func() {
			Expect(m.FindRobot("Robot 4")).To(BeNil())
		})
	})
})
