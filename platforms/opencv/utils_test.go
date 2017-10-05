package opencv

import (
	"path"
	"runtime"
	"testing"

	"github.com/hybridgroup/gocv"
	"gobot.io/x/gobot/gobottest"
)

func TestUtils(t *testing.T) {
	_, currentfile, _, _ := runtime.Caller(0)
	image := gocv.IMRead(path.Join(path.Dir(currentfile), "lena-256x256.jpg"), gocv.IMReadColor)
	rect := DetectFaces("haarcascade_frontalface_alt.xml", image)
	gobottest.Refute(t, len(rect), 0)
	DrawRectangles(image, rect, 0, 0, 0, 0)
}
