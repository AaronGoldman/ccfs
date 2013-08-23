package main

import (
	"errors"
	"time"
)

func GetCommit(h HKID) (commit, error) {
	commit_chan := make(chan commit)
	go func(commit_chan chan commit) {
		c, err := localfileservice_GetCommit(h)
		if err == nil {
			commit_chan <- c
		}
	}(commit_chan)
	select {
	case b := <-commit_chan:
		return b, nil
	case <-time.After(time.Second):
		return commit{}, errors.New("GetCommit Timeout")
	}
	//return commit{}, errors.New("select mux.go")
}

func GetTag(h HKID, namesegment string) (Tag, error) {
	tag_chan := make(chan Tag)
	go func(tag_chan chan Tag) {
		c, err := localfileservice_GetTag(h, namesegment)
		if err == nil {
			tag_chan <- c
		}
	}(tag_chan)
	select {
	case b := <-tag_chan:
		return b, nil
	case <-time.After(time.Second):
		return Tag{}, errors.New("GetTag Timeout")
	}
	//return Tag{}, errors.New("select mux.go")
}

func GetBlob(h HCID) (blob, error) {
	blob_chan := make(chan blob)
	go func(blob_chan chan blob) {
		c, err := localfileservice_GetBlob(h)
		if err == nil {
			blob_chan <- c
		}
	}(blob_chan)
	select {
	case b := <-blob_chan:
		return b, nil
	case <-time.After(time.Second):
		return nil, errors.New("GetBlob Timeout")
	}
	//return nil, errors.New("select mux.go")
}

func GetKey(h HKID) (blob, error) {
	blob_chan := make(chan blob)
	go func(blob_chan chan blob) {
		c, err := localfileservice_GetKey(h)
		if err == nil {
			blob_chan <- c
		}
	}(blob_chan)
	select {
	case b := <-blob_chan:
		return b, nil
	case <-time.After(time.Second):
		return nil, errors.New("GetBlob Timeout")
	}
	//return nil, errors.New("select mux.go")
}
