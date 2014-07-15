package services

import (
	"fmt"
	"log"
	"strings"

	"github.com/AaronGoldman/ccfs/objects"
)

func Post(objecthash objects.HKID, path string, post_bytes objects.Byteser) (hid objects.HID, err error) {
	hid, err = post(objects.HID(objecthash), path, "commit", post_bytes, "blob")
	return hid, err
}

func post(h objects.HID, path string, next_path_segment_type string,
	post_bytes objects.Byteser, post_type string) (hid objects.HID, err error) {
	//log.Printf(
	//	"\n\th: %v\n\tpath: %v\n\tnext_path_segment_type: %v\n\tpost_bytes: %s\n\tpost_type: %v\n",
	//	h, path, next_path_segment_type, post_bytes, post_type)
	if path == "" {
		//log.Printf("post_type: %s", post_type)
		err := PostBlob(post_bytes.(objects.Blob))
		return objects.HCID(post_bytes.(objects.Blob).Hash()), err
	}

	nameSegments := strings.SplitN(path, "/", 2)
	next_path_segment := nameSegments[0]
	rest_of_path := ""
	if len(nameSegments) > 1 {
		rest_of_path = nameSegments[1]
	}
	switch next_path_segment_type {
	default:
		return nil, fmt.Errorf(fmt.Sprintf("Invalid type %T", next_path_segment_type))
	case "blob":
		return nil, fmt.Errorf(fmt.Sprintf("only \"\" path can be blob"))
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

func list_helper(h objects.HID, next_path_segment string, rest_of_path string,
	post_bytes objects.Byteser, post_type string) (objects.HID, error) {
	geterr := fmt.Errorf("h in nil")
	l := objects.List(nil)
	if h != nil {
		l, geterr = GetList(h.(objects.HCID))
	}
	posterr := error(nil)
	if geterr == nil {
		//update and publish old list
		next_typeString, next_hash := l.Hash_for_namesegment(next_path_segment)

		next_path := rest_of_path
		if next_typeString == "" {
			next_typeString = "list"
		}
		if rest_of_path == "" {
			next_typeString = post_type
		}
		var hash_of_posted objects.HID
		if rest_of_path == "" && post_type != "blob" {
			hash_of_posted = post_bytes.(objects.HKID) //insrt reference by HKID
			//log.Printf("insrt reference by HKID\n\tnext_path_segment:%s\n\trest_of_path:%s\n", next_path_segment, rest_of_path)
		} else {
			hash_of_posted, posterr = post(next_hash, next_path,
				next_typeString, post_bytes, post_type) //post Data
		}
		l.Add(next_path_segment, hash_of_posted, next_typeString)
	} else {
		//build and publish new list
		next_hash := objects.HID(nil)
		next_path := rest_of_path
		next_typeString := "list"
		if rest_of_path == "" {
			next_typeString = post_type
		}
		var hash_of_posted objects.HID
		if rest_of_path == "" && post_type != "blob" {
			hash_of_posted = post_bytes.(objects.HKID) //insrt reference by HKID
			//log.Printf("insrt reference by HKID\n\tnext_path_segment:%s\n\trest_of_path:%s\n\t", next_path_segment, rest_of_path)
		} else {
			hash_of_posted, posterr = post(next_hash, next_path,
				next_typeString, post_bytes, post_type)
		}
		l = objects.NewList(hash_of_posted, next_typeString, next_path_segment)
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

func commit_helper(h objects.HKID, path string, post_bytes objects.Byteser, post_type string) (objects.HID, error) {
	c, geterr := GetCommit(h)
	posterr := error(nil)
	if geterr == nil {
		//A existing vertion was found
		next_typeString := "list"
		next_hash := c.ListHash
		next_path := path
		var hash_of_posted objects.HID
		hash_of_posted, posterr = post(next_hash, next_path,
			next_typeString, post_bytes, post_type)
		if posterr != nil {
			return nil, posterr
		}
		c = c.Update(hash_of_posted.Bytes())
	} else {
		//A existing vertion was NOT found
		_, err := GetPrivateKeyForHkid(h)
		if err == nil {
			next_hash := objects.HID(nil)
			next_path := path
			next_typeString := "list"
			var hash_of_posted objects.HID
			hash_of_posted, posterr = post(next_hash, next_path,
				next_typeString, post_bytes, post_type)
			c = objects.NewCommit(hash_of_posted.Bytes(), h)
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

func tag_helper(h objects.HID, next_path_segment string, rest_of_path string,
	post_bytes objects.Byteser, post_type string) (objects.HID, error) {
	t, geterr := GetTag(h.Bytes(), next_path_segment)
	posterr := error(nil)
	if geterr == nil {
		//the tag exists with that nameSegment
		next_typeString := t.TypeString
		next_hash := t.HashBytes
		next_path := rest_of_path
		var hash_of_posted objects.HID
		if rest_of_path == "" && post_type != "blob" {
			hash_of_posted = objects.HKID(post_bytes.Bytes()) //insrt reference by HKID
		} else {
			hash_of_posted, posterr = post(next_hash, next_path,
				next_typeString, post_bytes, post_type)
		}
		t = t.Update(hash_of_posted, t.TypeString)
	} else {
		//no tag exists with that nameSegment
		_, err := GetPrivateKeyForHkid(h.(objects.HKID))
		if err == nil {
			//you own the Domain

			next_hash := objects.HID(nil)
			next_path := rest_of_path
			next_typeString := "list"
			if rest_of_path == "" {
				next_typeString = post_type
			}
			log.Printf("tag pionting to a %s named %s", next_typeString, next_path_segment)
			var hash_of_posted objects.HID
			if rest_of_path == "" && post_type != "blob" {
				hash_of_posted = objects.HKID(post_bytes.Bytes()) //insrt reference by HKID
			} else {
				hash_of_posted, posterr = post(next_hash, next_path,
					next_typeString, post_bytes, post_type)
			}
			t = objects.NewTag(hash_of_posted, next_typeString, next_path_segment, nil, h.Bytes())
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
func InsertRepo(h objects.HKID, path string, foreign_hkid objects.HKID) error {
	//log.Printf("\n\trootHKID:%s\n\tPath:%s\n\tforeignHKID:%s", h, path, foreign_hkid)
	_, err := post(
		h,
		path,
		"commit",
		foreign_hkid,
		"commit")
	return err
}

//InsertDomain inserts a given foreign hkid to the local HKID at the path spesified
func InsertDomain(h objects.HKID, path string, foreign_hkid objects.HKID) error {
	//log.Printf("\n\trootHKID:%s\n\tPath:%s\n\tforeignHKID:%s", h, path, foreign_hkid)
	_, err := post(
		h,
		path,
		"commit",
		foreign_hkid,
		"tag")
	return err
}

//InitRepo creates a new repository and inserts it to the HKID at the path spesified
func InitRepo(h objects.HKID, path string) error {
	foreign_hkid := objects.GenHKID()
	err := InsertRepo(h, path, foreign_hkid)
	return err
}

//InitDomain creates a new domain and inserts it to the HKID at the path spesified
func InitDomain(h objects.HKID, path string) error {
	foreign_hkid := objects.GenHKID()
	err := InsertDomain(h, path, foreign_hkid)
	return err
}
