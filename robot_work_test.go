package gobot

import (
	"context"
	"testing"

	"time"

	"github.com/gobuffalo/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRobotWork(t *testing.T) {
	id, _ := uuid.NewV4()

	rw := &RobotWork{
		id:       id,
		kind:     EveryWorkKind,
		function: func() {},
	}

	duration := time.Second * 1
	ctx, cancelFunc := context.WithCancel(context.Background())

	rw.ctx = ctx
	rw.cancelFunc = cancelFunc
	rw.duration = duration

	t.Run("ID()", func(t *testing.T) {
		assert.Equal(t, rw.ID(), id)
	})

	t.Run("Ticker()", func(t *testing.T) {
		t.Skip()
	})

	t.Run("Duration()", func(t *testing.T) {
		assert.Equal(t, rw.duration, duration)
	})
}

func TestRobotWorkRegistry(t *testing.T) {
	robot := NewRobot("testbot")

	rw := robot.Every(context.Background(), time.Millisecond*250, func() {
		_ = 1 + 1
	})

	t.Run("Get retrieves", func(t *testing.T) {
		assert.Equal(t, robot.workRegistry.Get(rw.id), rw)
	})

	t.Run("delete deletes", func(t *testing.T) {
		robot.workRegistry.delete(rw.id)
		postDeleteKeys := collectStringKeysFromWorkRegistry(robot.workRegistry)
		assert.NotContains(t, postDeleteKeys, rw.id.String())
	})
}

func collectStringKeysFromWorkRegistry(rwr *RobotWorkRegistry) []string {
	keys := make([]string, len(rwr.r))
	for k, _ := range rwr.r {
		keys = append(keys, k)
	}
	return keys
}
