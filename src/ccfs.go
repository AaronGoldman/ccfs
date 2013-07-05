package main

import (
	"bufio"
	golist "container/list"
	"crypto/elliptic"
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

func InsertRopo(thereHKID, myHKID, path string) {

}

func InsertDomain(thereHKID, myHKID, path string) {

}

func InitRepo(hkid HKID, path string) HKID {
	return InitCommit()
}
func InitDomain(hkid HKID, path string) HKID {
	return GenHKID()
}

func GenHKID() HKID {
	privkey := KeyGen()
	PostBlob(elliptic.Marshal(privkey.PublicKey.Curve,
		privkey.PublicKey.X, privkey.PublicKey.Y))
	return GenerateHKID(privkey)
}

func regen(objectsToRegen *golist.List, objecthash HKID, b blob) error {
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
