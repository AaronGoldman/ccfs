//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file
package fuse

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const mountpoint = "../../mountpoint"

func TestPwd(t *testing.T) {
	fileInfos, _ := ioutil.ReadDir(mountpoint)
	t.Logf("pwd: %s", fileInfos)
}

func TestWrightFile(t *testing.T) {
	filename := mountpoint + "/TestFile.txt"
	data := []byte("Test File Data")
	perm := os.FileMode(0777)
	err := ioutil.WriteFile(filename, data, perm)
	if err != nil {
		t.Logf("got err %s", err)
		t.Fail()
	}
}

func TestReadFile(t *testing.T) {
	path := mountpoint + "/TestFile.txt"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Logf("got err %s", err)
		t.Fail()
	}
	expectedData := []byte("Test File Data")
	if !(bytes.Equal(data, expectedData)) {
		t.Logf("Expected:%s, Got:%s", expectedData, data)
		t.Fail()
	}
}

func TestFileFunctions(t *testing.T) {

	const testFile = "textfile.txt"

	fileInfos, err := ioutil.ReadDir(mountpoint)
	if err != nil {
		t.Errorf("Testing Directory Existence Error - %s /n", err)
	}
	for _, fileInfo := range fileInfos {
		if fileInfo == testFile {
			t.Errorf("File Creation Error - File Already Exists")
			return
		}
	}
	path := filepath.Join(mountpoint, testfile)
	ioutil.WriteFile(path, []byte(""), 0777)
	fileFound := false
	for _, fileInfo := range fileInfos {
		if fileInfo == testFile {
			fileFound = true
		}
	}
	if fileFound != true {
		t.Errorf("File Creation Error - Creation Failed")
		return
	}

	inBytes := []byte("Test Data")
	ioutil.WriteFile(path, inBytes, 0777)
	outBytes, err := ioutil.ReadFile(path, []byte{}, 0777)
	if err != nil {
		t.errorf("Read File Error - %s", err)
		return
	}
	if !bytes.Equal(inBytes, outBytes) {
		t.errorf("Data Write Error - Data Read Not Equal to Data Written")
		return
	}

	err = os.Remove(path)
	if err != nil {
		t.errorf("File Deletion Error - %s", err)
		return
	}

}
