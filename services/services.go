//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

// services.go
//Package services is the common function for all the ccfs services

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

//Registerblobgeter adds the pased in blobgeter to the blobgeters that are muxed
func Registerblobgeter(service blobgeter) { blobgeters = append(blobgeters, service) }

//Registercommitgeter adds the pased in commitgeter to the commitgeters that are muxed
func Registercommitgeter(service commitgeter) { commitgeters = append(commitgeters, service) }

//Registertaggeter adds the pased in taggeter to the taggeters that are muxed
func Registertaggeter(service taggeter) { taggeters = append(taggeters, service) }

//Registertagsgeter adds the pased in tagsgeter to the tagsgeters that are muxed
func Registertagsgeter(service tagsgeter) { tagsgeters = append(tagsgeters, service) }

//Registerkeygeter adds the pased in keygeter to the keygeters that are muxed
func Registerkeygeter(service keygeter) { keygeters = append(keygeters, service) }

//Registerblobposter adds the pased in blobposter to the blobposters that are muxed
func Registerblobposter(service blobposter) { blobposters = append(blobposters, service) }

//Registercommitposter adds the pased in commitposter to the commitposters that are muxed
func Registercommitposter(service commitposter) { commitposters = append(commitposters, service) }

//Registertagposter adds the pased in tagposter to the tagposters that are muxed
func Registertagposter(service tagposter) { tagposters = append(tagposters, service) }

//Registerkeyposter adds the pased in keyposter to the keyposters that are muxed
func Registerkeyposter(service keyposter) { keyposters = append(keyposters, service) }

//Registercontentservice Registers the service with all the content services
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
