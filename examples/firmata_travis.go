//go:build example
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
	"io"
	"net/http"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/firmata"
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
	if err := robot.Device(device).(*gpio.LedDriver).On(); err != nil {
		fmt.Println(err)
	}
}

func resetLeds(robot *gobot.Robot) {
	if err := robot.Device("red").(*gpio.LedDriver).Off(); err != nil {
		fmt.Println(err)
	}
	if err := robot.Device("green").(*gpio.LedDriver).Off(); err != nil {
		fmt.Println(err)
	}
	if err := robot.Device("blue").(*gpio.LedDriver).Off(); err != nil {
		fmt.Println(err)
	}
}

func checkTravis(robot *gobot.Robot) {
	resetLeds(robot)
	user := "hybridgroup"
	name := "gobot"
	// name := "broken-arrow"
	fmt.Printf("Checking repo %s/%s\n", user, name)
	turnOn(robot, "blue")
	resp, err := http.Get(fmt.Sprintf("https://api.travis-ci.org/repos/%s/%s.json", user, name))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var travis TravisResponse
	if err := json.Unmarshal(body, &travis); err != nil {
		fmt.Println(err)
	}
	resetLeds(robot)
	if travis.LastBuildStatus == 0 {
		turnOn(robot, "green")
	} else {
		turnOn(robot, "red")
	}
}

func main() {
	manager := gobot.NewManager()
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])
	red := gpio.NewLedDriver(firmataAdaptor, "7", gpio.WithName("red"))
	green := gpio.NewLedDriver(firmataAdaptor, "6", gpio.WithName("green"))
	blue := gpio.NewLedDriver(firmataAdaptor, "5", gpio.WithName("blue"))

	work := func() {
		checkTravis(manager.Robot("travis"))
		gobot.Every(10*time.Second, func() {
			checkTravis(manager.Robot("travis"))
		})
	}

	robot := gobot.NewRobot("travis",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{red, green, blue},
		work,
	)

	manager.AddRobot(robot)
	if err := manager.Start(); err != nil {
		panic(err)
	}
}
