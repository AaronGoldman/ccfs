//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file
package kademliadht

import (
	"bytes"
	"testing"

	"github.com/AaronGoldman/ccfs/objects"
)

func TestKademliaserviceBlob(t *testing.T) {
	indata := objects.Blob("TestPostData")
	Instance.PostBlob(indata)
	outdata, err := Instance.GetBlob(indata.Hash())
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Errorf("\nExpected:%s\nGot:%s\nErr:%s", indata, outdata, err)
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
		domainHkid,
	)
	Instance.PostTag(intag)
	outtag, err := Instance.GetTag(
		domainHkid,
		"TestPostBlob",
	)
	if err != nil || !outtag.Verify() /*|| !bytes.Equal(outtag.Hkid(), domainHkid)*/ {
		t.Errorf(
			"\nExpected:%s\nGot:%s\nErr:%v\nVerify:%v",
			intag,
			outtag,
			err,
			outtag.Verify(),
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
	if err != nil || !outcommit.Verify() {
		t.Errorf("\nExpected:%v\nGot:%v\nErr:%s\nVerify:%t", incommit, outcommit, err, outcommit.Verify())
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
