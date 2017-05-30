package mocks

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func TestFilesystemOpen(t *testing.T) {
	fs := NewFilesystem()
	f1 := fs.Add("foo")

	gobot.Assert(t, f1.Opened, false)
	f2, err := fs.OpenFile("foo", 0, 0666)
	gobot.Assert(t, f1, f2)
	gobot.Assert(t, err, nil)
}

func TestFilesystemWrite(t *testing.T) {
	fs := NewFilesystem()
	f1 := fs.Add("bar")

	f2, err := fs.OpenFile("bar", 0, 0666)
	gobot.Assert(t, err, nil)
	// Never been read or written.
	gobot.Assert(t, f1.Seq <= 0, true)

	f2.WriteString("testing")
	// Was written.
	gobot.Assert(t, f1.Seq > 0, true)
	gobot.Assert(t, f1.Contents, "testing")
}

func TestFilesystemRead(t *testing.T) {
	fs := NewFilesystem()
	f1 := fs.Add("bar")
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
}
