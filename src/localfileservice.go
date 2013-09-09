package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func localfileservice_GetBlob(hash HCID) (b blob, err error) {
	//ToDo Validate input
	filepath := fmt.Sprintf("../blobs/%s", hash.Hex())
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("\n\t%v\n", err)
	}
	//build object
	b = blob(data)
	//log.Printf("\n\t%v\n", string(b))
	return
}

func PostBlob(b blob) (err error) {
	filepath := fmt.Sprintf("../blobs/%s", b.Hash().Hex())
	err = os.MkdirAll("../blobs", 0764)
	err = ioutil.WriteFile(filepath, b.Bytes(), 0664)
	return
}

func localfileservice_GetTag(hkid HKID, nameSegment string) (t Tag, err error) {
	//ToDo Validate input
	matches, err := filepath.Glob(fmt.Sprintf("../tags/%s/%s/*",
		hex.EncodeToString(hkid), nameSegment))
	filepath := latestVersion(matches)
	data, err := ioutil.ReadFile(filepath)
	if err == nil {
		t, _ = TagFromBytes(data)
	}
	return
}

func PostTag(t Tag) (err error) {
	filepath := fmt.Sprintf("../tags/%s/%s/%d", hex.EncodeToString(t.hkid),
		t.nameSegment, t.version)
	dirpath := fmt.Sprintf("../tags/%s/%s", hex.EncodeToString(t.hkid),
		t.nameSegment)
	err = os.MkdirAll(dirpath, 0764)
	err = ioutil.WriteFile(filepath, t.Bytes(), 0664)
	return
}

func localfileservice_GetCommit(hkid HKID) (c commit, err error) {
	//Validate input
	matches, err := filepath.Glob(fmt.Sprintf("../commits/%s/*",
		hex.EncodeToString(hkid)))
	filepath := latestVersion(matches)

	data, err := ioutil.ReadFile(filepath)
	//log.Printf("%v\n", err)
	if err == nil {
		c, _ = CommitFromBytes(data)
	}

	return
}

func PostCommit(c commit) (err error) {
	filepath := fmt.Sprintf("../commits/%s/%d", hex.EncodeToString(c.hkid),
		c.version)
	dirpath := fmt.Sprintf("../commits/%s", hex.EncodeToString(c.hkid))
	err = os.MkdirAll(dirpath, 0764)
	err = ioutil.WriteFile(filepath, c.Bytes(), 0664)
	return
}

func localfileservice_GetKey(hkid []byte) (data blob, err error) {
	filepath := fmt.Sprintf("../keys/%s", hex.EncodeToString(hkid))
	filedata, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Panic(err)
	}
	return filedata, err
}

func PostKey(p *ecdsa.PrivateKey) (err error) {
	hkid := blob(elliptic.Marshal(p.PublicKey.Curve,
		p.PublicKey.X, p.PublicKey.Y)).Hash()
	err = os.MkdirAll("../keys", 0700)
	filepath := fmt.Sprintf("../keys/%s", hex.EncodeToString(hkid))
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
