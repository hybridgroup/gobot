package sysfs

import (
	"fmt"
	"io"
	"log"
	"sync"
	"syscall"
)

const (
	EPOLLET = 1 << 31
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

func (l *InterruptListener) add(fd int, handler InterruptListenerHandler) error {

	l.addHandler(fd, handler)

	e := syscall.EpollEvent{
		Events: syscall.EPOLLIN | syscall.EPOLLPRI | EPOLLET,
		Fd:     int32(fd),
	}

	return syscall.EpollCtl(l.epfd, syscall.EPOLL_CTL_ADD, fd, &e)
}

func (l *InterruptListener) Close() error {
	return syscall.Close(l.epfd)
}

// starts a go routine in the background that does the even looop
func (l *InterruptListener) Start() error {
	go func() {
		err := l.Run()
		if err != nil {
			log.Printf("error in InterruptListener.Run: %s", err)
		}
	}()
	return nil
}

func (l *InterruptListener) Run() error {
	const MaxEpollEvents = 16
	events := make([]syscall.EpollEvent, MaxEpollEvents)

	for {
		num, err := syscall.EpollWait(l.epfd, events, -1)
		if err != nil {
			if err == syscall.EINTR {
				continue
			}
			return err
		}

		for _, e := range events[0:num] {
			_, err = syscall.Seek(int(e.Fd), 0, io.SeekStart)
			if err != nil {
				return err
			}

			buf := make([]byte, 1)
			read, err := syscall.Read(int(e.Fd), buf)
			if err != nil {
				return err
			}

			if read == 0 {
				return fmt.Errorf("why did I read nothing!")
			}

			h := l.getHandler(int(e.Fd))
			h(buf[0])
		}
	}
}

func NewInterruptListener() (*InterruptListener, error) {
	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &InterruptListener{epfd: epfd, handlers: map[int]InterruptListenerHandler{}}, nil
}
