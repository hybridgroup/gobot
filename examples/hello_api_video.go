// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hybridgroup/mjpeg"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
	"gocv.io/x/gocv"
)

var (
	deviceID int
	err      error
	webcam   *gocv.VideoCapture
	stream   *mjpeg.Stream
)

func main() {
	// parse args
	deviceID := os.Args[1]

	master := gobot.NewMaster()

	a := api.NewAPI(master)

	// add the standard C3PIO API routes manually.
	a.AddC3PIORoutes()

	// starts the API without adding the Robeaux web interface. However, since the C3PIO API was
	// already added manually above, that will be available.
	a.StartRaw()

	hello := master.AddRobot(gobot.NewRobot("hello"))

	hello.AddCommand("hi_there", func(params map[string]interface{}) interface{} {
		return fmt.Sprintf("This command is attached to the robot %v", hello.Name)
	})

	// open webcam
	webcam, err = gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	// create the mjpeg stream
	stream = mjpeg.NewStream()
	http.Handle("/video", stream)

	// start capturing
	go mjpegCapture()

	master.Start()
}

func mjpegCapture() {
	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf)
	}
}
