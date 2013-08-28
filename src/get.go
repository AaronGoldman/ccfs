package main

import (
	"errors"
	"log"
	"strings"
)

func Get(objecthash HID, path string) (b blob, err error) {
	typeString := "commit"
	err = nil
	nameSegments := []string{"", path}
	for {
		//log.Printf("\n\tPath: %s\n\tType: %s\n", path, typeString)
		switch typeString {
		case "blob":
			if len(nameSegments) < 2 {
				b, err = GetBlob(objecthash.(HCID))
				//log.Printf("\n\t%v\n", string(b))
				if err != nil {
					if err != nil {
						log.Printf("\n\t%v\n", err)
					}
				}
				return
			}
		case "list":

			if len(nameSegments) > 1 {
				nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			}
			//log.Printf("\n\t%v\n\t%v\n", nameSegments, objecthash.Hex())
			l, err := GetList(objecthash.Bytes())
			if err != nil {
				log.Printf("\n\t%v\n", err)
			}
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
		case "tag":
			if len(nameSegments) > 1 {
				nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			}
			t, err := GetTag(objecthash.Bytes(), nameSegments[0])
			if err != nil {
				log.Printf("\n\t%v\n", err)
			}
			//log.Printf("\n\t%v\n", t)
			if !t.Verifiy() {
				b = nil
				err = errors.New("Tag Verifiy Failed")
				return b, err
			}
			typeString = t.TypeString
			objecthash = HCID(t.HashBytes)
			if len(nameSegments) == 1 {
				path = ""
			} else {
				path = nameSegments[1]
			}

		case "commit":
			if len(nameSegments) > 1 {
				nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			}
			c, err := GetCommit(objecthash.Bytes())
			if err != nil {
				log.Printf("\n\t%v\n", err)
			}
			if !c.Verifiy() {
				b = nil
				err = errors.New("Commit Verifiy Failed")
				return b, err
			}
			l, err := GetList(c.listHash)
			if err != nil {
				log.Printf("\n\t%v\n", err)
			}
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
			//log.Printf("\n\t%v\n", c)
			//log.Printf("\n\t%v\n", l)
			path = ""
			if len(nameSegments) > 1 {
				path = nameSegments[1]
			}
			if err != nil {
				log.Printf("%v\n", err)
			}

		case "":
			return nil, err
		}
	}
}
