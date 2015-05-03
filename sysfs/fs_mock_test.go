package sysfs

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func TestMockFilesystemOpen(t *testing.T) {
	fs := NewMockFilesystem([]string{"foo"})
	f1 := fs.Files["foo"]

	gobot.Assert(t, f1.Opened, false)
	f2, err := fs.OpenFile("foo", 0, 0666)
	gobot.Assert(t, f1, f2)
	gobot.Assert(t, err, nil)

	err = f2.Sync()
	gobot.Assert(t, err, nil)

	_, err = fs.OpenFile("bar", 0, 0666)
	gobot.Refute(t, err, nil)

	fs.Add("bar")
	f4, err := fs.OpenFile("bar", 0, 0666)
	gobot.Refute(t, f4.Fd(), f1.Fd())
}

func TestMockFilesystemWrite(t *testing.T) {
	fs := NewMockFilesystem([]string{"bar"})
	f1 := fs.Files["bar"]

	f2, err := fs.OpenFile("bar", 0, 0666)
	gobot.Assert(t, err, nil)
	// Never been read or written.
	gobot.Assert(t, f1.Seq <= 0, true)

	f2.WriteString("testing")
	// Was written.
	gobot.Assert(t, f1.Seq > 0, true)
	gobot.Assert(t, f1.Contents, "testing")
}

func TestMockFilesystemRead(t *testing.T) {
	fs := NewMockFilesystem([]string{"bar"})
	f1 := fs.Files["bar"]
	f1.Contents = "Yip"

	f2, err := fs.OpenFile("bar", 0, 0666)
	gobot.Assert(t, err, nil)
	// Never been read or written.
	gobot.Assert(t, f1.Seq <= 0, true)

	buffer := make([]byte, 20)
	n, err := f2.Read(buffer)

	// Was read.
	gobot.Assert(t, f1.Seq > 0, true)
	gobot.Assert(t, n, 3)
	gobot.Assert(t, string(buffer[:3]), "Yip")

	n, err = f2.ReadAt(buffer, 10)
	gobot.Assert(t, n, 3)
}
