// +build example
//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_travis.go /dev/ttyACM0
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

type TravisResponse struct {
	ID                  int    `json:"id"`
	Slug                string `json:"slug"`
	Description         string `json:"description"`
	PublicKey           string `json:"public_key"`
	LastBuildID         int    `json:"last_build_id"`
	LastBuildNumber     string `json:"last_build_number"`
	LastBuildStatus     int    `json:"last_build_status"`
	LastBuildResult     int    `json:"last_build_result"`
	LastBuildDuration   int    `json:"last_build_duration"`
	LastBuildLanguage   string `json:"last_build_language"`
	LastBuildStartedAt  string `json:"last_build_started_at"`
	LastBuildFinishedAt string `json:"last_build_finished_at"`
}

func turnOn(robot *gobot.Robot, device string) {
	robot.Device(device).(*gpio.LedDriver).On()
}

func resetLeds(robot *gobot.Robot) {
	robot.Device("red").(*gpio.LedDriver).Off()
	robot.Device("green").(*gpio.LedDriver).Off()
	robot.Device("blue").(*gpio.LedDriver).Off()
}

func checkTravis(robot *gobot.Robot) {
	resetLeds(robot)
	user := "hybridgroup"
	name := "gobot"
	//name := "broken-arrow"
	fmt.Printf("Checking repo %s/%s\n", user, name)
	turnOn(robot, "blue")
	resp, err := http.Get(fmt.Sprintf("https://api.travis-ci.org/repos/%s/%s.json", user, name))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var travis TravisResponse
	json.Unmarshal(body, &travis)
	resetLeds(robot)
	if travis.LastBuildStatus == 0 {
		turnOn(robot, "green")
	} else {
		turnOn(robot, "red")
	}
}

func main() {
	master := gobot.NewMaster()
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])
	red := gpio.NewLedDriver(firmataAdaptor, "7")
	red.SetName("red")
	green := gpio.NewLedDriver(firmataAdaptor, "6")
	green.SetName("green")
	blue := gpio.NewLedDriver(firmataAdaptor, "5")
	blue.SetName("blue")

	work := func() {
		checkTravis(master.Robot("travis"))
		gobot.Every(10*time.Second, func() {
			checkTravis(master.Robot("travis"))
		})
	}

	robot := gobot.NewRobot("travis",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{red, green, blue},
		work,
	)

	master.AddRobot(robot)
	master.Start()
}
