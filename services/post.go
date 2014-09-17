//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file
package services

import (
	"fmt"
	"log"
	"strings"

	"github.com/AaronGoldman/ccfs/objects"
)

//Post releases a content object and the necessary intermediate objects to storage
func Post(objecthash objects.HKID, path string, postBytes objects.Byteser) (hid objects.HID, err error) {
	hid, err = post(objects.HID(objecthash), path, "commit", postBytes, "blob")
	return hid, err
}

func post(h objects.HID, path string, nextPathSegmentType string,
	postBytes objects.Byteser, postType string) (hid objects.HID, err error) {
	//log.Printf(
	//	"\n\th: %v\n\tpath: %v\n\tnext_path_segment_type: %v\n\tpost_bytes: %s\n\tpost_type: %v\n",
	//	h, path, next_path_segment_type, post_bytes, post_type)
	if path == "" {
		//log.Printf("post_type: %s", post_type)
		err := PostBlob(postBytes.(objects.Blob))
		return objects.HCID(postBytes.(objects.Blob).Hash()), err
	}

	nameSegments := strings.SplitN(path, "/", 2)
	nextPathSegment := nameSegments[0]
	restOfPath := ""
	if len(nameSegments) > 1 {
		restOfPath = nameSegments[1]
	}
	switch nextPathSegmentType {
	default:
		return nil, fmt.Errorf(fmt.Sprintf("Invalid type %T", nextPathSegmentType))
	case "blob":
		return nil, fmt.Errorf(fmt.Sprintf("only \"\" path can be blob"))
	case "list":
		postedListHash, err := listHelper(h, nextPathSegment, restOfPath,
			postBytes, postType)
		return postedListHash, err
	case "commit":
		postedCommitHash, err := commitHelper(h.Bytes(), path, postBytes,
			postType)
		return postedCommitHash, err
	case "tag":
		postedTagHash, err := tagHelper(h, nextPathSegment, restOfPath,
			postBytes, postType)
		return postedTagHash, err
	}
}

func listHelper(h objects.HID, nextPathSegment string, restOfPath string,
	postBytes objects.Byteser, postType string) (objects.HID, error) {
	geterr := fmt.Errorf("h in nil")
	l := objects.List(nil)
	if h != nil {
		l, geterr = GetList(h.(objects.HCID))
	}
	posterr := error(nil)
	if geterr == nil {
		//update and publish old list
		nextTypeString, nextHash := l.HashForNamesegment(nextPathSegment)

		nextPath := restOfPath
		if nextTypeString == "" {
			nextTypeString = "list"
		}
		if restOfPath == "" {
			nextTypeString = postType
		}
		var hashOfPosted objects.HID
		if restOfPath == "" && postType != "blob" {
			hashOfPosted = postBytes.(objects.HKID) //insrt reference by HKID
			//log.Printf("insrt reference by HKID\n\tnext_path_segment:%s\n\trest_of_path:%s\n", next_path_segment, rest_of_path)
		} else {
			hashOfPosted, posterr = post(nextHash, nextPath,
				nextTypeString, postBytes, postType) //post Data
		}
		l.Add(nextPathSegment, hashOfPosted, nextTypeString)
	} else {
		//build and publish new list
		nextHash := objects.HID(nil)
		nextPath := restOfPath
		nextTypeString := "list"
		if restOfPath == "" {
			nextTypeString = postType
		}
		var hashOfPosted objects.HID
		if restOfPath == "" && postType != "blob" {
			hashOfPosted = postBytes.(objects.HKID) //insrt reference by HKID
			//log.Printf("insrt reference by HKID\n\tnext_path_segment:%s\n\trest_of_path:%s\n\t", next_path_segment, rest_of_path)
		} else {
			hashOfPosted, posterr = post(nextHash, nextPath,
				nextTypeString, postBytes, postType)
		}
		l = objects.NewList(hashOfPosted, nextTypeString, nextPathSegment)
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

func commitHelper(h objects.HKID, path string, postBytes objects.Byteser, postType string) (objects.HID, error) {
	c, geterr := GetCommit(h)
	posterr := error(nil)
	if geterr == nil {
		//A existing vertion was found
		nextTypeString := "list"
		nextHash := c.ListHash
		nextPath := path
		var hashOfPosted objects.HID
		hashOfPosted, posterr = post(nextHash, nextPath,
			nextTypeString, postBytes, postType)
		if posterr != nil {
			return nil, posterr
		}
		c = c.Update(hashOfPosted.Bytes())
	} else {
		//A existing vertion was NOT found
		_, err := GetPrivateKeyForHkid(h)
		if err == nil {
			nextHash := objects.HID(nil)
			nextPath := path
			nextTypeString := "list"
			var hashOfPosted objects.HID
			hashOfPosted, posterr = post(nextHash, nextPath,
				nextTypeString, postBytes, postType)
			c = objects.NewCommit(hashOfPosted.Bytes(), h)
		} else {
			log.Printf("You don't seem to own this repo\n\th=%v\n\terr=%v\n", h, err)
			return objects.HKID{}, fmt.Errorf("You don't seem to own this repo")
		}
	}
	if posterr != nil {
		return nil, posterr
	}
	//log.Print(c)
	err := PostCommit(c)
	if err == nil {
		return c.Hkid, nil
	}
	return nil, err
}

func tagHelper(h objects.HID, nextPathSegment string, restOfPath string,
	postBytes objects.Byteser, postType string) (objects.HID, error) {
	t, geterr := GetTag(h.Bytes(), nextPathSegment)
	posterr := error(nil)
	if geterr == nil {
		//the tag exists with that nameSegment
		nextTypeString := t.TypeString
		nextHash := t.HashBytes
		nextPath := restOfPath
		var hashOfPosted objects.HID
		if restOfPath == "" && postType != "blob" {
			hashOfPosted = objects.HKID(postBytes.Bytes()) //insrt reference by HKID
		} else {
			hashOfPosted, posterr = post(nextHash, nextPath,
				nextTypeString, postBytes, postType)
		}
		t = t.Update(hashOfPosted, t.TypeString)
	} else {
		//no tag exists with that nameSegment
		_, err := GetPrivateKeyForHkid(h.(objects.HKID))
		if err == nil {
			//you own the Domain

			nextHash := objects.HID(nil)
			nextPath := restOfPath
			nextTypeString := "list"
			if restOfPath == "" {
				nextTypeString = postType
			}
			log.Printf("tag pionting to a %s named %s", nextTypeString, nextPathSegment)
			var hashOfPosted objects.HID
			if restOfPath == "" && postType != "blob" {
				hashOfPosted = objects.HKID(postBytes.Bytes()) //insrt reference by HKID
			} else {
				hashOfPosted, posterr = post(nextHash, nextPath,
					nextTypeString, postBytes, postType)
			}
			t = objects.NewTag(hashOfPosted, nextTypeString, nextPathSegment, nil, h.Bytes())
		} else {
			log.Printf("You don't seem to own this Domain")
			return nil, fmt.Errorf("You dont own the Domain, meanie")
		}
	}
	if posterr != nil {
		return nil, posterr
	}
	//log.Print(t)
	err := PostTag(t)
	if err == nil {
		return t.Hkid, nil
	}
	return nil, err
}

//InsertRepo inserts a given foreign hkid to the local HKID at the path spesified
func InsertRepo(h objects.HKID, path string, foreignHkid objects.HKID) error {
	//log.Printf("\n\trootHKID:%s\n\tPath:%s\n\tforeignHKID:%s", h, path, foreign_hkid)
	_, err := post(
		h,
		path,
		"commit",
		foreignHkid,
		"commit")
	return err
}

//InsertDomain inserts a given foreign hkid to the local HKID at the path spesified
func InsertDomain(h objects.HKID, path string, foreignHkid objects.HKID) error {
	//log.Printf("\n\trootHKID:%s\n\tPath:%s\n\tforeignHKID:%s", h, path, foreign_hkid)
	_, err := post(
		h,
		path,
		"commit",
		foreignHkid,
		"tag")
	return err
}

//InitRepo creates a new repository and inserts it to the HKID at the path spesified
func InitRepo(h objects.HKID, path string) error {
	foreignHkid := objects.GenHKID()
	err := InsertRepo(h, path, foreignHkid)
	return err
}

//InitDomain creates a new domain and inserts it to the HKID at the path spesified
func InitDomain(h objects.HKID, path string) error {
	foreignHkid := objects.GenHKID()
	err := InsertDomain(h, path, foreignHkid)
	return err
}
