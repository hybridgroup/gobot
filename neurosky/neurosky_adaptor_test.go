package gobotNeurosky

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NeuroskyAdaptor", func() {
	var (
		adaptor *NeuroskyAdaptor
	)

	BeforeEach(func() {
		adaptor = new(NeuroskyAdaptor)
	})

	PIt("Must be able to Finalize", func() {
		Expect(adaptor.Finalize()).To(Equal(true))
	})
	PIt("Must be able to Connect", func() {
		Expect(adaptor.Connect()).To(Equal(true))
	})
	PIt("Must be able to Disconnect", func() {
		Expect(adaptor.Disconnect()).To(Equal(true))
	})
	PIt("Must be able to Reconnect", func() {
		Expect(adaptor.Reconnect()).To(Equal(true))
	})
})
