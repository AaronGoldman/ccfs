package main

import (
	"bufio"
	golist "container/list"
	"encoding/hex"
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

type HCID []byte

func (hcid HCID) Hex() string {
	return hex.EncodeToString(hcid)
}

type HKID HCID

func (hkid HKID) Hex() string {
	return hex.EncodeToString(hkid)
}

func regen(objectsToRegen *golist.List, objecthash HKID, b blob) error {
	return nil
}

func initRepo(objecthash HKID, path string) (repoHkid []byte) {
	return nil
}

func initDomain(objecthash HCID, path string) (domainHkid []byte) {
	return nil
}
