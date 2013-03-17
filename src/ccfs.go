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
	//Validate input
	filepath := fmt.Sprintf("../blobs/%s", hash)
	data, err := ioutil.ReadFile(filepath)
	return
}

func getTag(hkid [32]byte , nameSegment string) (data []byte, err error) {
	//Validate input
	filepath := fmt.Sprintf("../tags/%s/%s/%s",hkid,nameSegment,version)
	data, err := ioutil.ReadFile(filepath)
	return
}

func getCommit(hkid [32]byte) (data []byte, err error) {
	//Validate input
	filepath := fmt.Sprintf("../commits/%s/%s",hkid,version)
	data, err := ioutil.ReadFile(filepath)
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
