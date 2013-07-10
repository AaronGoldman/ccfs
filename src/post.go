package main

import (
	golist "container/list"
	"errors"
	"fmt"
	"strings"
)

func Post(objecthash Hexer, path string, b Byteser) (err error) {
	typeString := "commit"
	//objecthash := hkid
	err = nil
	nameSegments := strings.SplitN(path, "/", 2)
	regenlist := golist.New()
	regenpath := []string{}
	//for {
	fmt.Printf("%t %t %t\n", objecthash, path, b)
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
		if err == nil {
			regenlist.PushBack(l)
			regenpath = append(regenpath, nameSegments[0])
		} else {
			regen(regenlist, objecthash.(HID), b)
		}
	case "tag":
		nameSegments = strings.SplitN(nameSegments[1], "/", 2)
		t, _ := GetTag(objecthash.(HID), nameSegments[0])
		if !t.Verifiy() {
			return errors.New("Tag Verifiy Failed")
		}
		typeString = t.TypeString
		objecthash = HID(t.HashBytes)
		regenlist.Init()
		regenlist.PushBack(t)
	case "commit":
		c, _ := GetCommit(objecthash.(HKID))
		if !c.Verifiy() {
			return errors.New("Commit Verifiy Failed")
		}
		regenlist.Init()
		regenlist.PushBack(c)
		l, _ := GetList(c.listHash)

		typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
	}
	if err != nil {
		return err
	}
	if typeString != "list" {

		regenpath = []string{}
	}
	fmt.Printf("%t %t %t\n", objecthash, path, b)
	//}
	b = nil
	return
}
