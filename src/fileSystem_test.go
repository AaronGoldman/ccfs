// fileSystem_test.go
package main

import (
	"io/ioutil"
	//"path/filename"
	"fmt"
	"os"
	"os/exec"
	"path"
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
	cmd := exec.Command("./src", "-createDomain=true", "-path=\"TestPostNewTag\"")
	b, err := cmd.CombinedOutput()
	fmt.Printf("%s", b)
	if err != nil {
		t.Errorf("Create Domain Errored - %s \n", err)
	}
}
func TestCLCreateRepo(t *testing.T) {
	cmd := exec.Command("./src", "-createRepository=true", "-path=\"TestPostNewCommit\"")
	b, err := cmd.CombinedOutput()
	fmt.Printf("%s", b)
	if err != nil {
		t.Errorf("Create Repository Errored - %s \n", err)
	}
}
func TestCLInsertDomain(t *testing.T) {
	cmd := exec.Command("./src", "-insertDomain=true", "-path=\"TestPostTag\"", fmt.Sprintf("-hkid=\"%s\"", benchmarkTagHkid))
	b, err := cmd.CombinedOutput()
	fmt.Printf("%s", b)
	if err != nil {
		t.Errorf("Insert Domain Errored - %s \n", err)
	}
}
func TestCLInsertRepo(t *testing.T) {
	cmd := exec.Command("./src", "-insertRepository=true", "-path=\"TestPostCommit\"", fmt.Sprintf("-hkid=\"%s\"", benchmarkCommitHkid))
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
