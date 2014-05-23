package neurosky

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NeuroskyAdaptor", func() {
	var (
		n *NeuroskyAdaptor
	)

	BeforeEach(func() {
		n = NewNeuroskyAdaptor("bot", "/dev/null")
	})

	PIt("Must be able to Finalize", func() {
		Expect(n.Finalize()).To(Equal(true))
	})
	PIt("Must be able to Connect", func() {
		Expect(n.Connect()).To(Equal(true))
	})
})
