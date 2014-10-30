package arietta

import (
	"github.com/hybridgroup/gobot/internal"
	"os"
	"strconv"
	"strings"
)

func joinPath(path ...string) string {
	return strings.Join(path, string(os.PathSeparator))
}

// openOrDie opens the path or aborts the program.
func openOrDie(mode int, path ...string) internal.File {
	full := joinPath(path...)
	fi, err := internal.OpenFile(full, mode, 0666)
	if err != nil {
		panic(err)
	}
	return fi
}

// writeStr writes a string value to a path.
func writeStr(value string, path ...string) {
	fi := openOrDie(os.O_WRONLY|os.O_APPEND, path...)
	defer fi.Close()
	fi.WriteString(value)
}

// writeInt writes a integer value to a path.
func writeInt(value int, path ...string) {
	writeStr(strconv.Itoa(value), path...)
}
