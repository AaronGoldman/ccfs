// multicastservice.go
package main

import ()

type multicastservice struct{}

func (m multicastservice) getBlob(h HCID) (b blob, err error) {
	return b, err

}

func (m multicastservice) getCommit(h HKID) (c commit, err error) {
	return c, err
}

func (m multicastservice) getTag(h HKID, namesegment string) (t tag, err error) {
	return t, err

}

func (m multicastservice) getKey(h HKID) (b blob, err error) {
	return b, err
}

var multicastserviceInstance multicastservice = multicastservice{}
