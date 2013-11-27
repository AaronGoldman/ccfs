package main

import (
	"errors"
	//"time"
)

var blobgeters = []func(chan blob, chan error, HCID){
	localfileservice_blobgeter,
	timeoutservice_blobgeter,
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

var commitgeters = []func(chan commit, chan error, HKID){
	localfileservice_commitgeter,
	timeoutservice_commitgeter,
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

var taggeters = []func(chan Tag, chan error, HKID, string){
	localfileservice_taggeter,
	timeoutservice_taggeter,
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

var keygeters = []func(chan blob, chan error, HKID){
	localfileservice_keygeter,
	timeoutservice_keygeter,
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
