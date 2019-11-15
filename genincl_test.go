package genincl_test

import (
	"os"
	"testing"
	"time"

	. "github.com/dbulkow/genincl"
)

var files = map[string]File{
	"testdata": File{
		Modsec:  1573777556,
		Modnano: 901181155,
		Data: []byte{
			0x4e, 0x6f, 0x77, 0x20, 0x69, 0x73, 0x20, 0x74, 0x68,
			0x65, 0x20, 0x74, 0x69, 0x6d, 0x65, 0x20, 0x66, 0x6f,
			0x72, 0x20, 0x61, 0x6c, 0x6c, 0x20, 0x67, 0x6f, 0x6f,
			0x64, 0x20, 0x6d, 0x65, 0x6e, 0x20, 0x74, 0x6f, 0x20,
			0x74, 0x61, 0x6b, 0x65, 0x20, 0x74, 0x68, 0x65, 0x69,
			0x72, 0x20, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x20, 0x69,
			0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x63, 0x6f, 0x6e,
			0x66, 0x65, 0x64, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
			0x6e, 0x20, 0x6f, 0x66, 0x20, 0x70, 0x6c, 0x61, 0x6e,
			0x65, 0x74, 0x73,
		},
	},
	"bar": File{
		Modsec:  1573777556,
		Modnano: 901181155,
		Data: []byte{
			0x4e, 0x6f, 0x77, 0x20, 0x69, 0x73, 0x20, 0x74, 0x68,
			0x65, 0x20, 0x74, 0x69, 0x6d, 0x65, 0x20, 0x66, 0x6f,
			0x72, 0x20, 0x61, 0x6c, 0x6c, 0x20, 0x67, 0x6f, 0x6f,
			0x64, 0x20, 0x6d, 0x65, 0x6e, 0x20, 0x74, 0x6f, 0x20,
			0x74, 0x61, 0x6b, 0x65, 0x20, 0x74, 0x68, 0x65, 0x69,
			0x72, 0x20, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x20, 0x69,
			0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x63, 0x6f, 0x6e,
			0x66, 0x65, 0x64, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
			0x6e, 0x20, 0x6f, 0x66, 0x20, 0x70, 0x6c, 0x61, 0x6e,
			0x65, 0x74, 0x73,
		},
	},
}

func init() {
	Register(files)
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

	if n != len(files["testdata"].Data) {
		t.Fatal(err)
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

	if n != len(files["bar"].Data) {
		t.Fatal(err)
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

	file := files["bar"]

	if int(fi.Size()) != len(file.Data) {
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
