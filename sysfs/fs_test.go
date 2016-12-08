package sysfs

import (
	"os"
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func TestFilesystemOpen(t *testing.T) {
	SetFilesystem(&NativeFilesystem{})
	file, err := OpenFile(os.DevNull, os.O_RDONLY, 666)
	gobottest.Assert(t, err, nil)
	var _ File = file
}
