# Neurosky

NeuroSky delivers fully integrated, single chip EEG biosensors. NeuroSky enables its partners and developers to bring their brainwave application ideas to market with the shortest amount of time, and lowest end consumer price.

This package contains the Gobot adaptor and driver for the [Neurosky Mindwave Mobile EEG](http://store.neurosky.com/products/mindwave-mobile).

## How to Install
Installing Gobot with Neurosky support is pretty easy.

```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/neurosky
```

## How To Connect

### OSX

In order to allow Gobot running on your Mac to access the Mindwave, go to "Bluetooth > Open Bluetooth Preferences > Sharing Setup" and make sure that "Bluetooth Sharing" is checked.

Now you must pair with the Mindwave. Open System Preferences > Bluetooth. Now with the Bluetooth devices windows open, hold the On/Pair button on the Mindwave towards the On/Pair text until you see "Mindwave" pop up as available devices. Pair with that device. Once paired your Mindwave will be accessable through the serial device similarly named as `/dev/tty.MindWaveMobile-DevA`

### Ubuntu

Connecting to the Mindwave from Ubuntu or any other Linux-based OS can be done entirely from the command line using [Gort](https://github.com/hybridgroup/gort) CLI commands. Here are the steps.

Find the address of the Mindwave, by using:
```
gort scan bluetooth
```

Pair to Mindwave using this command (substituting the actual address of your Mindwave):
```
gort bluetooth pair <address>
```

Connect to the Mindwave using this command (substituting the actual address of your Mindwave):
```
gort bluetooth connect <address>
```

### Windows

You should be able to pair your Mindwave using your normal system tray applet for Bluetooth, and then connect to the COM port that is bound to the device, such as `COM3`.

## How to Use

This small program lets you connect the Neurosky an load data.

```go
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
```