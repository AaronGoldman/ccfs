package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	go BlobServerStart()
	//hashfindwalk()
	return
}

func GetBlob(hash []byte) (b blob, err error) {
	//ToDo Validate input
	filepath := fmt.Sprintf("../blobs/%s", hex.EncodeToString(hash))
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}
	//build object
	b = blob(data)
	return
}

func PostBlob(b blob) (err error) {
	filepath := fmt.Sprintf("../blobs/%s", hex.EncodeToString(b.Hash()))
	err = ioutil.WriteFile(filepath, b.Bytes(), 0664)
	return
}

func GetTag(hkid []byte, nameSegment string) (t Tag, err error) {
	//ToDo Validate input
	matches, err := filepath.Glob(fmt.Sprintf("../tags/%s/%s/*",
		hex.EncodeToString(hkid), nameSegment))
	filepath := latestVersion(matches)
	data, err := ioutil.ReadFile(filepath)
	//build object
	tagStrings := strings.Split(string(data), ",\n")
	tagHashBytes, _ := hex.DecodeString(tagStrings[0])
	tagTypeString := tagStrings[1]
	tagNameSegment := tagStrings[2]
	tagVersion, _ := strconv.ParseInt(tagStrings[3], 10, 64)
	tagHkid, _ := hex.DecodeString(tagStrings[4])
	tagSignature, _ := hex.DecodeString(tagStrings[5])
	t = Tag{tagHashBytes, tagTypeString, tagNameSegment, tagVersion, tagHkid,
		tagSignature}
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

func GetCommit(hkid []byte) (c commit, err error) {
	//Validate input
	matches, err := filepath.Glob(fmt.Sprintf("../commits/%s/*",
		hex.EncodeToString(hkid)))
	filepath := latestVersion(matches)
	data, err := ioutil.ReadFile(filepath)
	//build object
	commitStrings := strings.Split(string(data), ",\n")
	listHash, _ := hex.DecodeString(commitStrings[0])
	version, _ := strconv.ParseInt(commitStrings[1], 10, 64)
	chkid, _ := hex.DecodeString(commitStrings[2])
	signature, _ := hex.DecodeString(commitStrings[3])
	c = commit{listHash, version, chkid, signature}
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

func GetKey(hkid []byte) (data []byte, err error) {
	filepath := fmt.Sprintf("../keys/%s", hex.EncodeToString(hkid))
	filedata, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return filedata, err
}

func PostKey(p *ecdsa.PrivateKey) (err error) {
	hkid := blob(elliptic.Marshal(p.PublicKey.Curve,
		p.PublicKey.X, p.PublicKey.Y)).Hash()
	filepath := fmt.Sprintf("../keys/%s", hex.EncodeToString(hkid))
	err = ioutil.WriteFile(filepath, KeyBytes(*p), 0600)
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
