package ardrone

type testDrone struct{}

func (t testDrone) Takeoff() bool              { return true }
func (t testDrone) Land()                      {}
func (t testDrone) Up(a float64)               {}
func (t testDrone) Down(a float64)             {}
func (t testDrone) Left(a float64)             {}
func (t testDrone) Right(a float64)            {}
func (t testDrone) Forward(a float64)          {}
func (t testDrone) Backward(a float64)         {}
func (t testDrone) Clockwise(a float64)        {}
func (t testDrone) Counterclockwise(a float64) {}
func (t testDrone) Hover()                     {}
