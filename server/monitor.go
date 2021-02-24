package server

import (
	"os"
	"strings"
	"syscall"
)

const (
	sigTerm = syscall.Signal(15)
)

var (
	innerProcess = os.Getenv("MW_MONITORED") != ""
)

func childEnv() []string {
	var env []string
	for _, str := range os.Environ() {
		if strings.HasPrefix(str, "MW_NORESTART=") {
			continue
		}
		if strings.HasPrefix(str, "MW_MONITORED=") {
			continue
		}
		env = append(env, str)
	}
	env = append(env, "MW_MONITORED=yes")
	return env
}
