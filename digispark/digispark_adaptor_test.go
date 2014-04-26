package gobotDigispark

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Digispark", func() {
	var (
		adaptor *DigisparkAdaptor
	)

	BeforeEach(func() {
		adaptor = new(DigisparkAdaptor)
		connect = func() *LittleWire {
			return nil
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
})
