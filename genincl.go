package genincl

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type File struct {
	Data    []byte
	Modsec  int64
	Modnano int64
}

var genfiles map[string]File

func Register(files map[string]File) {
	genfiles = files
}

func ReadFile(filename string) ([]byte, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return ioutil.ReadFile(filename)
	}

	file, ok := genfiles[filename]
	if ok == false {
		return nil, fmt.Errorf("readfile %s: file not found", filename)
	}

	return file.Data, nil
}

func Stat(filename string) (os.FileInfo, error) {
	fi, err := os.Stat(filename)
	if err == nil {
		return fi, nil
	}

	return newFileInfo(filename)
}

type fileinfo struct {
	filename string
	size     int64
	modsec   int64
	modnano  int64
}

func newFileInfo(filename string) (os.FileInfo, error) {
	file, ok := genfiles[filename]
	if ok == false {
		return nil, os.ErrNotExist
	}

	return &fileinfo{
		filename: filename,
		size:     int64(len(file.Data)),
		modsec:   file.Modsec,
		modnano:  file.Modnano,
	}, nil
}

func (f *fileinfo) Name() string       { return f.filename }
func (f *fileinfo) Size() int64        { return f.size }
func (f *fileinfo) Mode() os.FileMode  { return 0444 }
func (f *fileinfo) ModTime() time.Time { return time.Unix(f.modsec, f.modnano) }
func (f *fileinfo) IsDir() bool        { return false }
func (f *fileinfo) Sys() interface{}   { return nil }
