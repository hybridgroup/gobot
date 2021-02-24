package sysfs

import (
	"fmt"
)

var (
	notSupportedError = fmt.Errorf("epoll/InterruptListenerHandler not supported on darwin")
)

func (l *InterruptListener) add(fd int, handler InterruptListenerHandler) error {
	return notSupportedError
}

func NewInterruptListener() (*InterruptListener, error) {
	return nil, notSupportedError
}

func (l *InterruptListener) Close() error {
	return nil
}

// starts a go routine in the background that does the even looop
func (l *InterruptListener) Start() error {
	return notSupportedError
}
