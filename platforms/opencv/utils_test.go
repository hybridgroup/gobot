package opencv

import (
	"path"
	"runtime"
	"testing"

	"gobot.io/x/gobot/gobottest"
	"gocv.io/x/gocv"
)

func TestUtils(t *testing.T) {
	_, currentfile, _, _ := runtime.Caller(0)
	image := gocv.IMRead(path.Join(path.Dir(currentfile), "lena-256x256.jpg"), gocv.IMReadColor)
	rect := DetectObjects("haarcascade_frontalface_alt.xml", image)
	gobottest.Refute(t, len(rect), 0)
	DrawRectangles(image, rect, 0, 0, 0, 0)
}
