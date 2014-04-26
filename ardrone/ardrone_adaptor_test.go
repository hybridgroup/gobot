package gobotArdrone

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArdroneAdaptor", func() {
	var (
		adaptor *ArdroneAdaptor
		ardrone drone
	)

	BeforeEach(func() {
		ardrone = new(testDrone)
		connect = func(me *ArdroneAdaptor) {
			me.ardrone = ardrone
		}
		adaptor = new(ArdroneAdaptor)
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
		Expect(adaptor.Drone()).To(Equal(ardrone))
	})
})
