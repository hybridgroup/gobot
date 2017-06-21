# Neurosky

NeuroSky delivers fully integrated, single chip EEG biosensors. NeuroSky enables its partners and developers to bring their brainwave application ideas to market with the shortest amount of time, and lowest end consumer price.

This package contains the Gobot adaptor and driver for the [Neurosky Mindwave Mobile EEG](http://store.neurosky.com/products/mindwave-mobile).

## How to Install
Installing Gobot with Neurosky support is pretty easy.

```
go get -d -u gobot.io/x/gobot/...
```

## How To Connect

### OSX

In order to allow Gobot running on your Mac to access the Mindwave, go to "Bluetooth > Open Bluetooth Preferences > Sharing Setup" and make sure that "Bluetooth Sharing" is checked.

Now you must pair with the Mindwave. Open System Preferences > Bluetooth. Now with the Bluetooth devices windows open, hold the On/Pair button on the Mindwave towards the On/Pair text until you see "Mindwave" pop up as available devices. Pair with that device. Once paired your Mindwave will be accessable through the serial device similarly named as `/dev/tty.MindWaveMobile-DevA`

### Ubuntu

Connecting to the Mindwave from Ubuntu or any other Linux-based OS can be done entirely from the command line using [Gort](https://gobot.io/x/gort) CLI commands. Here are the steps.

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
			fmt.Println("\n")
		})
	}

	robot := gobot.NewRobot("brainBot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{neuro},
		work,
	)

	robot.Start()
}
```
