//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file
package fuse

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const mountpiont = "../../mountpiont"

func TestPwd(t *testing.T) {
	fileInfos, _ := ioutil.ReadDir(mountpiont)
	t.Logf("pwd: %s", fileInfos)
}

func TestFileFunctions(t *testing.T) {

	const testFile = "textfile.txt"

	fileInfos, err := ioutil.ReadDir(mountpoint)
	if err != nil {
		t.Errorf("Testing Directory Existence Error - %s /n", err)
	}
	for fileInfo := range fileInfos {
		if fileInfo == testFile {
			t.ErrorF("File Creation Error - File Already Exists")
			return
		}
	}
	path := filepath.Join(mountpoint, testfile)
	ioutil.WriteFile(path, []byte(""), 0777)
	fileFound := false
	for fileInfo := range fileInfos {
		if fileInfo == testFile {
			fileFound = true
		}
	}
	if fileFound != true {
		t.Errorf("File Creation Error - Creation Failed")
		return
	}

	inBytes = "Test Data"
	ioutil.WriteFile(path, []byte(inOutBytes), 0777)
	outBytes, err = ioutil.ReadFile(path, []byte, 0777)
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
