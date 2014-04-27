package sphero

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SpheroAdaptor", func() {
	var (
		a *SpheroAdaptor
	)

	BeforeEach(func() {
		a = NewSpheroAdaptor()
		a.sp = sp{}
		a.connect = func(a *SpheroAdaptor) {}
	})

	It("Must be able to Finalize", func() {
		Expect(a.Finalize()).To(Equal(true))
	})
	It("Must be able to Connect", func() {
		Expect(a.Connect()).To(Equal(true))
	})
	It("Must be able to Disconnect", func() {
		Expect(a.Disconnect()).To(Equal(true))
	})
	It("Must be able to Reconnect", func() {
		Expect(a.Reconnect()).To(Equal(true))
	})
})
