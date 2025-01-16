# Migration guide

From time to time a breaking change of API can happen. Following to [SemVer](https://semver.org/), the gobot main version
should be increased. In such case all users needs to adjust there projects for the next update, although they not using
a driver or platform with changed API.

To prevent this scenario for most users, the main version will not always increased, but affected drivers and platforms
are listed here and a migration strategy is provided.

## Switch from version 2.4.0 (applications using the gpiod options affected)

### NanoPi NEO adaptor was moved to friendlyelec folder

With introduce of FriendlyELEC NanoPC-T6 a second adaptor from FriendlyELEC (formerly friendlarm) now exists. Please
search and replace to change the import path as follows.

```go
// old
...
  "gobot.io/x/gobot/v2/platforms/nanopi"
...

// new
...
  "gobot.io/x/gobot/v2/platforms/friendlyelec/nanopi"
...
```

### The term gpiod was renamed to cdev

Using the term "cdev" (short for character device Kernel ABI for GPIO access) is more suitable than using "gpiod" (the
name of the user space driver in Linux). Also it relates better to the term "sysfs" (the legacy sysfs Kernel ABI for
GPIO access). The former name was chosen so there would be no difference to the used go module "gpiod". Since also this
module is now renamed to "go-gpiocdev", we choose the better name "cdev" from now on.

This change affects all applications, which using the With... options of "gpiod" or "sysfs". A search and replace is
suitable:

```go
// old
...
  a := NewAdaptor(adaptors.WithGpiodAccess())
...

// new
...
  a := NewAdaptor(adaptors.WithGpioCdevAccess())
...
```

```go
// old
...
  a := NewAdaptor(adaptors.WithSysfsAccess())
...

// new
...
  a := NewAdaptor(adaptors.WithGpioSysfsAccess())
...
```

Also those findings needs to be replaced, which usually affects developers, but not users:

* `system.WithDigitalPinGpiodAccess()` --> `system.WithDigitalPinCdevAccess()`
* `IsGpiodDigitalPinAccess()` --> `HasDigitalPinCdevAccess()`

### PocketBeagle adaptor goes cdev

The beagle board "PocketBeagle" supports with latest images the Linux Kernel character device API, so the adaptor was
changed to use this as the default.

By "adaptors.WithSysfsAccess()" the old behavior can be forced. This is most likely needed if an old image with old
Kernel is used which does not support this new API.

A small renaming is also done, please search and replace:

```go
// old
...
"gobot.io/x/gobot/v2/platforms/beaglebone"
...
beaglebone.NewPocketBeagleAdaptor(...)
...

// new
...
"gobot.io/x/gobot/v2/platforms/beagleboard/pocketbeagle"
...
pocketbeagle.NewAdaptor(...)
...
```

## Switch from version 2.3.0 (ble and sphero adaptors affected)

### BLE drivers and client adaptor

All BLE drivers now can be found in the folder "drivers/ble". Formerly the drivers are located below "platforms/ble".
In addition the location of the BLE client adaptor was changed to "platforms/bleclient". Therefore a change for the
import paths is needed. The constructor function was also renamed, see below.

```go
// old
import(
  ...
  "gobot.io/x/gobot/v2/platforms/ble"
  ...
)

...
  bleAdaptor := ble.NewClientAdaptor(os.Args[1])
...

// new
import(
  ...
  "gobot.io/x/gobot/v2/drivers/ble"
  "gobot.io/x/gobot/v2/platforms/bleclient"
  ...
)
...
  bleAdaptor := bleclient.NewAdaptor(os.Args[1])
...
```

### BLE client adaptor changed signature for Subscribe()

Since introducing the usage of "github.com/muka/go-bluetooth" in 2020, the callback do not support the given error
parameter anymore. The switch to usage of "tinygo.org/x/bluetooth" has not changed this. Therefore it is removed now
from the function.

### BLE generic drivers changed signature for Get*() functions

All those functions log an error only or panic, so the caller gets no nice programmatic feedback. The error is now
returned instead and the log output needs to be done at caller side.

```go
// old
...
  devName := access.GetDeviceName()
  appearance := access.GetAppearance()
  modelNo := info.GetModelNumber()
  fwRev := info.GetFirmwareRevision()
  hwRev := info.GetHardwareRevision()
  manuName := info.GetManufacturerName()
  pid := info.GetPnPId()
  level := battery.GetBatteryLevel()
...

// new
...
  devName, err := access.GetDeviceName()
  if err != nil {
    fmt.Println(err)
  }
  appearance, err := access.GetAppearance()
  if err != nil {
    fmt.Println(err)
  }
  ...
...
```

### Sphero adaptor split off

The Serial Based Sphero adaptor was split off into a generic serial adaptor and the driver part. With this, the imports
needs to be adjusted. In addition all events now have a suffix "Event", see below.

```go
// old
import(
  ...
  "gobot.io/x/gobot/v2/platforms/sphero"
  ...
)

...
  adaptor := sphero.NewAdaptor("/dev/rfcomm0")
  spheroDriver := sphero.NewSpheroDriver(adaptor)
...
  _ = spheroDriver.On(sphero.Collision, func(data interface{}) {
...

// new
import(
  ...
  "gobot.io/x/gobot/v2/drivers/common/spherocommon"
  "gobot.io/x/gobot/v2/drivers/serial"
  "gobot.io/x/gobot/v2/platforms/serialport"
  ...
)
...
  adaptor := serialport.NewAdaptor("/dev/rfcomm0")
  spheroDriver := sphero.NewSpheroDriver(adaptor)
...
  _ = spheroDriver.On(sphero.CollisionEvent, func(data interface{}) {
...
```

### Neurosky adaptor split off

The Neurosky adaptor now use the generic serial adaptor. The driver part was moved. With this, the imports needs to be
adjusted. In addition all events now have a suffix "Event", see below.

```go
// old
import(
  ...
  "gobot.io/x/gobot/v2/platforms/neurosky"
  ...
)

...
  adaptor := neurosky.NewAdaptor("/dev/rfcomm0")
  neuro := neurosky.NewDriver(adaptor)
...
  _ = neuro.On(neurosky.Extended, func(data interface{}) {
...

// new
import(
  ...
  "gobot.io/x/gobot/v2/drivers/serial/neurosky"
  "gobot.io/x/gobot/v2/platforms/serialport"
  ...
)
...
  adaptor := serialport.NewAdaptor("/dev/rfcomm0", serialport.WithName("Neurosky"), serialport.WithBaudRate(57600))
  neuro := neurosky.NewMindWaveDriver(adaptor)
...
  _ = neuro.On(neurosky.ExtendedEvent, func(data interface{}) {
...
```

### MegaPi adaptor split off

The MegaPi adaptor now use the generic serial adaptor. The driver part was moved. With this, the imports needs to be
adjusted.

```go
// old
import(
  ...
  "gobot.io/x/gobot/v2/platforms/megapi"
  ...
)

...
  megaPiAdaptor := megapi.NewAdaptor("/dev/ttyS0")
  motor := megapi.NewMotorDriver(megaPiAdaptor, 1)
...

// new
import(
  ...
  "gobot.io/x/gobot/v2/drivers/serial/megapi"
  "gobot.io/x/gobot/v2/platforms/serialport"
  ...
)
...
  adaptor := serialport.NewAdaptor("/dev/ttyS0", serialport.WithName("MegaPi"))
  motor := megapi.NewMotorDriver(adaptor, 1)
...
```

## Switch from version 2.2.0 (gpio drivers affected)

### gpio.ButtonDriver, gpio.PIRMotionDriver: substitute parameter "v time.duration"

A backward compatible case is still included, but it is recommended to use "WithButtonPollInterval" instead, see example
below.

```go
// old
d := gpio.NewButtonDriver(adaptor, "1", 50*time.Millisecond)

// new
d := gpio.NewButtonDriver(adaptor, "1", gpio.WithButtonPollInterval(50*time.Millisecond))
```

### gpio.EasyDriver: optional pins

There is no need to use the direction, enable or sleep feature of the driver. Therefore the parameters are removed from
constructor. Please migrate according to the examples below. The order of the optional functions does not matter.

```go
// old
d0 := gpio.NewEasyDriver(adaptor, 0.80, "1", "", "", "")
d1 := gpio.NewEasyDriver(adaptor, 0.81, "11", "12", "", "")
d2 := gpio.NewEasyDriver(adaptor, 0.82, "21", "22", "23", "")
d3 := gpio.NewEasyDriver(adaptor, 0.83, "31", "32", "33", "34")

// new
d0 := gpio.NewEasyDriver(adaptor, 0.80, "1")
d1 := gpio.NewEasyDriver(adaptor, 0.81, "11", gpio.WithEasyDirectionPin("12"))
d2 := gpio.NewEasyDriver(adaptor, 0.82, "21", gpio.WithEasyDirectionPin("22"), gpio.WithEasyEnablePin("23"))
d3 := gpio.NewEasyDriver(adaptor, 0.83, "31", gpio.WithEasyDirectionPin("32"), gpio.WithEasyEnablePin("33"),
  gpio.WithEasySleepPin("34"))
```

### gpio.BuzzerDriver: unexport 'BPM' attribute

```go
d := gpio.NewBuzzerDriver(adaptor, "1")
// old
d.BPM = 120.0
fmt.Println("BPM:", d.BPM)

// new
d.SetBPM(120.0)
fmt.Println("BPM:", d.BPM())
```

### gpio.RelayDriver: unexport 'Inverted' attribute

Usually the relay is inverted or not, except be rewired. From now on the inverted behavior can only be changed on
initialization. If there is really a different use case, please file a new issue.

```go
// old
d := gpio.NewRelayDriver(adaptor, "1")
d.Inverted = true
fmt.Println("is inverted:", d.Inverted)

// new
d := gpio.NewRelayDriver(adaptor, "1", gpio.WithRelayInverted())
fmt.Println("is inverted:", d.IsInverted())
```

### gpio.HD44780Driver: make 'SetRWPin()' an option

```go
// old
d := gpio.NewHD44780Driver(adaptor, ...)
d.SetRWPin("10")

// new
d := gpio.NewHD44780Driver(adaptor, ..., gpio.WithHD44780RWPin("10"))
```

### gpio.ServoDriver: unexport 'CurrentAngle' and rename functions 'Min()', 'Max()', 'Center()'

```go
d := gpio.NewServoDriver(adaptor, "1")
// old
d.Min()
fmt.Println("current position:", d.CurrentAngle)
d.Center()
d.Max()

// new
d.ToMin()
fmt.Println("current position:", d.Angle())
d.ToCenter()
d.ToMax()
```

### gpio.MotorDriver: unexport pin and state attributes, rename functions

The motor driver was heavily revised - sorry for the inconveniences.

affected pins:

* SpeedPin
* SwitchPin (removed, was unused)
* DirectionPin
* ForwardPin
* BackwardPin

Usually the pins will not change without a hardware rewiring. All pins, except the speed pin are optionally, so options
are designed for that.

```go
// old
d := gpio.NewMotorDriver(adaptor, "1")
d.DirectionPin = "10"

// new
d := gpio.NewMotorDriver(adaptor, "1", gpio.WithMotorDirectionPin("10"))
```

```go
// old
d := gpio.NewMotorDriver(adaptor, "1")
d.ForwardPin = "10"
d.BackWardPin = "11"

// new
d := gpio.NewMotorDriver(adaptor, "1", gpio.WithMotorForwardPin("10"), gpio.WithMotorBackwardPin("11"))
```

affected functions:

* Speed() --> SetSpeed()
* Direction() --> SetDirection()
* Max() --> RunMax()
* Min() --> RunMin()

affected states:

* CurrentState
* CurrentSpeed
* CurrentMode
* CurrentDirection

Most of the attributes were used only for reading. If there is something missing, please file a new issue.

```go
d := gpio.NewMotorDriver(adaptor, "1")
// old
d.On()
fmt.Println("is on:", d.CurrentState==1)
fmt.Println("speed:", d.CurrentSpeed)
d.Off()
fmt.Println("is off:", d.CurrentState==0)
fmt.Println("mode is digital:", d.CurrentMode=="digital")
fmt.Println("direction:", d.CurrentDirection)

// new
d.On()
fmt.Println("is on:", d.IsOn())
d.Off()
fmt.Println("is on:", d.IsOff())
fmt.Println("speed:", d.Speed())
fmt.Println("mode is digital:", d.IsDigital())
fmt.Println("direction:", d.Direction())
```

```go
d := gpio.NewMotorDriver(adaptor, "1")
// old
d.Speed(123)
fmt.Println("is mode now analog?", d.CurrentMode!="digital")

// new
d.SetSpeed(123)
fmt.Println("is mode now analog?", d.IsAnalog())
```

Although, it is working like above, it will be more clear, if the mode is defined at the beginning, like so.

```go
// old
d := gpio.NewMotorDriver(adaptor, "1")
d.CurrentMode=="analog"
d.Max()

// new
d := gpio.NewMotorDriver(adaptor, "1", gpio.WithMotorAnalog())
d.RunMax()
```
