package system

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockFilesystemOpen(t *testing.T) {
	fs := newMockFilesystem([]string{"foo"})
	f1 := fs.Files["foo"]

	assert.False(t, f1.Opened)
	f2, err := fs.openFile("foo", 0, 0o666)
	assert.Equal(t, f2, f1)
	require.NoError(t, err)

	err = f2.Sync()
	require.NoError(t, err)

	_, err = fs.openFile("bar", 0, 0o666)
	require.ErrorContains(t, err, " : bar: no such file")

	fs.Add("bar")
	f4, _ := fs.openFile("bar", 0, 0o666)
	assert.NotEqual(t, f1.Fd(), f4.Fd())
}

func TestMockFilesystemStat(t *testing.T) {
	fs := newMockFilesystem([]string{"foo", "bar/baz"})

	fileStat, err := fs.stat("foo")
	require.NoError(t, err)
	assert.False(t, fileStat.IsDir())

	dirStat, err := fs.stat("bar")
	require.NoError(t, err)
	assert.True(t, dirStat.IsDir())

	_, err = fs.stat("plonk")
	require.ErrorContains(t, err, " : plonk: no such file")
}

func TestMockFilesystemFind(t *testing.T) {
	// arrange
	fs := newMockFilesystem([]string{"/foo", "/bar/foo", "/bar/foo/baz", "/bar/baz/foo", "/bar/foo/bak"})
	tests := map[string]struct {
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
			require.NoError(t, err)
			sort.Strings(dirs)
			assert.Equal(t, tt.want, dirs)
		})
	}
}

func TestMockFilesystemWrite(t *testing.T) {
	fs := newMockFilesystem([]string{"bar"})
	f1 := fs.Files["bar"]

	f2, err := fs.openFile("bar", 0, 0o666)
	require.NoError(t, err)
	// Never been read or written.
	assert.LessOrEqual(t, f1.Seq, 0)

	_, _ = f2.WriteString("testing")
	// Was written.
	assert.Positive(t, f1.Seq)
	assert.Equal(t, "testing", f1.Contents)
}

func TestMockFilesystemRead(t *testing.T) {
	fs := newMockFilesystem([]string{"bar"})
	f1 := fs.Files["bar"]
	f1.Contents = "Yip"

	f2, err := fs.openFile("bar", 0, 0o666)
	require.NoError(t, err)
	// Never been read or written.
	assert.LessOrEqual(t, f1.Seq, 0)

	buffer := make([]byte, 20)
	n, _ := f2.Read(buffer)

	// Was read.
	assert.Positive(t, f1.Seq)
	assert.Equal(t, 3, n)
	assert.Equal(t, "Yip", string(buffer[:3]))

	n, _ = f2.ReadAt(buffer, 10)
	assert.Equal(t, 3, n)
}
