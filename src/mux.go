package main

import (
	"errors"
	"log"
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
	//appsscriptserviceInstance,
}
var taggeters = []taggeter{
	timeoutserviceInstance,
	localfileserviceInstance,
	//googledriveserviceInstance,
	//appsscriptserviceInstance,
}
var keygeters = []keygeter{
	timeoutserviceInstance,
	localfileserviceInstance,
	//googledriveserviceInstance,
	//appsscriptserviceInstance,
}

func GetBlob(h HCID) (blob, error) {
	datach := make(chan blob, len(blobgeters))
	errorch := make(chan error, len(blobgeters))
	for _, rangeblobgeterInstance := range blobgeters {
		go func(blobgeterInstance blobgeter, datach chan blob, errorch chan error, h HCID) {
			b, err := blobgeterInstance.getBlob(h)
			if err == nil {
				datach <- b
				return
			} else {
				errorch <- err
				return
			}
		}(rangeblobgeterInstance, datach, errorch, h)
	}
	for {
		select {
		case b := <-datach:
			if b != nil && b.Hash().Hex() == h.Hex() {
				return b, nil
			}
			return nil, errors.New("Blob Verifiy Failed")
		case err := <-errorch:
			if err.Error() == "GetBlob Timeout" {
				return nil, err
			} else {
				log.Println(err)
			}
		}
	}
}

func GetCommit(h HKID) (commit, error) {
	datach := make(chan commit, len(commitgeters))
	errorch := make(chan error, len(commitgeters))
	for _, rangecommitgeterInstance := range commitgeters {
		go func(commitgeterInstance commitgeter, datach chan commit, errorch chan error, h HKID) {
			c, err := commitgeterInstance.getCommit(h)
			if err == nil {
				datach <- c
				return
			} else {
				errorch <- err
				return
			}
		}(rangecommitgeterInstance, datach, errorch, h)
	}
	for {
		select {
		case c := <-datach:
			if c.Verifiy() {
				return c, nil
			}
			return commit{}, errors.New("Commit Verifiy Failed")
		case err := <-errorch:
			if err.Error() == "GetCommit Timeout" {
				return commit{}, err
			} else {
				log.Println(err)
			}
		}
	}
}

func GetTag(h HKID, namesegment string) (tag, error) {
	datach := make(chan tag, len(taggeters))
	errorch := make(chan error, len(taggeters))
	for _, rangetaggeterInstance := range taggeters {
		go func(taggeterInstance taggeter, datach chan tag, errorch chan error, h HKID, namesegment string) {
			t, err := taggeterInstance.getTag(h, namesegment)
			if err == nil {
				datach <- t
				return
			} else {
				errorch <- err
				return
			}
		}(rangetaggeterInstance, datach, errorch, h, namesegment)
	}
	for {
		select {
		case t := <-datach:
			if t.Verifiy() {
				return t, nil
			}
			return tag{}, errors.New("Tag Verifiy Failed")
		case err := <-errorch:
			if err.Error() == "GetTag Timeout" {
				return tag{}, err
			} else {
				log.Println(err)
			}
		}
	}
}

func GetKey(h HKID) (blob, error) {
	datach := make(chan blob, len(keygeters))
	errorch := make(chan error, len(keygeters))
	for _, rangekeygeterInstance := range keygeters {
		go func(keygeterInstance keygeter, datach chan blob, errorch chan error, h HKID) {
			k, err := keygeterInstance.getKey(h)
			if err == nil {
				datach <- k
				return
			} else {
				errorch <- err
				return
			}
		}(rangekeygeterInstance, datach, errorch, h)
	}
	for {
		select {
		case b := <-datach:
			//if something { //How to Verifiy key?
			return b, nil
			//}
			//return nil, errors.New("Key Verifiy Failed")
		case err := <-errorch:
			if err.Error() == "GetKey Timeout" {
				return nil, err
			} else {
				log.Println(err)
			}
		}
	}
}
