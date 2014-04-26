package gobotSphero

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SpheroAdaptor", func() {
	var (
		adaptor *SpheroAdaptor
	)

	BeforeEach(func() {
		adaptor = new(SpheroAdaptor)
		adaptor.sp = sp{}
		connect = func(sa *SpheroAdaptor) {
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
