// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/neurosky"
)

func main() {
	adaptor := neurosky.NewAdaptor("/dev/rfcomm0")
	neuro := neurosky.NewDriver(adaptor)

	work := func() {
		neuro.On(neuro.Event("extended"), func(data interface{}) {
			fmt.Println("Extended", data)
		})
		neuro.On(neuro.Event("signal"), func(data interface{}) {
			fmt.Println("Signal", data)
		})
		neuro.On(neuro.Event("attention"), func(data interface{}) {
			fmt.Println("Attention", data)
		})
		neuro.On(neuro.Event("meditation"), func(data interface{}) {
			fmt.Println("Meditation", data)
		})
		neuro.On(neuro.Event("blink"), func(data interface{}) {
			fmt.Println("Blink", data)
		})
		neuro.On(neuro.Event("wave"), func(data interface{}) {
			fmt.Println("Wave", data)
		})
		neuro.On(neuro.Event("eeg"), func(data interface{}) {
			eeg := data.(neurosky.EEGData)
			fmt.Println("Delta", eeg.Delta)
			fmt.Println("Theta", eeg.Theta)
			fmt.Println("LoAlpha", eeg.LoAlpha)
			fmt.Println("HiAlpha", eeg.HiAlpha)
			fmt.Println("LoBeta", eeg.LoBeta)
			fmt.Println("HiBeta", eeg.HiBeta)
			fmt.Println("LoGamma", eeg.LoGamma)
			fmt.Println("MidGamma", eeg.MidGamma)
			fmt.Printf("\n\n")
		})
	}

	robot := gobot.NewRobot("brainBot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{neuro},
		work,
	)

	robot.Start()
}
