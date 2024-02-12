# Neurosky

NeuroSky delivers fully integrated, single chip EEG biosensors. NeuroSky enables its partners and developers to bring their
brainwave application ideas to market with the shortest amount of time, and lowest end consumer price.

This package contains the Gobot adaptor and driver for the [Neurosky MindWave Mobile EEG](http://store.neurosky.com/products/mindwave-mobile).

## How to Install

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

## How To Connect

### OSX

In order to allow Gobot running on your Mac to access the MindWave, go to "Bluetooth > Open Bluetooth Preferences > Sharing Setup"
and make sure that "Bluetooth Sharing" is checked.

Now you must pair with the MindWave. Open System Preferences > Bluetooth. Now with the Bluetooth devices windows open, hold
the On/Pair button on the MindWave towards the On/Pair text until you see "MindWave" pop up as available devices. Pair with
that device. Once paired your MindWave will be accessable through the serial device similarly named as `/dev/tty.MindWaveMobile-DevA`

### Ubuntu

Connecting to the MindWave from Ubuntu or any other Linux-based OS can be done entirely from the command line using [Gort](https://gobot.io/x/gort)
CLI commands. Here are the steps.

Find the address of the MindWave, by using:

```sh
gort scan bluetooth
```

Pair to MindWave using this command (substituting the actual address of your MindWave):

```sh
gort bluetooth pair <address>
```

Connect to the MindWave using this command (substituting the actual address of your MindWave):

```sh
gort bluetooth connect <address>
```

### Windows

You should be able to pair your MindWave using your normal system tray applet for Bluetooth, and then connect to the
COM port that is bound to the device, such as `COM3`.

## How to Use

Please refer to the provided example `examples/serialport_neurosky.go`.
