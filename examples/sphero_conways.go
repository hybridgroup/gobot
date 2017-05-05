// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/sphero"
)

type conway struct {
	alive    bool
	age      int
	contacts int
	cell     *sphero.SpheroDriver
}

func main() {
	master := gobot.NewMaster()

	spheros := []string{
		"/dev/rfcomm0",
		"/dev/rfcomm1",
		"/dev/rfcomm2",
	}

	for _, port := range spheros {
		spheroAdaptor := sphero.NewAdaptor(port)

		cell := sphero.NewSpheroDriver(spheroAdaptor)
		cell.SetName("Sphero" + port)

		work := func() {
			conway := new(conway)
			conway.cell = cell

			conway.birth()

			cell.On(sphero.Collision, func(data interface{}) {
				conway.contact()
			})

			gobot.Every(3*time.Second, func() {
				if conway.alive {
					conway.movement()
				}
			})

			gobot.Every(10*time.Second, func() {
				if conway.alive {
					conway.birthday()
				}
			})
		}

		robot := gobot.NewRobot("conway",
			[]gobot.Connection{spheroAdaptor},
			[]gobot.Device{cell},
			work,
		)

		master.AddRobot(robot)
	}

	master.Start()
}

func (c *conway) resetContacts() {
	c.contacts = 0
}

func (c *conway) contact() {
	c.contacts++
}

func (c *conway) rebirth() {
	fmt.Println("Welcome back", c.cell.Name(), "!")
	c.life()
}

func (c *conway) birth() {
	c.resetContacts()
	c.age = 0
	c.life()
	c.movement()
}

func (c *conway) life() {
	c.cell.SetRGB(0, 255, 0)
	c.alive = true
}

func (c *conway) death() {
	fmt.Println(c.cell.Name(), "died :(")
	c.alive = false
	c.cell.SetRGB(255, 0, 0)
	c.cell.Stop()
}

func (c *conway) enoughContacts() bool {
	if c.contacts >= 2 && c.contacts < 7 {
		return true
	}
	return false
}

func (c *conway) birthday() {
	c.age++

	fmt.Println("Happy birthday", c.cell.Name(), "you are", c.age, "and had", c.contacts, "contacts.")

	if c.enoughContacts() {
		if !c.alive {
			c.rebirth()
		}
	} else {
		c.death()
	}

	c.resetContacts()
}

func (c *conway) movement() {
	if c.alive {
		c.cell.Roll(100, uint16(gobot.Rand(360)))
	}
}
