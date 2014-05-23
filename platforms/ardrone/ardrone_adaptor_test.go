package ardrone

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArdroneAdaptor", func() {
	var (
		adaptor *ArdroneAdaptor
		drone   *testDrone
	)

	BeforeEach(func() {
		drone = &testDrone{}
		adaptor = NewArdroneAdaptor("drone")
		adaptor.connect = func(a *ArdroneAdaptor) {
			a.drone = drone
		}
	})

	It("Must be able to Finalize", func() {
		Expect(adaptor.Finalize()).To(Equal(true))
	})
	It("Must be able to Connect", func() {
		Expect(adaptor.Connect()).To(Equal(true))
	})
	It("Must be able to Disconnect", func() {
		Expect(adaptor.Disconnect()).To(Equal(true))
	})
	It("Must be able to Reconnect", func() {
		Expect(adaptor.Reconnect()).To(Equal(true))
	})
	It("Must be able to return a Drone", func() {
		adaptor.Connect()
		Expect(adaptor.Drone()).To(Equal(drone))
	})
})
