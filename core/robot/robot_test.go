package robot

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Robot", func() {

	var (
		r *Robot
	)

	Context("when valid", func() {
		BeforeEach(func() {
			r = newTestRobot("")
			r.master = NewMaster()
			r.master.trap = func(c chan os.Signal) {
				c <- os.Interrupt
			}
			r.Start()
		})

		It("should set random name when not set", func() {
			Expect(r.Name).NotTo(BeNil())
		})
		It("GetDevice should return nil if device doesn't exist", func() {
			Expect(r.GetDevice("Device 4")).To(BeNil())
		})
		It("GetDevice should return device", func() {
			Expect(r.GetDevice("Device 1").Name).To(Equal("Device 1"))
		})
		It("GetDevices should return devices", func() {
			Expect(len(r.GetDevices())).To(Equal(3))
		})
		It("GetConnection should return nil if connection doesn't exist", func() {
			Expect(r.GetConnection("Connection 4")).To(BeNil())
		})
		It("GetConnection should return connection", func() {
			Expect(r.GetConnection("Connection 1").Name).To(Equal("Connection 1"))
		})
		It("GetConnections should return connections", func() {
			Expect(len(r.GetConnections())).To(Equal(3))
		})
	})
})
