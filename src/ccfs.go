package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
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

func get(hkid []byte, path string) (b blob, err error) {
	typeString := "commit"
	objecthash := hkid
	err = nil
	nameSegments := []string{"", path}

	for {
		if typeString == "blob" || nameSegments[1] == "" {
			b, err = GetBlob(objecthash)
			return
		}
		if typeString == "list" {
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			l, _ := GetList(objecthash)
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
		}
		if typeString == "tag" {
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			t, err := GetTag(objecthash, nameSegments[0])
			if !t.Verifiy() {
				b = nil
				err = errors.New("Tag Verifiy Failed")
				return b, err
			}
			typeString = t.TypeString
			objecthash = t.HashBytes
		}
		if typeString == "commit" {
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			c, err := GetCommit(objecthash)
			if !c.Verifiy() {
				b = nil
				err = errors.New("Commit Verifiy Failed")
				return b, err
			}
			l, err := GetList(c.listHash)
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
		}
	}
	b = nil
	return
}

func post(hkid []byte, path string, b blob) (err error) {
	//follow get prosedres, push object as there pased, drop at tag/commit
	//isue the blob regen list isue list regen tag/commit isue tag/commit
	typeString := "commit"
	objecthash := hkid
	regenlist := [][]byte{}
	regenpath := []string{}
	err = nil
	nameSegments := []string{"", path}
	for {
		if typeString == "blob" || nameSegments[1] == "" {
			b, err = GetBlob(objecthash)
			if err == nil {

			}
			//if pulled blob maches no update nesisary
		}
		if typeString == "list" {
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			l, err := GetList(objecthash)
			if err == nil {
				regenlist = append(regenlist, l.Bytes())
				regenpath = append(regenpath, nameSegments[0])
			} else {
				regen(regenlist, hkid, b)
			}
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
		}
		if typeString == "tag" {
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			regenlist = nil
			regenpath = []string{}
			t, err := GetTag(objecthash, nameSegments[0])
			if !t.Verifiy() {
				b = nil
				err = errors.New("Tag Verifiy Failed")
				return err
			}
			typeString = t.TypeString
			objecthash = t.HashBytes
		}
		if typeString == "commit" {
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			regenlist = nil
			regenpath = []string{}
			c, err := GetCommit(objecthash)
			if !c.Verifiy() {
				b = nil
				err = errors.New("Commit Verifiy Failed")
				return err
			}
			l, err := GetList(c.listHash)
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
		}
	}
	b = nil
	return
}

func regen(objectsToRegen [][]byte, hkid []byte, b blob)(error) {
	return nil
}

func initRepo(hkid []byte, path string) (repoHkid []byte) {
	return nil
}

func initDomain(hkid []byte, path string) (domainHkid []byte) {
	return nil
}
