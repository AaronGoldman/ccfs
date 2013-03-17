package main

import (
	"encoding/pem"
	//	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func main() {
	go BlobServerStart()
	//commitgentest()
	//taggentest()
	//hashfindwalk()
	return
}

func getBlob(hash [32]byte) (data []byte, err error) {
	//Validate input
	filepath := fmt.Sprintf("../blobs/%s", hash)
	data, err = ioutil.ReadFile(filepath)
	return
}

func postBlob(data []byte) (err error) {
	return
}

func getTag(hkid [32]byte, nameSegment string) (data []byte, err error) {
	//Validate input
	//matches []string, err error
	matches, err := filepath.Glob(fmt.Sprintf("../tags/%s/%s/*", hkid, nameSegment))
	version := latestVersion(matches)
	filepath := fmt.Sprintf("../tags/%s/%s/%s", hkid, nameSegment, version)
	data, err = ioutil.ReadFile(filepath)
	return
}

func postTag(data []byte, nameSegment string, version string) (err error) {
	return
}

func getCommit(hkid [32]byte) (data []byte, err error) {
	//Validate input
	matches, err := filepath.Glob(fmt.Sprintf("../commits/%s/*", hkid))
	version := latestVersion(matches)
	filepath := fmt.Sprintf("../commits/%s/%s", hkid, version)
	data, err = ioutil.ReadFile(filepath)
	return
}

func postCommit(data []byte, version string) (err error) {
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

func postKey(data []byte) (err error) {
	return
}

func latestVersion(matches []string) (match string) {
	match = ""
	for _, element := range matches {
		if match < element {
			match = element
		}
	}
	return
}
