package genincl_test

import (
	"os"
	"testing"
	"time"

	. "github.com/dbulkow/genincl"
)

var genfiles = map[string]File{
	"testdata": File{
		Modsec:  1573777556,
		Modnano: 901181155,
		Length:  84,
		Format:  "zlib",
		Data: []byte{
			0x78, 0x9c, 0x1c, 0xcb, 0xe1, 0x0d, 0x80, 0x20, 0x0c,
			0x44, 0xe1, 0x55, 0x6e, 0x22, 0x77, 0x68, 0xe0, 0xd0,
			0x46, 0xe8, 0x19, 0x68, 0xe2, 0xfa, 0x06, 0x7f, 0xbf,
			0xf7, 0x1d, 0x7a, 0xe1, 0x0b, 0x79, 0x11, 0xe9, 0x83,
			0x68, 0x9a, 0xb0, 0xde, 0x71, 0x4a, 0x15, 0x83, 0x81,
			0x14, 0xd2, 0x6e, 0xee, 0xc3, 0x27, 0x9e, 0x6e, 0x85,
			0xf0, 0xf8, 0x41, 0x51, 0x34, 0x56, 0x4e, 0x4b, 0x57,
			0x40, 0x6d, 0xd7, 0x60, 0xae, 0x2f, 0x00, 0x00, 0xff,
			0xff, 0xee, 0x4c, 0x1e, 0x58,
		},
	},
	"bar": File{
		Modsec:  1573777556,
		Modnano: 901181155,
		Length:  84,
		Format:  "zlib",
		Data: []byte{
			0x78, 0x9c, 0x1c, 0xcb, 0xe1, 0x0d, 0x80, 0x20, 0x0c,
			0x44, 0xe1, 0x55, 0x6e, 0x22, 0x77, 0x68, 0xe0, 0xd0,
			0x46, 0xe8, 0x19, 0x68, 0xe2, 0xfa, 0x06, 0x7f, 0xbf,
			0xf7, 0x1d, 0x7a, 0xe1, 0x0b, 0x79, 0x11, 0xe9, 0x83,
			0x68, 0x9a, 0xb0, 0xde, 0x71, 0x4a, 0x15, 0x83, 0x81,
			0x14, 0xd2, 0x6e, 0xee, 0xc3, 0x27, 0x9e, 0x6e, 0x85,
			0xf0, 0xf8, 0x41, 0x51, 0x34, 0x56, 0x4e, 0x4b, 0x57,
			0x40, 0x6d, 0xd7, 0x60, 0xae, 0x2f, 0x00, 0x00, 0xff,
			0xff, 0xee, 0x4c, 0x1e, 0x58,
		},
	},
}

func init() {
	Register(genfiles)
}

func TestOpen(t *testing.T) {
	f, err := Open("testdata")
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 512)

	n, err := f.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	if n != int(genfiles["testdata"].Length) {
		t.Fatal("file length difference")
	}

	if string(buf[:n]) != "Now is the time for all good men to take their place in the confederation of planets" {
		t.Fatal("file contents mismatch")
	}

	f.Close()
}

func TestOpenLocal(t *testing.T) {
	f, err := Open("bar")
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 512)

	n, err := f.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	if n != int(genfiles["bar"].Length) {
		t.Fatal("file length difference")
	}

	if string(buf[:n]) != "Now is the time for all good men to take their place in the confederation of planets" {
		t.Fatal("file contents mismatch")
	}

	f.Close()
}

func TestOpenUnknownFile(t *testing.T) {
	_, err := Open("foobar")
	if err != nil {
		if err.Error() != "unknown filename foobar" {
			t.Fatalf("expected \"unknown filename foobar\" got \"%s\"", err.Error())
		}
		return
	}

	t.Error("expected a failure")
}

func TestReadFile(t *testing.T) {
	buf, err := ReadFile("testdata")
	if err != nil {
		t.Fatal(err)
	}

	if string(buf) != "Now is the time for all good men to take their place in the confederation of planets" {
		t.Fatal("file contents mismatch")
	}
}

func TestReadFileLocal(t *testing.T) {
	buf, err := ReadFile("bar")
	if err != nil {
		t.Fatal(err)
	}

	if string(buf) != "Now is the time for all good men to take their place in the confederation of planets" {
		t.Fatal("file contents mismatch")
	}
}

func TestReadFileUnknown(t *testing.T) {
	_, err := ReadFile("unknownfile")
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestStat(t *testing.T) {
	_, err := Stat("testdata")
	if err != nil {
		t.Fatal(err)
	}
}

func TestStatLocal(t *testing.T) {
	fi, err := Stat("bar")
	if err != nil {
		t.Fatal(err)
	}

	if fi.Name() != "bar" {
		t.Fatal("filename mismatch")
	}

	file := genfiles["bar"]

	if int(fi.Size()) != int(file.Length) {
		t.Fatal("file size mismatch")
	}

	if fi.Mode() != 0444 {
		t.Fatal("file mode mismatch")
	}

	if fi.ModTime() != time.Unix(file.Modsec, file.Modnano) {
		t.Fatal("file time mismatch")
	}

	if fi.IsDir() != false {
		t.Fatal("IsDir() is true")
	}

	if fi.Sys() != nil {
		t.Fatal("Sys is not nil")
	}
}

func TestStatUnknown(t *testing.T) {
	_, err := Stat("unknown")
	if err == nil {
		t.Fatal("expected an error")
	}

	if os.IsNotExist(err) == false {
		t.Fatalf("expected os.ErrNotExist got \"%v\"", err)
	}
}
