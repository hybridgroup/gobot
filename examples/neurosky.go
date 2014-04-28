package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/neurosky"
)

func main() {

	adaptor := neurosky.NewNeuroskyAdaptor()
	adaptor.Name = "neurosky"
	adaptor.Port = "/dev/rfcomm0"

	neuro := neurosky.NewNeuroskyDriver(adaptor)
	neuro.Name = "neuro"

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

	robot := gobot.Robot{
		Connections: []gobot.Connection{adaptor},
		Devices:     []gobot.Device{neuro},
		Work:        work,
	}

	robot.Start()
}
