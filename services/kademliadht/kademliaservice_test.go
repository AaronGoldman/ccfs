//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file
package kademliadht

import (
	"bytes"
	"log"
	"testing"

	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
	"github.com/AaronGoldman/ccfs/services/localfile"
	"github.com/AaronGoldman/ccfs/services/timeout"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	localfile.Start()
	timeout.Start()
	Start()
	objects.RegisterGeterPoster(
		services.GetPublicKeyForHkid,
		services.GetPrivateKeyForHkid,
		services.PostKey,
		services.PostBlob,
	)
}

func TestKademliaserviceBlob(t *testing.T) {
	indata := objects.Blob("TestPostData")
	Instance.PostBlob(indata)
	outdata, err := Instance.GetBlob(indata.Hash())

	if err != nil {
		t.Errorf("\nErr:%s", err)
	}

	if !bytes.Equal(indata, outdata) {
		t.Errorf("\nExpected:%s\nGot:%s", indata, outdata)
	}
}

func TestKademliaserviceTag(t *testing.T) {
	b := objects.Blob("TestPostData")
	domainHkid := objects.HkidFromDString("36353776900433612923412235249557"+
		"5547801975514185453610009798109590341135752484880821676711220739"+
		"4029990114056594565164287698180880449563881968956877896844137", 10)
	intag := objects.NewTag(
		b.Hash(),
		"blob",
		"TestPostBlob",
		[]objects.HCID{},
		domainHkid,
	)
	Instance.PostTag(intag)
	outtag, err := Instance.GetTag(
		domainHkid,
		"TestPostBlob",
	)
	if err != nil {
		t.Errorf("Get Tag Err: %s", err)
	}

	if !outtag.Verify() /*|| !bytes.Equal(outtag.Hkid(), domainHkid)*/ {
		t.Errorf(
			"\nVerify:%v",
			outtag.Verify(),
		)
	}

	if !bytes.Equal(intag.Bytes(), outtag.Bytes()) {
		t.Errorf(
			"\nExpected:%s\nGot:%s",
			intag,
			outtag,
		)
	}
}

func TestKademliaserviceCommit(t *testing.T) {
	b := objects.Blob("TestPostData")
	l := objects.NewList(
		b.Hash(),
		"blob",
		"TestPostBlob",
	)
	repoHkid := objects.HkidFromDString("64171129167204289916774847858432"+
		"1039643124642934014944704416438487015947986633802511102841255411"+
		"2620702113155684804978643917650455537680636225253952875765474", 10)
	incommit := objects.NewCommit(l.Hash(), repoHkid)
	Instance.PostCommit(incommit)
	outcommit, err := Instance.GetCommit(repoHkid)

	if err != nil {
		t.Fatalf("\nGet Commit Err:%s\n", err)
	}
	if !outcommit.Verify() {
		t.Fatalf("\nVerify:%t", outcommit.Verify())
	}
	if !bytes.Equal(incommit.Bytes(), outcommit.Bytes()) {
		t.Fatalf("\nExpected:%v\nGot:%v", incommit, outcommit)
	}
}

//func TestKademliaserviceKey(t *testing.T) {
//	//t.Skipf("")
//	repoHkid := hkidFromDString("64171129167204289916774847858432"+
//		"1039643124642934014944704416438487015947986633802511102841255411"+
//		"2620702113155684804978643917650455537680636225253952875765474", 10)
//	inkey, _ := GetKey(repoHkid)
//	kademliaserviceInstance.PostKey(inkey)
//	outkey, err := kademliaserviceInstance.GetKey(repoHkid)
//	if err != nil || !bytes.Equal(repoHkid, outkey.Hash()) {
//		t.Errorf("\nExpected:%s\nGot:%s", inkey, outkey)
//	}
//}
