# UP2 (Squared)

The UP2 Board is a single board SoC computer based on the Intel Apollo Lake processor. It has built-in GPIO, PWM, SPI, and I2C interfaces.

For more info about the UP2 Board, go to [http://www.up-board.org/upsquared/](http://www.up-board.org/upsquared/).

## How to Install

We recommend updating to the latest Ubuntu when using the UP2.

You would normally install Go and Gobot on your workstation. Once installed, cross compile your program on your workstation, transfer the final executable to your UP2, and run the program on the UP2 as documented here.

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use

The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.

```go
r := up2.NewAdaptor()
led := gpio.NewLedDriver(r, "13")
```

## How to Connect

### Compiling

Compile your Gobot program on your workstation like this:

```bash
$ GOARCH=386 GOOS=linux go build examples/up2_blink.go
```

Once you have compiled your code, you can you can upload your program and execute it on the UP2 from your workstation using the `scp` and `ssh` commands like this:

```bash
$ scp up2_blink ubuntu@192.168.1.xxx:/home/ubuntu/
$ ssh -t ubuntu@192.168.1.xxx "./up2_blink"
```
