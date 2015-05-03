package sysfs

import (
	"os"
	"testing"

	"github.com/hybridgroup/gobot"
)

func TestFilesystemOpen(t *testing.T) {
	SetFilesystem(&NativeFilesystem{})
	file, err := OpenFile(os.DevNull, os.O_RDONLY, 666)
	gobot.Assert(t, err, nil)
	var _ File = file
}
