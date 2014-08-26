//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file
// services.go
package services

import (
	"github.com/AaronGoldman/ccfs/objects"
)

var blobgeters = []blobgeter{}
var commitgeters = []commitgeter{}
var taggeters = []taggeter{}
var tagsgeters = []tagsgeter{}
var keygeters = []keygeter{}
var blobposters = []blobposter{}
var commitposters = []commitposter{}
var tagposters = []tagposter{}
var keyposters = []keyposter{}

func Registerblobgeter(service blobgeter)       { blobgeters = append(blobgeters, service) }
func Registercommitgeter(service commitgeter)   { commitgeters = append(commitgeters, service) }
func Registertaggeter(service taggeter)         { taggeters = append(taggeters, service) }
func Registertagsgeter(service tagsgeter)       { tagsgeters = append(tagsgeters, service) }
func Registerkeygeter(service keygeter)         { keygeters = append(keygeters, service) }
func Registerblobposter(service blobposter)     { blobposters = append(blobposters, service) }
func Registercommitposter(service commitposter) { commitposters = append(commitposters, service) }
func Registertagposter(service tagposter)       { tagposters = append(tagposters, service) }
func Registerkeyposter(service keyposter)       { keyposters = append(keyposters, service) }
func Registercontentservice(service contentservice) {
	Registerblobgeter(service)
	Registercommitgeter(service)
	Registertaggeter(service)
	Registertagsgeter(service)
	Registerkeygeter(service)
	Registerblobposter(service)
	Registercommitposter(service)
	Registertagposter(service)
	Registerkeyposter(service)
}

type contentservice interface {
	blobgeter
	commitgeter
	taggeter
	tagsgeter
	keygeter
	blobposter
	commitposter
	tagposter
	keyposter
}

type blobgeter interface {
	GetBlob(objects.HCID) (objects.Blob, error)
}
type commitgeter interface {
	GetCommit(objects.HKID) (objects.Commit, error)
}
type taggeter interface {
	GetTag(h objects.HKID, namesegment string) (objects.Tag, error)
}
type tagsgeter interface {
	GetTags(h objects.HKID) ([]objects.Tag, error)
}
type keygeter interface {
	GetKey(objects.HKID) (objects.Blob, error)
}

type blobposter interface {
	PostBlob(b objects.Blob) error
}
type commitposter interface {
	PostCommit(c objects.Commit) error
}
type tagposter interface {
	PostTag(t objects.Tag) error
}
type keyposter interface {
	PostKey(p *objects.PrivateKey) error
}
