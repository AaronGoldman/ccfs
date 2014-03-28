// fileSystem_test.go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
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
	//t.Skip("skip create domain")
	log.SetFlags(log.Lshortfile)
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "../mountpoint")
	os.MkdirAll(fmt.Sprintf("%s/TestPostNewTag", path), 0777)
	list, _ := ioutil.ReadDir(fmt.Sprintf("%s/TestPostNewTag", path))
	if len(list) != 0 {
		t.Errorf("Folder not empty")
	}
	cmd := exec.Command("./src", "-createDomain=true", fmt.Sprintf("-path=%s/TestPostNewTag", path))
	b, err := cmd.CombinedOutput()
	fmt.Printf("%s", b)
	if err != nil {
		t.Errorf("Create Domain Errored - %s \n", err)
	}
}
func TestCLCreateRepo(t *testing.T) {
	//t.Skip("skip create repo")
	log.SetFlags(log.Lshortfile)
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "../mountpoint")
	os.MkdirAll(fmt.Sprintf("%s/TestPostNewCommit", path), 0777)
	list, _ := ioutil.ReadDir(fmt.Sprintf("%s/TestPostNewCommit", path))
	if len(list) != 0 {
		t.Errorf("Folder not empty")
	}
	cmd := exec.Command("./src", "-createRepository=true", fmt.Sprintf("-path=%s/TestPostNewCommit", path))
	b, err := cmd.CombinedOutput()
	fmt.Printf("%s", b)
	if err != nil {
		t.Errorf("Create Repository Errored - %s \n", err)
	}
}
func TestCLInsertDomain(t *testing.T) {
	t.Skip("skip insert domain")
	log.SetFlags(log.Lshortfile)
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "../mountpoint")
	os.MkdirAll(fmt.Sprintf("%s/TestPostTag", path), 0777)
	list, _ := ioutil.ReadDir(fmt.Sprintf("%s/TestPostTag", path))
	if len(list) != 0 {
		t.Errorf("Folder not empty")
	}
	domainHKID := hkidFromDString("2990018983336786774600773215435487572040278176087795322342464389288172846099779527029312056191767811453586805184323598252008160483472900619326359336945638850", 10)
	cmd := exec.Command("./src", "-insertDomain=true", fmt.Sprintf("-path=%s/TestPostTag", path), fmt.Sprintf("-hkid=\"%s\"", domainHKID))
	b, err := cmd.CombinedOutput()
	fmt.Printf("%s", b)
	if err != nil {
		t.Errorf("Insert Domain Errored - %s \n", err)
	}
}
func TestCLInsertRepo(t *testing.T) {
	t.Skip("skip insert repo")
	log.SetFlags(log.Lshortfile)
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "../mountpoint")
	os.MkdirAll(fmt.Sprintf("%s/TestPostCommit", path), 0777)
	list, _ := ioutil.ReadDir(fmt.Sprintf("%s/TestPostCommit", path))
	if len(list) != 0 {
		t.Errorf("Folder not empty")
	}
	repoHKID := hkidFromDString("5824205648082772934729637225579799788842383870921308642349398134394915270944497186356984254449560747108115423811117570014383411154383531617434061770576416540", 10)
	cmd := exec.Command("./src", "-insertRepository=true", fmt.Sprintf("-path=%s/TestPostCommit", path), fmt.Sprintf("-hkid=\"%s\"", repoHKID))
	b, err := cmd.CombinedOutput()
	fmt.Printf("%s", b)
	if err != nil {
		t.Errorf("Insert Repository Errored - %s \n", err)
	}
}
func TestWriteFileSystemInterface(t *testing.T) {
	log.SetFlags(log.Lshortfile)
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
	log.SetFlags(log.Lshortfile)
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
func BenchmarkReadFileSystemInterface(b *testing.B) {
	var answerKey = []struct {
		fileName    string
		fileContent string
	}{
		{"../mountpoint/TestPostBlob", "TestPostData"},
		{"../mountpoint/TestPostCommit/TestPostBlob", "TestPostCommitBlobData"},
		{"../mountpoint/TestPostTag/TestPostBlob", "TestPostTagBlobData"},
		{"../mountpoint/TestPostList/TestPostList2/TestPostBlob", "TestPostListListBlobData"},
	}
	for i := 0; i < b.N; i++ {
		for _, answer := range answerKey {
			data, err := ioutil.ReadFile(answer.fileName)
			if err != nil {
				b.Errorf("Benchmark File Interface Failed - %s \n", err)
			}
			if string(data) != answer.fileContent {
				b.Errorf("Filepath: %s\n Expected: %s\n Actual: %s", answer.fileName, answer.fileContent, data)
			}
		}
	}
}
