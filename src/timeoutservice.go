package main

import (
	"fmt"
	"time"
)

type timeoutservice struct{}

func (timeoutservice) getBlob(HCID) (blob, error) {
	time.Sleep(time.Second)
	return blob{}, fmt.Errorf("GetBlob Timeout")
}
func (timeoutservice) getCommit(HKID) (commit, error) {
	time.Sleep(time.Second)
	return commit{}, fmt.Errorf("GetCommit Timeout")
}
func (timeoutservice) getTag(h HKID, namesegment string) (tag, error) {
	time.Sleep(time.Second)
	return tag{}, fmt.Errorf("GetTag Timeout")
}
func (timeoutservice) getKey(HKID) (blob, error) {
	time.Sleep(time.Second)
	return blob{}, fmt.Errorf("GetKey Timeout")
}

var timeoutserviceInstance timeoutservice
