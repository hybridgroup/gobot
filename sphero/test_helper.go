package gobotSphero

type sp struct{}

func (me sp) Write(b []byte) (int, error) {
	return len(b), nil
}
func (me sp) Read(b []byte) (int, error) {
	return len(b), nil
}
func (me sp) Close() error {
	return nil
}
