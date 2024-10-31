//go:build gocv
// +build gocv

package opencv

import (
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"gocv.io/x/gocv"
)

func TestUtils(t *testing.T) {
	_, currentfile, _, _ := runtime.Caller(0)
	image := gocv.IMRead(path.Join(path.Dir(currentfile), "lena-256x256.jpg"), gocv.IMReadColor)
	rect := DetectObjects("haarcascade_frontalface_alt.xml", image)
	assert.NotEqual(t, 0, len(rect))
	DrawRectangles(image, rect, 0, 0, 0, 0)
}
