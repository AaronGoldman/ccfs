package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

//localfileservice is an
type localfileservice struct{}

func (lfs localfileservice) PostBlob(b blob) (err error) {
	filepath := fmt.Sprintf("../blobs/%s", b.Hash().Hex())
	err = os.MkdirAll("../blobs", 0764)
	err = ioutil.WriteFile(filepath, b.Bytes(), 0664)
	return
}
func (lfs localfileservice) PostTag(t tag) (err error) {
	filepath := fmt.Sprintf("../tags/%s/%s/%d", t.hkid.Hex(),
		t.nameSegment, t.version)
	dirpath := fmt.Sprintf("../tags/%s/%s", t.hkid.Hex(),
		t.nameSegment)
	err = os.MkdirAll(dirpath, 0764)
	err = ioutil.WriteFile(filepath, t.Bytes(), 0664)
	return
}
func (lfs localfileservice) PostCommit(c commit) (err error) {
	filepath := fmt.Sprintf("../commits/%s/%d", c.hkid.Hex(),
		c.version)
	dirpath := fmt.Sprintf("../commits/%s", c.hkid.Hex())
	err = os.MkdirAll(dirpath, 0764)
	err = ioutil.WriteFile(filepath, c.Bytes(), 0664)
	return
}
func (lfs localfileservice) PostKey(p *PrivateKey) (err error) {
	err = os.MkdirAll("../keys", 0700)
	filepath := fmt.Sprintf("../keys/%s", p.Hkid().Hex())
	err = ioutil.WriteFile(filepath, PrivateKey(*p).Bytes(), 0600)
	return
}

func (lfs localfileservice) GetBlob(h HCID) (b blob, err error) {
	//ToDo Validate input
	if h == nil {
		return nil, fmt.Errorf("[localfileservice] GetBlob() HCID is nil")
	}
	filepath := fmt.Sprintf("../blobs/%s", h.Hex())
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("\n\t%v\n", err)
	}
	//build object
	b = blob(data)
	//log.Printf("\n\t%v\n", string(b))
	return b, err
}
func (lfs localfileservice) GetCommit(h HKID) (c commit, err error) {
	//Validate input
	matches, err := filepath.Glob(fmt.Sprintf("../commits/%s/*", h.Hex()))
	filepath := lfs.latestVersion(matches)

	data, err := ioutil.ReadFile(filepath)
	//log.Printf("%v\n", err)
	if err == nil {
		c, _ = CommitFromBytes(data)
	}
	return c, err
}
func (lfs localfileservice) GetTag(h HKID, namesegment string) (t tag, err error) {
	//ToDo Validate input
	matches, err := filepath.Glob(fmt.Sprintf("../tags/%s/%s/*",
		h.Hex(), namesegment))
	filepath := lfs.latestVersion(matches)
	data, err := ioutil.ReadFile(filepath)
	if err == nil {
		t, _ = TagFromBytes(data)
	}
	return t, err
}
func (lfs localfileservice) GetKey(h HKID) (blob, error) {
	filepath := fmt.Sprintf("../keys/%s", h.Hex())
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

var localfileserviceInstance localfileservice = localfileservice{}

func init() {
	Registerblobgeter(localfileserviceInstance)
	Registerblobposter(localfileserviceInstance)
	Registercommitgeter(localfileserviceInstance)
	Registercommitposter(localfileserviceInstance)
	Registertaggeter(localfileserviceInstance)
	Registertagposter(localfileserviceInstance)
	Registerkeygeter(localfileserviceInstance)
	Registerkeyposter(localfileserviceInstance)
}
