package main

import (
	"fmt"
	"time"
)

type timeoutservice struct{}

func (timeoutservice) GetBlob(HCID) (blob, error) {
	time.Sleep(time.Second)
	return blob{}, fmt.Errorf("GetBlob Timeout")
}
func (timeoutservice) GetCommit(HKID) (commit, error) {
	time.Sleep(time.Second)
	return commit{}, fmt.Errorf("GetCommit Timeout")
}
func (timeoutservice) GetTag(h HKID, namesegment string) (tag, error) {
	time.Sleep(time.Second)
	return tag{}, fmt.Errorf("GetTag Timeout")
}
func (timeoutservice) GetKey(HKID) (blob, error) {
	time.Sleep(time.Second)
	return blob{}, fmt.Errorf("GetKey Timeout")
}

var timeoutserviceInstance timeoutservice

func init() {
	Registerblobgeter(timeoutserviceInstance)
	Registercommitgeter(timeoutserviceInstance)
	Registertaggeter(timeoutserviceInstance)
	Registerkeygeter(timeoutserviceInstance)
}
