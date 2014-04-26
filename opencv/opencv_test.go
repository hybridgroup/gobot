package gobotOpencv

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Opencv", func() {
	var (
		adaptor *Opencv
	)

	BeforeEach(func() {
		adaptor = new(Opencv)
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
