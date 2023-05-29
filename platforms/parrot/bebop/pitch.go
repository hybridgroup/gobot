package bebop

import "math"

// ValidatePitch helps validate pitch values such as those created by
// a joystick to values between 0-100 that are required as
// params to Parrot Bebop PCMDs
func ValidatePitch(data float64, offset float64) int {
	value := math.Abs(data) / offset
	if value >= 0.1 {
		if value <= 1.0 {
			return int((float64(int(value*100)) / 100) * 100)
		}
		return 100
	}
	return 0
}
