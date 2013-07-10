package main

import (
	"bufio"
	golist "container/list"
	"crypto/elliptic"
	"crypto/sha256"
	"fmt"
	"os"
)

func main() {
	//go BlobServerStart()
	//hashfindwalk()
	in := bufio.NewReader(os.Stdin)
	input, err := in.ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Println(input)
	return
}

func InsertRopo(thereHKID HKID, myHKID HKID, path string) {

}

func InsertDomain(thereHKID HKID, myHKID HKID, path string) {

}

func InitRepo(hkid HKID) error {
	//InitCommit()
	c := NewCommit(sha256.New().Sum(nil), hkid)
	err := PostCommit(c)
	return err
}
func InitDomain(hkid HKID, nameSegment string) error {
	//GenHKID()
	t := NewTag(sha256.New().Sum(nil), "blob", nameSegment, hkid)
	err := PostTag(t)
	return err
}

func GenHKID() HKID {
	privkey := KeyGen()
	PostBlob(elliptic.Marshal(privkey.PublicKey.Curve,
		privkey.PublicKey.X, privkey.PublicKey.Y))
	return GenerateHKID(privkey)
}

func regen(objectsToRegen *golist.List, objecthash Hexer, b Byteser) error {
	return nil
}

/*func initRepo(objecthash HKID, path string) (repoHkid HKID) {
	c, privkey := InitCommit()
	err := PostKey(privkey)
	if err != nil {
		panic(err)
	}
	return c.hkid
}*/
