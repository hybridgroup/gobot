package spark

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Spark", func() {
	var (
		s *SparkCoreAdaptor
	)

	BeforeEach(func() {
		s = NewSparkCoreAdaptor("bot", "", "")
	})

	It("Must be able to Finalize", func() {
		Expect(s.Finalize()).To(Equal(true))
	})
	It("Must be able to Connect", func() {
		Expect(s.Connect()).To(Equal(true))
	})
})
