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
		//log.Printf("\n\tPath: %s\n\tType: %v\n\tobjecthash: %v\n",
		//	nameSegments,
		//	typeString,
		//	objecthash)
		switch typeString {
		case "blob":
			//if len(nameSegments) < 2 {
			b, err = GetBlob(objecthash.Bytes())
			//log.Printf("\n\t%v\n", string(b))
			if err != nil {
				if err != nil {
					log.Printf("\n\t%v\n", err)
				}
			}
			return b, err
			//}
		case "list":

			if len(nameSegments) > 1 {
				nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			}
			//log.Printf("\n\t%v\n\t%v\n", nameSegments, objecthash)
			l, err := GetList(objecthash.Bytes())
			if err != nil {
				log.Printf("\n\t%v\n", err)
			}
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
			if len(nameSegments) == 1 && typeString != "blob" {
				//log.Printf("\n\t%v\n", l)
				return blob(l.Bytes()), err
			}
		case "tag":
			if len(nameSegments) > 1 {
				nameSegments = strings.SplitN(nameSegments[1], "/", 2)
			}
			//log.Printf("\n\t%v\n\t%v\n", objecthash, nameSegments)
			t, err := GetTag(objecthash.Bytes(), nameSegments[0])
			if err != nil {
				log.Printf("\n\t%v\n", err)
				return nil, err
			}
			if !t.Verifiy() {
				b = nil
				err = errors.New("Tag Verifiy Failed")
				return b, err
			}
			typeString = t.TypeString
			objecthash = HCID(t.HashBytes)
			if len(nameSegments) == 1 && typeString != "blob" {
				return blob(t.Bytes()), err
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
			if err != nil {
				log.Printf("%v\n", err)
			}
			if len(nameSegments) == 1 && typeString != "blob" {
				return blob(l.Bytes()), err
			}

		case "":
			log.Printf("\n\t%v\n", objecthash)
			b, err = GetBlob(objecthash.Bytes())
			return b, err
		}
	}
}
