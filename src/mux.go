package main

import (
	"fmt"
	"log"
)

//GetBlob looks up blobs by their HCIDs.
func GetBlob(h HCID) (blob, error) {
	datach := make(chan blob, len(blobgeters))
	errorch := make(chan error, len(blobgeters))
	for _, rangeblobgeterInstance := range blobgeters {
		go func(blobgeterInstance blobgeter, datach chan blob, errorch chan error, h HCID) {
			b, err := blobgeterInstance.GetBlob(h)
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
			return nil, fmt.Errorf("Blob Verifiy Failed")
		case err := <-errorch:
			if err.Error() == "GetBlob Timeout" {
				return nil, err
			} else {
				log.Println(err)
			}
		}
	}
}

func PostList(l list) (err error) {
	return PostBlob(blob(l.Bytes()))
}

//GetCommit retreves the newest commit for a given HKID
func GetCommit(h HKID) (commit, error) {
	datach := make(chan commit, len(commitgeters))
	errorch := make(chan error, len(commitgeters))
	for _, rangecommitgeterInstance := range commitgeters {
		go func(commitgeterInstance commitgeter, datach chan commit, errorch chan error, h HKID) {
			c, err := commitgeterInstance.GetCommit(h)
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
			if c.Verify() {
				return c, nil
			}
			return commit{}, fmt.Errorf("Commit Verifiy Failed")
		case err := <-errorch:
			if err.Error() == "GetCommit Timeout" {
				return commit{}, err
			} else {
				log.Println(err)
			}
		}
	}
}

//GetTag retreves the newest tag for a given HKID and name segment
func GetTag(h HKID, namesegment string) (tag, error) {
	datach := make(chan tag, len(taggeters))
	errorch := make(chan error, len(taggeters))
	for _, rangetaggeterInstance := range taggeters {
		go func(taggeterInstance taggeter, datach chan tag, errorch chan error, h HKID, namesegment string) {
			t, err := taggeterInstance.GetTag(h, namesegment)
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
			if t.Verify() {
				return t, nil
			}
			return tag{}, fmt.Errorf("Tag Verifiy Failed")
		case err := <-errorch:
			if err.Error() == "GetTag Timeout" {
				return tag{}, err
			} else {
				log.Println(err)
			}
		}
	}
}

//GetKey uses the HKID to lookup the PrivateKey.
func GetKey(h HKID) (*PrivateKey, error) {
	datach := make(chan blob, len(keygeters))
	errorch := make(chan error, len(keygeters))
	for _, rangekeygeterInstance := range keygeters {
		go func(keygeterInstance keygeter, datach chan blob, errorch chan error, h HKID) {
			k, err := keygeterInstance.GetKey(h)
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
			privkey := PrivteKeyFromBytes(b)
			if privkey.Hkid().String() == h.String() && privkey.Verify() {
				return PrivteKeyFromBytes(b), nil
			} else {
				log.Println("Key Verifiy Failed")
			}
		case err := <-errorch:
			if err.Error() == "GetKey Timeout" {
				return nil, err
			} else {
				log.Println(err)
			}
		}
	}
}

//release blob to storage
func PostBlob(b blob) (err error) {
	var firsterr error
	for _, service := range blobposters {
		err := service.PostBlob(b)
		if err != nil {
			firsterr = err
		}
	}
	return firsterr
	//return localfileserviceInstance.PostBlob(b)
}

//release commit to storage
func PostCommit(c commit) (err error) {
	var firsterr error
	for _, service := range commitposters {
		err := service.PostCommit(c)
		if err != nil {
			firsterr = err
		}
	}
	return firsterr
	//return localfileserviceInstance.PostCommit(c)
}

//release key to storage
func PostKey(p *PrivateKey) (err error) {
	var firsterr error
	for _, service := range keyposters {
		err := service.PostKey(p)
		if err != nil {
			firsterr = err
		}
	}
	return firsterr
	//return localfileserviceInstance.PostKey(p)
}

//release tag to storage
func PostTag(t tag) (err error) {
	var firsterr error
	for _, service := range tagposters {
		err := service.PostTag(t)
		if err != nil {
			firsterr = err
		}
	}
	return firsterr
	//return localfileserviceInstance.PostTag(t)
}
