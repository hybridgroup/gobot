/*
Package neurosky contains the Gobot adaptor and driver for the Neurosky Mindwave Mobile EEG.

Installing:

	Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

Example:

	package main

	import (
		"fmt"

		"gobot.io/x/gobot/v2"
		"gobot.io/x/gobot/v2/platforms/neurosky"
	)

	func main() {
		adaptor := neurosky.NewAdaptor("/dev/rfcomm0")
		neuro := neurosky.NewDriver(adaptor)

		work := func() {
			_ = neuro.On(neuro.Event("extended"), func(data interface{}) {
				fmt.Println("Extended", data)
			})
			_ = neuro.On(neuro.Event("signal"), func(data interface{}) {
				fmt.Println("Signal", data)
			})
			_ = neuro.On(neuro.Event("attention"), func(data interface{}) {
				fmt.Println("Attention", data)
			})
			_ = neuro.On(neuro.Event("meditation"), func(data interface{}) {
				fmt.Println("Meditation", data)
			})
			_ = neuro.On(neuro.Event("blink"), func(data interface{}) {
				fmt.Println("Blink", data)
			})
			_ = neuro.On(neuro.Event("wave"), func(data interface{}) {
				fmt.Println("Wave", data)
			})
			_ = neuro.On(neuro.Event("eeg"), func(data interface{}) {
				eeg := data.(neurosky.EEGData)
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

		if err := robot.Start(); err != nil {
			panic(err)
		}
	}

For further information refer to neuroky README:
https://github.com/hybridgroup/gobot/blob/master/platforms/neurosky/README.md
*/
package neurosky // import "gobot.io/x/gobot/v2/platforms/neurosky"
