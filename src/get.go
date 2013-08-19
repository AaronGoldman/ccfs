package main

import (
	"errors"
	"fmt"
	"strings"
)

func Get(objecthash Hexer, path string) (b blob, err error) {
	typeString := "commit"
	//objecthash := hkid
	err = nil
	nameSegments := strings.SplitN(path, "/", 2)
	for {
		fmt.Printf("Path: %s Type: %s\n", path, typeString)
		switch typeString {
		case "blob":
			if len(nameSegments) < 2 {
				b, err = GetBlob(objecthash.(HCID))
				if err != nil {
					panic(err)
				}
				return
			}
		case "list":
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			l, _ := GetList(objecthash.(HCID))
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
		case "tag":
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			t, err := GetTag(HKID(objecthash.(HID)), nameSegments[0])
			if !t.Verifiy() {
				b = nil
				err = errors.New("Tag Verifiy Failed")
				return b, err
			}
			typeString = t.TypeString
			objecthash = HCID(t.HashBytes)
			if len(nameSegments) < 2 {
				path = ""
			} else {
				path = nameSegments[1]
			}

		case "commit":
			c, err := GetCommit(objecthash.(HKID))
			if !c.Verifiy() {
				b = nil
				err = errors.New("Commit Verifiy Failed")
				return b, err
			}
			l, err := GetList(c.listHash)
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
			//fmt.Printf("%v\n", c)
			path = nameSegments[1]
			_ = err

		case "":
			return nil, err
		}

		if err != nil {
			return nil, err
		}
	}
	b = nil
	return
}
