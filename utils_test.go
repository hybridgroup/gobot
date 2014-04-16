package gobot

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Utils", func() {

	var (
		testInterface interface{}
	)

	Context("when valid", func() {
		It("should execute function at every interval", func() {
			var i = 0
			Every("500ms", func() {
				i = i + 1
			})
			time.Sleep(600 * time.Millisecond)
			Expect(i).To(Equal(1))
			time.Sleep(600 * time.Millisecond)
			Expect(i).To(Equal(2))
		})
		It("should execute function after specific interval", func() {
			var i = 0
			After("500ms", func() {
				i = i + 1
			})
			time.Sleep(600 * time.Millisecond)
			Expect(i).To(Equal(1))
			time.Sleep(600 * time.Millisecond)
			Expect(i).To(Equal(1))
		})
		It("should Publish message to channel without blocking", func() {
			c := make(chan interface{}, 1)
			Publish(c, 1)
			Publish(c, 2)
			i := <-c
			Expect(i.(int)).To(Equal(1))
		})
		It("should execution function on event", func() {
			c := make(chan interface{})
			var i int
			On(c, func(data interface{}) {
				i = data.(int)
			})
			c <- 10
			Expect(i).To(Equal(10))
		})
		It("should scale the value between 0...1", func() {
			Expect(FromScale(5, 0, 10)).To(Equal(0.5))
		})
		It("should scale the 0...1 to scale ", func() {
			Expect(ToScale(500, 0, 10)).To(Equal(float64(10)))
			Expect(ToScale(-1, 0, 10)).To(Equal(float64(0)))
			Expect(ToScale(0.5, 0, 10)).To(Equal(float64(5)))
		})
		It("should return random int", func() {
			a := Rand(100)
			b := Rand(100)
			Expect(a).NotTo(Equal(b))
		})
		It("should return the Field", func() {
			testInterface = *newTestStruct()
			Expect(FieldByName(testInterface, "i").Int()).To(Equal(int64(10)))
		})
		It("should return the Field from ptr", func() {
			testInterface = newTestStruct()
			Expect(FieldByNamePtr(testInterface, "f").Float()).To(Equal(0.2))
		})
		It("should call function on interface", func() {
			testInterface = newTestStruct()
			Expect(Call(testInterface, "Hello", "Human", "How are you?")[0].String()).To(Equal("Hello Human! How are you?"))
		})
	})
})
