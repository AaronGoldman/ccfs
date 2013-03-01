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

func getBlob(hash [32]byte) (data []byte, err error) {
	filepath := fmt.Sprintf("../blobs/%s", hash)
	filedata, err := ioutil.ReadFile(filepath)
	if err == nil {
		dataBlock, _ := pem.Decode(filedata)
		data = dataBlock.Bytes
	} else {
		data = []byte("")
	}
	return
	//	data = []byte{0xfa,
	//	0x84,0xff,0xaa,0xe4,0xd6,0x5f,0x49,0x67,0x67,0x8c,0x95,0xb7,0xf9,0x6d,0x61,
	//	0xe0,0x67,0x2b,0x77,0xe2,0x67,0x26,0x78,0x44,0x44,0x95,0x24,0x55,0x07,0x56,0xb4}
	//	err = errors.New("Not yet implimented")
	//	return
}

func getTag(hkid [32]byte) (data []byte, err error) {
	filepath := fmt.Sprintf("../tags/%s/%s/%s", hkid, name, version)
	filedata, err := ioutil.ReadFile(filepath)
	if err == nil {
		dataBlock, _ := pem.Decode(filedata)
		data = dataBlock.Bytes
	} else {
		data = []byte("")
	}
	return
}

func getCommit(hkid [32]byte) (data []byte, err error) {
	filepath := fmt.Sprintf("../commits/%s/%s", hkid, version)
	filedata, err := ioutil.ReadFile(filepath)
	if err == nil {
		dataBlock, _ := pem.Decode(filedata)
		data = dataBlock.Bytes
	} else {
		data = []byte("")

	}
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
