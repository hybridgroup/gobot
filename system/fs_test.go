package system

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilesystemOpen(t *testing.T) {
	fs := &nativeFilesystem{}
	file, err := fs.openFile(os.DevNull, os.O_RDONLY, 0o666)
	require.NoError(t, err)
	var _ File = file
}

func TestFilesystemStat(t *testing.T) {
	fs := &nativeFilesystem{}
	fileInfo, err := fs.stat(os.DevNull)
	require.NoError(t, err)
	assert.NotNil(t, fileInfo)
}
