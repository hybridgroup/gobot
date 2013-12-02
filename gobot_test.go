package gobot

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Robot", func() {

var (
        someRobot Robot
    )

    BeforeEach(func() {
        someRobot = Robot{
          Work: func() {
          },
        }
    })

    Describe("Robots", func() {
        Context("when true", func() {
            It("should be true", func() {
                Expect(true)
            })
        })
    })
})
