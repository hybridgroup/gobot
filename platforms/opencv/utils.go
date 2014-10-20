package opencv

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
)

// loadHaarClassifierCascade returns open cv HaarCascade loaded
func loadHaarClassifierCascade(haar string) *cv.HaarCascade {
	return cv.LoadHaarClassifierCascade(haar)
}

// DetectFaces loads Haar cascade to detect face objects in image
func DetectFaces(haar string, image *cv.IplImage) []*cv.Rect {
	return loadHaarClassifierCascade(haar).DetectObjects(image)
}

// DrawRectangles uses Rect array values to return image with rectangles drawn.
func DrawRectangles(image *cv.IplImage, rect []*cv.Rect, r int, g int, b int, thickness int) *cv.IplImage {
	for _, value := range rect {
		cv.Rectangle(image,
			cv.Point{value.X() + value.Width(), value.Y()},
			cv.Point{value.X(), value.Y() + value.Height()},
			cv.NewScalar(b, g, r), thickness, 1, 0)
	}
	return image
}
