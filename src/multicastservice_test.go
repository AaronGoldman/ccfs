// multicastservice_test
package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

func TestMulticastservice_GetBlob(t *testing.T) {
	//t.Skipf("Come back to this test")
	go BlobServerStart()
	AnswerKey := []struct {
		hcid     HCID
		response blob
	}{
		{blob([]byte("blob found")).Hash(), blob([]byte("blob found"))},
	}
	for _, answer := range AnswerKey {
		go func() {
			time.Sleep(1 * time.Millisecond)
			mcaddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:8000")
			multicastserviceInstance.receivemessage("{\"type\":\"blob\", \"hcid\": \"42cc3a4c4a9d9d3ee7de9322b45acb0e5a5c33550d9ad4791df6ae937a869e12\", \"URL\": \"/b/42cc3a4c4a9d9d3ee7de9322b45acb0e5a5c33550d9ad4791df6ae937a869e12\"}", mcaddr)
		}()

		output, err := multicastserviceInstance.GetBlob(answer.hcid)

		if err != nil {
			t.Errorf("Get Blob Failed \nError:%t", err)
		}

		if !bytes.Equal(output.Hash(), answer.hcid) || err != nil {
			t.Errorf("Make URL Failed \nExpected:%s \nGot: %s", answer.response, output)
		}

	}
}

func TestMulticastservice_GetCommit(t *testing.T) {
	//t.Skipf("Come back to this test")
	hkid := hkidFromDString("5198719439877464148627795433286736285873678110640040333794349799294848737858561643942881983506066042818105864129178593001327423646717446545633525002218361750", 10)
	log.Printf("The HKID is, %s", hkid)

	b := blob([]byte("blob found"))
	l := NewList(b.Hash(), "blob", "Blobinlist")
	c := NewCommit(l.Hash(), hkid)
	localfileserviceInstance.PostCommit(c)
	go func() {
		time.Sleep(1 * time.Millisecond)
		mcaddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:8000")
		multicastserviceInstance.receivemessage(fmt.Sprintf("{\"type\":\"commit\", \"hkid\": \"9bd1b3c9aeda7025068319c0a4af1d2b7b644066c9820d247b19f1b9bf40840c\", \"URL\": \"/c/9bd1b3c9aeda7025068319c0a4af1d2b7b644066c9820d247b19f1b9bf40840c/%d\"}", c.version), mcaddr)
	}()

	output, err := multicastserviceInstance.GetCommit(c.Hkid())

	if err != nil {
		t.Errorf("Get Commit Failed \nError:%t", err)
	}

	if !bytes.Equal(output.Hash(), c.Hash()) || err != nil {
		t.Errorf("Make URL Failed \nExpected:%s \nGot: %s", c, output)
	}

}

func TestMulticastservice_GetTag(t *testing.T) {
	//t.Skipf("Come back to this test")
	//log.Printf("The key generated is, %d", KeyGen().D)
	//hkid := HKID{}
	hkid := hkidFromDString("6450698573071574057685373503239926609554390924514830851922442833127942726436428023022500281659846836919706975681006884631876585143520956760217923400876937896", 10)
	log.Printf("The HKID is, %s", hkid)
	b := blob([]byte("blob found"))
	log.Printf("Blob is, %s", b)
	tag_t := NewTag(b.Hash(), "blob", "BlobinTag", hkid)
	localfileserviceInstance.PostTag(tag_t)
	go func() {
		time.Sleep(1 * time.Millisecond)
		mcaddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:8000")
		multicastserviceInstance.receivemessage(fmt.Sprintf("{\"type\":\"tag\", \"hkid\": \"%s\", \"namesegment\": \"%s\", \"URL\": \"/t/%s/%s/%d\"}", hkid, tag_t.nameSegment, hkid, tag_t.nameSegment, tag_t.version), mcaddr)
	}()

	output, err := multicastserviceInstance.GetTag(tag_t.Hkid(), tag_t.nameSegment)

	if err != nil {
		t.Errorf("Get Tag Failed \nError:%s", err)
	} else if !bytes.Equal(output.Hash(), tag_t.Hash()) {
		t.Errorf("Make URL Failed \nExpected:%s \nGot: %s", tag_t, output)
	}

}
