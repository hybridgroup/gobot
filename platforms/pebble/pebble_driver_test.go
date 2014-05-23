package pebble

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("PebbleDriver", func() {
  var (
    driver *PebbleDriver
  )

  BeforeEach(func() {
    driver = NewPebble(new(PebbleAdaptor))
  })

  It("Must be able to Start", func() {
    Expect(driver.Start()).To(Equal(true))
  })
  It("Must be able to Init", func() {
    Expect(driver.Init()).To(Equal(true))
  })
  It("Must be able to Halt", func() {
    Expect(driver.Halt()).To(Equal(true))
  })
})
