//go:build windows

package upgrader

func dupFd(fd uintptr, name fileName) (*file, error) {
	return newFile(fd, name), nil
}
