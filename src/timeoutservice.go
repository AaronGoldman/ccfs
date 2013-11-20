package main

import (
	"errors"
	"time"
)

func timeoutservice_blobgeter(datach chan blob, errorch chan error, h HCID) {
	time.Sleep(time.Second)
	errorch <- errors.New("GetBlob Timeout")
}

func timeoutservice_taggeter(datach chan Tag, errorch chan error, h HKID, namesegment string) {
	time.Sleep(time.Second)
	errorch <- errors.New("GetTag Timeout")
}

func timeoutservice_commitgeter(datach chan commit, errorch chan error, h HKID) {
	time.Sleep(time.Second)
	errorch <- errors.New("GetCommit Timeout")
}

func timeoutservice_keygeter(datach chan blob, errorch chan error, h HKID) {
	time.Sleep(time.Second)
	errorch <- errors.New("GetKey Timeout")
}
