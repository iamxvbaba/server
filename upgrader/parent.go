package upgrader

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const (
	sentinelEnvVar = "mw_upgrade_parent"
	notifyReady    = 42
)

type parent struct {
	wr     *os.File
	result <-chan error
	exited <-chan struct{}
}

func newParent(env *env) (*parent, map[fileName]*file, error) {
	if env.getenv(sentinelEnvVar) == "" {
	//	log.Println("没有子进程!!!!!")
		return nil, make(map[fileName]*file), nil
	}

	wr := env.newFile(3, "write")
	rd := env.newFile(4, "read")

	var names [][]string
	dec := gob.NewDecoder(rd)
	if err := dec.Decode(&names); err != nil {
		return nil, nil, fmt.Errorf("can't decode names from parent process: %s", err)
	}
	//log.Println("names:",names)
	files := make(map[fileName]*file)
	for i, parts := range names {
		var key fileName
		copy(key[:], parts)

		// Start at 5 to account for stdin, etc. and write
		// and read pipes.
		fd := 5 + i
		// 设置关闭文件只对当前进程有效
		env.closeOnExec(fd)

		files[key] = &file{
			env.newFile(uintptr(fd), key.String()),
			uintptr(fd),
		}
	}

	result := make(chan error, 1)
	exited := make(chan struct{})
	go func() {
		defer rd.Close()

		n, err := io.Copy(ioutil.Discard, rd)
		if n != 0 {
			err = errors.New("unexpected data from parent process")
		} else if err != nil {
			err = fmt.Errorf("unexpected error while waiting for parent to exit: %s", err)
		}
		result <- err
		close(exited)
	}()

	return &parent{
		wr:     wr,
		result: result,
		exited: exited,
	}, files, nil
}

func (ps *parent) sendReady() error {
	defer ps.wr.Close()
	if _, err := ps.wr.Write([]byte{notifyReady}); err != nil {
		return fmt.Errorf("can't notify parent process: %s", err)
	}
	return nil
}
