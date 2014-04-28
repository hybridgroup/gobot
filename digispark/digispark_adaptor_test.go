package digispark

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Digispark", func() {
	var (
		d *DigisparkAdaptor
	)

	BeforeEach(func() {
		d = NewDigisparkAdaptor()
		d.connect = func(d *DigisparkAdaptor) {}
	})

	It("Must be able to Finalize", func() {
		Expect(d.Finalize()).To(Equal(true))
	})
	It("Must be able to Connect", func() {
		Expect(d.Connect()).To(Equal(true))
	})
	It("Must be able to Disconnect", func() {
		Expect(d.Disconnect()).To(Equal(true))
	})
	It("Must be able to Reconnect", func() {
		Expect(d.Reconnect()).To(Equal(true))
	})
})
