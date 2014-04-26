package gobotOpencv

import (
	"github.com/hybridgroup/go-opencv/opencv"
)

func loadHaarClassifierCascade(haar string) *opencv.HaarCascade {
	return opencv.LoadHaarClassifierCascade(haar)
}

func DetectFaces(haar string, image *opencv.IplImage) []*opencv.Rect {
	cascade := loadHaarClassifierCascade(haar)
	return cascade.DetectObjects(image)
}

func DrawRectangles(image *opencv.IplImage, rect []*opencv.Rect, r int, g int, b int, thickness int) *opencv.IplImage {
	for _, value := range rect {
		opencv.Rectangle(image,
			opencv.Point{value.X() + value.Width(), value.Y()},
			opencv.Point{value.X(), value.Y() + value.Height()},
			opencv.NewScalar(b, g, r), thickness, 1, 0)
	}
	return image
}
