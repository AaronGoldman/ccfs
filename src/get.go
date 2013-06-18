package main

import (
	"errors"
	"strings"
)

func get(objecthash HKID, path string) (b blob, err error) {
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
		if err != nil {
			return nil, err
		}
	}
	b = nil
	return
}
