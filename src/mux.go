package main

import (
	"errors"
	//"time"
	//"log"
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
	getTag(h HKID, namesegment string) (tag, error)
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
	postTag(t tag) error
}
type keyposter interface {
	postKey(p PrivateKey) error
}

var blobgeters = []blobgeter{
	timeoutserviceInstance,
	localfileserviceInstance,
	//googledriveserviceInstance,
	//appsscriptserviceInstance,
}
var commitgeters = []commitgeter{
	timeoutserviceInstance,
	localfileserviceInstance,
	//googledriveserviceInstance,
}
var taggeters = []taggeter{
	timeoutserviceInstance,
	localfileserviceInstance,
	//googledriveserviceInstance,
}
var keygeters = []keygeter{
	timeoutserviceInstance,
	localfileserviceInstance,
	//googledriveserviceInstance,
}

//var commitgeters = []func(chan commit, chan error, HKID){
//	localfileservice_commitgeter,
//	timeoutservice_commitgeter,
//	//googledriveservice_commitgeter,
//}

//var taggeters = []func(chan tag, chan error, HKID, string){
//	localfileservice_taggeter,
//	timeoutservice_taggeter,
//	//googledriveservice_taggeter,
//}

//var keygeters = []func(chan blob, chan error, HKID){
//	localfileservice_keygeter,
//	timeoutservice_keygeter,
//	//googledriveservice_keygeter,
//}

func GetBlob(h HCID) (blob, error) {
	datach := make(chan blob, len(blobgeters))
	errorch := make(chan error, len(blobgeters))
	for _, blobgeterInstance := range blobgeters {
		go func(datach chan blob, errorch chan error, h HCID) {
			b, err := blobgeterInstance.getBlob(h)
			if err == nil {
				datach <- b
				return
			} else {
				errorch <- err
				return
			}
		}(datach, errorch, h)
	}
	for {
		select {
		case b := <-datach:
			if b != nil && b.Hash().Hex() == h.Hex() {
				return b, nil
			}
			return nil, errors.New("Blob Verifiy Failed")
		case err := <-errorch:
			return nil, err
		}
	}
}

func GetCommit(h HKID) (commit, error) {
	datach := make(chan commit, len(commitgeters))
	errorch := make(chan error, len(commitgeters))
	for _, commitgeterInstance := range commitgeters {
		//go commitgeter(datach, errorch, h)
		go func(datach chan commit, errorch chan error, h HKID) {
			c, err := commitgeterInstance.getCommit(h)
			if err == nil {
				datach <- c
				return
			} else {
				errorch <- err
				return
			}
		}(datach, errorch, h)
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

func GetTag(h HKID, namesegment string) (tag, error) {
	datach := make(chan tag, len(taggeters))
	errorch := make(chan error, len(taggeters))
	for _, taggeterInstance := range taggeters {
		//go taggeter(datach, errorch, h, namesegment)
		go func(datach chan tag, errorch chan error, h HKID, namesegment string) {
			t, err := taggeterInstance.getTag(h, namesegment)
			if err == nil {
				datach <- t
				return
			} else {
				errorch <- err
				return
			}
		}(datach, errorch, h, namesegment)
	}
	for {
		select {
		case t := <-datach:
			if t.Verifiy() {
				return t, nil
			}
			return tag{}, errors.New("Tag Verifiy Failed")
		case err := <-errorch:
			return tag{}, err
		}
	}
}

func GetKey(h HKID) (blob, error) {
	datach := make(chan blob, len(keygeters))
	errorch := make(chan error, len(keygeters))
	for _, keygeterInstance := range keygeters {
		//log.Printf("Instance %T\n", keygeterInstance)
		//go keygeter(datach, errorch, h)
		go func(datach chan blob, errorch chan error, h HKID) {
			//log.Printf("begin %T\n", keygeterInstance)
			k, err := keygeterInstance.getKey(h)
			//log.Printf("end %T\n", keygeterInstance)
			if err == nil {
				datach <- k
				return
			} else {
				errorch <- err
				return
			}
		}(datach, errorch, h)
	}
	for {
		select {
		case b := <-datach:
			//log.Printf("select loop data%v\n", b)
			//if something { //How to Verifiy key?
			return b, nil
			//}
			//return nil, errors.New("Key Verifiy Failed")
		case err := <-errorch:
			//log.Printf("select loop error%v\n", err)
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
