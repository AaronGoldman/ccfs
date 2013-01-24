package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
)

func main() {
	err := filepath.Walk("./", visit)
	if err != nil {
		panic(err)
	}
}

//fromchannaltofile()

func visit(path string, f os.FileInfo, err error) error {
	if f.IsDir() == false {
		hbuf := hashfile(path)
		fmt.Printf("%v %s\n", hex.EncodeToString(hbuf), path)
	}
	return nil
}

func hashfile(filepath string) []byte {
	fi, err := os.Open(filepath)

	if err != nil {
		fmt.Printf("%v", err)
		return nil
		//panic(err)
	}
	defer fi.Close()

	buf := make([]byte, 1024)
	var h hash.Hash = sha256.New()
	for {
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Printf("%v", err)
			return nil
			//panic(err)
		}
		if n == 0 {
			break
		}
		h.Write(buf[:n])
	}
	sbuf := make([]byte, 0)
	sbuf = h.Sum(sbuf)
	return sbuf
}
