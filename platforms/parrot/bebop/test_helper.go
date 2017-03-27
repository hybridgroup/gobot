package bebop

type testDrone struct{}

func (t testDrone) TakeOff() error                    { return nil }
func (t testDrone) Land() error                       { return nil }
func (t testDrone) Up(n int) error                    { return nil }
func (t testDrone) Down(n int) error                  { return nil }
func (t testDrone) Left(n int) error                  { return nil }
func (t testDrone) Right(n int) error                 { return nil }
func (t testDrone) Forward(n int) error               { return nil }
func (t testDrone) Backward(n int) error              { return nil }
func (t testDrone) Clockwise(n int) error             { return nil }
func (t testDrone) CounterClockwise(n int) error      { return nil }
func (t testDrone) Stop() error                       { return nil }
func (t testDrone) Connect() error                    { return nil }
func (t testDrone) Video() chan []byte                { return nil }
func (t testDrone) StartRecording() error             { return nil }
func (t testDrone) StopRecording() error              { return nil }
func (t testDrone) HullProtection(protect bool) error { return nil }
func (t testDrone) Outdoor(outdoor bool) error        { return nil }
func (t testDrone) VideoEnable(enable bool) error     { return nil }
func (t testDrone) VideoStreamMode(mode int8) error   { return nil }
