package main

import (
	"errors"
	//"time"
)

type contentservice interface {
	contentgeter
	contentposter
}

type contentgeter interface {
	blobgeter
	commitgeter
	taggeter
	keygeter
}

type blobgeter interface {
	getBlob(HCID) (blob, error)
}
type commitgeter interface {
	getCommit(HKID) (commit, error)
}
type taggeter interface {
	//	getTag(h HKID, namesegment string) (tag, error)
}
type keygeter interface {
	getKey(HKID) (blob, error)
}

type contentposter interface {
	blobpostter
	commitposter
	tagposter
	keyposter
}

type blobpostter interface {
	postBlob(b blob) error
}
type commitposter interface {
	postCommit(c commit) error
}
type tagposter interface {
	//	postTag(t tag) error
}
type keyposter interface {
	//	postKey(p *ecdsa.PrivateKey) error
}

//type Contentservice interface {
//	blobgeter(datach chan blob, errorch chan error, h HCID)
//	commitgeter(datach chan commit, errorch chan error, h HKID)
//	keygeter(datach chan blob, errorch chan error, h HKID)
//	taggeter(datach chan Tag, errorch chan error, h HKID, namesegment string)
//}

//var services = []contentservice{
//	localfileservice{},
//	timeoutservice{},
//	googledriveservice{},
//}

var blobgeters = []func(chan blob, chan error, HCID){
	localfileservice_blobgeter,
	timeoutservice_blobgeter,
	//googledriveservice_blobgeter,
}

var commitgeters = []func(chan commit, chan error, HKID){
	localfileservice_commitgeter,
	timeoutservice_commitgeter,
	//googledriveservice_commitgeter,
}

var taggeters = []func(chan Tag, chan error, HKID, string){
	localfileservice_taggeter,
	timeoutservice_taggeter,
	//googledriveservice_taggeter,
}

var keygeters = []func(chan blob, chan error, HKID){
	localfileservice_keygeter,
	timeoutservice_keygeter,
	//googledriveservice_keygeter,
}

func GetBlob(h HCID) (blob, error) {
	datach := make(chan blob, len(blobgeters))
	errorch := make(chan error, len(blobgeters))
	for _, blobgeter := range blobgeters {
		go blobgeter(datach, errorch, h)
	}
	for {
		select {
		case b := <-datach:
			if b != nil && b.Hash().Hex() == h.Hex() {
				return b, nil
			}
			return nil, errors.New("Commit Verifiy Failed")
		case err := <-errorch:
			return nil, err
		}
	}
}

func GetCommit(h HKID) (commit, error) {
	datach := make(chan commit, len(commitgeters))
	errorch := make(chan error, len(commitgeters))
	for _, commitgeter := range commitgeters {
		go commitgeter(datach, errorch, h)
	}
	for {
		select {
		case c := <-datach:
			if c.Verifiy() {
				return c, nil
			}
			return commit{}, errors.New("Commit Verifiy Failed")
		case err := <-errorch:
			return commit{}, err
		}
	}
}

func GetTag(h HKID, namesegment string) (Tag, error) {
	datach := make(chan Tag, len(taggeters))
	errorch := make(chan error, len(taggeters))
	for _, taggeter := range taggeters {
		go taggeter(datach, errorch, h, namesegment)
	}
	for {
		select {
		case t := <-datach:
			if t.Verifiy() {
				return t, nil
			}
			return Tag{}, errors.New("Tag Verifiy Failed")
		case err := <-errorch:
			return Tag{}, err
		}
	}
}

func GetKey(h HKID) (blob, error) {
	datach := make(chan blob, len(keygeters))
	errorch := make(chan error, len(keygeters))
	for _, keygeter := range keygeters {
		go keygeter(datach, errorch, h)
	}
	for {
		select {
		case b := <-datach:
			//if something { //How to Verifiy key?
			return b, nil
			//}
			//return nil, errors.New("Key Verifiy Failed")
		case err := <-errorch:
			return nil, err
		}
	}
}

//func GetCommit(h HKID) (commit, error) {
//	commit_chan := make(chan commit)
//	go func(commit_chan chan commit) {
//		c, err := localfileservice_GetCommit(h)
//		if err == nil {
//			commit_chan <- c
//		}
//	}(commit_chan)
//	select {
//	case b := <-commit_chan:
//		return b, nil
//	case <-time.After(time.Second):
//		return commit{}, errors.New("GetCommit Timeout")
//	}
//}

//func GetTag(h HKID, namesegment string) (Tag, error) {
//	tag_chan := make(chan Tag)
//	go func(tag_chan chan Tag) {
//		c, err := localfileservice_GetTag(h, namesegment)
//		if err == nil {
//			tag_chan <- c
//		}
//	}(tag_chan)
//	select {
//	case b := <-tag_chan:
//		return b, nil
//	case <-time.After(time.Second):
//		return Tag{}, errors.New("GetTag Timeout")
//	}
//}

//func GetKey(h HKID) (blob, error) {
//	blob_chan := make(chan blob)
//	go func(blob_chan chan blob) {
//		c, err := localfileservice_GetKey(h)
//		if err == nil {
//			blob_chan <- c
//		}
//	}(blob_chan)
//	select {
//	case b := <-blob_chan:
//		return b, nil
//	case <-time.After(time.Second):
//		return nil, errors.New("GetBlob Timeout")
//	}
//}
