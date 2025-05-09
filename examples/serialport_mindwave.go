//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/serial/neurosky"
	"gobot.io/x/gobot/v2/platforms/serialport"
)

func main() {
	adaptor := serialport.NewAdaptor("/dev/rfcomm0", serialport.WithName("Neurosky"), serialport.WithBaudRate(57600))
	neuro := neurosky.NewMindWaveDriver(adaptor)

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
			eeg := data.(neurosky.MindWaveEEGData)
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

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
