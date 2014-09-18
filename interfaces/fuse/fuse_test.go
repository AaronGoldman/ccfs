package fuse

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

const mountpiont = "../../mountpiont"

func TestPwd(t *testing.T) {
	fileInfos, _ := ioutil.ReadDir(mountpiont)
	t.Logf("pwd: %s", fileInfos)
}

func TestWrightFile(t *testing.T) {
	filename := mountpiont + "/TestFile.txt"
	data := []byte("Test File Data")
	perm := os.FileMode(0777)
	err := ioutil.WriteFile(filename, data, perm)
	if err != nil {
		t.Logf("got err %s", err)
		t.Fail()
	}
}

func TestReadFile(t *testing.T) {
	path := mountpiont + "/TestFile.txt"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Logf("got err %s", err)
		t.Fail()
	}
	expectedData := []byte("Test File Data")
	if !bytes.Equal(data, expectedData) {
		t.Logf("Expected:%s, Got:%s", expectedData, data)
		t.Fail()
	}
}
