package server

import (
	"log"
	"os"
	"os/exec"
)

func daemon() {
	var (
		childEnv = childEnv()
		args     = os.Args
		err      error
	)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = childEnv

	log.Println("start with daemon......")
	if err = cmd.Start(); err != nil {
		panic(err)
	}
}
