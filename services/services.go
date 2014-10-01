//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

// services.go
//Package services is the common function for all the ccfs services

package services

import (
	"github.com/AaronGoldman/ccfs/objects"
)

var blobgeters = map[string]blobgeter{}
var commitgeters = map[string]commitgeter{}
var taggeters = map[string]taggeter{}
var tagsgeters = map[string]tagsgeter{}
var keygeters = map[string]keygeter{}
var blobposters = map[string]blobposter{}
var commitposters = map[string]commitposter{}
var tagposters = map[string]tagposter{}
var keyposters = map[string]keyposter{}

//Registerblobgeter adds the pased in blobgeter to the blobgeters that are muxed
func Registerblobgeter(service blobgeter) {
	if blobgeters == nil {
		blobgeters = make(map[string]blobgeter)
	}
	blobgeters[service.GetId()] = service
}

func DeRegisterblobgeter(service blobgeter) {
	delete(blobgeters, service.GetId())
}

//Registercommitgeter adds the pased in commitgeter to the commitgeters that are muxed
func Registercommitgeter(service commitgeter) {
	if commitgeters == nil {
		commitgeters = make(map[string]commitgeter)
	}
	commitgeters[service.GetId()] = service
}

func DeRegistercommitgeter(service commitgeter) {
	delete(commitgeters, service.GetId())
}

//Registertaggeter adds the pased in taggeter to the taggeters that are muxed
func Registertaggeter(service taggeter) {
	if taggeters == nil {
		taggeters = make(map[string]taggeter)
	}
	taggeters[service.GetId()] = service
}

func DeRegistertaggeter(service taggeter) {
	delete(taggeters, service.GetId())
}

//Registertagsgeter adds the pased in tagsgeter to the tagsgeters that are muxed
func Registertagsgeter(service tagsgeter) {
	if tagsgeters == nil {
		tagsgeters = make(map[string]tagsgeter)
	}
	tagsgeters[service.GetId()] = service
}

func DeRegistertagsgeter(service tagsgeter) {
	delete(tagsgeters, service.GetId())
}

//Registerkeygeter adds the pased in keygeter to the keygeters that are muxed
func Registerkeygeter(service keygeter) {
	if keygeters == nil {
		keygeters = make(map[string]keygeter)
	}
	keygeters[service.GetId()] = service
}

func DeRegisterkeygeter(service keygeter) {
	delete(keygeters, service.GetId())
}

//Registerblobposter adds the pased in blobposter to the blobposters that are muxed
func Registerblobposter(service blobposter) {
	if blobposters == nil {
		blobposters = make(map[string]blobposter)
	}
	blobposters[service.GetId()] = service
}

func DeRegisterblobposter(service blobposter) {
	delete(blobposters, service.GetId())
}

//Registercommitposter adds the pased in commitposter to the commitposters that are muxed
func Registercommitposter(service commitposter) {
	if commitposters == nil {
		commitposters = make(map[string]commitposter)
	}
	commitposters[service.GetId()] = service
}

func DeRegistercommitposter(service commitposter) {
	delete(commitposters, service.GetId())
}

//Registertagposter adds the pased in tagposter to the tagposters that are muxed
func Registertagposter(service tagposter) {
	if tagposters == nil {
		tagposters = make(map[string]tagposter)
	}
	tagposters[service.GetId()] = service
}

func DeRegistertagposter(service tagposter) {
	delete(tagposters, service.GetId())
}

//Registerkeyposter adds the pased in keyposter to the keyposters that are muxed
func Registerkeyposter(service keyposter) {
	if keyposters == nil {
		keyposters = make(map[string]keyposter)
	}
	keyposters[service.GetId()] = service
}

func DeRegisterkeyposter(service keyposter) {
	delete(keyposters, service.GetId())
}

//Registercontentservice Registers the service with all the content services
//func Registercontentservice(service contentservice) {
//	Registerblobgeter(service)
//	Registercommitgeter(service)
//	Registertaggeter(service)
//	Registertagsgeter(service)
//	Registerkeygeter(service)
//	Registerblobposter(service)
//	Registercommitposter(service)
//	Registertagposter(service)
//	Registerkeyposter(service)
//}

//type contentservice interface {
//	idgeter
//	blobgeter
//	commitgeter
//	taggeter
//	tagsgeter
//	keygeter
//	blobposter
//	commitposter
//	tagposter
//	keyposter
//}

type idgeter interface {
	GetId() string
}

type blobgeter interface {
	GetBlob(objects.HCID) (objects.Blob, error)
	idgeter
}
type commitgeter interface {
	GetCommit(objects.HKID) (objects.Commit, error)
	idgeter
}
type taggeter interface {
	GetTag(h objects.HKID, namesegment string) (objects.Tag, error)
	idgeter
}
type tagsgeter interface {
	GetTags(h objects.HKID) ([]objects.Tag, error)
	idgeter
}
type keygeter interface {
	GetKey(objects.HKID) (objects.Blob, error)
	idgeter
}

type blobposter interface {
	PostBlob(b objects.Blob) error
	idgeter
}
type commitposter interface {
	PostCommit(c objects.Commit) error
	idgeter
}
type tagposter interface {
	PostTag(t objects.Tag) error
	idgeter
}
type keyposter interface {
	PostKey(p *objects.PrivateKey) error
	idgeter
}
