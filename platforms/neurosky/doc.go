/*
Package neurosky contains the Gobot adaptor and driver for the Neurosky Mindwave Mobile EEG.

Installing:

	go get github.com/hybridgroup/gobot/platforms/neurosky

Example:

	package main

	import (
		"fmt"

		"github.com/hybridgroup/gobot"
		"github.com/hybridgroup/gobot/platforms/neurosky"
	)

	func main() {
		gbot := gobot.NewGobot()

		adaptor := neurosky.NewNeuroskyAdaptor("neurosky", "/dev/rfcomm0")
		neuro := neurosky.NewNeuroskyDriver(adaptor, "neuro")

		work := func() {
			gobot.On(neuro.Event("extended"), func(data interface{}) {
				fmt.Println("Extended", data)
			})
			gobot.On(neuro.Event("signal"), func(data interface{}) {
				fmt.Println("Signal", data)
			})
			gobot.On(neuro.Event("attention"), func(data interface{}) {
				fmt.Println("Attention", data)
			})
			gobot.On(neuro.Event("meditation"), func(data interface{}) {
				fmt.Println("Meditation", data)
			})
			gobot.On(neuro.Event("blink"), func(data interface{}) {
				fmt.Println("Blink", data)
			})
			gobot.On(neuro.Event("wave"), func(data interface{}) {
				fmt.Println("Wave", data)
			})
			gobot.On(neuro.Event("eeg"), func(data interface{}) {
				eeg := data.(neurosky.EEG)
				fmt.Println("Delta", eeg.Delta)
				fmt.Println("Theta", eeg.Theta)
				fmt.Println("LoAlpha", eeg.LoAlpha)
				fmt.Println("HiAlpha", eeg.HiAlpha)
				fmt.Println("LoBeta", eeg.LoBeta)
				fmt.Println("HiBeta", eeg.HiBeta)
				fmt.Println("LoGamma", eeg.LoGamma)
				fmt.Println("MidGamma", eeg.MidGamma)
				fmt.Println("\n")
			})
		}

		robot := gobot.NewRobot("brainBot",
			[]gobot.Connection{adaptor},
			[]gobot.Device{neuro},
			work,
		)

		gbot.AddRobot(robot)
		gbot.Start()
	}

For further information refer to neuroky README:
https://github.com/hybridgroup/gobot/blob/master/platforms/neurosky/README.md
*/
package neurosky
