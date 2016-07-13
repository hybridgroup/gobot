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
    - Conform to the [cppp.io](https://github.com/hybridgroup/cppp-io) spec
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
