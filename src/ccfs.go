package main

import (
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
)

func main() {
	go BlobServerStart()
	//commitgentest()
	//taggentest()
	//hashfindwalk()

}

func getBlob(hkid [32]byte) (data []byte, err error) {
	data = []byte("testing")
	err = errors.New("Not yet implimented")
	return
}

func getKey(hkid [32]byte) (data []byte, err error) {
	filepath := fmt.Sprintf("../keys/%s", hkid)
	filedata, err := ioutil.ReadFile(filepath)
	if err == nil {
		dataBlock, _ := pem.Decode(filedata)
		data = dataBlock.Bytes
	} else {
		data = []byte("")
	}
	return
}
