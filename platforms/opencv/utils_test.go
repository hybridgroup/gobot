package opencv

import (
	"path"
	"runtime"
	"testing"

	"gobot.io/x/gobot/gobottest"
	cv "github.com/lazywei/go-opencv/opencv"
)

func TestUtils(t *testing.T) {
	_, currentfile, _, _ := runtime.Caller(0)
	image := cv.LoadImage(path.Join(path.Dir(currentfile), "lena-256x256.jpg"))
	rect := DetectFaces("haarcascade_frontalface_alt.xml", image)
	gobottest.Refute(t, len(rect), 0)
	gobottest.Refute(t, DrawRectangles(image, rect, 0, 0, 0, 0), nil)
}
