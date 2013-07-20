package main

import (
	golist "container/list"
	"errors"
	"fmt"
	"strings"
)

func Post(objecthash Hexer, path string, b Byteser) (hid HID, err error) {
	typeString := "commit"
	//objecthash := hkid
	err = nil
	if path == "" {
		switch b.(type) {
		case blob:
			PostBlob(b.(blob))
			return b.(blob).Hash(), nil
		case list:
			PostList(b.(list))
			return b.(list).Hash(), nil
		case commit:
			PostCommit(b.(commit))
			return b.(commit).Hash(), nil
		case Tag:
			PostTag(b.(Tag))
			return b.(Tag).Hash(), nil
		}
	}
	nameSegments := strings.SplitN(path, "/", 2)
	regenlist := golist.New()
	//regenpath := []string{}
	//for {
	fmt.Printf("%t %t %t %t\n", objecthash, path, typeString, b)
	switch typeString {
	case "blob":
		if len(nameSegments) < 2 {
			b, err = GetBlob(objecthash.(HCID))
			err = PostBlob(b.Bytes())
			return BlobFromBytes(b.Bytes()).Hash(), err
		}
	case "list":
		nameSegments = strings.SplitN(nameSegments[1], "/", 2)
		l, _ := GetList(objecthash.(HCID))
		typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
		hash_of_new_entry, err := Post(objecthash, nameSegments[1], b)
		PostList(l.add(nameSegments[0], hash_of_new_entry, "list"))
		return l.Hash(), err
		//if err == nil {
		//	regenlist.PushBack(l)
		//	regenpath = append(regenpath, nameSegments[0])
		//} else {
		//	regen(regenlist, objecthash.(HID), b)
		//}
	case "tag":
		nameSegments = strings.SplitN(nameSegments[1], "/", 2)
		t, _ := GetTag(objecthash.(HKID), nameSegments[0])
		if !t.Verifiy() {
			return nil, errors.New("Tag Verifiy Failed")
		}
		typeString = t.TypeString
		objecthash = HID(t.HashBytes)
		regenlist.Init()
		regenlist.PushBack(t)
	case "commit":
		c, _ := GetCommit(objecthash.(HKID))
		if !c.Verifiy() {
			return nil, errors.New("Commit Verifiy Failed")
		}
		regenlist.Init()
		regenlist.PushBack(c)
		l, _ := GetList(c.listHash)
		fmt.Printf("c.listHash: %s\n", c.listHash.Hex())
		typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
	}
	if err != nil {
		return nil, err
	}

	fmt.Printf("%t %t %t %t\n", objecthash, path, typeString, b)
	//}
	b = nil
	return
}
