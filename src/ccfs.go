package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	//"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	go BlobServerStart()
	//hashfindwalk()
	return
}

func GetBlob(hash []byte) (data []byte, err error) {
	//Validate input
	filepath := fmt.Sprintf("../blobs/%s", hex.EncodeToString(hash))
	data, err = ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func PostBlob(b blob) (err error) {
	filepath := fmt.Sprintf("../blobs/%s", hex.EncodeToString(b.Hash()))
	//err = os.MkdirAll(filepath, 0664)
	err = ioutil.WriteFile(filepath, b.Bytes(), 0664)
	return
	//	data = []byte{0xfa,
	//	0x84,0xff,0xaa,0xe4,0xd6,0x5f,0x49,0x67,0x67,0x8c,0x95,0xb7,0xf9,0x6d,0x61,
	//	0xe0,0x67,0x2b,0x77,0xe2,0x67,0x26,0x78,0x44,0x44,0x95,0x24,0x55,0x07,0x56,0xb4}
	//	err = errors.New("Not yet implimented")
	//	return
}

func GetTag(hkid []byte, nameSegment string) (data []byte, err error) {
	//Validate input
	//matches []string, err error
	matches, err := filepath.Glob(fmt.Sprintf("../tags/%s/%s/*", hkid,
		nameSegment))
	filepath := latestVersion(matches)
	//filepath := fmt.Sprintf("../tags/%s/%s/%s", hkid, nameSegment, version)
	data, err = ioutil.ReadFile(filepath)
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

func GetCommit(hkid []byte) (data []byte, err error) {
	//Validate input
	matches, err := filepath.Glob(fmt.Sprintf("../commits/%s/*",
		hex.EncodeToString(hkid)))
	filepath := latestVersion(matches)
	data, err = ioutil.ReadFile(filepath)
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
	//dataBlock, _ := pem.Decode(filedata)
	//	data = dataBlock.Bytes
	return filedata, err
}

func PostKey(p *ecdsa.PrivateKey) (err error) {
	hkid := blob(elliptic.Marshal(p.PublicKey.Curve,
		p.PublicKey.X, p.PublicKey.Y)).Hash()
	filepath := fmt.Sprintf("../keys/%s", hex.EncodeToString(hkid))
	//err = os.MkdirAll(filepath, 0600)
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
