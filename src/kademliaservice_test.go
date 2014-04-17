package main

import (
	"bytes"
	//"log"
	"testing"
)

//TestKademliaServiceIsContentService tests if kademliaservice is contentservice
func TestKademliaServiceIsContentService(t *testing.T) {
	func(contentservice) {}(kademliaserviceInstance)
}

func TestKademliaserviceBlob(t *testing.T) {
	t.Skip()
	indata := blob("TestPostData")
	kademliaserviceInstance.PostBlob(indata)
	outdata, err := kademliaserviceInstance.GetBlob(indata.Hash())
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Errorf("\nExpected:%s\nGot:%s\nErr:%s", indata, outdata, err)
	}
}

func TestKademliaserviceTag(t *testing.T) {
	t.Skip()
	b := blob("TestPostData")
	domainHkid := hkidFromDString("36353776900433612923412235249557"+
		"5547801975514185453610009798109590341135752484880821676711220739"+
		"4029990114056594565164287698180880449563881968956877896844137", 10)
	intag := NewTag(
		b.Hash(),
		"blob",
		"TestPostBlob",
		domainHkid,
	)
	kademliaserviceInstance.PostTag(intag)
	outtag, err := kademliaserviceInstance.GetTag(
		domainHkid,
		"TestPostBlob",
	)
	if err != nil || !outtag.Verify() /*|| !bytes.Equal(outtag.Hkid(), domainHkid)*/ {
		t.Errorf("\nExpected:%s\nGot:%s\nErr:%v\nVerify:%v", intag, outtag, err, outtag.Verify())
	}
}

func TestKademliaserviceCommit(t *testing.T) {
	t.Skipf("")
	b := blob("TestPostData")
	l := NewList(
		b.Hash(),
		"blob",
		"TestPostBlob",
	)
	repoHkid := hkidFromDString("64171129167204289916774847858432"+
		"1039643124642934014944704416438487015947986633802511102841255411"+
		"2620702113155684804978643917650455537680636225253952875765474", 10)
	incommit := NewCommit(l.Hash(), repoHkid)
	kademliaserviceInstance.PostCommit(incommit)
	outcommit, err := kademliaserviceInstance.GetCommit(repoHkid)
	if err != nil || !outcommit.Verify() {
		t.Errorf("\nExpected:%s\nGot:%s\nErr:%s\nVerify:%s", incommit, outcommit, err, outcommit.Verify())
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
