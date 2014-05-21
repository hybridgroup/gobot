package ardrone

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArdroneDriver", func() {
	var (
		driver *ArdroneDriver
	)

	BeforeEach(func() {
		adaptor := NewArdroneAdaptor("drone")
		adaptor.connect = func(a *ArdroneAdaptor) {
			a.drone = &testDrone{}
		}
		driver = NewArdroneDriver(adaptor, "drone")
		adaptor.Connect()
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
	It("Must be able to TakeOff", func() {
		driver.TakeOff()
	})
	It("Must be able to Land", func() {
		driver.Land()
	})
	It("Must be able to go Up", func() {
		driver.Up(1)
	})
	It("Must be able to go Down", func() {
		driver.Down(1)
	})
	It("Must be able to go Left", func() {
		driver.Left(1)
	})
	It("Must be able to go Right", func() {
		driver.Right(1)
	})
	It("Must be able to go Forward", func() {
		driver.Forward(1)
	})
	It("Must be able to go Backward", func() {
		driver.Backward(1)
	})
	It("Must be able to go Clockwise", func() {
		driver.Clockwise(1)
	})
	It("Must be able to go CounterClockwise", func() {
		driver.CounterClockwise(1)
	})
	It("Must be able to Hover", func() {
		driver.Hover()
	})
})
