package neurosky

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NeuroskyDriver", func() {
	var (
		n *NeuroskyDriver
	)

	BeforeEach(func() {
		n = NewNeuroskyDriver(NewNeuroskyAdaptor())
	})

	PIt("Must be able to Start", func() {
		Expect(n.Start()).To(Equal(true))
	})
	PIt("Must be able to Init", func() {
		Expect(n.Init()).To(Equal(true))
	})
	PIt("Must be able to Halt", func() {
		Expect(n.Halt()).To(Equal(true))
	})
})
