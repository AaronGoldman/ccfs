package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
)

type hashfinder struct{}

func (hashfinder) walk() {
	err := filepath.Walk("./", hashfinder{}.visit)
	if err != nil {
		log.Panic(err)
	}
}

//fromchannaltofile()

func (hashfinder) visit(path string, f os.FileInfo, err error) error {
	if f.IsDir() == false {
		hbuf := hashfinder{}.hashfile(path)
		log.Printf("%v %s\n", hex.EncodeToString(hbuf), path)
	}
	return nil
}

func (hashfinder) hashfile(filepath string) []byte {
	fi, err := os.Open(filepath)

	if err != nil {
		log.Printf("%v", err)
		return nil
		//log.Panic(err)
	}
	defer fi.Close()

	buf := make([]byte, 1024)
	var h = sha256.New()
	for {
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("%v", err)
			return nil
			//log.Panic(err)
		}
		if n == 0 {
			break
		}
		h.Write(buf[:n])
	}
	var sbuf []byte
	sbuf = h.Sum(sbuf)
	return sbuf
}
