// +build example
//
// Do not build by default.

package main

import (
	"log"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	stationName := "DlSnIpEr FM Demo"
	rdsMessage := "DlSnIpEr in the mix"

	adaptor := raspi.NewAdaptor()

	radioConfig := i2c.AdafruitSi4713Config{
		TransmitFrequency: 8850,
		TransmitPower:     115,
		ResetPin:          "29",
		DebugMode:         false,
		HasRDS:            true,
		RDSProgramID:      0x3104,
		RDSStationName:    stationName,
		RDSMessage:        rdsMessage,
		Log:               log.Printf,
		DebugLog:          nil,
	}
	fmRadioBot, err := i2c.NewSi4713Driver(adaptor, radioConfig)
	if err != nil {
		log.Fatalln(err)
	}

	work := func() {
		if err = fmRadioBot.SetRDSMessage(rdsMessage); err != nil {
			log.Fatalln(err)
		}

		gobot.Every(1*time.Second, func() {
			if err = fmRadioBot.CheckDeviceStatus(); err != nil {
				log.Fatalln(err)
			}
		})
	}

	robot := gobot.NewRobot("FM Transmitter Station demo",
		[]gobot.Connection{adaptor},
		[]gobot.Device{fmRadioBot},
		work,
	)

	if err = robot.Start(); err != nil {
		log.Fatalln(err)
	}
}
