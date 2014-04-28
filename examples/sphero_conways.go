package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sphero"
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

	for s := range spheros {
		spheroAdaptor := sphero.NewSpheroAdaptor()
		spheroAdaptor.Name = "Sphero"
		spheroAdaptor.Port = spheros[s]

		cell := sphero.NewSpheroDriver(spheroAdaptor)
		cell.Name = "Sphero" + spheros[s]

		work := func() {

			conway := new(conway)
			conway.cell = cell

			conway.birth()

			gobot.On(sphero.Events["Collision"], func(data interface{}) {
				conway.contact()
			})

			gobot.Every("3s", func() {
				if conway.alive == true {
					conway.movement()
				}
			})

			gobot.Every("10s", func() {
				if conway.alive == true {
					conway.birthday()
				}
			})
		}

		master.Robots = append(master.Robots, &gobot.Robot{
			Connections: []gobot.Connection{spheroAdaptor},
			Devices:     []gobot.Device{cell},
			Work:        work,
		})
	}

	master.Start()
}

func (c *conway) resetContacts() {
	c.contacts = 0
}

func (c *conway) contact() {
	c.contacts += 1
}

func (c *conway) rebirth() {
	fmt.Println("Welcome back", c.cell.Name, "!")
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
	fmt.Println(c.cell.Name, "died :(")
	c.alive = false
	c.cell.SetRGB(255, 0, 0)
	c.cell.Stop()
}

func (c *conway) enoughContacts() bool {
	if c.contacts >= 2 && c.contacts < 7 {
		return true
	} else {
		return false
	}
}

func (c *conway) birthday() {
	c.age += 1

	fmt.Println("Happy birthday", c.cell.Name, "you are", c.age, "and had", c.contacts, "contacts.")

	if c.enoughContacts() == true {
		if c.alive == false {
			c.rebirth()
		}
	} else {
		c.death()
	}

	c.resetContacts()
}

func (c *conway) movement() {
	if c.alive == true {
		c.cell.Roll(100, uint16(gobot.Rand(360)))
	}
}
