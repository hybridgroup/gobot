package beaglebone

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Beaglebone", func() {
	var (
		b *BeagleboneAdaptor
	)

	BeforeEach(func() {
		b = NewBeagleboneAdaptor("bot")
	})

	It("Must be able to Finalize", func() {
		Expect(b.Finalize()).To(Equal(true))
	})
	It("Must be able to Connect", func() {
		Expect(b.Connect()).To(Equal(true))
	})
	It("Must be able to Disconnect", func() {
		Expect(b.Disconnect()).To(Equal(true))
	})
	It("Must be able to Reconnect", func() {
		Expect(b.Reconnect()).To(Equal(true))
	})
})
