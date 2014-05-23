package ardrone

type testDrone struct{}

func (me testDrone) Takeoff() bool              { return true }
func (me testDrone) Land()                      {}
func (me testDrone) Up(a float64)               {}
func (me testDrone) Down(a float64)             {}
func (me testDrone) Left(a float64)             {}
func (me testDrone) Right(a float64)            {}
func (me testDrone) Forward(a float64)          {}
func (me testDrone) Backward(a float64)         {}
func (me testDrone) Clockwise(a float64)        {}
func (me testDrone) Counterclockwise(a float64) {}
func (me testDrone) Hover()                     {}
