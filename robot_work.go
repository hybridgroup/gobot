// RobotWork and the RobotWork registry represent units of executing computation
// managed at the Robot level. Unlike the utility functions gobot.After and gobot.Every,
// RobotWork units require a context.Context, and can be cancelled externally by calling code.
package gobot

import (
	"context"
	"fmt"
	"time"

	"sync"

	"github.com/gobuffalo/uuid"
)

// RobotWorkRegistry contains all the work units registered on a Robot
type RobotWorkRegistry struct {
	sync.RWMutex

	r map[string]*RobotWork
}

const (
	EveryWorkKind = "every"
	AfterWorkKind = "after"
)

// RobotWork represents a unit of work (in the form of an arbitrary Go function)
// to be done once or on a recurring basis. It encapsulations notions of duration,
// context, count of successful runs, etc.
type RobotWork struct {
	id         uuid.UUID
	kind       string
	tickCount  int
	ctx        context.Context
	cancelFunc context.CancelFunc
	function   func()
	ticker     *time.Ticker
	duration   time.Duration
}

// ID returns the UUID of the RobotWork
func (rw *RobotWork) ID() uuid.UUID {
	return rw.id
}

// CancelFunc returns the context.CancelFunc used to cancel the work
func (rw *RobotWork) CancelFunc() context.CancelFunc {
	return rw.cancelFunc
}

// Ticker returns the time.Ticker used in an Every so that calling code can sync on the same channel
func (rw *RobotWork) Ticker() *time.Ticker {
	if rw.kind == AfterWorkKind {
		return nil
	}
	return rw.ticker
}

// Duration returns the timeout until an After fires or the period of an Every
func (rw *RobotWork) Duration() time.Duration {
	return rw.duration
}

func (rw *RobotWork) String() string {
	format := `ID: %s
Kind: %s
TickCount: %d

`
	return fmt.Sprintf(format, rw.id, rw.kind, rw.tickCount)
}

// WorkRegistry returns the Robot's WorkRegistry
func (r *Robot) WorkRegistry() *RobotWorkRegistry {
	return r.workRegistry
}

// Every calls the given function for every tick of the provided duration.
func (r *Robot) Every(ctx context.Context, d time.Duration, f func()) *RobotWork {
	rw := r.workRegistry.registerEvery(ctx, d, f)
	go func() {
	EVERYWORK:
		for {
			select {
			case <-rw.ctx.Done():
				r.workRegistry.delete(rw.id)
				rw.ticker.Stop()
				break EVERYWORK
			case <-rw.ticker.C:
				rw.tickCount++
				f()
			}
		}
	}()

	return rw
}

// After calls the given function after the provided duration has elapsed
func (r *Robot) After(ctx context.Context, d time.Duration, f func()) *RobotWork {
	rw := r.workRegistry.registerAfter(ctx, d, f)
	ch := time.After(d)
	go func() {
	AFTERWORK:
		for {
			select {
			case <-rw.ctx.Done():
				r.workRegistry.delete(rw.id)
				break AFTERWORK
			case <-ch:
				f()
			}
		}
	}()
	return rw
}

// Get returns the RobotWork specified by the provided ID. To delete something from the registry, it's
// necessary to call its context.CancelFunc, which will perform a goroutine-safe delete on the underlying
// map.
func (rwr *RobotWorkRegistry) Get(id uuid.UUID) *RobotWork {
	rwr.Lock()
	defer rwr.Unlock()
	return rwr.r[id.String()]
}

// Delete returns the RobotWork specified by the provided ID
func (rwr *RobotWorkRegistry) delete(id uuid.UUID) {
	rwr.Lock()
	defer rwr.Unlock()
	delete(rwr.r, id.String())
}

// registerAfter creates a new unit of RobotWork and sets up its context/cancellation
func (rwr *RobotWorkRegistry) registerAfter(ctx context.Context, d time.Duration, f func()) *RobotWork {
	rwr.Lock()
	defer rwr.Unlock()

	id, _ := uuid.NewV4()
	rw := &RobotWork{
		id:       id,
		kind:     AfterWorkKind,
		function: f,
		duration: d,
	}

	rw.ctx, rw.cancelFunc = context.WithCancel(ctx)
	rwr.r[id.String()] = rw
	return rw
}

// registerEvery creates a new unit of RobotWork and sets up its context/cancellation
func (rwr *RobotWorkRegistry) registerEvery(ctx context.Context, d time.Duration, f func()) *RobotWork {
	rwr.Lock()
	defer rwr.Unlock()

	id, _ := uuid.NewV4()
	rw := &RobotWork{
		id:       id,
		kind:     EveryWorkKind,
		function: f,
		duration: d,
		ticker:   time.NewTicker(d),
	}

	rw.ctx, rw.cancelFunc = context.WithCancel(ctx)

	rwr.r[id.String()] = rw
	return rw
}
