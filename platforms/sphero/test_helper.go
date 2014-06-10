package sphero

type sp struct{}

func (s sp) Write(b []byte) (int, error) {
	return len(b), nil
}
func (s sp) Read(b []byte) (int, error) {
	return len(b), nil
}
func (s sp) Close() error {
	return nil
}
