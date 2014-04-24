package localfile

import (
	"fmt"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

//localfileservice is an
type localfileservice struct{}

func (lfs localfileservice) PostBlob(b objects.Blob) (err error) {
	filepath := fmt.Sprintf("bin/blobs/%s", b.Hash().Hex())
	err = os.MkdirAll("bin/blobs", 0764)
	err = ioutil.WriteFile(filepath, b.Bytes(), 0664)
	return
}
func (lfs localfileservice) PostTag(t objects.Tag) (err error) {
	filepath := fmt.Sprintf("bin/tags/%s/%s/%d", t.Hkid().Hex(),
		t.NameSegment, t.Version)
	dirpath := fmt.Sprintf("bin/tags/%s/%s", t.Hkid().Hex(),
		t.NameSegment)
	err = os.MkdirAll(dirpath, 0764)
	err = ioutil.WriteFile(filepath, t.Bytes(), 0664)
	return
}
func (lfs localfileservice) PostCommit(c objects.Commit) (err error) {
	filepath := fmt.Sprintf("bin/commits/%s/%d", c.Hkid().Hex(),
		c.Version)
	dirpath := fmt.Sprintf("bin/commits/%s", c.Hkid().Hex())
	err = os.MkdirAll(dirpath, 0764)
	err = ioutil.WriteFile(filepath, c.Bytes(), 0664)
	return
}
func (lfs localfileservice) PostKey(p *objects.PrivateKey) (err error) {
	err = os.MkdirAll("bin/keys", 0700)
	filepath := fmt.Sprintf("bin/keys/%s", p.Hkid().Hex())
	err = ioutil.WriteFile(filepath, objects.PrivateKey(*p).Bytes(), 0600)
	return
}

func (lfs localfileservice) GetBlob(h objects.HCID) (b objects.Blob, err error) {
	//ToDo Validate input
	if h == nil {
		return nil, fmt.Errorf("[localfileservice] GetBlob() HCID is nil")
	}
	filepath := fmt.Sprintf("bin/blobs/%s", h.Hex())
	//log.Printf("Filepath: %v", filepath)
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("\n\t%v\n", err)
	}
	//build object
	b = objects.Blob(data)
	return b, err
}
func (lfs localfileservice) GetCommit(h objects.HKID) (c objects.Commit, err error) {
	//Validate input
	matches, err := filepath.Glob(fmt.Sprintf("bin/commits/%s/*", h.Hex()))
	filepath := lfs.latestVersion(matches)
	//log.Printf("Filepath: %v", filepath)
	data, err := ioutil.ReadFile(filepath)
	if err == nil {
		c, _ = objects.CommitFromBytes(data)
	}
	return c, err
}
func (lfs localfileservice) GetTag(h objects.HKID, namesegment string) (t objects.Tag, err error) {
	//ToDo Validate input
	matches, err := filepath.Glob(fmt.Sprintf("bin/tags/%s/%s/*",
		h.Hex(), namesegment))
	filepath := lfs.latestVersion(matches)
	log.Printf("Filepath: %v", filepath)
	data, err := ioutil.ReadFile(filepath)
	if err == nil {
		t, _ = objects.TagFromBytes(data)
	}
	return t, err
}
func (lfs localfileservice) GetKey(h objects.HKID) (objects.Blob, error) {
	filepath := fmt.Sprintf("bin/keys/%s", h.Hex())
	//log.Printf("Filepath: %v", filepath)
	filedata, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Println(err)
	}
	return filedata, err
}

func (lfs localfileservice) latestVersion(matches []string) string {
	match := ""
	for _, element := range matches {
		if match < element {
			match = element
		}
	}
	return match
}

var Instance localfileservice = localfileservice{}

func init() {
	services.Registerblobgeter(Instance)
	services.Registerblobposter(Instance)
	services.Registercommitgeter(Instance)
	services.Registercommitposter(Instance)
	services.Registertaggeter(Instance)
	services.Registertagposter(Instance)
	services.Registerkeygeter(Instance)
	services.Registerkeyposter(Instance)
}