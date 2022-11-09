package sysfs

import (
	"os"
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func TestFilesystemOpen(t *testing.T) {
	fs := &nativeFilesystem{}
	file, err := fs.openFile(os.DevNull, os.O_RDONLY, 666)
	gobottest.Assert(t, err, nil)
	var _ File = file
}

func TestFilesystemStat(t *testing.T) {
	fs := &nativeFilesystem{}
	fileInfo, err := fs.stat(os.DevNull)
	gobottest.Assert(t, err, nil)
	gobottest.Refute(t, fileInfo, nil)
}
