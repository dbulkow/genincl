package main

import (
	"bytes"
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"os"
)

const Outfile = "generated.go"

func addfile(filename string, out io.Writer) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "\t\"%s\": File{\n", filename)
	fmt.Fprintf(out, "\t\tModsec:  %d,\n", fi.ModTime().Unix())
	fmt.Fprintf(out, "\t\tModnano: %d,\n", fi.ModTime().Nanosecond())
	fmt.Fprintf(out, "\t\tLength: %d,\n", fi.Size())
	fmt.Fprintf(out, "\t\tFormat: \"zlib\",\n")
	fmt.Fprintf(out, "\t\tData: []byte{\n")

	var b bytes.Buffer

	w := zlib.NewWriter(&b)
	io.Copy(w, file)
	w.Close()

	var p = []byte{0}
	var offset = 0
	for {
		n, err := b.Read(p)
		if err != nil || n == 0 {
			if err == io.EOF {
				fmt.Fprintln(out)
				break
			}

			return err
		}

		switch offset {
		case 9:
			fmt.Fprintln(out)
			offset = 0
			fallthrough
		case 0:
			fmt.Fprint(out, "\t\t\t")
		default:
			fmt.Fprint(out, " ")
		}

		fmt.Fprintf(out, "0x%2.2x,", p[0])

		offset++
	}

	fmt.Fprintln(out, "\t\t},")
	fmt.Fprintln(out, "\t},")

	return nil
}

var header = `package main

import . "github.com/dbulkow/genincl"

func init() {
	Register(genfiles)
}

var genfiles = map[string]File {
`

func main() {
	var outfile = os.Getenv("GENINCL_OUTFILE")

	if outfile == "" {
		outfile = Outfile
	}

	flag.StringVar(&outfile, "outfile", outfile, "Output file name")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "usage: %s [-outfile] <filename or glob>...\n\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	out, err := os.Create(outfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer out.Close()

	fmt.Fprint(out, header)

	names := flag.Args()
	for _, file := range names {
		if err := addfile(file, out); err != nil {
			os.Remove(Outfile)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	fmt.Fprint(out, "}\n")
}
