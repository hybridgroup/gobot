package ardrone

import "math"

// ValidatePitch helps validate pitch values such as those created by
// a joystick to values between 0-1.0 that are required as
// params to Parrot ARDrone PCMDs
func ValidatePitch(data float64, offset float64) float64 {
	value := math.Abs(data) / offset
	if value >= 0.1 {
		if value <= 1.0 {
			return float64(int(value*100)) / 100
		}
		return 1.0
	}
	return 0.0
}
