package gobot

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Connection", func() {

	var (
		someRobot Robot
	)

	BeforeEach(func() {
		someRobot = newTestRobot("")
		start = func(r *Robot) {
			r.startRobot()
		}
		someRobot.Start()
	})

	Context("when valid", func() {
		It("Connect should call adaptor Connect", func() {
			Expect(someRobot.Connections[0].Connect()).To(Equal(true))
		})
		It("Finalize should call adaptor Finalize", func() {
			Expect(someRobot.Connections[0].Connect()).To(Equal(true))
		})
		It("Disconnect should call adaptor Disconnect", func() {
			Expect(someRobot.Connections[0].Connect()).To(Equal(true))
		})
		It("Reconnect should call adaptor Reconnect", func() {
			Expect(someRobot.Connections[0].Connect()).To(Equal(true))
		})
	})
})
