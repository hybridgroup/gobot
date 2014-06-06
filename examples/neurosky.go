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
		gobot.On(neuro.Events["Extended"], func(data interface{}) {
			fmt.Println("Extended", data)
		})
		gobot.On(neuro.Events["Signal"], func(data interface{}) {
			fmt.Println("Signal", data)
		})
		gobot.On(neuro.Events["Attention"], func(data interface{}) {
			fmt.Println("Attention", data)
		})
		gobot.On(neuro.Events["Meditation"], func(data interface{}) {
			fmt.Println("Meditation", data)
		})
		gobot.On(neuro.Events["Blink"], func(data interface{}) {
			fmt.Println("Blink", data)
		})
		gobot.On(neuro.Events["Wave"], func(data interface{}) {
			fmt.Println("Wave", data)
		})
		gobot.On(neuro.Events["EEG"], func(data interface{}) {
			eeg := data.(gobotNeurosky.EEG)
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

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("brainBot", []gobot.Connection{adaptor}, []gobot.Device{neuro}, work))

	gbot.Start()
}
