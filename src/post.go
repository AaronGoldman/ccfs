package main

import (
	//golist "container/list"
	"errors"
	"fmt"
	//"log"
	"strings"
)

func Post(objecthash HKID, path string, post_bytes Byteser) (hid HID, err error) {
	hid, err = post(HID(objecthash), path, "commit", post_bytes)
	return hid, err
}

func post(h HID, path string, next_path_segment_type string,
	post_bytes Byteser) (hid HID, err error) {
	if path == "" {
		err := PostBlob(post_bytes.(blob))
		return HID(post_bytes.(blob).Hash()), err
		//switch btype := b.(type) {
		//default:
		//	return nil, errors.New(fmt.Sprintf("Can not post type %T", btype))
		//case blob:
		//	err := PostBlob(post_bytes.(blob))
		//	return HID(post_bytes.(blob).Hash()), err
		//case list:
		//	return nil, errors.New("posting lists must have a path")
		//case commit:
		//case Tag:
		//}
	}
	nameSegments := strings.SplitN(path, "/", 2)
	next_path_segment := nameSegments[0]
	rest_of_path := ""
	if len(nameSegments) < 1 {
		rest_of_path = nameSegments[1]
	}
	switch next_path_segment_type {
	default:
		return nil, errors.New(fmt.Sprintf("Invalid type %T", next_path_segment_type))
	case "blob":
		return nil, errors.New(fmt.Sprintf("only \"\" path can be blob"))
	case "list":
		//posted_hash, err := list_helper(h, next_path_segment, rest_of_path, post_bytes)
		//return posted_hash, err
		l, geterr := GetList(h.Bytes())
		posterr := error(nil)
		if geterr == nil {
			//update and publish old list
			next_typeString, next_hash := l.hash_for_namesegment(next_path_segment)
			next_path := rest_of_path
			var hash_of_posted HID
			hash_of_posted, posterr = post(next_hash, next_path,
				next_typeString, post_bytes)
			l.add(next_path_segment, hash_of_posted, "list")
		} else {
			//build and publish new list
			next_hash := HID(nil)
			next_path := rest_of_path
			next_typeString := "list"
			if rest_of_path == "" {
				next_typeString = "blob"
			}
			var hash_of_posted HID
			hash_of_posted, posterr = post(next_hash, next_path,
				next_typeString, post_bytes)
			l = NewList(hash_of_posted, "list", next_path_segment)
		}
		if posterr != nil {
			return nil, posterr
		} else {
			err := PostList(l)
			if err == nil {
				return l.Hash(), nil
			}
			return nil, err
		}
	case "commit":

	case "tag":

	}
	return nil, errors.New("failed to return before end of function")
}

func list_helper(h HCID, next_path_segment string, rest_of_path string, post_bytes Byteser)

/*
func Post(objecthash Hexer, path string, b Byteser) (hid HID, err error) {
	typeString := "commit"
	//objecthash := hkid
	err = nil
	if path == "" {
		switch btype := b.(type) {
		default:
			return nil, errors.New(fmt.Sprintf("Can not post type %T", btype))
		case blob:
			err = PostBlob(b.(blob))
			return HID(b.(blob).Hash()), err
		case list:
			err = PostList(b.(list))
			return HID(b.(list).Hash()), err
		case commit:
			err = PostCommit(b.(commit))
			return HID(b.(commit).Hash()), err
		case Tag:
			err = PostTag(b.(Tag))
			return HID(b.(Tag).Hash()), err

		}
	}
	nameSegments := strings.SplitN(path, "/", 2)
	regenlist := golist.New()
	//regenpath := []string{}
	//for {
	log.Printf("\n\tobjecthash: %v\n\tpath: %v\n\ttypeString: %v\n\tblob %t\n",
		objecthash, path, typeString, b)
	switch typeString {
	case "blob":
		if len(nameSegments) < 2 {
			b, err = GetBlob(objecthash.(HCID))
			err = PostBlob(b.Bytes())
			return HID(BlobFromBytes(b.Bytes()).Hash()), err
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
		c, err := GetCommit(objecthash.(HKID))
		if err != nil {
			log.Printf("\n\tGetCommit err: %v\n", err)
			//commit_to_post = 
			//err = PostCommit(commit_to_post)
			return nil, err
		}
		log.Printf("%v %v", c, err != nil)
		if !c.Verifiy() {
			return nil, errors.New("Commit Verifiy Failed")
		}
		regenlist.Init()
		regenlist.PushBack(c)
		l, _ := GetList(c.listHash)
		log.Printf("c.listHash: %s\n", c.listHash.Hex())
		typeString, objecthash = l.hash_for_namesegment(nameSegments[0])
	}
	if err != nil {
		return nil, err
	}

	log.Printf("%t %t %t %t\n", objecthash, path, typeString, b)
	//}
	b = nil
	return
}*/
