package main

import (
	"errors"
	"strings"
)

func Get(objecthash Hexer, path string) (b blob, err error) {
	typeString := "commit"
	//objecthash := hkid
	err = nil
	nameSegments := strings.SplitN(path, "/", 2)
	for {
		switch typeString {
		case "blob":
			if len(nameSegments) < 2 {
				b, err = GetBlob(objecthash.(HID))
				return
			}
		case "list":
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			l, _ := GetList(objecthash.(HID))
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
		case "tag":
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			t, err := GetTag(objecthash.(HID), nameSegments[0])
			if !t.Verifiy() {
				b = nil
				err = errors.New("Tag Verifiy Failed")
				return b, err
			}
			typeString = t.TypeString
			objecthash = HKID(t.HashBytes)
		case "commit":
			c, err := GetCommit(objecthash.(HKID))
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
