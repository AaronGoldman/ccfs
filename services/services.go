//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

//Package services implements the common function for all the ccfs services
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

//Registerblobgeter adds a blobgeter to blobgeters
func Registerblobgeter(service blobgeter) {
	if blobgeters == nil {
		blobgeters = make(map[string]blobgeter)
	}
	blobgeters[service.ID()] = service
}

//DeRegisterblobgeter removes a blobgeter from blobgeters
func DeRegisterblobgeter(service blobgeter) {
	delete(blobgeters, service.ID())
}

//Registercommitgeter adds a commitgeter to commitgeters
func Registercommitgeter(service commitgeter) {
	if commitgeters == nil {
		commitgeters = make(map[string]commitgeter)
	}
	commitgeters[service.ID()] = service
}

//DeRegistercommitgeter removes a commitgeter from commitgeters
func DeRegistercommitgeter(service commitgeter) {
	delete(commitgeters, service.ID())
}

//Registertaggeter adds a taggeter to taggeters
func Registertaggeter(service taggeter) {
	if taggeters == nil {
		taggeters = make(map[string]taggeter)
	}
	taggeters[service.ID()] = service
}

//DeRegistertaggeter removes a taggeter from taggeters
func DeRegistertaggeter(service taggeter) {
	delete(taggeters, service.ID())
}

//Registertagsgeter adds a tagsgeter to tagsgeters
func Registertagsgeter(service tagsgeter) {
	if tagsgeters == nil {
		tagsgeters = make(map[string]tagsgeter)
	}
	tagsgeters[service.ID()] = service
}

//DeRegistertagsgeter removes a tagsgeter from tagsgeters
func DeRegistertagsgeter(service tagsgeter) {
	delete(tagsgeters, service.ID())
}

//Registerkeygeter adds a keygeter the keygeters
func Registerkeygeter(service keygeter) {
	if keygeters == nil {
		keygeters = make(map[string]keygeter)
	}
	keygeters[service.ID()] = service
}

//DeRegisterkeygeter removes a keygeter from keygeters
func DeRegisterkeygeter(service keygeter) {
	delete(keygeters, service.ID())
}

//Registerblobposter adds a blobposter to blobposters
func Registerblobposter(service blobposter) {
	if blobposters == nil {
		blobposters = make(map[string]blobposter)
	}
	blobposters[service.ID()] = service
}

//DeRegisterblobposter removes a blobposter to blobposters
func DeRegisterblobposter(service blobposter) {
	delete(blobposters, service.ID())
}

//Registercommitposter adds a commitposter to commitposters
func Registercommitposter(service commitposter) {
	if commitposters == nil {
		commitposters = make(map[string]commitposter)
	}
	commitposters[service.ID()] = service
}

//DeRegistercommitposter removes a commitposter from commitposters
func DeRegistercommitposter(service commitposter) {
	delete(commitposters, service.ID())
}

//Registertagposter adds a tagposter to tagposters
func Registertagposter(service tagposter) {
	if tagposters == nil {
		tagposters = make(map[string]tagposter)
	}
	tagposters[service.ID()] = service
}

//DeRegistertagposter removes a tagposter from tagposters
func DeRegistertagposter(service tagposter) {
	delete(tagposters, service.ID())
}

//Registerkeyposter adds a keyposter to keyposters
func Registerkeyposter(service keyposter) {
	if keyposters == nil {
		keyposters = make(map[string]keyposter)
	}
	keyposters[service.ID()] = service
}

//DeRegisterkeyposter removes a keyposter from keyposters
func DeRegisterkeyposter(service keyposter) {
	delete(keyposters, service.ID())
}

type idgeter interface {
	ID() string
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
