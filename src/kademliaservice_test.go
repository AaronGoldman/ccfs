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

func TestKademliaservicePostBlob(t *testing.T) {
	t.Skipf("")
	b := blob("TestPostData")
	kademliaserviceInstance.PostBlob(b)
}
func TestKademliaservicePostTag(t *testing.T) {
	t.Skipf("")
	b := blob("TestPostData")
	domainHkid := hkidFromDString("36353776900433612923412235249557"+
		"5547801975514185453610009798109590341135752484880821676711220739"+
		"4029990114056594565164287698180880449563881968956877896844137", 10)
	tt := NewTag(
		b.Hash(),
		"blob",
		"TestPostBlob",
		domainHkid,
	)
	kademliaserviceInstance.PostTag(tt)
}
func TestKademliaservicePostCommit(t *testing.T) {
	t.Skipf("")
	//kademliaserviceInstance.PostCommit()
}
func TestKademliaservicePostKey(t *testing.T) {
	t.Skipf("")
	//kademliaserviceInstance.PostKey()
}
func TestKademliaserviceGetBlob(t *testing.T) {
	t.Skipf("")
	indata := blob("TestPostData")
	outdata, err := kademliaserviceInstance.GetBlob(indata.Hash())
	if !bytes.Equal(indata, outdata) || err == nil {
		t.Errorf("Expected:%s\nGot:%s", string(indata), string(outdata))
	}
}
func TestKademliaserviceGetCommit(t *testing.T) {
	t.Skipf("")
	domainHkid := hkidFromDString("36353776900433612923412235249557"+
		"5547801975514185453610009798109590341135752484880821676711220739"+
		"4029990114056594565164287698180880449563881968956877896844137", 10)
	c, err := kademliaserviceInstance.GetCommit(domainHkid)
	if err != nil || !c.Verify() {
		t.Errorf("fail")
	}
}
func TestKademliaserviceGetTag(t *testing.T) {
	t.Skipf("")
	//kademliaserviceInstance.GetTag()
}
func TestKademliaserviceGetKey(t *testing.T) {
	t.Skipf("")
	//kademliaserviceInstance.GetKey()
}
