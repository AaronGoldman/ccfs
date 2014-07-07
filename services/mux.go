package services

import (
	"bytes"
	"fmt"
	"log"

	"github.com/AaronGoldman/ccfs/objects"
)

//GetBlob looks up blobs by their HCIDs.
func GetBlob(h objects.HCID) (objects.Blob, error) {
	if h == nil {
		log.Printf("GetBlob(nil)")
		return nil, fmt.Errorf("nil pased in to GetBlob")
	}
	datach := make(chan objects.Blob, len(blobgeters))
	errorch := make(chan error, len(blobgeters))
	for _, rangeblobgeterInstance := range blobgeters {
		go func(
			blobgeterInstance blobgeter,
			datach chan objects.Blob,
			errorch chan error,
			h objects.HCID,
		) {
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
			if b != nil && bytes.Equal(b.Hash(), h) {
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

func PostList(l objects.List) (err error) {
	return PostBlob(objects.Blob(l.Bytes()))
}

//GetCommit retreves the newest commit for a given HKID
func GetCommit(h objects.HKID) (objects.Commit, error) {
	datach := make(chan objects.Commit, len(commitgeters))
	errorch := make(chan error, len(commitgeters))
	for _, rangecommitgeterInstance := range commitgeters {
		go func(
			commitgeterInstance commitgeter,
			datach chan objects.Commit,
			errorch chan error,
			h objects.HKID,
		) {
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
			return objects.Commit{}, fmt.Errorf("Commit Verifiy Failed")
		case err := <-errorch:
			if err.Error() == "GetCommit Timeout" {
				return objects.Commit{}, err
			} else {
				log.Println(err)
			}
		}
	}
}

//GetTag retreves the newest tag for a given HKID and name segment
func GetTag(h objects.HKID, namesegment string) (objects.Tag, error) {
	datach := make(chan objects.Tag, len(taggeters))
	errorch := make(chan error, len(taggeters))
	for _, rangetaggeterInstance := range taggeters {
		go func(
			taggeterInstance taggeter,
			datach chan objects.Tag,
			errorch chan error,
			h objects.HKID,
			namesegment string,
		) {
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
			return objects.Tag{}, fmt.Errorf("Tag Verifiy Failed")
		case err := <-errorch:
			if err.Error() == "GetTag Timeout" {
				return objects.Tag{}, err
			} else {
				log.Println(err)
			}
		}
	}
}

//GetTags retreves the newest tag for each name segment for a given HKID
func GetTags(h objects.HKID) (tags []objects.Tag, err error) {
	datach := make(chan objects.Tag, len(tagsgeters))
	errorch := make(chan error, len(tagsgeters))
	for _, rangetagsgeterInstance := range tagsgeters {
		go func(
			tagsgeterInstance tagsgeter,
			datach chan objects.Tag,
			errorch chan error,
			h objects.HKID,
		) {
			t, err := tagsgeterInstance.GetTags(h)
			if err == nil {
				for _, tag := range t {
					datach <- tag
				}
				return
			} else {
				errorch <- err
				return
			}
		}(rangetagsgeterInstance, datach, errorch, h)
	}
	for {
		select {
		case t := <-datach:
			if t.Verify() {
				tags = append(tags, t)
			} else {
				fmt.Println("Tag Verifiy Failed")
			}
		case err := <-errorch:
			if err.Error() == "GetTags Timeout" {
				if len(tags) > 0 {
					return tags, nil
				} else {
					return tags, err
				}
			} else {
				log.Println(err)
			}
		}
	}
}

//GetKey uses the HKID to lookup the PrivateKey.
func GetKey(h objects.HKID) (*objects.PrivateKey, error) {
	datach := make(chan objects.Blob, len(keygeters))
	errorch := make(chan error, len(keygeters))
	for _, rangekeygeterInstance := range keygeters {
		go func(
			keygeterInstance keygeter,
			datach chan objects.Blob,
			errorch chan error,
			h objects.HKID,
		) {
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
			privkey, err := objects.PrivteKeyFromBytes(b)
			if bytes.Equal(privkey.Hkid(), h) && privkey.Verify() {
				return privkey, err
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
func PostBlob(b objects.Blob) (err error) {
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
func PostCommit(c objects.Commit) (err error) {
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
func PostKey(p *objects.PrivateKey) (err error) {
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
func PostTag(t objects.Tag) (err error) {
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
