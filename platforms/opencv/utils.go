package opencv

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
)

func loadHaarClassifierCascade(haar string) *cv.HaarCascade {
	return cv.LoadHaarClassifierCascade(haar)
}

func DetectFaces(haar string, image *cv.IplImage) []*cv.Rect {
	return loadHaarClassifierCascade(haar).DetectObjects(image)
}

func DrawRectangles(image *cv.IplImage, rect []*cv.Rect, r int, g int, b int, thickness int) *cv.IplImage {
	for _, value := range rect {
		cv.Rectangle(image,
			cv.Point{value.X() + value.Width(), value.Y()},
			cv.Point{value.X(), value.Y() + value.Height()},
			cv.NewScalar(b, g, r), thickness, 1, 0)
	}
	return image
}
