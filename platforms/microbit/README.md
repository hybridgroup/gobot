# Microbit

The Microbit is a tiny computer with built-in Bluetooth LE aka Bluetooth 4.0.

## How to Install
```
go get -d -u gobot.io/x/gobot/... && go install gobot.io/x/gobot/platforms/microbit
```

## How to Use
```go
// code here...
```

## How to Connect

The Microbit is a Bluetooth LE device.

You need to know the BLE ID of the Microbit that you want to connect to.

### OSX

To run any of the Gobot BLE code you must use the `GODEBUG=cgocheck=0` flag in order to get around some of the issues in the CGo-based implementation.

For example:

    GODEBUG=cgocheck=0 go run examples/microbit_blink.go "BBC micro:bit"

OSX uses its own Bluetooth ID system which is different from the IDs used on Linux. The code calls thru the XPC interfaces provided by OSX, so as a result does not need to run under sudo.

### Ubuntu

On Linux the BLE code will need to run as a root user account. The easiest way to accomplish this is probably to use `go build` to build your program, and then to run the requesting executable using `sudo`.

For example:

    go build examples/microbit_blink.go
    sudo ./microbit_blink "BBC micro:bit"

### Windows

Hopefully coming soon...
