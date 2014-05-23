package main

import (
	"encoding/json"
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"io/ioutil"
	"net/http"
	"time"
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

func resetLeds(robot *gobot.Robot) {
	gobot.Call(robot.GetDevice("red").Driver, "Off")
	gobot.Call(robot.GetDevice("green").Driver, "Off")
	gobot.Call(robot.GetDevice("blue").Driver, "Off")
}

func checkTravis(robot *gobot.Robot) {
	resetLeds(robot)
	user := "hybridgroup"
	name := "gobot"
	//name := "broken-arrow"
	fmt.Printf("Checking repo %s/%s\n", user, name)
	gobot.Call(robot.GetDevice("blue").Driver, "On")
	resp, err := http.Get(fmt.Sprintf("https://api.travis-ci.org/repos/%s/%s.json", user, name))
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var travis TravisResponse
	json.Unmarshal(body, &travis)
	resetLeds(robot)
	if travis.LastBuildStatus == 0 {
		gobot.Call(robot.GetDevice("green").Driver, "On")
	} else {
		gobot.Call(robot.GetDevice("red").Driver, "On")
	}
}

func main() {
	master := gobot.NewGobot()
	firmataAdaptor := firmata.NewFirmataAdaptor("firmata", "/dev/ttyACM0")
	red := gpio.NewLedDriver(firmata, "red", "7")
	green := gpio.NewLedDriver(firmata, "green", "6")
	blue := gpio.NewLedDriver(firmata, "blue", "5")

	work := func() {
		checkTravis(master.FindRobot("travis"))
		gobot.Every(10*time.Second, func() {
			checkTravis(master.FindRobot("travis"))
		})
	}

	master.Robots = append(master.Robots,
		gobot.NewRobot("travis", []gobot.Connection{firmataAdaptor}, []gobot.Device{red, green, blue}, work))

	master.Start()
}
