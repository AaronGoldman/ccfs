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

type hkid []byte

func get(objecthash hkid, path string) (b blob, err error) {
	typeString := "commit"
	//objecthash := hkid
	err = nil
	nameSegments := strings.SplitN(path, "/", 2)
	for {
		switch typeString {
		case "blob":
			if len(nameSegments) < 2 {
				b, err = GetBlob(objecthash)
				return
			}
		case "list":
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			l, _ := GetList(objecthash)
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
		case "tag":
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			t, err := GetTag(objecthash, nameSegments[0])
			if !t.Verifiy() {
				b = nil
				err = errors.New("Tag Verifiy Failed")
				return b, err
			}
			typeString = t.TypeString
			objecthash = t.HashBytes
		case "commit":
			c, err := GetCommit(objecthash)
			if !c.Verifiy() {
				b = nil
				err = errors.New("Commit Verifiy Failed")
				return b, err
			}
			l, err := GetList(c.listHash)
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
		}
		if err != nil{
			return nil, err
		}
	}
	b = nil
	return
}

func post(objecthash hkid, path string, b blob) (err error) {
	typeString := "commit"
	//objecthash := hkid
	err = nil
	nameSegments := strings.SplitN(path, "/", 2)
	regenlist := []blob{}
	regenpath := []string{}
	for {
		switch typeString {
		case "blob":
			if len(nameSegments) < 2 {
				b, err = GetBlob(objecthash)
				return
			}
		case "list":
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			l, _ := GetList(objecthash)
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
			if err == nil {
				regenlist = append(regenlist, l.Bytes())
				regenpath = append(regenpath, nameSegments[0])
			} else {
				regen(regenlist, objecthash, b)
			}
		case "tag":
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			t, _ := GetTag(objecthash, nameSegments[0])
			if !t.Verifiy() {
				return errors.New("Tag Verifiy Failed")
			}
			typeString = t.TypeString
			objecthash = t.HashBytes
		case "commit":
			c, _ := GetCommit(objecthash)
			if !c.Verifiy() {
				return errors.New("Commit Verifiy Failed")
			}
			l, _ := GetList(c.listHash)
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
		}
		if err != nil{
			return err
		}
		if typeString != "list"{
			regenlist = nil
			regenpath = []string{}
		}
	}
	b = nil
	return
}


func regen(objectsToRegen [][]byte, objecthash hkid, b blob) error {
	return nil
}

func initRepo(objecthash hkid, path string) (repoHkid []byte) {
	return nil
}

func initDomain(objecthash hkid, path string) (domainHkid []byte) {
	return nil
}
