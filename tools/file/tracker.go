package file

import (
	"io"
	"os"
)

type FileTracker struct {
	target    **os.File
	oldValue  *os.File
	tempValue *os.File
}

func NewFileTracker(target **os.File) *FileTracker {
	return &FileTracker{target: target, oldValue: *target}
}

func (f *FileTracker) Begin() {
	file, err := os.CreateTemp("", "go_test_*.txt")
	if err != nil {
		panic(err)
	}
	*f.target = file
	f.tempValue = file
}

func (f *FileTracker) End() {
	*f.target = f.oldValue
	if err := f.tempValue.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.tempValue.Name()); err != nil {
		panic(err)
	}
}

func (f *FileTracker) Value() string {
	if _, err := f.tempValue.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}
	buff, err := io.ReadAll(f.tempValue)
	if err != nil {
		panic(err)
	}
	return string(buff)
}
