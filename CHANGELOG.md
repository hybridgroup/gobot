# CHANGELOG

## [Unreleased](https://github.com/hybridgroup/gobot/compare/v2.5.0...HEAD)

## [v2.5.0](https://github.com/hybridgroup/gobot/compare/v2.4.0...v2.5.0) (2025-02-15)

### Build

* **deps:** update modules ([#1131](https://github.com/hybridgroup/gobot/issues/1131))
* **linter:** update linter to v1.64.5 and fix issues ([#1130](https://github.com/hybridgroup/gobot/issues/1130))

### Core

* **system:** split system accesser features ([#1121](https://github.com/hybridgroup/gobot/issues/1121))

### Doc

* **core:** prepare for release v2.5.0

### Drivers

* fix some function names in comment ([#1103](https://github.com/hybridgroup/gobot/issues/1103))

### Examples

* correct all usage of ps3 joystick reads, which are now int instead of int16 ([#1085](https://github.com/hybridgroup/gobot/issues/1085))

### Gpio

* rename gpiod to gpiocdev ([#1116](https://github.com/hybridgroup/gobot/issues/1116))
* **adaptors:** fix so now gpiodev is used as default ([#1112](https://github.com/hybridgroup/gobot/issues/1112))
* **core:** fix sporadic hang on finish for discrete edge polling ([#1107](https://github.com/hybridgroup/gobot/issues/1107))

### Nanopct6

* introduce adaptor for FriendlyELEC NanoPC-T6 ([#1126](https://github.com/hybridgroup/gobot/issues/1126))

### Onewire

* **ds18b20:** introduce 1-wire device access by sysfs and temp driver ([#1091](https://github.com/hybridgroup/gobot/issues/1091))

### Orangepi5pro

* introduce adaptor for OrangePi 5 Pro ([#1129](https://github.com/hybridgroup/gobot/issues/1129))

### Platforms

* file name and folder unification ([#1127](https://github.com/hybridgroup/gobot/issues/1127))

### Pocketbeagle

* introduce and use cdev by default ([#1118](https://github.com/hybridgroup/gobot/issues/1118))

### Rock64

* introduce adaptor for PINE64-ROCK64 ([#1122](https://github.com/hybridgroup/gobot/issues/1122))

### Sphero

* fix panic by removing unnecessary go routine ([#1117](https://github.com/hybridgroup/gobot/issues/1117))

### Tests

* fix for arm64 environment and stabilize ([#1115](https://github.com/hybridgroup/gobot/issues/1115))

### Tinkerboard2

* introduce adapter for Asus Tinker Board 2 ([#1108](https://github.com/hybridgroup/gobot/issues/1108))

### Zero

* introduce adaptor for Radxa Zero ([#1128](https://github.com/hybridgroup/gobot/issues/1128))

## [v2.4.0](https://github.com/hybridgroup/gobot/compare/v2.3.0...v2.4.0) (2024-11-05)

### Bebop

* fix concurrent map writes ([#1063](https://github.com/hybridgroup/gobot/issues/1063))

### Ble

* add support for functional options, add tests ([#1059](https://github.com/hybridgroup/gobot/issues/1059))
* introduce in drivers folder ([#1057](https://github.com/hybridgroup/gobot/issues/1057))
* **client:** add scan timout ([#1051](https://github.com/hybridgroup/gobot/issues/1051))
* **module:** update tinygo.org/x/bluetooth to v0.10 and adapt code ([#1084](https://github.com/hybridgroup/gobot/issues/1084))

### Build

* **go, deps:** switch to Go 1.22 and update modules, linter v1.61.0 and adapt code ([#1093](https://github.com/hybridgroup/gobot/issues/1093),[#1092](https://github.com/hybridgroup/gobot/issues/1092))
* **go, deps:** switch to Go 1.20 and update modules ([#1067](https://github.com/hybridgroup/gobot/issues/1067))
* **linter:** update linter to v1.56.1 and fix issues ([#1068](https://github.com/hybridgroup/gobot/issues/1068))

### Doc

* update links to release or tagged branch ([#1069](https://github.com/hybridgroup/gobot/issues/1069))
* **core:** prepare for release v2.4.0

### Examples

* fix missing checks of return values ([#1060](https://github.com/hybridgroup/gobot/issues/1060))

### Gobot

* rename Master to Manager ([#1070](https://github.com/hybridgroup/gobot/issues/1070))

### Megapi

* use serialport adaptor and move driver to drivers/serial ([#1062](https://github.com/hybridgroup/gobot/issues/1062))

### Neurosky

* use serialport adaptor and move driver to drivers/serial ([#1061](https://github.com/hybridgroup/gobot/issues/1061))

### Test

* try to stabilize eventer tests ([#1066](https://github.com/hybridgroup/gobot/issues/1066))
* try to stabilize firmata tests ([#1097](https://github.com/hybridgroup/gobot/issues/1097))

## [v2.3.0](https://github.com/hybridgroup/gobot/compare/v2.2.0...v2.3.0) (2024-01-06)

### Adaptors

* **pwm:** introduce scale option for servo ([#1046](https://github.com/hybridgroup/gobot/issues/1046))
* **analogpins:** add a generic analog pin adaptor ([#1041](https://github.com/hybridgroup/gobot/issues/1041))

### Aio

* fix data race in AnalogSensorDriver ([#1024](https://github.com/hybridgroup/gobot/issues/1024))
* **all:** introduce functional options ([#1039](https://github.com/hybridgroup/gobot/issues/1039))
* **analog sensor:** fix deadlock in cyclic reading ([#1042](https://github.com/hybridgroup/gobot/issues/1042))
* **thermalzone:** add driver for read a thermalzone from system ([#1040](https://github.com/hybridgroup/gobot/issues/1040))

### Build

* **go, deps:** update modules ([#1047](https://github.com/hybridgroup/gobot/issues/1047), [#1052](https://github.com/hybridgroup/gobot/issues/1052))

### Doc

* **test:** use -race for tests by default ([#1035](https://github.com/hybridgroup/gobot/issues/1035))

### Gpio

* fix data race in StepperDriver ([#1029](https://github.com/hybridgroup/gobot/issues/1029))
* fix data race in PIRMotionDriver ([#1028](https://github.com/hybridgroup/gobot/issues/1028))
* fix data race in ButtonDriver ([#1027](https://github.com/hybridgroup/gobot/issues/1027))
* fix data race in EasyDriver ([#1025](https://github.com/hybridgroup/gobot/issues/1025))
* **all:** introduce functional options ([#1045](https://github.com/hybridgroup/gobot/issues/1045))

### I2c

* **core:** fix problems with usage of uintptr ([#1033](https://github.com/hybridgroup/gobot/issues/1033))

### Lint

* **all:** fix issues of errorlint etc ([#1037](https://github.com/hybridgroup/gobot/issues/1037))
* **all:** switch to 1.55.2 and adjust linter issues ([#1036](https://github.com/hybridgroup/gobot/issues/1036))

### Ollie

* **test:** fix data race in test ([#1034](https://github.com/hybridgroup/gobot/issues/1034))

### Raspi

* **pwm:** add support for sysfs and fix pi-blaster ([#1048](https://github.com/hybridgroup/gobot/issues/1048))

## [v2.2.0](https://github.com/hybridgroup/gobot/compare/v2.1.1...v2.2.0) (2023-10-29)

### Adaptors

* **PWM:** fix wrong duty cycle after kill program ([#994](https://github.com/hybridgroup/gobot/issues/994))

### Beaglebone

* **doc:** fix preceding typo ([#985](https://github.com/hybridgroup/gobot/issues/985))

### Build

* **deps:** module update ([#992](https://github.com/hybridgroup/gobot/issues/992))
* **go, deps:** switch to Go 1.19 and update modules ([#1008](https://github.com/hybridgroup/gobot/issues/1008))
* **style:** switch to gofumpt and add linters ([#1009](https://github.com/hybridgroup/gobot/issues/1009))

### Doc

* **roadmap:** remove file ROADMAP.md after creating issues ([#1005](https://github.com/hybridgroup/gobot/issues/1005))

### Dragonboard

* fix example and documentation ([#977](https://github.com/hybridgroup/gobot/issues/977))

### Gpio

* **hcsr04:** add driver for ultrasonic ranging module ([#1012](https://github.com/hybridgroup/gobot/issues/1012))

### I2c

* **PCA9685, adafruit, adafruit2327, adafruit2348:** clean up architecture and fix init ([#1021](https://github.com/hybridgroup/gobot/issues/1021))

### Jetson

* **PWM:** fix set period ([#1019](https://github.com/hybridgroup/gobot/issues/1019))

### Joystick

* **core:** replace sdl with 0xcafed00d/joystick  package ([#988](https://github.com/hybridgroup/gobot/issues/988))

### Sphero

* Add support for calibration

### System

* **gpio:** add edge polling function ([#1015](https://github.com/hybridgroup/gobot/issues/1015))

### Test

* **all:** substitude assert.Nil by assert.NoError if useful ([#1016](https://github.com/hybridgroup/gobot/issues/1016))
* **all:** substitude assert.Error by assert.ErrorContains ([#1014](https://github.com/hybridgroup/gobot/issues/1014), [#1011](https://github.com/hybridgroup/gobot/issues/1011))
* **all:** switch to test package stretchr testify ([#1006](https://github.com/hybridgroup/gobot/issues/1006))
* **gpio, aio:** cleanup helper_test ([#1018](https://github.com/hybridgroup/gobot/issues/1018))

## [v2.1.1](https://github.com/hybridgroup/gobot/compare/v2.1.0...v2.1.1) (2023-07-07)

### All

* upgrade modules
* substitute deprecated ioutil methods ([#923](https://github.com/hybridgroup/gobot/issues/923))
* **linter:** activate linters "errcheck", "ineffassign", "unused" and fix issues ([#950](https://github.com/hybridgroup/gobot/issues/950))
* **linter, format:** format code and list of linter todo's ([#962](https://github.com/hybridgroup/gobot/issues/962))
* **linter:** activate linter "makezero" and fix issue ([#965](https://github.com/hybridgroup/gobot/issues/965))

### Ble

* simplify and substitute errors.Wrap() ([#958](https://github.com/hybridgroup/gobot/issues/958))

### Build

* **go:** switch to Go 1.18 ([#940](https://github.com/hybridgroup/gobot/issues/940))

### Core

* CLI removed ([#946](https://github.com/hybridgroup/gobot/issues/946))

### Doc

* fix and improve documentation regarding to installation ([#946](https://github.com/hybridgroup/gobot/issues/946), [#970](https://github.com/hybridgroup/gobot/issues/970))

### Mavlink

* fix linter issues of errcheck ([#959](https://github.com/hybridgroup/gobot/issues/959))

### System

* **syscall:** switch to x/sys ([#963](https://github.com/hybridgroup/gobot/issues/963))

### Tello

* fix wifiMessage and lightMessage ([#957](https://github.com/hybridgroup/gobot/issues/957))
* fix partially [#793](https://github.com/hybridgroup/gobot/issues/793) by initialize doneCh in NewDriverWithIP ([#931](https://github.com/hybridgroup/gobot/issues/931))

## [v2.1.0](https://github.com/hybridgroup/gobot/compare/v2.0.3...v2.1.0) (2023-05-29)

### Build

* **v2:** revert of [#927](https://github.com/hybridgroup/gobot/pull/927), no usage of a v2 subfolder anymore (issue [#920](https://github.com/hybridgroup/gobot/issues/920))

## [v2.0.3](https://github.com/hybridgroup/gobot/compare/v2.0.2...v2.0.3) (2023-05-24)

* accidentally created release without any changes

## [v2.0.2](https://github.com/hybridgroup/gobot/compare/v2.0.1...v2.0.2) (2023-05-22)

### Build

* **v2:** fix usage by moving code to a v2 subfolder ([#927](https://github.com/hybridgroup/gobot/pull/927))

## [v2.0.1](https://github.com/hybridgroup/gobot/compare/v2.0.0...v2.0.1) (2023-05-21)

### Build

* **style:** add golangci-lint workflow configuration ([#916](https://github.com/hybridgroup/gobot/issues/916))
* **style:** fix linter findings of "gosimple", "govet" and "staticcheck" ([#917](https://github.com/hybridgroup/gobot/issues/917))

### Bump

* periph.io/x/conn/v3 from 3.6.10 to 3.7.0 ([#913](https://github.com/hybridgroup/gobot/issues/913))
* github.com/gofrs/uuid from 4.3.0+incompatible to 4.4.0+incompatible ([#914](https://github.com/hybridgroup/gobot/issues/914))
* golang.org/x/net from 0.1.0 to 0.10.0 ([#915](https://github.com/hybridgroup/gobot/issues/915))
* github.com/nats-io/nats-server/v2 from 2.1.0 to 2.7.4 ([#906](https://github.com/hybridgroup/gobot/issues/906))

### Core

* fix Semantic Import Versioning for v2 ([#921](https://github.com/hybridgroup/gobot/issues/921))

### Docs

* **core:** adjust changelog generation ([#924](https://github.com/hybridgroup/gobot/issues/924))

## [v2.0.0](https://github.com/hybridgroup/gobot/compare/v1.16.0...v2.0.0) (2023-05-15)

### ble

* update to TinyGo Bluetooth package v0.6.0 release

### build

* update appveyor for go 1.19
* switch to new cimg with golang 1.17
* new home path for cimg
* check examples in CI ([#884](https://github.com/hybridgroup/gobot/issues/884))
* add tests of more platforms to CI
* add configuration file for dependabot ([#907](https://github.com/hybridgroup/gobot/issues/907))
* add PR template

### core

* use base driver for all I2C devices
* rename package "sysfs" to "system"
* go.mod to 1.17 and all modules incl. code upgrades

## digispark

* add example for generic i2c.Driver
* fix i2c.ReadBlockData(), Read_Data() and some small other fixes

### dji tello

* Halt does not terminate all the related goroutines and may wait forever when it is called multiple times

### docs

* README for gpio, pwm, i2c and add example
* document fields for flight data

### aio

* analog sensor driver to prevent ReadValue() to get float64

### gopigo3

* fix examples and driver

### gpio

* add advanced digital pin options (pull, bias, drive, debounce, event)
* add support for new character device Kernel ABI for GPIO
* add read firmware version and DHT sensors for grovepi

### i2c

* add generic i2c driver
* fix I2C connection-bus caching and multiple device usage
* introduce I2cBusAdaptor for composition in platforms
* **Adafruit1109:** fix driver shows bad characters after Halt()
* **ads1x15:** fix driver not working stable when reading multiple inputs
* **ADXL345:** use ReadBlockData()
* **bmxy8z:** use ReadBlockData
* **BMP180, BMP280 BMP388 BME280:** use ReadBlockData()
* **CCS811:** use ReadBlockData()
* **HMC5883L:** fix I2C driver typo: change from HMC8553L
* **HMC5883L:** fix driver returns wrong values
* **L3GD20H:** fix full scale range usage
* **MPL115A2:** use ReadBlockData(), WriteByteData()
* **MPU6050:** fix wrong initialize and reduced temperature resolution
* **PCA9501:** add driver
* **PCA953x:** add driver
* **PCF8583:** add driver
* **TH02:** fix wrong register usage for read heater

### jetson nano

* add Jetson Nano adpator
* fix pwm feature

### joystick

* add Xbox-One controller
* add configuration for Nintendo Switch controllers ([#903](https://github.com/hybridgroup/gobot/issues/903))
* add Dualsense joystick (PlayStation 5) ([#880](https://github.com/hybridgroup/gobot/issues/880))

### nanopi neo

* add platform

### piblaster

* add unused but missing interface implementation

### radxa rock pi 4(c+)

* add platform ([#902](https://github.com/hybridgroup/gobot/issues/902))

### raspi

* fix  pwm cache
* fix Stopping and Starting Robot and LED Driver/LED does not toggle on restart

### spi

* fix spi.SpiConnection is not gobot.Connection: missing method Connect
* using GPIO's is now possible
* **MFRC522:** add driver

### test

* increase some timings to make tests a little less fragile
* skip test TestNatsAdaptorFailedConnect when flaky
* stabilize "every"-test
* stabilize flaky utils_test
* stabilize firmata tests
* fix tests with sysfs mocks, ReadBlockData, WriteBlockData
* fix keyboard tests and exclude opencv
* fix PWM related read/write tests
* add check for examples in Makefile

### tinkerboard

* fix new pwm behaviour

### BREAKING CANGES

* some interfaces moved, see folder system and adaptor.go

## [v1.16.0](https://github.com/hybridgroup/gobot/compare/v1.15.0...v1.16.0) (2022-05-02)

### bugfix

* failing leftovers after usage of PR #569
* Fix servo and DC motors presence
* FIX the bug #568 without further impact, heavy improvements of tests
* fixed PinMode, SetPullUp and SetPolarity, unit tests activated
* ReadGPIO fixed with #576, failing leftovers for PinMode, SetPullUp and SetPolarity
* helper_test ReadByteData, ReadWordData to use reg

### core

* update uuid package and directly access it; remove archived uuid package

### digispark

* fix ReadByte & WriteByte, rework and add i2c tests
* remove useless code in i2c test

### drivers

* add AnalogActuatorDriver, analog temperature sensor, driver for PCF8591 (with 400kbit stabilization), driver for YL-40
* Adding support for hmc5883l compass
* bmp388 fix missing address write byte in test of Measurements
* drv2605l fix missing address write byte in test of Halt()
* introduce adafruit1109 2x16 LCD with 5 keys
* mcp23017: add mutex for write, hd44780: fix mutexes
* MCP3004: correct number of channels

### raspi

* fix raspi PWMPin.SetDutyCycle (#800)

### tello

* Guards Dji Tello Halt against nil dereference

### test

* don't panic on 'With*' allow simpler wrapping of drivers

### tinkerboard

* fix tinkerboard i2c0 to i2c4, improve comments in pin map, improve README

## [v1.15.0](https://github.com/hybridgroup/gobot/compare/v1.14.0...v1.15.0) (2020-11-30)

### build

* Switch to CircleCI

### ble

* replace go-ble with tinygo bluetooth package, restore macOS functionality

### gpio

* Update RelayDriver to invert value written on Inverted
* Add tests for DigitalWrite value
* Add support for HD44780 LCD controller
* Add delay for Run function of StepperDriver

### spi

* fixes #700 * Avoid to close the connection.

### i2c

* add SHT2x device
* add BMP388 Barometric Pressure/Temperature/Altitude Sensor

### pwm

* Resolve issue with PWM for PWMWrite

### mqtt

* Add method to publish MQTT messages with retain flag

### tello

* Add graceful halt for Tello driver
* Add Tello EDU driver

### keyboard

* add symbol keys for platform/keyboard

### examples

* Update ffmpeg command to decrease latency in tello example

## [v1.14.0](https://github.com/hybridgroup/gobot/compare/v1.13.0...v1.14.0) (2019-10-15)

### core

* migrating from dep to go modules
* update codegangsta to urfave (#690)

### docs

* Fix a link in package docs' example code.

### examples

* fixed broken imports due to changed path causing go get to fail

### gpio

* Added ability to make a relay driver inverted (#674)

### opencv

* Update to GoCV 0.21.0

### spi

* Apa102 use default brightness (#671)

### tello

* Updated videoPort for DJI Tello to 11111

## [v1.13.0](https://github.com/hybridgroup/gobot/compare/v1.12.0...v1.13.0) (2019-05-22)

### api

* Initial stab at Robot-based work

### build

* correct package version as suggested by @dlisin thanks
* only build last 2 versions of Go plus tip for CI
* Update dep script for AppVeyor
* update deps to latest versions of dependencies for GoCV and others
* Update Gopkg and add test dep to Travis YML
* update OpenCV build script for OpenCV 4.1.0

### docs

* update to remove Gitter and replace with Slack, and update copyright dates

### example

* add missing nobuild header

### gpio

* Add SparkFun’s EasyDriver (and BigEasyDriver)
* Add unit tests for TH02 & Minor improvement
* Added rudiementary support for TH02  Grove Sensor
* pwm_pin * Fix DutyCycle() parse error, need to trim off trailing '\n' before calling strconv.Atoi(), as other functions in this package do
* Simplify code as suggested in #617

### grovepi

* add mutex to control transactionality of the device communication

### i2c

* add 128x32 and 96x16 sizes to the i2c ssd1306 driver
* build out the ccs811 driver
* update PCA9685 driver to use same protocol as Adafruit Python lib

### leapmotion

* Parser error in Pointable.Bases: Write test and fix
* Update gobot leap platform to support Leap Motion API v6

### mavlink

* fix mavlink README to use correct example code

### mqtt

* Add some new MQTT adaptor functions with QOS
* Allow setting QoS on MTT adaptor
* make tests run correctly even when a local MQTT server is in fact running
* Do not skip verification of root CA certificates by default InsecureSkipVerify

### nats

* Update Go NATS client library import

### opencv

* minor updates to opencv README
* update to OpenCV 4.1.0

### sphero

* Added methods to read Sphero Power States
* Added some new features to the sphero ollie, bb-8 and sprkplus

### spi

* correct param used for APA102 Draw() method
* Stop using Red parameter for brightness value

### tello

* add direct vector access
* add example with keyboard
* Change fps to 60
* Check for error immediately and skip publish if error occurred
* update FlightData struct

### up2

* add support for built-in LEDs
* correct i2c default bus information to match correct values
* finalize docs for UP2 config steps
* update README to include more complete setup information
* useful constant values to access the built-in LEDs

## [v1.12.0](https://github.com/hybridgroup/gobot/compare/1.11.1...v1.12.0) (2018-08-27)

### api

* further improvement of the modular API changes
* modify Start() for more modular initialization, and add StartRaw() for completely custom API implementations
* settled on StartWithoutDefaults() as the method to start API without default routes

### core

* add Rescale utility function for straight linear rescaling

### digispark

* add examples using digispark with i2c devices blinkm and mlp115a2
* Added i2c to digispark, but not working yet
* Added some tests for digispark i2c connector
* Digispark i2c fixes, added Test for checking available addresses
* remove test method that should not be in adaptor
* remove test that is expected to ofail, but passes when digispark board is actually connected

### docs

* add GrovePi to README
* adjust order of badges in README
* Fixing broken link

### examples

* add example that uses both the API and also a custom handler with MJPEG streaming from an attached camera
* small improvements to Tello examples
* update Tello examples for main thread friendly macOS/Windows, add Tello face tracker

### i2c

* add commands to JHD1313MDriver
* add commands to PCA9685Driver
* add missing methods so the GrovePi fully implements the Adaptor interface
* add ShowImage() function to ssd1306 driver based on @mikegleasonjr suggestion
* GrovePi digitalwrite implemented
* implemented DigitalRead, DigitalWrite, and AnalogRead for GrovePi
* improve godocs for PCA9685
* mention that GrovePi requires running firmware 1.3.0
* update GrovePi to v1.3.0 firmware
* work in progress on GrovePi plus driver

### joystick

* add config file for Magicsee R1 contributed by @carl-ranson
* add some additional test coverage for file-based config
* added error handling for config loading in joystick driver
* mention need to be running a Linux kernel v4.14+ for controller mappings to work as expected
* provide constant values for existing joystick configurations

### raspi

* export PiBlasterPeriod in Adaptor

### spi

* add ShowImage() function to ssd1306 driver based on @mikegleasonjr suggestion

### tello

* specify end of msgType position
* add handleResponse testing
* Add motion cessation commands to Tello
* handleResponse only needs an io.Reader
* handleResponse should not send commands
* rename reqConn to cmdConn
* reqConn is only an io.WriteCloser
* send Land() command to drone on Halt() to avoid floating mid-air

## [1.11.1](https://github.com/hybridgroup/gobot/compare/1.11.0...1.11.1) (2018-07-10)

### build

* exclude vendor and other previously excluded subpackages
* update Travis build to use OpenCV 3.4.2 release
* update deps for GoCV to v0.14.0 release
* Bump periph.io/x/periph to v3.0.0
* update to Go 1.10.3 and 1.9.7 for Travis builds

### docs

* Fix Leap Motion package link

### i2c

* fix write/read gpio on mcp23017, and cleaned up some comments
* correct pca9685 SetPWMFreq function scaling

### gopigo3

* update with default spi values, cleanup

## [1.11.0](https://github.com/hybridgroup/gobot/compare/1.10.2...1.11.0) (2018-05-31)

### build

* correct profile file location for codecov upload
* Make Go Lint happier by adding some explicit type conversions and ignoring unused error returns
* single quotes needed to upload any .cov file to codecov for reporting
* update deps to latest versions for Paho MQTT, go-sdl, and gocv
* upload any .cov file to codecov for reporting
* use go 1.10.2 and 1.9.6 for Travis builds
* add step to call dep ensure before contributing #524

### examples

* correct events used by XBox360 joystick example

### firmata

* Update the Firmata homepage in platform README

### gpio

* Improve Stepper Driver
* Initial support for MAX7219 (gpio) led driver

### joystick

* full corrected ds3 and ds4 mappings plus examples to match for latest sdl 2.0.8
* add instructions to README on how to install SDL on Linux from source
* add missing type conversion
* add new contributions to README
* Add T-Flight Hotas X flight controoller
* add xbox360 rock band drums controller
* Correct Dualshock4 controller mappings and add ps/left/right buttons
* correct test issue
* exclude scanner from test builds
* Fix joystick_driver to detect dpad input for xbox controllers
* Update dualshock4.json to match joystick_dualshock4.go
* update scanner to match go-sdl 0.3 API changes
* Update the joystick driver test to read DPAD properly

### leapmotion

* change timestamp to uint64 to fix #516

### tello

* slow/fast mode switch function
* StopLanding feature
* Add Bounce() and PalmLand() funcs and their associated events.
* bug fix
* Change several fields in FlightData struct from int16 to bool
* Export the FlightData fields (see Issue #531)

## [1.10.2](https://github.com/hybridgroup/gobot/compare/1.10.1...1.10.2) (2018-04-24)

### opencv

* update GoCV to latest version

## [1.10.1](https://github.com/hybridgroup/gobot/compare/1.10.0...1.10.1) (2018-04-24)

### tello

* improve support for DJI Tello drone, especially video

## [1.10.0](https://github.com/hybridgroup/gobot/compare/v1.9.0...1.10.0) (2018-04-20)

### docs

* add gitter badge to readme

### gpio

* AIP1640 led driver, used in Wemos D1 mini's matrix LED shield

### spi

* switch to using periph.io for SPI interfaces
* add support for ssd1306
* add optional params such as bus/chip to all current drivers
* complete refactoring to spi.Connection
* remove unneeded code as suggested by @maruel
* remove unneeded type and cleanup GoDocs

### ble

* correct spelling error in function name

### build

* update to latest version of Go 1.10 for Travis build

### cli

* remove extra newline

### docs

* add recently contributed GPIO devices to README

### joystick

* able to configure joysticks without external json file
* correct error in scanning script
* correct events used by gamepad-style up/down/left/right buttons
* correct scanner error from ID
* removed double release event

### tello

* add support for DJI Tello drone

## [v1.9.0](https://github.com/hybridgroup/gobot/compare/v1.8.0...v1.9.0) (2018-02-14)

### beaglebone

* update pin naming, docs, and examples for the latest Debian OS releases

### opencv

* update build settings needed to build OpenCV/GoCV as part of test suite
* deps for latest GoCV v0.9.0

### build

* update Travis build to use very latest Go versions

### docs

* add references to new drivers for ADXL345, BH1750, and TM1638.
* improve docs for installation and use of OpenCV/GoCV from Gobot
* update copyright date to 2018

### gpio

* Initial support for TM1638 modules

### i2c

* Added basic driver for BH1750 (light sensor), board GY-302
* support for accel ADXL345

### bb8/ollie/sprkplus

* add Boost command
* add Set Back LED Output command
* add Set Raw Motor Values command
* add Set Rotation Rate command
* add Set Stabilization command

### test

* Refactor TestAdaptorDigitalPinConcurrency test

## [v1.8.0](https://github.com/hybridgroup/gobot/compare/v1.7.1...v1.8.0) (2017-12-21)

### sysfs

* pause briefly to allow udev rules to apply when exporting PWMPin

### beaglebone

* correct uboot installation instructions
* add SPI support
* no more slots, add docs on configuring u-boot overlays
* handle gpio pinmux without relying on specific pre-existing setup

### pocketbeagle

* add support for PocketBeagle
* use universal io cape manager to initialize board setup
* improve docs for latest Debian OS

### build

* Add dep, change how tests run in CI
* update dependencies to latest GoCV version

### spi

* Add MCP3002, MCP3202, MCP3204, MCP3208, MCP3304, MCP3004, and MCP3008 A/D converter drivers
* adding initial support for APA102 LEDs, thanks to code sample from @rakyll
* extract shared SPI init code into spi package

### up2

* initial work on support for UP2 board

### gopigo3

* fixed set/get bug with motor dps

### gpio

* Adding stepper motor module

### firmata

* handle cases where out of sync data is read from serial port on first connecting

### i2c

* Change init payload sequence within jhd1313m1 driver Start() func.

## [v1.7.1](https://github.com/hybridgroup/gobot/compare/v1.7.0...v1.7.1) (2017-11-05)

### sprkplus

* add new platform for Sphero SPRK+

### firmata

* correct problem where last analog pin(s) were being ignored from capabilities query

### ble

* use go-ble/ble fork for BLE interactions

### build

* update to use latest OpenCV version
* update to use latest Golang versions

## [v1.7.0](https://github.com/hybridgroup/gobot/compare/v1.6.1...v1.7.0) (2017-10-23)

### curie

* Add Linux specific step to Intel Curie docs

### mqtt

* Added SetCleanSession

### build

* add go1.9 to versions tested in Travis CI
* add missing OpenCV lib dependency
* Update build to use latest Golang versions
* Travis build will now require sudo to install due to OpenCV

### docs

* some helpful edits for the initial spi implementation

### gopigo3

* integration of recent GoPiGo3 contributions
* Added grove support, and more gopigo3 examples

### gpio

* Add ButtonDriver.DefaultState to allow for 'reverse' buttons (ones that go from HIGH to LOW)

### holystone

* Add initial support for HS-200

### i2c

* SSD1306.WithDisplayHeight() and SSD1306.WithDisplayWidth() for SSD1306 that use different display ratios

### joystick

* add CLI utilty to scan display events to make it easier to add new joyticks
* update README to address #441

### opencv

* Switchover to use GoCV and OpenCV 3.3
* Switch to use custom domain for GoCV package
* all examples using new GoCV based code
* correct formatting in face detect example
* OpenCV face detector that is much more concurrent
* update interface and examples to indicate multipurpose

## [v1.6.1](https://github.com/hybridgroup/gobot/compare/v1.6.0...v1.6.1) (2017-07-15)

### core

* log failure errors on Robot Start()

### build

* run test coverage with covermode=set
* update build to use Golang 1.7.6 and 1.8.3

### docs

* work on ROADMAP doc

### sysfs

* increase test coverage

### bb8

* use updated ble adaptor interface for tests

### ble

* allow for characteristic writes both with and without a response
* allow override of specific HCI device to use
* eliminate race conditions from response handling

### curie

* Implement Accelerometer, Gyroscope, and Temperature sensors implemented
* motion detect implemented
* shock detect implemented
* step count implemented
* tap detect implemented

### digispark

* update blink example to display error message on Start()
* update README with latest development info

### edison

* auto-discovery of Edison board option
* removed commented lines

### firmata

* expose WriteSysex to external callers
* adjust client test timeout values
* cleanup error handling for connection code
* client tests don't need so many goroutines
* expose WriteSysex to external callers
* improve connection code to use a proper timeout
* increase test coverage
* make it possible to test external devices that use firmata adaptor
* refactoring firmata client
* remove circular import in test
* remove unused code, increase test coverage
* return connect errors to client
* switch to using go-serial package
* Sysex response events now being handled as expected

### bme280

* fix signed/unsigned bug
* Fixed incorrect error condition check when reading the 'ctrl_hum' register.
* Expanded the BME280 unit test for TestBME280DriverStart() to support reading from the 'ctrl_hum' register.
* Enables humidity readings in the BME280 driver by enforcing the write to the 'ctrl_meas' register, as per Section 5.4.3 of the BME280 data sheet

### chip

* Fixed PWM duty cycle calculation for C.H.I.P ServoWrite
* Fixed PWM init bug for C.H.I.P
* C.H.I.P PWM init robust for already enabled state

### i2c

* remove unused test code
* write config register in little endian

### joystick

* add needed constants for all PS3 buttons

### littlewire

* littlewire.cc links changed to littlewire.github.io

### mavlink

* switch to using go-serial package

### megapi

* switch to using go-serial package

### microbit

* use updated ble adaptor interface for tests

### minidrone

* add example for Parrot Mambo
* add support for Mambo external accessories
* increase test coverage
* never expect responses for characteristic writes
* remove unneeded code, increase test coverage
* separate flight status processing and add test coverage

### neurosky

* switch to using go-serial package

### ollie

* use updated ble adaptor interface for tests

### sphero

* switch to using go-serial package

### tinkerboard

* Updated Tinkerboard and sysfs tests to updated PWM polarity contract

## [v1.6.0](https://github.com/hybridgroup/gobot/compare/v1.5.0...v1.6.0) (2017-06-15)

### Bb8

* use updated ble adaptor interface for tests

### Ble

* eliminate race conditions from response handling
* allow for characteristic writes both with and without a response
* allow override of specific HCI device to use

### Build

* update build to use Golang 1.8.3
* update build to use Golang 1.7.6 and 1.8.2
* run test coverage with covermode=set

### Core

* log failure errors on Robot Start()

### Curie

* update docs formatting
* add Curie to main README platform list
* more improvements for README setup instructions
* improve README info
* improve tests and examples
* increase test coverage
* increase test coverage
* motion detect implemented
* tap detect implemented
* step count implemented
* shock detect implemented
* Accelerometer, Gyroscope, and Temperature sensors implemented
* WIP on adding support for Intel Curie IMU

### Digispark

* update blink example to display error message on Start()
* update README with latest development info

### Docs

* add more wishlist to ROADMAP
* add helpful information to examples themselves
* correct installation instructions to match latest versions
* more wishlish items for roadmap
* update BLE connect info to latest
* more work on ROADMAP doc
* add first attempt at roadmap document

### Edison

* refactor auto-discovery of Edison board option
* removed commented lines

### Firmata

* remove circular import in test
* make it possible to test external devices that use firmata adaptor
* Sysex response events now being handled as expected
* expose WriteSysex to external callers
* expose WriteSysex to external callers
* cleanup error handling for connection code
* improve connection code to use a proper timeout
* remove unused code, increase test coverage
* increase test coverage
* switch to using go-serial package
* return connect errors to client
* client tests don't need so many goroutines
* adjust client test timeout values
* refactoring firmata client

### Fix

* signed/unsigned bug

### Fixed

* incorrect error condition check when reading the 'ctrl_hum' register. Expanded the BME280 unit test for TestBME280DriverStart() to support reading from the 'ctrl_hum' register.
* PWM duty cycle calculation for C.H.I.P ServoWrite
* PWM init bug for C.H.I.P

### I2c

* remove unused test code

### Improved

* BME280 humidity initialisation so that it does not override existing oversampling rates that have been set up for the ctrl_meas register.

### Issue

* [#424](https://github.com/hybridgroup/gobot/issues/424): Enables humidity readings in the BME280 driver by enforcing the write to the 'ctrl_meas' register, as per Section 5.4.3 of the BME280 data sheet

### Joystick

* add needed constants for all PS3 buttons

### Made

* C.H.I.P PWM init robust for already enabled state

### Mavlink

* switch to using go-serial package

### Megapi

* switch to using go-serial package

### Microbit

* use updated ble adaptor interface for tests

### Minidrone

* never expect responses for characteristic writes
* add example for Parrot Mambo
* separate flight status processing and add test coverage
* add support for Mambo external accessories
* remove unneeded code, increase test coverage
* increase test coverage

### Neurosky

* switch to using go-serial package

### Ollie

* use updated ble adaptor interface for tests

### Prepare

* for v1.6.0 release

### Release

* correct changelog incorrect section titles

### Sphero

* switch to using go-serial package

### Sysfs

* increase test coverage

### Updated

* Tinkerboard and sysfs tests to updated PWM polarity contract

### Write

* config register in little endian

## [v1.5.0](https://github.com/hybridgroup/gobot/compare/v1.4.0...v1.5.0) (2017-05-10)

### core

* Add Running() methods for Manager and Robot and increase test coverage accordingly

### sysfs

* define DigitalPinnerProvider and PWMPinnerProvider interfaces
* add Chip to be able to change pwmchip, and some related refactoring
* add file read/write testing for failure conditions
* proper handling of busy state vs. other errors
* return sensible result when no valid data read

### test

* increase coverage on test helpers

### build

* switching to Travis builds using Ubuntu 14.04 Trusty

### aio

* only need to support AnalogReader interface
* avoid test race conditions
* ensure that AnalogSensor event Data is always int

### gpio

* only need to support DigitalReader/DigitalWriter interface

### i2c

* Added support for the ADS1015 and ADS1115 ADCs
* Add INA3221 Voltage Monitor
* Ensure lock of i2c bus for each individual operation
* Small refactoring and increase test coverage for BMP180

### beaglebone

* implement DigitalPinner and PWMPinner interfaces
* protect against pin map races
* increase test coverage

### chip

* add preliminary support for C.H.I.P. Pro
* add back ServoWrite implementation
* implement DigitalPinnerProvider and PWMPinnerProvider interfaces
* protect against pin map races

### dragonboard

* export DigitalPin and PWMPin adaptor methods
* protect against pin map races
* increase test coverage

### edison

* auto-detect arduino breakout board, if no specific board is expected
* ensure that we initialize tristate if arduino breakout board
* export DigitalPin and PWMPin adaptor methods
* implement DigitalPinnerProvider and PWMPinnerProvider interfaces
* protect against pin map races
* refactoring to reduce code duplication

### firmata

* remove processing that might have been eating test events, increase test coverage

### joule

* implement DigitalPinnerProvider and PWMPinnerProvider interfaces
* protect against pin map races
* remove incorrect pin assignment and improve test coverage
* add examples using Joule with ADS1015 ADC
* naming system changes
* correct pin mappings and add PWM example

### mavlink

* add a Mavlink-over-UDP adaptor.

### microbit

* Add DigitalWriter, DigitalReader, and AnalogReader support using IOPinDriver
* Handle start error and increase test coverage

### mqtt

* Add a (topic, payload) event type
* change the On handler to take mqtt.Message
* increase test coverage
* update examples that use mqtt for updated notification signature

### nats

* change the On() handler to take the subject as an argument
* increase test coverage

### raspi

* implement DigitalPinnerProvider and PWMPinnerProvider interfaces
* add implementation for PWMPinner interface that wraps pi blaster
* fix adaptor race conditions
* increase test coverage

### tinkerboard

* Add support for ASUS Tinker Board

## [v1.4.0](https://github.com/hybridgroup/gobot/compare/v1.3.0...v1.4.0) (2017-04-12)

### core

* Use 10-buffered chans for events, see #374

### i2c

* Many refactors and increases in test coverage
* Eliminate race conditions introduced by tests
* Adds Altitude() function to BMP280/BME280
* bme280 driver Humidity compensation formula
* ssd1306 driver implementation

### aio

* Eliminate race conditions introduced by tests

### gpio

* Fix motor mode change when speed is set
* Eliminate race conditions introduced by tests
* Reduce test side effects

### ardrone

* Increase test coverage

### audio

* Increase test coverage

### bb8

* Refactoring to use BLEConnector interface and provide tests

### bebop

* Increase test coverage

### beaglebone

* Increase test coverage

### ble

* Increase test coverage for battery, device information, and generic access drivers
* Refactoring drivers to use BLEConnector interface and provide tests

### chip

* Added PWM0 support
* Increase test coverage

### digispark

* Increase test coverage

### dragonboard

* Increase test coverage

### edison

* Remove pointless error checking code
* Refactor digital pin creation process method
* Increase test coverage

### firmata

* Eliminate race conditions introduced by tests
* Increase test coverage for i2c commands

### joule

* Increase test coverage

### joystick

* Increase test coverage

### keyboard

* Increase test coverage

### mavlink

* Eliminate race conditions introduced by tests
* Increase test coverage

### mavlink

* Increase test coverage

### microbit

* Refactoring to use BLEConnector interface and provide tests
* Address #404 by adding info about required magnetometer calibration step to README
* Increase test coverage

### minidrone

* Refactoring to use BLEConnector interface and provide tests

### mqtt

* Increase test coverage

### nats

* Increase test coverage

### neurosky

* Update neurosky README & example
* Eliminate race conditions introduced by tests
* Increase test coverage

### ollie

* Refactoring to use BLEConnector interface and provide tests
* Correct race condition error on seq
* Increase test coverage

### opencv

* Increase test coverage

### particle

* Increase test coverage

### raspi

* Address #391 by providing more details about normal development workflow
* Increase test coverage

### sphero

* Eliminate race conditions
* Increase test coverage

### sysfs

* Address race condition from udev rules when exporting GPIO pins
* Increase test coverage

### docs

* Improve explanations for scp/ssh workflow on SoC boards
* Include entire Apache 2.0 license in the license text

### test

* Add crude travis check for gofmt; format all sources
* Significantly speed up travis and make runs
* Remove test code no longer being called
* Update Travis to run tests using Golang 1.8.1
* Increase gobottest test coverage

## [v1.3.0](https://github.com/hybridgroup/gobot/compare/v1.2.1...v1.3.0) (2017-03-22)

### microbit

* Add new platform support

### dragonboard

* Add new platform support

### gpio

* Increase test coverage

### i2c

* Update list of supported i2c devices
* Minor adjustments and test coverage improvements
* Added more capabilities checks for I2C
* Removed smbus block operations

### core

* Increase test coverage

### test

* Improvements to run tests much faster thanks @maruel
* Use codecov.io for code coverage reporting

### docs

* Update CoC based on Contributor Covenant

## [v1.2.1](https://github.com/hybridgroup/gobot/compare/v1.2.0...v1.2.1) (2017-02-16)

### Allow

* NATS options to pass in the NATS adaptor for TLS support.

### Chip

* correct docs to describe valid pin mappings

### Update

* version to 1.2.1 for point release

## [v1.2.0](https://github.com/hybridgroup/gobot/compare/v1.1.0...v1.2.0) (2017-02-16)

### core

* Use new improved default namer to avoid API conflicts

### gpio

* Removed scaling function from servo driver
* Correct servo driver to pass along angle to adaptor to sort out implementation

### i2c

* Refactored platforms and drivers to new I2C interfaces
* Change to make I2C support more than one bus
* Refactor drivers to support new optional params

### bb8

* Added collision detection support and example

### beaglebone

* Correct i2c buses to match actual mapping

### ble

* Switch to using [ble](https://github.com/currantlabs/ble) package for Bluetooth LE
* Basic serial over BLE working with Arduino101 with StandardFirmataBLE
* WIP on multiple simultaneous ble devices

### chip

* Fixed chip XIO base address lookup

### digispark

* Fix #288 by using pkg-config to locate libusb-compat includes

### firmata

* Remove race conditions identified in Firmata client
* Correct error in I2C reads not listening to board events

### mqtt

* Add driver for syntactical sugar around virtual devices
* Add SSL/TLS client options support
* Fix #277 by adding SetAutoReconnect method to set Paho MQTT client
* Change both 'On' and 'Publish' method function signatures to match Eventer interface

### nats

* Add driver to make it easier to create virtual devices

### ollie

* Added collision detection support and example

### parrot

* Add ValidatePitch helper function for Parrot Minidrone, Parrot Bebop & ARDrone 2.0 to package

### docs

* Fix #363 by using atomic.Value to protect current values used by multiple goroutines in drone examples

### test

* Remove Golang 1.5 from TravisCI tests in prep for Golang 1.8 release

## [v1.1.0](https://github.com/hybridgroup/gobot/compare/v1.0.0...v1.1.0) (2017-01-09)

### core

* use canonical import path for sysfs package

### i2c

* Add a driver for the SHT3X chip
* Add a driver for BMP180
* Add support for L3GD20H gyroscope

### firmata

* Add support for TCPFirmata connections, allowing ESP8266 and other WiFi-connected controllers
* Add mention to README to use 'tty.' serial port on OSX
* Add mention of A4 and A5 normally unavailable on Firmata

### raspi

* Correct README build instructions with missing 'go build' command

### snapcraft

* Add the packaging metadata to build the gobot snap for Ubuntu Snappy

### particle

* Update examples to take key params via command line
* Address #160 by adding support for tinker-servo sketch if installed on Particle device

### esp8266 add experimental ESP8266 support to list of supported platforms

### sysfs

* Should fix #272 by using first byte of data as command register for I2C reads
* Some additional cleanup suggested by golint

### ble

* Add generic access service driver
* Update docs to include reference to included drivers
* Move various test code to test file

### ollie

* Refactoring so no need to expose internal implementation details

### bebop

* Add support/example of RTP video
* Enable video on firmware 3.3+
* Update ps3 and video example to enable the video stream
* Update README for brief explanation of how to get drone video
* Corrected import paths for client examples

### bb8

* Correct NewDriver params and set name
* Add missing constructor to wrap Ollie implementation

### minidrone

* Update README with example and which specific models are currently supported
* Add all piloting flying state events
* Adds Emergency() and TakePicture() commands

### test

* Add Golang 1.8beta2 to Travis builds
* Correct aio references for AnalogRead tests

## [v1.0.0](https://github.com/hybridgroup/gobot/compare/v0.13.0...v1.0.0) (2016-12-21)

### core

* Refactoring to allow 'Metal' development using Gobot packages
* Able to run robots without being part of a Manager.
* Now running all work in separate goroutines
* Rename internal name of Manager type
* Refactor events to use channels all the way down.
* Eliminate potential race conditions from Events and Every functions
* Add Unsubscribe() to Eventer, now Once() works as expected
* DeleteEvent function added to Eventer interface
* Ranges over event channels instead of using select
* No longer return non-standard slices of errors, instead use hashicorp/go-multierror
* Ensure that all drivers have default names
* Now both Robot and Manager operate using AutoRun as expected
* Use canonical import domain of gobot.io for all code
* Use time.Sleep unless waiting for a timeout in a select
* Uses time.NewTimer() instead of time.After() to be more efficient

### test

* Add deps tasks to Makefile
* Add golang 1.7 to Travis CI tests
* Add golang 1.8beta1 to build matrix for Travis
* Reduce Travis builds to golang 1.4+ since it is late 2016 already
* Complete move of test interfaces into the test files where they belong
* Adds Parrot Minidrone and Sphero Ollie to Travis tests

### Add missing godocs for everything

### i2c

* Move I2C drivers into appropriately named 'drivers/i2c' directory
* Add support for Adafruit Servo/PWM HAT

### gpio

* Move GPIO drivers into appropriately named 'drivers/gpio' directory
* Add support for PIR motion detector

### beaglebone

* auto-detect Linux kernel version
* map usr LEDs to match all kernels

### ble

* Rename drivers to make them more obvious
* Add test placeholders

### chip

* Auto-detect OS version to adjust pin mappings
* Correct base for new 4.4 GPIO

### edison

* Support for other breakout boards besides Arduino

### firmata

* Use io.ReadFull in platforms/firmata/client
* Update tarm/goserial to tarm/serial

### joule

* Add support for Intel Joule

### megapi

* Adding support for MakeBlock megapi

### nats

* Add support for NATS server

### particle

* Complete renaming Spark platform to Particle


### parrot
* Move Parrot Minidrone into own platform
* Move both ARDrone and Bebop under Parrot package

### raspi

* Add missing godocs and small refactors for platform

### sphero

* Add initial support for Sphero BB-8 platform
* Move Sphero Ollie into own platform

## [v0.13.0](https://github.com/hybridgroup/gobot/compare/v0.12.1...v0.13.0) (2016-10-10)

### Add

* PinMode test case
* PinMode func for MCP23017
* example for Edison blink demo without gobot initialization.
* ServoConfig to the FirmataAdaptor
* unit tests for ServoConfig

### Adding

* support for MakeBlock megapi
* tests for the Adafruit driver, and corresponding minor driver changes.
* support for MakeBlock megapi
* a Servo example program for the Adafruit Servo Hat driver code.
* support for the Adafruit Servo/PWM HAT.  This required a slight refactor to the existing Motor HAT code to support multiple I2C addresses.
* two API funcs for the Adafruit driver with respect to the DC Motor, fleshing out the raspi_adafruit example with a runner function.
* the initial NATS platform support

### Another

* attempt at correct Travis syntax for gnatsd -[#5](https://github.com/hybridgroup/gobot/issues/5)
* attempt at correct Travis syntax for gnatsd -[#4](https://github.com/hybridgroup/gobot/issues/4)
* attempt at correct Travis syntax for gnatsd -[#3](https://github.com/hybridgroup/gobot/issues/3)
* attempt at correct Travis syntax for gnatsd -[#2](https://github.com/hybridgroup/gobot/issues/2)
* attempt at correct Travis syntax for gnatsd

### Ble

* fix unused var
* populate descriptors after descovering characterisitcs

### Bug

* fix in the Adafruit stepper code, specifically with respect to the AdafruitDouble step-style.

### Code

* cleanups suggested by gosimple

### Core

* update README with an example of 'Metal' Gobot
* should correct occasional test errors due to event overlap with test
* correct behavior in Mavlink driver, correct tests to match
* Add Unsubscribe() to eventer, now Once() works as expected
* Add further tests for Eventer
* cleanup comments on Eventer interface
* function DeleteEvent added to Eventer interface
* Refactor tests to allow 'metal' development using Gobot adaptors/drivers.
* Refactor tests to allow 'metal' development using Gobot adaptors/drivers.
* Refactor tests to allow 'metal' development using Gobot adaptors/drivers.
* Refactor examples to allow 'metal' development using Gobot adaptors/drivers.
* Refactoring to allow 'metal' development using Gobot adaptors/drivers.
* Continue refactoring to allow 'metal' development using Gobot libs.
* Refactor events to use channels all the way down. Allows 'metal' development using Gobot libs.
* update README with an example of 'Metal' Gobot
* should correct occasional test errors due to event overlap with test
* correct behavior in Mavlink driver, correct tests to match
* Add Unsubscribe() to eventer, now Once() works as expected
* Add further tests for Eventer
* cleanup comments on Eventer interface
* function DeleteEvent added to Eventer interface
* Refactor tests to allow 'metal' development using Gobot adaptors/drivers.
* Refactor tests to allow 'metal' development using Gobot adaptors/drivers.
* Refactor tests to allow 'metal' development using Gobot adaptors/drivers.
* Refactor examples to allow 'metal' development using Gobot adaptors/drivers.
* Refactoring to allow 'metal' development using Gobot adaptors/drivers.
* Continue refactoring to allow 'metal' development using Gobot libs.
* Refactor events to use channels all the way down. Allows 'metal' development using Gobot libs.

### Docs

* go fmt files that needed it from recent changes
* go fmt files that needed it from recent changes

### File

* rename, adding a test file for the Adafruit driver, and slight func naming changes.
* rename, adding a test file for the Adafruit driver, and slight func naming changes.

### Fix

* a typo and update the doc comment for FirmataAdaptor.ServoConfig
* the ServoConfig byte order
* issues flagged by 'go vet'
* misspellings

### Fixing

* some code and finally have Travis building
* tests, adding a few more, adding nats server to Travis CI for testing

### Initial

* significant changes to the Adafruit Motor HAT driver to support Stepper Motors.
* commit of driver code, with accompanying example, for the Adafruit_MotorHat.

### Joule

* add i2c example and notes to README about pullup resistors
* adds pin mappings from the second header
* add pin mapping info to README
* go fmt the multi-LED example

### Merge

* branch 'dev' of github.com:jfinken/gobot into dev
* branch 'dev' of github.com:jfinken/gobot into dev

### Misc

* update all LICENSE files for current year

### More

* explicit initialization in Start, slight refactor, and separate DC Motor and Stepper Motor examples.

### Move

* interface assertions to test files.

### Release

* update to version 0.13.0

### Remove

* debug message from i2c_device.go

### Removing

* the raspi_adafruit program as it has been split into three separate programs, removing my Makefile for the raspi adafruit programs, and fixing up a few comments.
* my fork from the Adafruit tests.

### Starting

* support for Intel Joule with the built-in LEDs and more

### Test

* add golang 1.7 to Travis CI tests
* add golang 1.7 to Travis CI tests

### Tests

* complete move of test interfaces into the test files where they belong
* refactor test interfaces out of implementations and into the tests where they belong

### Update

* READMEs with up to date info for Edison/Joule

### Updating

* platform support info


## [v0.12.1](https://github.com/hybridgroup/gobot/compare/v0.12.0...v0.12.1) (2016-07-13)

### A

* little more WIP, can open a connection to a specific peripheral

### Add

* MQTT authentication support
* MQTT authentication support
* Go Reportcard badge for fun
* Go Reportcard badge for fun

### Added

* example of how to use temp36 temperature sensor with firmata
* example of how to use temp36 temperature sensor with firmata

### Adds

* support for Dualshock4 wireless gamepad
* support for Dualshock4 wireless gamepad

### Allow

* failures in Travis builds for Golang 1.3 due to SDL changes

### Almost

* reading battery info

### BLE

* seems to require Golang 1.4+

### Can

* see BLE devices, and connect to a specific one

### Change

* default value for PCMD flag to match the Bebop 2.0.57+ expectations
* default value for PCMD flag to match the Bebop 2.0.57+ expectations
* test delay to 50ms
* test delay to 50ms

### Code

* cleanup, improve go report card
* cleanup, improve go report card

### Fix

* specs
* specs
* for analog (quick changes lag)
* for analog (quick changes lag)
* [#201](https://github.com/hybridgroup/gobot/issues/201) by add 'make examples' command to Makefile
* mavlink link typo

### Fixes

* failing test
* failing test

### Go

* fmt the code

### Increase

* hover time and remove cruft from simple Bebop drone example
* hover time and remove cruft from simple Bebop drone example

### Introduce

* `gobottest` package with test helpers
* `gobottest` package with test helpers

### Make

* dev branch target more explicit

### Making

* sure tests pass

### Merge

* branch 'feature/audio' into dev
* branch 'bugfix/gpio-button-tests' into dev
* branch 'feature/ble' into feature/ble-wip

### More

* WIP on reading characteristics

### Pin

* 229 value left out of test fixture on edison

### Refactor

* to use `gobottest` test helpers
* to use `gobottest` test helpers

### Remove

* fmt no longer used here
* commented lines
* test code

### Resolve

* merge conflicts
* merge conflict in Travis CI file

### Simple

* implementation that can read data

### Support

* gpio pin turn on and off
* gpio pin turn on and off

### Switching

* to currantlabs fork of gatt, and some related refactoring

### Test

* generated error messages as well
* generated error messages as well

### Tests

* also need to be pointed to [@veandco](https://github.com/veandco) go-sdl2 fork

### Update

* to 0.12.1
* missing changelog entries
* missing changelog entries
* ARDrone face tracking example to use main go-opencv fork

### Use

* main go-sdl fork from [@veandco](https://github.com/veandco) to pickup any upstream changes
* OpenCV 2.4, as well as switch to main fork of go-opencv
* Seek to speed up read/write in sysfs

### WIP

* on BLE

## [v0.12.0](https://github.com/hybridgroup/gobot/compare/v.0.11.1...v0.12.0) (2016-07-13)

### Refactor Gobot test helpers into separate package

### Improve Gobot.Every method to return channel, allowing it to be halted

### Refactor of sysfs adds substantial speed improvements

### ble

* Experimental support for Bluetooth LE.
* Initial support for Battery & Device Information services
* Initial support for Sphero BLE robots such as Ollie
* Initial support for Parrot Minidrone

### audio

* Add new platform for Audio playback

### gpio

* Support added for new GPIO device:
* RGB LED
* Bugfixes:
* Correct analog to better handle quick changes
* Correct handling of errors and buffering for Wiichuk

### mqtt

* Add support for MQTT authentication

### opencv

* Switching to use main fork of OpenCV
* Some minor bugfixes related to face tracking

## [v.0.11.1](https://github.com/hybridgroup/gobot/compare/v0.11.0...v.0.11.1) (2016-02-17)

### Add

* support for 'hand' and 'gesture' Leap Motion events
* MMA7660 accelerometer example for C.H.I.P.
* C.H.I.P. to supported platforms
* support for the CHIP platform
* MCP23017 write and read functionality to GPIO

### Adds

* MCP23017 i2c device to README
* additional examples for C.H.I.P.

### Better

* I2C device descriptions in README

### Correct

* the release command sent to pi-blaster.
* Intel Edison docs location thanks to [@seanmarcia](https://github.com/seanmarcia)

### Default

* the new MQTT 'AutoReconnect' to false

### Failure

* is no longer an option for Go 1.6

### Fix

* [#236](https://github.com/hybridgroup/gobot/issues/236) & fix [#239](https://github.com/hybridgroup/gobot/issues/239) by correcting initialization and temperature conversion for MPU-6050

### Fixed

* event race condition

### Get

* I2C functionality before doing SMBus block I/O

### Golang

* 1.3.3 still works, adding back to build

### Increase

* button delay hack for test suite
* test delay hack for button tests

### Name

* C.H.I.P. pins according to printed names

### Need

* to explicitly set content type to text/html for Robeaux main page

### No

* coveralls repo token for provate repos?

### Remove

* coveralls badge

### Run

* builds against the latest major releases

### The

* take-off-before-event-handling bug again

### Trying

* to remove coveralls based code coverage
* conditional build before_install
* conditional build

### Update

* version to v.0.11.1
* version to 0.11
* MQTT README for latest info
* targeted golang versions to include 1.6, and to begin deprecating 1.3.3 and earlier
* coveralls badge in README
* API example

### Use

* newer naming system for C.H.I.P. pins

### What

* about -v

### Why

* do this twice?

## [v0.11.0](https://github.com/hybridgroup/gobot/compare/0.10.0...v0.11.0) (2016-02-17)

### Support for Golang 1.6

### Determine I2C adaptor capabilities dynamically to avoid use of block I/O when unavailable

### chip

* Add support for GPIO & I2C interfaces on C.H.I.P. $9 computer

### leap motion

* Add support additional "hand" and "gesture" events

### mqtt

* Support latest update to Eclipse Paho MQTT client library

### raspberry pi

* Proper release of Pi Blaster for PWM pins

### bebop

* Prevent event race conditions on takeoff/landing

### i2c

* Support added for new i2c device:
* MCP23017 Port Expander
* Bugfixes:
* Correct init and data parsing for MPU-6050
* Correct handling of errors and buffering for Wiichuk

## [0.10.0](https://github.com/hybridgroup/gobot/compare/0.8.2...0.10.0) (2015-10-27)

### Refactor core to cleanup robot initialization and shutdown

### Remove unnecessary goroutines spawned by NewEvent

### api

* Update Robeaux to v0.5.0

### bebop

* Add support for the Parrot Bebop drone

### keyboard

* Add support for keyboard control

### gpio

* Support added for 10 new Grove GPIO devices:
* Grove Touch Sensor
* Grove Sound Sensor
* Grove Button
* Grove Buzzer
* Grove Led
* Grove Light Sensor
* Grove Vibration Sensor
* Grove Rotary
* Grove Relay
* Grove Temperature Sensor

### i2c

* Support added for 2 new Grove i2c devices:
* Grove Accelerometer
* Grove LCD with RGB backlit display

### docs

* Many useful fixes and updates for docs, mostly contributed by our wonderful community.

## [0.8.2](https://github.com/hybridgroup/gobot/compare/0.8.1...0.8.2) (2015-06-30)

### firmata

* Refactor firmata adaptor and split firmata protocol implementation into sub `client` package

### gpio

* Add support for LIDAR-Lite

### raspi

* Add PWM support via pi-blaster

### sphero

* Add `ConfigureLocator`, `ReadLocator` and `SetRotationRate`

## [0.8.1](https://github.com/hybridgroup/gobot/compare/0.8...0.8.1) (2014-12-28)

### spark

* Add support for spark Events, Functions and Variables

### sphero

* Add `SetDataStreaming` and `ConfigureCollisionDetection` methods

## [0.8](https://github.com/hybridgroup/gobot/compare/0.7.1...0.8) (2014-12-24)

### Refactor core, gpio, and i2c interfaces

### Correctly pass errors throughout packages and remove all panics

### Numerous bug fixes and performance improvements

### api

* Update robeaux to v0.3.0

### firmata

* Add optional io.ReadWriteCloser parameter to FirmataAdaptor
* Fix `thread exhaustion` error

### cli

* generator

*  Update generator for new adaptor and driver interfaces

*  Add driver, adaptor and project generators

*  Add optional package name parameter

## [0.7.1](https://github.com/hybridgroup/gobot/compare/0.7...0.7.1) (2014-11-17)

### opencv

* Fix pthread_create issue on Mac OS

## [0.7](https://github.com/hybridgroup/gobot/compare/0.6.3...0.7) (2014-11-10)

### Dramatically increased test coverage and documentation

### api

* Conform to the [cppp.io](https://gobot.io/x/cppp-io) spec
* Add support for basic middleware
* Add support for custom routes
* Add SSE support

### ardrone

* Add optional parameter to specify the drones network address

### core

* Add `Once(e *Event, f func(s interface{})` Event function
* Rename `Expect` to `Assert` and add `Refute` test helper function

### i2c

* Add support for MPL115A2
* Add support for MPU6050

### mavlink

* Add support for `common` mavlink messages

### mqtt

* Add support for mqtt

### raspi

* Add support for the Raspberry Pi

### sphero

* Enable stop on sphero disconnect
* Add `Collision` data struct

### sysfs

* Add generic linux filesystem gpio implementation

## [0.6.3](https://github.com/hybridgroup/gobot/compare/0.6.2...0.6.3) (2014-09-24)

### Add support for the Intel Edison

## [0.6.2](https://github.com/hybridgroup/gobot/compare/0.6.1...0.6.2) (2014-07-28)

### cli

* Fix typo in generator

### leap

* Fix incorrect Port reference
* Fix incorrect Event name

### neurosky

* Fix incorrect Event names

### sphero

* Correctly format output of GetRGB

## [0.6.1](https://github.com/hybridgroup/gobot/compare/0.6...0.6.1) (2014-07-12)

### cli

* Fix template error in generator

## [0.6](https://github.com/hybridgroup/gobot/compare/0.5.2...0.6) (2014-07-11)

### api

* Add robeaux support

### core

* Refactor `Connection` and `Device`
* Connections are now a collection of Adaptors
* Devices are now a collection of Drivers
* Add `Event(string)` function instead of `Events[string]` for retrieving Driver event
* Add `AddEvent(string)` function to register an event on a Driver

### firmata

* Fix slice bounds out of range error

### sphero

* Fix issue where the driver would not halt correctly on OSX

## [0.5.2](https://github.com/hybridgroup/gobot/compare/0.5.1...0.5.2) (2014-06-30)

### beaglebone

* Add `DirectPinDriver`
* Ensure slots are properly loaded

## [0.5.1](https://github.com/hybridgroup/gobot/compare/0.5...0.5.1) (2014-06-28)

### core

* Add `Version()` function for Gobot version retrieval

### firmata

* Fix issue with reading analog inputs
* Add `data` event for `AnalogSensorDriver`

## [0.5](https://github.com/hybridgroup/gobot/compare/0.4...0.5) (2014-06-17)

### Idomatic clean up

* Removed reflections throughout packages
* All officially supported platforms are now in ./platforms
* API is now a new package ./api
* All platforms examples are in ./examples
* Replaced martini with net/http
* Replaced ginkgo/gomega with system testing package
* Refactor gobot/robot/device commands
* Added Event type
* Replaced Manager type with Gobot type
* Every` and `After` now accept `time.Duration`
* Removed reflection helper methods

## [0.4](https://github.com/hybridgroup/gobot/compare/0.3...0.4) (2014-06-12)

### API

* commands now return an array of results

### Add

* cors support
* basic auth support to api
* Joystick & Neurosky platforms to README
* utils tests
* coveralls

### Allow

* user to set Host and Port when starting up Api

### Change

* README image source to gobot-site repo
* startApi to private function

### Display

* warning when using API without SSL

### Fixed

* the logo link

### Format

* device and connection type

### Green

* tests

### More

* api test coverage
* tests

### Refactor

* tests

### Remove

* Travis build from IRC
* ConnectToSerial
* ConnectToTcp util
* Reconnect and Disconnect from AdaptorInterface

### Robot

* is now a pointer

### SSL

* support in Api

### Update

* README.md
* README for new API security features
* api robeaux api compatibility
* .travis.yml
* generated driver
* examples
* coveralls badge
* platforms and drivers

### Use

* go-martini/martini

### WIP

* for API host/port params
* api tests


## [0.3](https://github.com/hybridgroup/gobot/compare/0.2...0.3) (2014-04-07)

### Add

* Godeps file
* IRC notifications to Travis builds
* tests for generated projects
* Init function to DriverInterface
* Halt function to DriverInterface
* more GPIO devices to README
* scale functions

### All

* updates for new gonuts/commander api

### Fix

* typo in generator

### Merge

* branch 'master' of github.com:hybridgroup/gobot

### Update

* generator

## [0.2](https://github.com/hybridgroup/gobot/compare/0.1...0.2) (2014-02-04)

### Add

* robeaux submodule
* Finalize on SIGINT
* Publish function for driver events
* device test coverage
* manager and robot test coverage

### Clean

* up tests

### Do

* not run tests on gobot.io branch

### JSON

* compatibility with cylon and artoo

### More

* test coverage

### Refactor

* robot and manager

### Remove

* robeaux submodule

### Update

* README.md
* examples

### Use

* golang log

### WIP

* robeaux support

## 0.1 (2013-12-30)

### Accept

* POST and GET for commands

### Adaptor

* and driver generator

### Add

* support for additional parameters
* serialport support
* Travis banner to README
* api commands
* POST command
* manager example
* robot manager
* Sphero example
* Digispark to list of supported platforms
* helper functions
* Driver channel for events
* port to adaptor

### Alter

* structure

### Beaglebone

* Black GPIO

### Clean

* up files
* up variables
* up some comments

### Correctly

* start drivers

### DRY

* up On function

### Dots

* for ignoring imports

### Drop

* unnecessary api parameters

### Expose

* robot functions via api

### Fix

* example

### Go

* fmt examples

### Initial

* GETs for api
* commit

### Install

* ginkgo and gomega dependencies

### Merge

* branch 'examples'
* branch 'master' into ginkgo
* branch 'master' into ginkgo

### More

* WIP of base structs

### Now

* using Connection.Connect()

### Pending

* tests for Robot

### Proper

* formatting for README example

### Properly

* set default interval

### Refactor

* robot name assignment func, and tests to prove it

### Reformat

* examples using gofmt
* source using gofmt

### Remove

* Params from driver struct
* extra nesting

### Rename

* Gobot struct to Manager

### Set

* GOMAXPROCS property in GobotManager

### Skeleton

* for ginkgo/gomega testing

### Small

* refactor
* robot refactor

### StartDriver

* is now optional

### Switch

* to adaptor, driver, connection and device interfaces

### Travis

* lang build

### Tweak

* json output

### Update

* examples
* README.md
* timers and fix issues

### WIP

* multiple robot support
* connections and devices

### Work

* is optional
