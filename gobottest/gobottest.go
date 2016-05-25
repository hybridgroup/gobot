package gobottest

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var errFunc = func(t *testing.T, message string) {
	t.Errorf(message)
}

func logFailure(t *testing.T, message string) {
	_, file, line, _ := runtime.Caller(2)
	s := strings.Split(file, "/")
	errFunc(t, fmt.Sprintf("%v:%v: %v", s[len(s)-1], line, message))
}

// Assert checks if a and b are equal, emits a t.Errorf if they are not equal.
func Assert(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		logFailure(t, fmt.Sprintf("%v - \"%v\", should equal,  %v - \"%v\"",
			a, reflect.TypeOf(a), b, reflect.TypeOf(b)))
	}
}

// Refute checks if a and b are equal, emits a t.Errorf if they are equal.
func Refute(t *testing.T, a interface{}, b interface{}) {
	if reflect.DeepEqual(a, b) {
		logFailure(t, fmt.Sprintf("%v - \"%v\", should not equal,  %v - \"%v\"",
			a, reflect.TypeOf(a), b, reflect.TypeOf(b)))
	}
}

func ExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}
