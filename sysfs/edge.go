package sysfs

import (
	"sync"
)

type InterruptListenerHandler func(b byte)

type InterruptListener struct {
	epfd         int // the file descriptor of the epoll event loop
	handlersLock sync.Mutex
	handlers     map[int]InterruptListenerHandler
}

func (l *InterruptListener) addHandler(fd int, handler InterruptListenerHandler) {
	l.handlersLock.Lock()
	defer l.handlersLock.Unlock()
	l.handlers[fd] = handler
}

func (l *InterruptListener) getHandler(fd int) InterruptListenerHandler {
	l.handlersLock.Lock()
	defer l.handlersLock.Unlock()
	return l.handlers[fd]
}
