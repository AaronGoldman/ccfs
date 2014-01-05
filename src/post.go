package main

import (
	//golist "container/list"
	"errors"
	"fmt"
	"log"
	"strings"
)

func Post(objecthash HKID, path string, post_bytes Byteser) (hid HID, err error) {
	hid, err = post(HID(objecthash), path, "commit", post_bytes, "blob")
	return hid, err
}

func post(h HID, path string, next_path_segment_type string,
	post_bytes Byteser, post_type string) (hid HID, err error) {
	//log.Printf(
	//	"\n\th: %v\n\tpath: %v\n\tnext_path_segment_type: %v\n\tpost_bytes: %v\n\tpost_type: %v\n",
	//	h, path, next_path_segment_type, post_bytes, post_type)
	if path == "" {
		err := PostBlob(post_bytes.(blob))
		return HCID(post_bytes.(blob).Hash()), err
	}

	nameSegments := strings.SplitN(path, "/", 2)
	next_path_segment := nameSegments[0]
	rest_of_path := ""
	if len(nameSegments) > 1 {
		rest_of_path = nameSegments[1]
	}
	switch next_path_segment_type {
	default:
		return nil, errors.New(fmt.Sprintf("Invalid type %T", next_path_segment_type))
	case "blob":
		return nil, errors.New(fmt.Sprintf("only \"\" path can be blob"))
	case "list":
		posted_hash, err := list_helper(h, next_path_segment, rest_of_path,
			post_bytes, post_type)
		return posted_hash, err
	case "commit":
		posted_hash, err := commit_helper(h.Bytes(), path, post_bytes,
			post_type)
		return posted_hash, err
	case "tag":
		posted_hash, err := tag_helper(h, next_path_segment, rest_of_path,
			post_bytes, post_type)
		return posted_hash, err
	}
}

func list_helper(h HID, next_path_segment string, rest_of_path string,
	post_bytes Byteser, post_type string) (HID, error) {
	geterr := errors.New("h in nil")
	l := list(nil)
	if h != nil {
		l, geterr = GetList(h.Bytes())
	}
	posterr := error(nil)
	if geterr == nil {
		//update and publish old list
		next_typeString, next_hash := l.hash_for_namesegment(next_path_segment)

		next_path := rest_of_path
		if next_typeString == "" {
			next_typeString = "list"
		}
		if rest_of_path == "" {
			next_typeString = post_type
		}
		var hash_of_posted HID
		if rest_of_path == "" && post_type != "blob" {
			hash_of_posted = HKID(post_bytes.Bytes()) //insrt reference by HKID
		} else {
			hash_of_posted, posterr = post(next_hash, next_path,
				next_typeString, post_bytes, post_type) //post Data
		}
		l.add(next_path_segment, hash_of_posted, next_typeString)
	} else {
		//build and publish new list
		next_hash := HID(nil)
		next_path := rest_of_path
		next_typeString := "list"
		if rest_of_path == "" {
			next_typeString = post_type
		}
		var hash_of_posted HID
		if rest_of_path == "" && post_type != "blob" {
			hash_of_posted = HKID(post_bytes.Bytes()) //insrt reference by HKID
		} else {
			hash_of_posted, posterr = post(next_hash, next_path,
				next_typeString, post_bytes, post_type)
		}
		l = NewList(hash_of_posted, next_typeString, next_path_segment)
	}
	if posterr != nil {
		return nil, posterr
	}
	err := PostList(l)
	if err == nil {
		return l.Hash(), nil
	}
	return nil, err
}

func commit_helper(h HCID, path string, post_bytes Byteser, post_type string) (HID, error) {
	c, geterr := GetCommit(h.Bytes())
	posterr := error(nil)
	if geterr == nil {
		next_typeString := "list"
		next_hash := c.listHash
		next_path := path
		var hash_of_posted HID
		hash_of_posted, posterr = post(next_hash, next_path,
			next_typeString, post_bytes, post_type)
		c = c.Update(hash_of_posted.Bytes())
	} else {
		_, err := getPrivateKeyForHkid(h.Bytes())
		if err == nil {
			next_hash := HID(nil)
			next_path := path
			next_typeString := "list"
			var hash_of_posted HID
			hash_of_posted, posterr = post(next_hash, next_path,
				next_typeString, post_bytes, post_type)
			c = NewCommit(hash_of_posted.Bytes(), h.Bytes())
		} else {
			log.Printf("You don't seem to own this repo\nh=%v\nerr=%v\n", h, err)
			return HKID{}, fmt.Errorf("You don't seem to own this repo")
		}
	}
	if posterr != nil {
		return nil, posterr
	}
	//log.Print(c)
	err := PostCommit(c)
	if err == nil {
		return c.hkid, nil
	}
	return nil, err
}

func tag_helper(h HID, next_path_segment string, rest_of_path string,
	post_bytes Byteser, post_type string) (HID, error) {
	t, geterr := GetTag(h.Bytes(), next_path_segment)
	posterr := error(nil)
	if geterr == nil {
		//the tag exists with that nameSegment
		next_typeString := t.TypeString
		next_hash := t.HashBytes
		next_path := rest_of_path
		var hash_of_posted HID
		if rest_of_path == "" && post_type != "blob" {
			hash_of_posted = HKID(post_bytes.Bytes()) //insrt reference by HKID
		} else {
			hash_of_posted, posterr = post(next_hash, next_path,
				next_typeString, post_bytes, post_type)
		}
		t = t.Update(hash_of_posted.Bytes())
	} else {
		//no tag exists with that nameSegment
		_, err := getPrivateKeyForHkid(h.Bytes())
		if err == nil {
			//you own the Domain
			next_hash := HID(nil)
			next_path := rest_of_path
			next_typeString := "list"
			if rest_of_path == "" {
				next_typeString = post_type
			}
			var hash_of_posted HID
			if rest_of_path == "" && post_type != "blob" {
				hash_of_posted = HKID(post_bytes.Bytes()) //insrt reference by HKID
			} else {
				hash_of_posted, posterr = post(next_hash, next_path,
					next_typeString, post_bytes, post_type)
			}
			t = NewTag(hash_of_posted, next_typeString, next_path_segment, h.Bytes())
		} else {
			log.Panic("You don't seem to own this Domain")
		}
	}
	if posterr != nil {
		return nil, posterr
	}
	//log.Print(t)
	err := PostTag(t)
	if err == nil {
		return t.Hkid(), nil
	}
	return nil, err
}

func InsertRepo(h HKID, path string, foreign_hkid HKID) error {
	_, err := post(
		h,
		path,
		"commit",
		foreign_hkid,
		"commit")
	return err
}

func InsertDomain(h HKID, path string, foreign_hkid HKID) error {
	_, err := post(
		h,
		path,
		"commit",
		foreign_hkid,
		"tag")
	return err
}

func InitRepo(h HKID, path string) error {
	foreign_hkid := GenHKID()
	err := InsertRepo(h, path, foreign_hkid)
	return err
}

func InitDomain(h HKID, path string) error {
	foreign_hkid := GenHKID()
	err := InsertDomain(h, path, foreign_hkid)
	return err
}
