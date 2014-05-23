package pebble

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("PebbleAdaptor", func() {
  var (
    adaptor *PebbleAdaptor
  )

  BeforeEach(func() {
    adaptor = new(PebbleAdaptor)
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
