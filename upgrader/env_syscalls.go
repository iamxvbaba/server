// +build !windows

package upgrader

import (
	"os"
	"syscall"
)

var stdEnv = &env{
	newProc:     newOSProcess,
	newFile:     os.NewFile,
	environ:     os.Environ,
	getenv:      os.Getenv,
	closeOnExec: syscall.CloseOnExec,
}
