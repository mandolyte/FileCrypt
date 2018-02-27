package main

import (
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var key = flag.String("key", "", "encryption key; required")
var eflag = flag.Bool("e", false, "encrypt flag")
var dflag = flag.Bool("d", false, "decrypt flag")
var input = flag.String("i", "", "input file to be encrypted; required")
var output = flag.String("o", "", "output file to be decrypted; required")
var help = flag.Bool("help", false, "show usage message")

func main() {
	flag.Parse() // Scans the arg list and sets up flags

	if *help {
		usage("Usage:")
	}

	if *key == "" {
		usage("Encryption key required\n")
	}

	if (*eflag || *dflag) == false {
		usage("Either -e or -d must be specified\n")
	}

	if *output == "" {
		usage("Output filename to write encrypted or decrypted file is required")
	}

	if *input == "" {
		usage("Input filename to encrypt or decrypt is required")
	}

	keydata := make([]byte, 32)
	copy(keydata, (*key)[:])
	if *eflag {
		enzip(*input, keydata, *output)
	} else {
		dezip(*input, keydata, *output)
	}
}

func dezip(srcfile string, key []byte, dstfile string) {

	var iv = make([]byte, 16)

	var f *os.File
	var err error
	if f, err = os.Open(srcfile); err != nil {
		log.Fatal("os.Open() Error:" + err.Error())
	}
	defer f.Close()

	var c cipher.Block
	if c, err = aes.NewCipher(key); err != nil {
		log.Fatal("NewCipher() Error:" + err.Error())
	}

	var r io.Reader
	r = &cipher.StreamReader{S: cipher.NewOFB(c, iv), R: f}
	if r, err = gzip.NewReader(r); err != nil {
		log.Fatal("NewReader() Error:" + err.Error())
	}

	var w *os.File
	if w, err = os.Create(dstfile); err != nil {
		log.Fatal("os.Create() Error:" + err.Error())
	}

	_, err = io.Copy(w, r)
	if err != nil {
		log.Fatal("io.Copy() error:" + err.Error())
	}

	err = w.Close()
	if err != nil {
		log.Fatal("os.File.Close() error:" + err.Error())
	}
}

func enzip(srcfile string, key []byte, dstfile string) {

	var iv = make([]byte, 16)

	var r *os.File
	var err error
	if r, err = os.Open(srcfile); err != nil {
		log.Fatal("os.Open() Error:" + err.Error())
	}

	var w io.Writer
	if w, err = os.Create(dstfile); err != nil {
		log.Fatal("os.Create() Error:" + err.Error())
	}

	var c cipher.Block
	if c, err = aes.NewCipher(key); err != nil {
		log.Fatal("NewCipher() Error:" + err.Error())
	}

	w = &cipher.StreamWriter{S: cipher.NewOFB(c, iv), W: w}

	w2 := gzip.NewWriter(w)

	_, err = io.Copy(w2, r)
	if err != nil {
		log.Fatal("io.Copy() error:" + err.Error())
	}

	err = w2.Close()
	if err != nil {
		log.Fatal("gzip.Close() error:" + err.Error())
	}
}

func usage(msg string) {
	fmt.Printf("%v\n", msg)
	flag.PrintDefaults()
	fmt.Print("This code will encrypt and write to a gzip file; or decrypt from same\n")
	fmt.Print("Usage: FileCrypt -key keystring -e -i filename_to_encrypt -o encypted_filename\n")
	fmt.Print("      or\n")
	fmt.Print("Usage: FileCrypt -key keystring -d -i filename_to_decrypt -o decripted_filename\n")
	fmt.Print("where -e means encrypt and -d means decrypt\n")
	os.Exit(0)
}
