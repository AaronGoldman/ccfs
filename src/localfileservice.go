package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

//func localfileservice_GetBlob(hash HCID) (b blob, err error) {
//	//ToDo Validate input
//	filepath := fmt.Sprintf("../blobs/%s", hash.Hex())
//	data, err := ioutil.ReadFile(filepath)
//	if err != nil {
//		log.Printf("\n\t%v\n", err)
//	}
//	//build object
//	b = blob(data)
//	//log.Printf("\n\t%v\n", string(b))
//	return b, err
//}

//func localfileservice_blobgeter(datach chan blob, errorch chan error, h HCID) {
//	b, err := localfileservice_GetBlob(h)
//	if err == nil {
//		datach <- b
//	} else {
//		errorch <- err
//	}
//}

func PostBlob(b blob) (err error) {
	filepath := fmt.Sprintf("../blobs/%s", b.Hash().Hex())
	err = os.MkdirAll("../blobs", 0764)
	err = ioutil.WriteFile(filepath, b.Bytes(), 0664)
	return
}

//func localfileservice_GetTag(hkid HKID, nameSegment string) (t tag, err error) {
//	//ToDo Validate input
//	matches, err := filepath.Glob(fmt.Sprintf("../tags/%s/%s/*",
//		hex.EncodeToString(hkid), nameSegment))
//	filepath := latestVersion(matches)
//	data, err := ioutil.ReadFile(filepath)
//	if err == nil {
//		t, _ = TagFromBytes(data)
//	}
//	return
//}

//func localfileservice_taggeter(datach chan tag, errorch chan error, h HKID, namesegment string) {
//	t, err := localfileservice_GetTag(h, namesegment)
//	if err == nil {
//		datach <- t
//	} else {
//		errorch <- err
//	}
//}

func PostTag(t tag) (err error) {
	filepath := fmt.Sprintf("../tags/%s/%s/%d", t.hkid.Hex(),
		t.nameSegment, t.version)
	dirpath := fmt.Sprintf("../tags/%s/%s", t.hkid.Hex(),
		t.nameSegment)
	err = os.MkdirAll(dirpath, 0764)
	err = ioutil.WriteFile(filepath, t.Bytes(), 0664)
	return
}

//func localfileservice_GetCommit(hkid HKID) (c commit, err error) {
//	//Validate input
//	matches, err := filepath.Glob(fmt.Sprintf("../commits/%s/*",
//		hex.EncodeToString(hkid)))
//	filepath := latestVersion(matches)

//	data, err := ioutil.ReadFile(filepath)
//	//log.Printf("%v\n", err)
//	if err == nil {
//		c, _ = CommitFromBytes(data)
//	}

//	return
//}

//func localfileservice_commitgeter(datach chan commit, errorch chan error, h HKID) {
//	c, err := localfileservice_GetCommit(h)
//	if err == nil {
//		datach <- c
//	} else {
//		errorch <- err
//	}
//}

func PostCommit(c commit) (err error) {
	filepath := fmt.Sprintf("../commits/%s/%d", c.hkid.Hex(),
		c.version)
	dirpath := fmt.Sprintf("../commits/%s", c.hkid.Hex())
	err = os.MkdirAll(dirpath, 0764)
	err = ioutil.WriteFile(filepath, c.Bytes(), 0664)
	return
}

//func localfileservice_GetKey(h HKID) (data blob, err error) {
//	//log.Println("localfileservice", h)
//	filepath := fmt.Sprintf("../keys/%s", hex.EncodeToString(h))
//	filedata, err := ioutil.ReadFile(filepath)
//	if err != nil {
//		log.Println(err)
//	}
//	return filedata, err
//}

//func localfileservice_keygeter(datach chan blob, errorch chan error, h HKID) {
//	b, err := localfileservice_GetKey(h.Bytes())
//	if err == nil {
//		datach <- b
//	} else {
//		errorch <- err
//	}
//}

func PostKey(p *PrivateKey) (err error) {
	//hkid := blob(elliptic.Marshal(p.PublicKey.Curve,
	//	p.PublicKey.X, p.PublicKey.Y)).Hash()
	err = os.MkdirAll("../keys", 0700)
	filepath := fmt.Sprintf("../keys/%s", p.Hkid().Hex())
	err = ioutil.WriteFile(filepath, PrivateKey(*p).Bytes(), 0600)
	return
}

func latestVersion(matches []string) string {
	match := ""
	for _, element := range matches {
		if match < element {
			match = element
		}
	}
	return match
}

type localfileservice struct{}

func (localfileservice) getBlob(h HCID) (b blob, err error) {
	//ToDo Validate input
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
func (localfileservice) getCommit(h HKID) (c commit, err error) {
	//Validate input
	matches, err := filepath.Glob(fmt.Sprintf("../commits/%s/*", h.Hex()))
	filepath := latestVersion(matches)

	data, err := ioutil.ReadFile(filepath)
	//log.Printf("%v\n", err)
	if err == nil {
		c, _ = CommitFromBytes(data)
	}
	return c, err
}
func (localfileservice) getTag(h HKID, namesegment string) (t tag, err error) {
	//ToDo Validate input
	matches, err := filepath.Glob(fmt.Sprintf("../tags/%s/%s/*",
		h.Hex(), namesegment))
	filepath := latestVersion(matches)
	data, err := ioutil.ReadFile(filepath)
	if err == nil {
		t, _ = TagFromBytes(data)
	}
	return t, err
}
func (localfileservice) getKey(h HKID) (blob, error) {
	filepath := fmt.Sprintf("../keys/%s", h.Hex())
	filedata, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Println(err)
	}
	return filedata, err
}

var localfileserviceInstance localfileservice = localfileservice{}

//func init() {
//	localfileserviceInstance = localfileservice{}
//}
