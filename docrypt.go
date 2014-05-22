/*
 * learning how to use the command line user interface package

example run:
$ go run docrypt.go  -k xxxx -e docrypt.go encrypted
$ ls
docrypt.go  encrypted
$ go run docrypt.go  -k xxxx -d encrypted decrypted
$ diff decrypted docrypt.go
$ # no differences


*/
package main

import (
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"flag" // command line option parser
	"fmt"
	"io"
	"os"
)

func main() {
	// setup command line arg requirements
	var kflag *string = flag.String("k", "", "encryption key required")
	var eflag *bool = flag.Bool("e", false, "encrypt flag must be true")
	var dflag *bool = flag.Bool("d", false, "decrypt flag must be true")
	// parse command line
	flag.Parse() // Scans the arg list and sets up flags
	if flag.NArg() < 1 || flag.NArg() > 2 {
		fmt.Print("Error: an input filename is required\n")
		usage()
		os.Exit(1)
	}
	//fmt.Printf("(kflag,eflag,dflag,filename)=(%s,%t,%t,%s)\n",
	//  *kflag,*eflag,*dflag,flag.Arg(0))
	if *kflag == "" || *kflag == "-e" || *kflag == "-d" {
		fmt.Print("encryption key required\n")
		usage()
		os.Exit(2)
	}
	if (*eflag || *dflag) == false {
		fmt.Print("Either -e or -d must be specified\n")
		os.Exit(3)
	}
	var key []byte = make([]byte, 16)
	copy(key, (*kflag)[:])
	if *eflag {
		if flag.NArg() == 2 {
			enzip(flag.Arg(0), key, flag.Arg(1))
		} else {
			enzip(flag.Arg(0), key, flag.Arg(0)+".enzip")
		}
	} else if *dflag {
		if flag.NArg() == 2 {
			dezip(flag.Arg(0), key, flag.Arg(1))
		} else {
			dezip(flag.Arg(0), key, flag.Arg(0)+".dezip")
		}
	}
}

func dezip(srcfile string, key []byte, dstfile string) {

	var iv = make([]byte, 16)

	var f *os.File
	var err error
	if f, err = os.Open(srcfile); err != nil {
		println("os.Open() Error:" + err.Error())
		os.Exit(-1)
	}
	defer f.Close()

	var c cipher.Block
	if c, err = aes.NewCipher(key); err != nil {
		println("NewCipher() Error:" + err.Error())
		os.Exit(-1)
	}

	var r io.Reader
	r = &cipher.StreamReader{S: cipher.NewOFB(c, iv), R: f}
	if r, err = gzip.NewReader(r); err != nil {
		println("NewReader() Error:" + err.Error())
		os.Exit(-1)
	}

	var w *os.File
	if w, err = os.Create(dstfile); err != nil {
		println("os.Create() Error:" + err.Error())
		os.Exit(-1)
	}

	defer w.Close()
	io.Copy(w, r)
}

func enzip(srcfile string, key []byte, dstfile string) {

	var iv = make([]byte, 16)

	var r *os.File
	var err error
	if r, err = os.Open(srcfile); err != nil {
		println("os.Open() Error:" + err.Error())
		os.Exit(-1)
	}

	var w io.Writer
	if w, err = os.Create(dstfile); err != nil {
		println("os.Create() Error:" + err.Error())
		os.Exit(-1)
	}

	var c cipher.Block
	if c, err = aes.NewCipher(key); err != nil {
		println("NewCipher() Error:" + err.Error())
		os.Exit(-1)
	}

	w = &cipher.StreamWriter{S: cipher.NewOFB(c, iv), W: w}

	var w2 *gzip.Writer

	/*
		if w2, err = gzip.NewWriter(w); err != nil {
			println("NewWriter() Error:" + err.Error())
			os.Exit(-1)
		}
	*/
	w2 = gzip.NewWriter(w)

	io.Copy(w2, r)
	w2.Close()
}

func usage() {
	fmt.Print("Usage: docrypt -k keystring -e filename [output filename]\n")
	fmt.Print("      or\n")
	fmt.Print("Usage: docrypt -k keystring -d filename [output filename]\n")
	fmt.Print("where -e means encrypt and -d means decrypt\n")
	fmt.Print("If output filename omitted, \n")
	fmt.Print("then .enzip and .dezip are appended to input filename\n")
}
