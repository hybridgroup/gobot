package serialport

import "fmt"

type nullReadWriteCloser struct {
	written          []byte
	simulateReadErr  bool
	simulateWriteErr bool
	simulateCloseErr bool
}

func newNullReadWriteCloser() *nullReadWriteCloser {
	return &nullReadWriteCloser{}
}

func (rwc *nullReadWriteCloser) Write(data []byte) (int, error) {
	if rwc.simulateWriteErr {
		return 0, fmt.Errorf("write error")
	}
	rwc.written = append(rwc.written, data...)
	return len(data), nil
}

func (rwc *nullReadWriteCloser) Read(p []byte) (int, error) {
	if rwc.simulateReadErr {
		return 0, fmt.Errorf("read error")
	}
	return len(p), nil
}

func (rwc *nullReadWriteCloser) Close() error {
	if rwc.simulateCloseErr {
		return fmt.Errorf("close error")
	}
	return nil
}
