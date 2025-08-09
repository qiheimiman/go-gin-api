package lock

import (
	"os"
	"syscall"
)

type FileLock struct {
	file *os.File
	path string
}

func GetFileLock(filename string) (*FileLock, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	path := dir + "/" + filename

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		file.Close()
		return nil, err
	}

	return &FileLock{file: file, path: path}, nil
}

func (l *FileLock) Release() {
	if l.file != nil {
		syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
		l.file.Close()
		os.Remove(l.path)
	}
}
