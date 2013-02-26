package main

import (
	"errors"
)

func main() {
	go BlobServerStart()
	//commitgentest()
	//taggentest()
	//hashfindwalk()

}

func getBlob(hkid []byte) (data []byte, err error) {
	data = hkid
	err = errors.New("Not yet implimented")
	return hkid, err
}
