package main

import (
	"errors"
	"log"
	"strings"
)

func Get(objecthash HID, path string) (b blob, err error) {
	typeString := "commit"
	//objecthash := hkid
	err = nil
	nameSegments := strings.SplitN(path, "/", 2)
	for {
		log.Printf("\n\tPath: %s\n\tType: %s\n", path, typeString)
		switch typeString {
		case "blob":
			if len(nameSegments) < 2 {
				b, err = GetBlob(objecthash.(HCID))
				if err != nil {
					if err != nil {
						log.Printf("\n\t%v\n", err)
					}
				}
				return
			}
		case "list":
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			l, err := GetList(objecthash.(HCID))
			if err != nil {
				log.Printf("\n\t%v\n", err)
			}
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
		case "tag":
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			t, err := GetTag(objecthash.Bytes(), nameSegments[0])
			if err != nil {
				log.Printf("\n\t%v\n", err)
			}
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
			if err != nil {
				log.Printf("\n\t%v\n", err)
			}
			if !c.Verifiy() {
				b = nil
				err = errors.New("Commit Verifiy Failed")
				return b, err
			}
			l, err := GetList(c.listHash)
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
			//log.Printf("%v\n", c)
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
