package upgrader

import "errors"

func dupFd(fd uintptr, name fileName) (*file, error) {
	return nil, errors.New("upgrader: duplicating file descriptors is not supported on this platform")
}
