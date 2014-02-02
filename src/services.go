// services.go
package main

import ()

var blobgeters = []blobgeter{}
var commitgeters = []commitgeter{}
var taggeters = []taggeter{}
var keygeters = []keygeter{}
var blobposters = []blobposter{}
var commitposters = []commitposter{}
var tagposters = []tagposter{}
var keyposters = []keyposter{}

func Registerblobgeter(service blobgeter)       { blobgeters = append(blobgeters, service) }
func Registercommitgeter(service commitgeter)   { commitgeters = append(commitgeters, service) }
func Registertaggeter(service taggeter)         { taggeters = append(taggeters, service) }
func Registerkeygeter(service keygeter)         { keygeters = append(keygeters, service) }
func Registerblobposter(service blobposter)     { blobposters = append(blobposters, service) }
func Registercommitposter(service commitposter) { commitposters = append(commitposters, service) }
func Registertagposter(service tagposter)       { tagposters = append(tagposters, service) }
func Registerkeyposter(service keyposter)       { keyposters = append(keyposters, service) }

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
	GetBlob(HCID) (blob, error)
}
type commitgeter interface {
	GetCommit(HKID) (commit, error)
}
type taggeter interface {
	GetTag(h HKID, namesegment string) (tag, error)
}
type keygeter interface {
	GetKey(HKID) (blob, error)
}

type contentposter interface {
	blobposter
	commitposter
	tagposter
	keyposter
}
type blobposter interface {
	PostBlob(b blob) error
}
type commitposter interface {
	PostCommit(c commit) error
}
type tagposter interface {
	PostTag(t tag) error
}
type keyposter interface {
	PostKey(p *PrivateKey) error
}
