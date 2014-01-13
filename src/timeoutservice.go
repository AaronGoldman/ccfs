package main

import (
	"errors"
	"time"
)

type timeoutservice struct{}

func (timeoutservice) getBlob(HCID) (blob, error) {
	time.Sleep(time.Second)
	return blob{}, errors.New("GetBlob Timeout")
}
func (timeoutservice) getCommit(HKID) (commit, error) {
	time.Sleep(time.Second)
	return commit{}, errors.New("GetCommit Timeout")
}
func (timeoutservice) getTag(h HKID, namesegment string) (tag, error) {
	time.Sleep(time.Second)
	return tag{}, errors.New("GetTag Timeout")
}
func (timeoutservice) getKey(HKID) (blob, error) {
	time.Sleep(time.Second)
	return blob{}, errors.New("GetKey Timeout")
}

var timeoutserviceInstance timeoutservice
