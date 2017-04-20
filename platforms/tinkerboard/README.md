# Tinker Board

The ASUS Tinker Board is a single board SoC computer based on the Rockchip RK3288 processor. It has built-in GPIO, PWM, SPI, and I2C interfaces.

For more info about the Tinker Board, go to [https://www.asus.com/uk/Single-Board-Computer/Tinker-Board/](https://www.asus.com/uk/Single-Board-Computer/Tinker-Board/).

## How to Install

We recommend updating to the latest Debian TinkerOS when using the Tinker Board.

You would normally install Go and Gobot on your workstation. Once installed, cross compile your program on your workstation, transfer the final executable to your Tinker Board, and run the program on the Tinker Board as documented here.

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use

The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.

```go
// code here...
```

## How to Connect

### Compiling

Compile your Gobot program on your workstation like this:

```bash
$ GOARM=7 GOARCH=arm GOOS=linux go build examples/tinkerboard_blink.go
```

Once you have compiled your code, you can you can upload your program and execute it on the Tinkerboard from your workstation using the `scp` and `ssh` commands like this:

```bash
echo "todo"
```
