package main

import (
	"fmt"
	"log"
	"strings"
)

//Get retrieves the content objects using HID and path
func Get(objecthash HID, path string) (b blob, err error) {
	typeString := "commit"
	err = nil
	nameSegments := []string{"", path}
	for {
		if len(nameSegments) > 1 {
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
		} else {
			nameSegments = []string{""}
		}
		//log.Printf("\n\tPath: %s\n\tType: %v\n\tobjecthash: %v\n",
		//	nameSegments, typeString, objecthash)
		switch typeString {
		case "blob":
			b, err = GetBlob(objecthash.Bytes())
			if err != nil {
				log.Printf("\n\t%v\n", err)
			}
			return b, err
		case "list":
			var l list
			l, err = GetList(objecthash.Bytes())
			if err != nil {
				log.Printf("\n\t%v\n", err)
			}
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
			if objecthash == nil && nameSegments[0] != "" {
				err = fmt.Errorf("Blob not found")
			}
			b = l.Bytes()
		case "tag":
			var t tag
			if nameSegments[0] == "" {
				log.Printf("\n\tNo Path\n")
			}
			t, err = GetTag(objecthash.Bytes(), nameSegments[0])
			if err != nil {
				//log.Printf("\n\t%v\n", err)
				return nil, err
			}
			if !t.Verify() {
				return nil, fmt.Errorf("Tag Verifiy Failed")
			}
			typeString = t.TypeString
			objecthash = HCID(t.HashBytes)
			b = t.Bytes()
		case "commit":
			var c commit
			c, err = GetCommit(objecthash.Bytes())
			if err != nil {
				log.Printf("\n\t%v\n", err)
			}
			if !c.Verify() {
				return nil, fmt.Errorf("Commit Verifiy Failed")
			}
			var l list
			l, err = GetList(c.listHash)
			if err != nil {
				log.Printf("\n\t%v\n", err)
			}
			typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
			if objecthash == nil && nameSegments[0] != "" {
				err = fmt.Errorf("Blob not found")
			}
			//if err != nil {
			//	log.Printf("%v\n", err)
			//}
			b = l.Bytes()
		default:
			log.Printf("\n\t%v\n", err)
			panic(err)
		}
		//if len(nameSegments) == 1 && typeString != "blob" {
		if objecthash == nil {
			return b, err
		}
	}
}
