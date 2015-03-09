//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package fuse

import (
	"bytes"
	"io/ioutil"
	"os"
	//"path/filepath"
	"testing"
)

const mountpoint = "../../mountpoint"

/*
func TestList(t *testing.T) {
	t.Logf("List of Available Tests\n")
	t.Logf("-------------------------\n\n")
	t.Logf("TestPwd         - Verify Working Directory\n")
	t.Logf("TestCreateFile  - Creates an empty file in the root directory\n")
	t.Logf("TestWriteFile1   - Write to a file using ioutil\n")
	t.Logf("                - The file will be created in root if it does not exist\n")
	t.Logf("TestWriteFile2  - Same as TestWriteFile except creation of file is\n")
	t.Logf("                - done seperate from Write")
	t.Logf("-------------------------\n\n")
	t.Error("False failure to force list print")
}
*/
func TestPwd(t *testing.T) {
	fileInfos, _ := ioutil.ReadDir(mountpoint)
	t.Logf("pwd: %s", fileInfos)
}

func TestCreateFile(t *testing.T) {
	filename := mountpoint + "/TestCreateFile.txt"
	_, err := os.Create(filename)
	if err != nil {
		t.Logf("Could not create file:", err)
	}
}

//Test file create and write using ioutil
func TestWriteFile1(t *testing.T) {
	filename := mountpoint + "/TestFile1.txt"
	data := []byte("TestFile1 Data")
	perm := os.FileMode(0777)
	err := ioutil.WriteFile(filename, data, perm)
	if err != nil {
		t.Errorf("got err %s", err)
	}
}

//Test file read using ioutil
func TestReadFile1(t *testing.T) {
	path := mountpoint + "/TestFile1.txt"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Logf("got err %s", err)
		t.Fail()
	}
	expectedData := []byte("TestFile1 Data")
	if !(bytes.Equal(data, expectedData)) {
		t.Logf("Expected:%s, Got:%s", expectedData, data)
		t.Fail()
	}
}

//Test file write by creating file using os and writing using ioutil
func TestWriteFile2(t *testing.T) {
	filename := mountpoint + "/TestFile2.txt"
	_, err := os.Create(filename)
	if err != nil {
		t.Logf("Could not create file:", err)
	}
	data := []byte("TestFile2 Data")
	perm := os.FileMode(0777)
	err = ioutil.WriteFile(filename, data, perm)
	if err != nil {
		t.Errorf("got err %s", err)
	}
}

//Test file read using ioutil
func TestReadFile2(t *testing.T) {
	path := mountpoint + "/TestFile2.txt"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Logf("got err %s", err)
		t.Fail()
	}
	expectedData := []byte("TestFile2 Data")
	if !(bytes.Equal(data, expectedData)) {
		t.Logf("Expected:%s, Got:%s", expectedData, data)
		t.Fail()
	}
}

//Test file write using os
func TestWriteFile3(t *testing.T) {
	filename := mountpoint + "/TestFile3.txt"
	file, err := os.Create(filename) //Open(filename)
	if err != nil {
		t.Errorf("Could Not Create File - %s", err)
	}
	data := []byte("TestFile3 Data")
	_, err = file.Write(data)
	if err != nil {
		t.Errorf("Could Not Write To File - %s", err)
	}
}

//Test file read using ioutil
func TestReadFile3(t *testing.T) {
	path := mountpoint + "/TestFile3.txt"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Logf("got err %s", err)
		t.Fail()
	}
	expectedData := []byte("TestFile3 Data")
	if !(bytes.Equal(data, expectedData)) {
		t.Logf("Expected:%s, Got:%s", expectedData, data)
		t.Fail()
	}
}

//file, err = os.Open(filename)
//if err != nil {
//	t.Errorf("Could Not Open File - %s", err)
//}else{
//	fileInfo, fileInfoError:= file.Stat()
//	if(fileInfoError != nil){
//		t.Errorf("Error retrieving file information - %s", err)
//	} else{
//		t.Logf("File Name: %s", fileInfo.Name)
//		t.Logf("File Size: %v", fileInfo.Size)
//		t.Logf("File Mode: %v", fileInfo.Mode)
//	}
//}

/*
func TestFileFunctions(t *testing.T) {

	const testFile = "textfile.txt"

	fileInfos, err := ioutil.ReadDir(mountpoint)
	if err != nil {
		t.Errorf("Testing Directory Existence Error - %s /n", err)
	}
	for _, fileInfo := range fileInfos {
		if fileInfo.Name() == testFile {
			t.Errorf("File Creation Error - File Already Exists")
			return
		}
	}
	path := filepath.Join(mountpoint, testFile)
	ioutil.WriteFile(path, []byte(""), 0777)
	fileFound := false
	for _, fileInfo := range fileInfos {
		if fileInfo.Name() == testFile {
			fileFound = true
		}
	}
	if fileFound != true {
		t.Errorf("File Creation Error - Creation Failed")
		return
	}

	inBytes := []byte("Test Data")
	ioutil.WriteFile(path, inBytes, 0777)
	outBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("Read File Error - %s", err)
		return
	}
	if !bytes.Equal(inBytes, outBytes) {
		t.Errorf("Data Write Error - Data Read Not Equal to Data Written")
		return
	}

	err = os.Remove(path)
	if err != nil {
		t.Errorf("File Deletion Error - %s", err)
		return
	}

}

func TestDeleteFile(t *testing.T) {
	filename := mountpoint + "/TestFile.txt"
	err := os.Remove(filename)
	if err != nil {
		t.Errorf("File could not be deleted:", err)
	}
}*/
