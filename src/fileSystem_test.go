// fileSystem_test.go
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
)

var answerKey = []struct {
	fileName    string
	fileContent string
}{
	{"../mountpoint/TestPostBlob", "TestPostData"},
	{"../mountpoint/TestPostCommit/TestPostBlob", "TestPostCommitBlobData"},
	{"../mountpoint/TestPostList/TestPostList2/TestPostBlob", "TestPostListListBlobData"},
	{"../mountpoint/TestPostTag/TestPostBlob", "TestPostTagBlobData"},
}

func TestMountRepo(t *testing.T) {
	t.Skip()
}

func TestCLCreateDomain(t *testing.T) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "../mountpoint")
	cmd := exec.Command("./src", "-createDomain=true", fmt.Sprintf("-path=%s/TestPostNewTag", path))
	b, err := cmd.CombinedOutput()
	fmt.Printf("%s", b)
	if err != nil {
		t.Errorf("Create Domain Errored - %s \n", err)
	}
}
func TestCLCreateRepo(t *testing.T) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "../mountpoint")
	cmd := exec.Command("./src", "-createRepository=true", fmt.Sprintf("-path=%s/TestPostNewCommit", path))
	b, err := cmd.CombinedOutput()
	fmt.Printf("%s", b)
	if err != nil {
		t.Errorf("Create Repository Errored - %s \n", err)
	}
}
func TestCLInsertDomain(t *testing.T) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "../mountpoint")
	cmd := exec.Command("./src", "-insertDomain=true", fmt.Sprintf("-path=%s/TestPostTag", path), fmt.Sprintf("-hkid=\"%s\"", benchmarkTagHkid))
	b, err := cmd.CombinedOutput()
	fmt.Printf("%s", b)
	if err != nil {
		t.Errorf("Insert Domain Errored - %s \n", err)
	}
}
func TestCLInsertRepo(t *testing.T) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "../mountpoint")
	cmd := exec.Command("./src", "-insertRepository=true", fmt.Sprintf("-path=%s/TestPostCommit", path), fmt.Sprintf("-hkid=\"%s\"", benchmarkCommitHkid))
	b, err := cmd.CombinedOutput()
	fmt.Printf("%s", b)
	if err != nil {
		t.Errorf("Insert Repository Errored - %s \n", err)
	}
}
func TestWriteFileSystemInterface(t *testing.T) {
	for _, answer := range answerKey {
		err := os.MkdirAll(path.Dir(answer.fileName), 0764)
		if err != nil {
			t.Errorf("Test File Interface Creation Failed - %s \n", err)
		}
		err = ioutil.WriteFile(answer.fileName, []byte(answer.fileContent), 0664)
		if err != nil {
			t.Errorf("Test File Interface Creation Failed - %s \n", err)
		}
	}
}

func TestReadFileSystemInterface(t *testing.T) {

	for _, answer := range answerKey {
		data, err := ioutil.ReadFile(answer.fileName)
		if err != nil {
			t.Errorf("Test File Interface Failed - %s \n", err)
		}
		if string(data) != answer.fileContent {
			t.Errorf("Filepath: %s\n Expected: %s\n Actual: %s", answer.fileName, answer.fileContent, data)
		}
	}
}

//func TestParseMessage(t *testing.T) {
//	parseMessage("{\"hkid\":\"herp\", \"hcid\":\"derppp\", \"namesegment\":\"idk\", \"type\":\"....\"}")
//}
