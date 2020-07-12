package genincl

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type File struct {
	Data    []byte
	Modsec  int64
	Modnano int64
	Length  int64
	Format  string
}

var genfiles map[string]File

func Register(files map[string]File) {
	genfiles = files
}

type Reader struct {
	file     File
	data     []byte
	location int
}

func Open(filename string) (io.ReadCloser, error) {
	_, err := os.Stat(filename)
	if err == nil {
		f, err := os.Open(filename)
		return f, err
	}

	file, ok := genfiles[filename]

	if ok == false {
		return nil, fmt.Errorf("unknown filename %s", filename)
	}

	buf := &bytes.Buffer{}

	if file.Format == "zlib" {
		var b = bytes.NewReader(file.Data)

		r, err := zlib.NewReader(b)
		if err != nil {
			return nil, fmt.Errorf("data format error: %v", err)
		}

		if _, err := io.Copy(buf, r); err != nil {
			return nil, fmt.Errorf("data read error: %v", err)
		}
	} else {
		buf = bytes.NewBuffer(file.Data)
	}

	return &Reader{
		file:     file,
		data:     buf.Bytes(),
		location: 0,
	}, nil
}

func (r *Reader) Read(buf []byte) (int, error) {
	n := copy(buf, r.data[r.location:])
	r.location += n
	return n, nil
}

func (r *Reader) Close() error {
	return nil
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

	var out = &bytes.Buffer{}
	var in = bytes.NewBuffer(file.Data)

	r, err := zlib.NewReader(in)
	if err != nil {
		return nil, fmt.Errorf("readfile %s: zlib error %v", filename, err)
	}
	if _, err := io.Copy(out, r); err != nil {
		return nil, fmt.Errorf("readfile %s: zlib copy %v", filename, err)
	}
	r.Close()

	return out.Bytes(), nil
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
		size:     file.Length,
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
