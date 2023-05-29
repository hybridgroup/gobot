package system

import (
	"sort"
	"testing"

	"gobot.io/x/gobot/v2/gobottest"
)

func TestMockFilesystemOpen(t *testing.T) {
	fs := newMockFilesystem([]string{"foo"})
	f1 := fs.Files["foo"]

	gobottest.Assert(t, f1.Opened, false)
	f2, err := fs.openFile("foo", 0, 0666)
	gobottest.Assert(t, f1, f2)
	gobottest.Assert(t, err, nil)

	err = f2.Sync()
	gobottest.Assert(t, err, nil)

	_, err = fs.openFile("bar", 0, 0666)
	gobottest.Assert(t, err.Error(), " : bar: no such file")

	fs.Add("bar")
	f4, _ := fs.openFile("bar", 0, 0666)
	gobottest.Refute(t, f4.Fd(), f1.Fd())
}

func TestMockFilesystemStat(t *testing.T) {
	fs := newMockFilesystem([]string{"foo", "bar/baz"})

	fileStat, err := fs.stat("foo")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fileStat.IsDir(), false)

	dirStat, err := fs.stat("bar")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, dirStat.IsDir(), true)

	_, err = fs.stat("plonk")
	gobottest.Assert(t, err.Error(), " : plonk: no such file")
}

func TestMockFilesystemFind(t *testing.T) {
	// arrange
	fs := newMockFilesystem([]string{"/foo", "/bar/foo", "/bar/foo/baz", "/bar/baz/foo", "/bar/foo/bak"})
	var tests = map[string]struct {
		baseDir string
		pattern string
		want    []string
	}{
		"flat":                  {baseDir: "/", pattern: "foo", want: []string{"/foo"}},
		"in directory no slash": {baseDir: "/bar", pattern: "foo", want: []string{"/bar/foo", "/bar/foo", "/bar/foo"}},
		"file":                  {baseDir: "/bar/baz/", pattern: "foo", want: []string{"/bar/baz/foo"}},
		"file pattern":          {baseDir: "/bar/foo/", pattern: "ba.?", want: []string{"/bar/foo/bak", "/bar/foo/baz"}},
		"empty":                 {baseDir: "/", pattern: "plonk", want: nil},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// act
			dirs, err := fs.find(tt.baseDir, tt.pattern)
			// assert
			gobottest.Assert(t, err, nil)
			sort.Strings(dirs)
			gobottest.Assert(t, dirs, tt.want)
		})
	}
}

func TestMockFilesystemWrite(t *testing.T) {
	fs := newMockFilesystem([]string{"bar"})
	f1 := fs.Files["bar"]

	f2, err := fs.openFile("bar", 0, 0666)
	gobottest.Assert(t, err, nil)
	// Never been read or written.
	gobottest.Assert(t, f1.Seq <= 0, true)

	f2.WriteString("testing")
	// Was written.
	gobottest.Assert(t, f1.Seq > 0, true)
	gobottest.Assert(t, f1.Contents, "testing")
}

func TestMockFilesystemRead(t *testing.T) {
	fs := newMockFilesystem([]string{"bar"})
	f1 := fs.Files["bar"]
	f1.Contents = "Yip"

	f2, err := fs.openFile("bar", 0, 0666)
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
