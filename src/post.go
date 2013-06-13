package main

import (
	golist "container/list"
	"errors"
	"strings"
)

func post(objecthash hkid, path string, b blob) (err error) {
	typeString := "commit"
	//objecthash := hkid
	err = nil
	nameSegments := strings.SplitN(path, "/", 2)
	regenlist := golist.New()
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
				regenlist.PushBack(l)
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
			regenlist.Init()
			regenlist.PushBack(t)
		case "commit":
			c, _ := GetCommit(objecthash)
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
	}
	b = nil
	return
}
