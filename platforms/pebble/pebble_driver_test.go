package pebble

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PebbleDriver", func() {
	var (
		driver  *PebbleDriver
		adaptor *PebbleAdaptor
	)

	BeforeEach(func() {
		adaptor = NewPebbleAdaptor("pebble")
		driver = NewPebbleDriver(adaptor, "pebble")
	})

	It("Must be able to Start", func() {
		Expect(driver.Start()).To(Equal(true))
	})

	It("Must be able to Halt", func() {
		Expect(driver.Halt()).To(Equal(true))
	})

	It("Adds message when sending notification", func() {
		driver.SendNotification("Hello")
		Expect(driver.Messages[0]).To(Equal("Hello"))
	})

	It("Retrieves pending messages", func() {
		driver.SendNotification("Hello")
		driver.SendNotification("World")

		Expect(driver.PendingMessage()).To(Equal("Hello"))
		Expect(driver.PendingMessage()).To(Equal("World"))
		Expect(driver.PendingMessage()).To(Equal(""))
	})
})
