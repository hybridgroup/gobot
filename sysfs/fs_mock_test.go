package sysfs

import (
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func TestMockFilesystemOpen(t *testing.T) {
	fs := NewMockFilesystem([]string{"foo"})
	f1 := fs.Files["foo"]

	gobottest.Assert(t, f1.Opened, false)
	f2, err := fs.OpenFile("foo", 0, 0666)
	gobottest.Assert(t, f1, f2)
	gobottest.Assert(t, err, nil)

	err = f2.Sync()
	gobottest.Assert(t, err, nil)

	_, err = fs.OpenFile("bar", 0, 0666)
	gobottest.Refute(t, err, nil)

	fs.Add("bar")
	f4, _ := fs.OpenFile("bar", 0, 0666)
	gobottest.Refute(t, f4.Fd(), f1.Fd())
}

func TestMockFilesystemStat(t *testing.T) {
	fs := NewMockFilesystem([]string{"foo", "bar/baz"})

	fileStat, err := fs.Stat("foo")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fileStat.IsDir(), false)

	dirStat, err := fs.Stat("bar")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, dirStat.IsDir(), true)

	_, err = fs.Stat("plonk")
	gobottest.Refute(t, err, nil)
}

func TestMockFilesystemWrite(t *testing.T) {
	fs := NewMockFilesystem([]string{"bar"})
	f1 := fs.Files["bar"]

	f2, err := fs.OpenFile("bar", 0, 0666)
	gobottest.Assert(t, err, nil)
	// Never been read or written.
	gobottest.Assert(t, f1.Seq <= 0, true)

	f2.WriteString("testing")
	// Was written.
	gobottest.Assert(t, f1.Seq > 0, true)
	gobottest.Assert(t, f1.Contents, "testing")
}

func TestMockFilesystemRead(t *testing.T) {
	fs := NewMockFilesystem([]string{"bar"})
	f1 := fs.Files["bar"]
	f1.Contents = "Yip"

	f2, err := fs.OpenFile("bar", 0, 0666)
	gobottest.Assert(t, err, nil)
	// Never been read or written.
	gobottest.Assert(t, f1.Seq <= 0, true)

	buffer := make([]byte, 20)
	n, _ := f2.Read(buffer)

	// Was read.
	gobottest.Assert(t, f1.Seq > 0, true)
	gobottest.Assert(t, n, 3)
	gobottest.Assert(t, string(buffer[:3]), "Yip")

	n, _ = f2.ReadAt(buffer, 10)
	gobottest.Assert(t, n, 3)
}
