1.15.0
---
* **build**
    * Switch to CircleCI
* **ble**
    * replace go-ble with tinygo bluetooth package, restore macOS functionality
* **gpio**
    * Update RelayDriver to invert value written on Inverted
    * Add tests for DigitalWrite value
    * Add support for HD44780 LCD controller
    * Add delay for Run function of StepperDriver
* **spi**
    * fixes #700 - Avoid to close the connection.
* **i2c**
    * add SHT2x device
    * add BMP388 Barometric Pressure/Temperature/Altitude Sensor
* **pwm**
    * Resolve issue with PWM for PWMWrite
* **mqtt**
    * Add method to publish MQTT messages with retain flag
* **tello**
    * Add graceful halt for Tello driver
    * Add Tello EDU driver
* **keyboard**
    * add symbol keys for platform/keyboard
* **examples**
    * Update ffmpeg command to decrease latency in tello example

1.14.0
---
* **core**
    * migrating from dep to go modules
    * update codegangsta to urfave (#690)
* **docs**
    * Fix a link in package docs' example code.
* **examples**
    * fixed broken imports due to changed path causing go get to fail
* **gpio**
    * Added ability to make a relay driver inverted (#674)
* **opencv**
    * Update to GoCV 0.21.0
* **spi**
    * Apa102 use default brightness (#671)
* **tello**
    * Updated videoPort for DJI Tello to 11111

1.13.0
---
* **api**
    * Initial stab at Robot-based work
* **build**
    * correct package version as suggested by @dlisin thanks
    * only build last 2 versions of Go plus tip for CI
    * Update dep script for AppVeyor
    * update deps to latest versions of dependencies for GoCV and others
    * Update Gopkg and add test dep to Travis YML
    * update OpenCV build script for OpenCV 4.1.0
* **docs**
    * update to remove Gitter and replace with Slack, and update copyright dates
* **example**
    * add missing nobuild header
* **gpio**
    * Add SparkFun’s EasyDriver (and BigEasyDriver)
    * Add unit tests for TH02 & Minor improvement
    * Added rudiementary support for TH02  Grove Sensor
    * pwm_pin - Fix DutyCycle() parse error, need to trim off trailing '\n' before calling strconv.Atoi(), as other functions in this package do
    * Simplify code as suggested in #617
* **grovepi**
    * add mutex to control transactionality of the device communication
* **i2c**
    * add 128x32 and 96x16 sizes to the i2c ssd1306 driver
    * build out the ccs811 driver
    * update PCA9685 driver to use same protocol as Adafruit Python lib
* **leapmotion**
    * Parser error in Pointable.Bases: Write test and fix
    * Update gobot leap platform to support Leap Motion API v6
* **mavlink**
    * fix mavlink README to use correct example code
* **mqtt**
    * Add some new MQTT adaptor functions with QOS
    * Allow setting QoS on MTT adaptor
    * make tests run correctly even when a local MQTT server is in fact running
    * Do not skip verification of root CA certificates by default InsecureSkipVerify
* **nats**
    * Update Go NATS client library import
* **opencv**
    * minor updates to opencv README
    * update to OpenCV 4.1.0
* **sphero**
    * Added methods to read Sphero Power States
    * Added some new features to the sphero ollie, bb-8 and sprkplus
* **spi**
    * correct param used for APA102 Draw() method
    * Stop using Red parameter for brightness value
* **tello**
    * add direct vector access
    * add example with keyboard
    * Change fps to 60
    * Check for error immediately and skip publish if error occurred
    * update FlightData struct
* **up2**
    * add support for built-in LEDs
    * correct i2c default bus information to match correct values
    * finalize docs for UP2 config steps
    * update README to include more complete setup information
    * useful constant values to access the built-in LEDs

1.12.0
---
* **api**
  * further improvement of the modular API changes
  * modify Start() for more modular initialization, and add StartRaw() for completely custom API implementations
  * settled on StartWithoutDefaults() as the method to start API without default routes
* **core**
  * add Rescale utility function for straight linear rescaling
* **digispark**
  * add examples using digispark with i2c devices blinkm and mlp115a2
  * Added i2c to digispark, but not working yet
  * Added some tests for digispark i2c connector
  * Digispark i2c fixes, added Test for checking available addresses
  * remove test method that should not be in adaptor
  * remove test that is expected to ofail, but passes when digispark board is actually connected
* **docs**
  * add GrovePi to README
  * adjust order of badges in README
  * Fixing broken link
* **examples**
  * add example that uses both the API and also a custom handler with MJPEG streaming from an attached camera
  * small improvements to Tello examples
  * update Tello examples for main thread friendly macOS/Windows, add Tello face tracker
* **i2c**
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
* **joystick**
  * add config file for Magicsee R1 contributed by @carl-ranson
  * add some additional test coverage for file-based config
  * added error handling for config loading in joystick driver
  * mention need to be running a Linux kernel v4.14+ for controller mappings to work as expected
  * provide constant values for existing joystick configurations
* **raspi**
  * export PiBlasterPeriod in Adaptor
* **spi**
  * add ShowImage() function to ssd1306 driver based on @mikegleasonjr suggestion
* **tello**
  * specify end of msgType position
  * add handleResponse testing
  * Add motion cessation commands to Tello
  * handleResponse only needs an io.Reader
  * handleResponse should not send commands
  * rename reqConn to cmdConn
  * reqConn is only an io.WriteCloser
  * send Land() command to drone on Halt() to avoid floating mid-air

1.11.1
---
* **build**
    * exclude vendor and other previously excluded subpackages
    * update Travis build to use OpenCV 3.4.2 release
    * update deps for GoCV to v0.14.0 release
    * Bump periph.io/x/periph to v3.0.0
    * update to Go 1.10.3 and 1.9.7 for Travis builds
* **docs** 
    * Fix Leap Motion package link
* **i2c** 
    * fix write/read gpio on mcp23017, and cleaned up some comments
    * correct pca9685 SetPWMFreq function scaling
* **gopigo3** 
    * update with default spi values, cleanup

1.11.0
---
* **build**
    * correct profile file location for codecov upload
    * Make Go Lint happier by adding some explicit type conversions and ignoring unused error returns
    * single quotes needed to upload any .cov file to codecov for reporting
    * update deps to latest versions for Paho MQTT, go-sdl, and gocv
    * upload any .cov file to codecov for reporting
    * use go 1.10.2 and 1.9.6 for Travis builds
    * add step to call dep ensure before contributing #524
* **examples**
    * correct events used by XBox360 joystick example
* **firmata**
    * Update the Firmata homepage in platform README
* **gpio**
    * Improve Stepper Driver
    * Initial support for MAX7219 (gpio) led driver
* **joystick**
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
* **leapmotion**
    * change timestamp to uint64 to fix #516
* **tello**
    * slow/fast mode switch function
    * StopLanding feature
    * Add Bounce() and PalmLand() funcs and their associated events.
    * bug fix
    * Change several fields in FlightData struct from int16 to bool
    * Export the FlightData fields (see Issue #531)

1.10.2
---
* **opencv**
    * update GoCV to latest version

1.10.1
---
* **tello**
    * improve support for DJI Tello drone, especially video

1.10.0
---
* **docs**
    * add gitter badge to readme
* **gpio**
    * AIP1640 led driver, used in Wemos D1 mini's matrix LED shield
* **spi**
    * switch to using periph.io for SPI interfaces
    * add support for ssd1306 
    * add optional params such as bus/chip to all current drivers
    * complete refactoring to spi.Connection
    * remove unneeded code as suggested by @maruel
    * remove unneeded type and cleanup GoDocs
* **ble**
    * correct spelling error in function name
* **build**
    * update to latest version of Go 1.10 for Travis build
* **cli**
    * remove extra newline
* **docs**
    * add recently contributed GPIO devices to README
* **joystick**
    * able to configure joysticks without external json file
    * correct error in scanning script
    * correct events used by gamepad-style up/down/left/right buttons
    * correct scanner error from ID
    * removed double release event
* **tello**
    * add support for DJI Tello drone

1.9.0
---
* **beaglebone**
  * update pin naming, docs, and examples for the latest Debian OS releases
* **opencv**
  * update build settings needed to build OpenCV/GoCV as part of test suite
  * deps for latest GoCV v0.9.0
* **build**
  * update Travis build to use very latest Go versions
* **docs**
  * add references to new drivers for ADXL345, BH1750, and TM1638.
  * improve docs for installation and use of OpenCV/GoCV from Gobot
  * update copyright date to 2018
* **gpio**
  * Initial support for TM1638 modules
* **i2c**
  * Added basic driver for BH1750 (light sensor), board GY-302
  * support for accel ADXL345
* **bb8/ollie/sprkplus**
  * add Boost command
  * add Set Back LED Output command
  * add Set Raw Motor Values command
  * add Set Rotation Rate command
  * add Set Stabilization command
* **test** 
  * Refactor TestAdaptorDigitalPinConcurrency test

1.8.0
---
* **sysfs** 
  * pause briefly to allow udev rules to apply when exporting PWMPin
* **beaglebone** 
  * correct uboot installation instructions
  * add SPI support
  * no more slots, add docs on configuring u-boot overlays
  * handle gpio pinmux without relying on specific pre-existing setup
* **pocketbeagle** 
  * add support for PocketBeagle
  * use universal io cape manager to initialize board setup
  * improve docs for latest Debian OS
* **build** 
  * Add dep, change how tests run in CI
  * update dependencies to latest GoCV version
* **spi** 
  * Add MCP3002, MCP3202, MCP3204, MCP3208, MCP3304, MCP3004, and MCP3008 A/D converter drivers
  * adding initial support for APA102 LEDs, thanks to code sample from @rakyll
  * extract shared SPI init code into spi package
* **up2** 
  * initial work on support for UP2 board
* **gopigo3**
  * fixed set/get bug with motor dps
* **gpio**
  * Adding stepper motor module
* **firmata** 
  * handle cases where out of sync data is read from serial port on first connecting
* **i2c**
  * Change init payload sequence within jhd1313m1 driver Start() func.

1.7.1
---
* **sprkplus**
  * add new platform for Sphero SPRK+
* **firmata** 
  * correct problem where last analog pin(s) were being ignored from capabilities query
* **ble** 
  * use go-ble/ble fork for BLE interactions
* **build** 
  * update to use latest OpenCV version
  * update to use latest Golang versions

1.7.0
---
* **curie**
  * Add Linux specific step to Intel Curie docs
* **mqtt**
  * Added SetCleanSession
* **build**
  * add go1.9 to versions tested in Travis CI
  * add missing OpenCV lib dependency
  * Update build to use latest Golang versions
  * Travis build will now require sudo to install due to OpenCV
* **docs**
  * some helpful edits for the initial spi implementation
* **gopigo3**
  * integration of recent GoPiGo3 contributions
  * Added grove support, and more gopigo3 examples
* **gpio**
  * Add ButtonDriver.DefaultState to allow for 'reverse' buttons (ones that go from HIGH to LOW)
* **holystone**
  * Add initial support for HS-200
* **i2c**
  * SSD1306.WithDisplayHeight() and SSD1306.WithDisplayWidth() for SSD1306 that use different display ratios
* **joystick**
  * add CLI utilty to scan display events to make it easier to add new joyticks
  * update README to address #441
* **opencv**
  * Switchover to use GoCV and OpenCV 3.3
  * Switch to use custom domain for GoCV package
  * all examples using new GoCV based code
  * correct formatting in face detect example
  * OpenCV face detector that is much more concurrent
  * update interface and examples to indicate multipurpose

1.6.0
---
* **core**
  * log failure errors on Robot Start()
* **build**
  * run test coverage with covermode=set
  * update build to use Golang 1.7.6 and 1.8.3
* **docs**
  * work on ROADMAP doc
* **sysfs**
  * increase test coverage
* **bb8**
  * use updated ble adaptor interface for tests
* **ble**
  * allow for characteristic writes both with and without a response
  * allow override of specific HCI device to use
  * eliminate race conditions from response handling
* **curie**
  * Implement Accelerometer, Gyroscope, and Temperature sensors implemented
  * motion detect implemented
  * shock detect implemented
  * step count implemented
  * tap detect implemented
* **digispark**
  * update blink example to display error message on Start()
  * update README with latest development info
* **edison**
  * auto-discovery of Edison board option
  * removed commented lines
* **firmata**
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
* **bme280**
  * fix signed/unsigned bug
  * Fixed incorrect error condition check when reading the 'ctrl_hum' register.
  * Expanded the BME280 unit test for TestBME280DriverStart() to support reading from the 'ctrl_hum' register.
  * Enables humidity readings in the BME280 driver by enforcing the write to the 'ctrl_meas' register, as per Section 5.4.3 of the BME280 data sheet
* **chip**
  * Fixed PWM duty cycle calculation for C.H.I.P ServoWrite
  * Fixed PWM init bug for C.H.I.P
  * C.H.I.P PWM init robust for already enabled state
* **i2c**
  * remove unused test code
  * write config register in little endian
* **joystick**
  * add needed constants for all PS3 buttons
* **littlewire**
  * littlewire.cc links changed to littlewire.github.io
* **mavlink**
  * switch to using go-serial package
* **megapi**
  * switch to using go-serial package
* **microbit**
  * use updated ble adaptor interface for tests
* **minidrone**
  * add example for Parrot Mambo
  * add support for Mambo external accessories
  * increase test coverage
  * never expect responses for characteristic writes
  * remove unneeded code, increase test coverage
  * separate flight status processing and add test coverage
* **neurosky**
  * switch to using go-serial package
* **ollie**
  * use updated ble adaptor interface for tests
* **sphero**
  * switch to using go-serial package
* **tinkerboard**
  * Updated Tinkerboard and sysfs tests to updated PWM polarity contract

1.5.0
---
* **core**
  * Add Running() methods for Master and Robot and increase test coverage accordingly
* **sysfs**
  * define DigitalPinnerProvider and PWMPinnerProvider interfaces
  * add Chip to be able to change pwmchip, and some related refactoring
  * add file read/write testing for failure conditions
  * proper handling of busy state vs. other errors
  * return sensible result when no valid data read
* **test**
  * increase coverage on test helpers
* **build**
  * switching to Travis builds using Ubuntu 14.04 Trusty
* **aio**
  * only need to support AnalogReader interface
  * avoid test race conditions
  * ensure that AnalogSensor event Data is always int
* **gpio**
  * only need to support DigitalReader/DigitalWriter interface
* **i2c**
  * Added support for the ADS1015 and ADS1115 ADCs
  * Add INA3221 Voltage Monitor
  * Ensure lock of i2c bus for each individual operation
  * Small refactoring and increase test coverage for BMP180
* **beaglebone**
  * implement DigitalPinner and PWMPinner interfaces
  * protect against pin map races
  * increase test coverage
* **chip**
  * add preliminary support for C.H.I.P. Pro
  * add back ServoWrite implementation
  * implement DigitalPinnerProvider and PWMPinnerProvider interfaces
  * protect against pin map races
* **dragonboard**
  * export DigitalPin and PWMPin adaptor methods
  * protect against pin map races
  * increase test coverage
* **edison**
  * auto-detect arduino breakout board, if no specific board is expected
  * ensure that we initialize tristate if arduino breakout board
  * export DigitalPin and PWMPin adaptor methods
  * implement DigitalPinnerProvider and PWMPinnerProvider interfaces
  * protect against pin map races
  * refactoring to reduce code duplication
* **firmata**
  * remove processing that might have been eating test events, increase test coverage
* **joule**
  * implement DigitalPinnerProvider and PWMPinnerProvider interfaces
  * protect against pin map races
  * remove incorrect pin assignment and improve test coverage
  * add examples using Joule with ADS1015 ADC
  * naming system changes
  * correct pin mappings and add PWM example    
* **mavlink**
  * add a Mavlink-over-UDP adaptor.
* **microbit**
  * Add DigitalWriter, DigitalReader, and AnalogReader support using IOPinDriver
  * Handle start error and increase test coverage
* **mqtt**
  * Add a (topic, payload) event type
  * change the On handler to take mqtt.Message
  * increase test coverage
  * update examples that use mqtt for updated notification signature
* **nats**
  * change the On() handler to take the subject as an argument
  * increase test coverage
* **raspi**
  * implement DigitalPinnerProvider and PWMPinnerProvider interfaces
  * add implementation for PWMPinner interface that wraps pi blaster
  * fix adaptor race conditions
  * increase test coverage
* **tinkerboard**
  * Add support for ASUS Tinker Board

1.4.0
---
* **core**
  * Use 10-buffered chans for events, see #374
* **i2c**
  * Many refactors and increases in test coverage
  * Eliminate race conditions introduced by tests
  * Adds Altitude() function to BMP280/BME280
  * bme280 driver Humidity compensation formula
  * ssd1306 driver implementation
* **aio**
  * Eliminate race conditions introduced by tests
* **gpio**
  * Fix motor mode change when speed is set
  * Eliminate race conditions introduced by tests
  * Reduce test side effects
* **ardrone**
  * Increase test coverage
* **audio**
  * Increase test coverage
* **bb8**
  * Refactoring to use BLEConnector interface and provide tests
* **bebop**
  * Increase test coverage
* **beaglebone**
  * Increase test coverage
* **ble**
  * Increase test coverage for battery, device information, and generic access drivers
  * Refactoring drivers to use BLEConnector interface and provide tests
* **chip**
  * Added PWM0 support
  * Increase test coverage
* **digispark**
  * Increase test coverage
* **dragonboard**
  * Increase test coverage
* **edison**
  * Remove pointless error checking code
  * Refactor digital pin creation process method
  * Increase test coverage
* **firmata**
  * Eliminate race conditions introduced by tests
  * Increase test coverage for i2c commands
* **joule**
  * Increase test coverage
* **joystick**
  * Increase test coverage
* **keyboard**
  * Increase test coverage
* **mavlink**
  * Eliminate race conditions introduced by tests
  * Increase test coverage
* **mavlink**
  * Increase test coverage
* **microbit**
  * Refactoring to use BLEConnector interface and provide tests
  * Address #404 by adding info about required magnetometer calibration step to README
  * Increase test coverage
* **minidrone**
  * Refactoring to use BLEConnector interface and provide tests
* **mqtt**
  * Increase test coverage
* **nats**
  * Increase test coverage
* **neurosky**
  * Update neurosky README & example
  * Eliminate race conditions introduced by tests
  * Increase test coverage
* **ollie**
  * Refactoring to use BLEConnector interface and provide tests
  * Correct race condition error on seq
  * Increase test coverage
* **opencv**
  * Increase test coverage
* **particle**
  * Increase test coverage
* **raspi**
  * Address #391 by providing more details about normal development workflow
  * Increase test coverage
* **sphero**
  * Eliminate race conditions
  * Increase test coverage
* **sysfs**
  * Address race condition from udev rules when exporting GPIO pins
  * Increase test coverage
* **docs**
  * Improve explanations for scp/ssh workflow on SoC boards
  * Include entire Apache 2.0 license in the license text
* **test**
  * Add crude travis check for gofmt; format all sources
  * Significantly speed up travis and make runs
  * Remove test code no longer being called
  * Update Travis to run tests using Golang 1.8.1
  * Increase gobottest test coverage

1.3.0
---
* **microbit**
  * Add new platform support
* **dragonboard**
  * Add new platform support
* **gpio**
  * Increase test coverage
* **i2c**
  * Update list of supported i2c devices
  * Minor adjustments and test coverage improvements
  * Added more capabilities checks for I2C
  * Removed smbus block operations
* **core**
  * Increase test coverage
* **test**
  * Improvements to run tests much faster thanks @maruel
  * Use codecov.io for code coverage reporting
* **docs**
  * Update CoC based on Contributor Covenant

1.2.0
---
* **core**
  * Use new improved default namer to avoid API conflicts
* **gpio**
  * Removed scaling function from servo driver
  * Correct servo driver to pass along angle to adaptor to sort out implementation
* **i2c**
  * Refactored platforms and drivers to new I2C interfaces
  * Change to make I2C support more than one bus
  * Refactor drivers to support new optional params
* **bb8**
  * Added collision detection support and example
* **beaglebone**
  * Correct i2c buses to match actual mapping
* **ble**
  * Switch to using [ble](https://github.com/currantlabs/ble) package for Bluetooth LE
  * Basic serial over BLE working with Arduino101 with StandardFirmataBLE
  * WIP on multiple simultaneous ble devices
* **chip**
  * Fixed chip XIO base address lookup  
* **digispark**
 * Fix #288 by using pkg-config to locate libusb-compat includes
* **firmata**
  * Remove race conditions identified in Firmata client
  * Correct error in I2C reads not listening to board events
* **mqtt**
  * Add driver for syntactical sugar around virtual devices
  * Add SSL/TLS client options support
  * Fix #277 by adding SetAutoReconnect method to set Paho MQTT client
  * Change both 'On' and 'Publish' method function signatures to match Eventer interface
* **nats**
  * Add driver to make it easier to create virtual devices
* **ollie**
  * Added collision detection support and example
* **parrot**
  * Add ValidatePitch helper function for Parrot Minidrone, Parrot Bebop & ARDrone 2.0 to package
* **docs**
  * Fix #363 by using atomic.Value to protect current values used by multiple goroutines in drone examples
* **test**
  * Remove Golang 1.5 from TravisCI tests in prep for Golang 1.8 release

1.1.0
---
* **core**
  * use canonical import path for sysfs package
* **i2c**
  * Add a driver for the SHT3X chip
  * Add a driver for BMP180
  * Add support for L3GD20H gyroscope
* **firmata**
  * Add support for TCPFirmata connections, allowing ESP8266 and other WiFi-connected controllers
  * Add mention to README to use 'tty.' serial port on OSX
  * Add mention of A4 and A5 normally unavailable on Firmata
* **raspi**
  * Correct README build instructions with missing 'go build' command
* **snapcraft**
  * Add the packaging metadata to build the gobot snap for Ubuntu Snappy
* **particle**
  * Update examples to take key params via command line
  * Address #160 by adding support for tinker-servo sketch if installed on Particle device
* **esp8266** add experimental ESP8266 support to list of supported platforms
* **sysfs**
  * Should fix #272 by using first byte of data as command register for I2C reads
  * Some additional cleanup suggested by golint
* **ble**
  * Add generic access service driver
  * Update docs to include reference to included drivers
  * Move various test code to test file
* **ollie**
  * Refactoring so no need to expose internal implementation details
* **bebop**
  * Add support/example of RTP video
  * Enable video on firmware 3.3+
  * Update ps3 and video example to enable the video stream
  * Update README for brief explanation of how to get drone video
  * Corrected import paths for client examples
* **bb8**
  * Correct NewDriver params and set name
  * Add missing constructor to wrap Ollie implementation
* **minidrone**
  * Update README with example and which specific models are currently supported
  * Add all piloting flying state events
  * Adds Emergency() and TakePicture() commands
* **test**
  * Add Golang 1.8beta2 to Travis builds
  * Correct aio references for AnalogRead tests

1.0.0
---
* **core**
  * Refactoring to allow 'Metal' development using Gobot packages
  * Able to run robots without being part of a Master.
  * Now running all work in separate goroutines
  * Rename internal name of Master type
  * Refactor events to use channels all the way down.
  * Eliminate potential race conditions from Events and Every functions
  * Add Unsubscribe() to Eventer, now Once() works as expected
  * DeleteEvent function added to Eventer interface
  * Ranges over event channels instead of using select
  * No longer return non-standard slices of errors, instead use hashicorp/go-multierror
  * Ensure that all drivers have default names
  * Now both Robot and Master operate using AutoRun as expected
  * Use canonical import domain of gobot.io for all code
  * Use time.Sleep unless waiting for a timeout in a select
  * Uses time.NewTimer() instead of time.After() to be more efficient

* **test**
  * Add deps tasks to Makefile
  * Add golang 1.7 to Travis CI tests
  * Add golang 1.8beta1 to build matrix for Travis
  * Reduce Travis builds to golang 1.4+ since it is late 2016 already
  * Complete move of test interfaces into the test files where they belong
  * Adds Parrot Minidrone and Sphero Ollie to Travis tests

* **Add missing godocs for everything**

* **i2c**
  * Move I2C drivers into appropriately named 'drivers/i2c' directory
  * Add support for Adafruit Servo/PWM HAT

* **gpio**
  * Move GPIO drivers into appropriately named 'drivers/gpio' directory
  * Add support for PIR motion detector

* **beaglebone**
  * auto-detect Linux kernel version
  * map usr LEDs to match all kernels

* **ble**
  * Rename drivers to make them more obvious
  * Add test placeholders

* **chip**
  * Auto-detect OS version to adjust pin mappings
  * Correct base for new 4.4 GPIO

* **edison**
  * Support for other breakout boards besides Arduino

* **firmata**
  * Use io.ReadFull in platforms/firmata/client
  * Update tarm/goserial to tarm/serial

* **joule**
  * Add support for Intel Joule

* **megapi**
  * Adding support for MakeBlock megapi

* **nats**
  * Add support for NATS server

* **particle**
  * Complete renaming Spark platform to Particle

* **parrot**
  * Move Parrot Minidrone into own platform
  * Move both ARDrone and Bebop under Parrot package

* **raspi**
  * Add missing godocs and small refactors for platform

* **sphero**
  * Add initial support for Sphero BB-8 platform
  * Move Sphero Ollie into own platform

0.12.0
---
* **Refactor Gobot test helpers into separate package**
* **Improve Gobot.Every method to return channel, allowing it to be halted**
* **Refactor of sysfs adds substantial speed improvements**
* **ble**
  * Experimental support for Bluetooth LE.
  * Initial support for Battery & Device Information services
  * Initial support for Sphero BLE robots such as Ollie
  * Initial support for Parrot Minidrone
* **audio**
  * Add new platform for Audio playback
* **gpio**
  * Support added for new GPIO device:
    * RGB LED
  * Bugfixes:
    * Correct analog to better handle quick changes
    * Correct handling of errors and buffering for Wiichuk
* **mqtt**
  * Add support for MQTT authentication
* **opencv**
  * Switching to use main fork of OpenCV
  * Some minor bugfixes related to face tracking

0.11.0
---
* **Support for Golang 1.6**
* **Determine I2C adaptor capabilities dynamically to avoid use of block I/O when unavailable**
* **chip**
  * Add support for GPIO & I2C interfaces on C.H.I.P. $9 computer
* **leap motion**
  * Add support additional "hand" and "gesture" events
* **mqtt**
  * Support latest update to Eclipse Paho MQTT client library
* **raspberry pi**
  * Proper release of Pi Blaster for PWM pins
* **bebop**
  * Prevent event race conditions on takeoff/landing
* **i2c**
  * Support added for new i2c device:
    * MCP23017 Port Expander
  * Bugfixes:
    * Correct init and data parsing for MPU-6050
    * Correct handling of errors and buffering for Wiichuk

0.10.0
---
* **Refactor core to cleanup robot initialization and shutdown**
* **Remove unnecessary goroutines spawned by NewEvent**
* **api**
  * Update Robeaux to v0.5.0
* **bebop**
  * Add support for the Parrot Bebop drone
* **keyboard**
  * Add support for keyboard control
* **gpio**
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
* **i2c**
  * Support added for 2 new Grove i2c devices:
    * Grove Accelerometer
    * Grove LCD with RGB backlit display
* **docs**
  * Many useful fixes and updates for docs, mostly contributed by our wonderful community.

0.8.2
---
  - firmata
    - Refactor firmata adaptor and split firmata protocol implementation into sub `client` package
  - gpio
    - Add support for LIDAR-Lite
  - raspi
    - Add PWM support via pi-blaster
  - sphero
    - Add `ConfigureLocator`, `ReadLocator` and `SetRotationRate`  

0.8.1
---
  - spark
    - Add support for spark Events, Functions and Variables
  - sphero
    - Add `SetDataStreaming` and `ConfigureCollisionDetection` methods

0.8
---
  - Refactor core, gpio, and i2c interfaces
  - Correctly pass errors throughout packages and remove all panics
  - Numerous bug fixes and performance improvements
  - api
    - Update robeaux to v0.3.0
  - firmata
    - Add optional io.ReadWriteCloser parameter to FirmataAdaptor
    - Fix `thread exhaustion` error
  - cli
    - generator
      - Update generator for new adaptor and driver interfaces
      - Add driver, adaptor and project generators
      - Add optional package name parameter

0.7.1
---
  - opencv
    - Fix pthread_create issue on Mac OS

0.7
---
  - Dramatically increased test coverage and documentation
  - api
    - Conform to the [cppp.io](https://gobot.io/x/cppp-io) spec
    - Add support for basic middleware
    - Add support for custom routes
    - Add SSE support
  - ardrone
    - Add optional parameter to specify the drones network address
  - core
    - Add `Once(e *Event, f func(s interface{})` Event function
    - Rename `Expect` to `Assert` and add `Refute` test helper function
  - i2c
    - Add support for MPL115A2
    - Add support for MPU6050
  - mavlink
    - Add support for `common` mavlink messages
  - mqtt
    - Add support for mqtt
  - raspi
    - Add support for the Raspberry Pi
  - sphero
    - Enable stop on sphero disconnect
    - Add `Collision` data struct  
  - sysfs
    - Add generic linux filesystem gpio implementation

0.6.3
---
- Add support for the Intel Edison

0.6.2
---
- cli
  - Fix typo in generator
- leap
  - Fix incorrect Port reference
  - Fix incorrect Event name
- neurosky
  - Fix incorrect Event names
- sphero
  - Correctly format output of GetRGB

0.6.1
---
- cli
  - Fix template error in generator

0.6  
---  
- api
  - Add robeaux support
- core
  - Refactor `Connection` and `Device`
  - Connections are now a collection of Adaptors
  - Devices are now a collection of Drivers
  - Add `Event(string)` function instead of `Events[string]` for retrieving Driver event
  - Add `AddEvent(string)` function to register an event on a Driver
- firmata
  - Fix slice bounds out of range error
- sphero
  - Fix issue where the driver would not halt correctly on OSX

0.5.2  
---  
- beaglebone
  - Add `DirectPinDriver`
  - Ensure slots are properly loaded

0.5.1  
---  
- core
  - Add `Version()` function for Gobot version retrieval
- firmata
  - Fix issue with reading analog inputs
  - Add `data` event for `AnalogSensorDriver`

0.5      
---  
- Idomatic clean up
- Removed reflections throughout packages
- All officially supported platforms are now in ./platforms
- API is now a new package ./api
- All platforms examples are in ./examples
- Replaced martini with net/http
- Replaced ginkgo/gomega with system testing package
- Refactor gobot/robot/device commands
- Added Event type
- Replaced Master type with Gobot type
- Every` and `After` now accept `time.Duration`
- Removed reflection helper methods
