// multicastservice_test
package multicast

import (
	"bytes"
	"fmt"
	"github.com/AaronGoldman/ccfs/interfaces/web"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
	"github.com/AaronGoldman/ccfs/services/localfile"
	"github.com/AaronGoldman/ccfs/services/timeout"
	"log"
	"net"
	"testing"
	"time"
)

var benchmarkRepo objects.HKID

func init() {
	benchmarkRepo = objects.HkidFromDString("44089867384569081480066871308647"+
		"4832666868594293316444099156169623352946493325312681245061254048"+
		"6538169821270508889792789331438131875225590398664679212538621", 10)
	services.Registerblobgeter(Instance)
	services.Registerblobgeter(timeout.Instance)
}

func TestMulticastservice_GetBlob(t *testing.T) {
	//t.Skipf("Come back to this test")
	web.Start()
	//go BlobServerStart()

	AnswerKey := []struct {
		hcid     objects.HCID
		response objects.Blob
	}{
		{objects.Blob([]byte("blob found")).Hash(), objects.Blob([]byte("blob found"))},
	}
	for _, answer := range AnswerKey {

		localfile.Instance.PostBlob(answer.response)
		go func() {
			time.Sleep(1 * time.Millisecond)
			mcaddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:1234")
			Instance.receivemessage("{\"type\":\"blob\", \"hcid\": \"42cc3a4c4a9d9d3ee7de9322b45acb0e5a5c33550d9ad4791df6ae937a869e12\", \"URL\": \"/b/42cc3a4c4a9d9d3ee7de9322b45acb0e5a5c33550d9ad4791df6ae937a869e12\"}", mcaddr)
		}()

		output, err := Instance.GetBlob(answer.hcid)

		if err != nil {
			t.Errorf("Get Blob Failed \nError:%s", err)
		} else if !bytes.Equal(output.Hash(), answer.hcid) {
			t.Errorf("Get Blob Failed \nExpected:%s \nGot: %s", answer.response, output)
		}

	}
}

func TestMulticastservice_GetCommit(t *testing.T) {
	//t.Skipf("Come back to this test")
	hkid := objects.HkidFromDString("5198719439877464148627795433286736285873678110640040333794349799294848737858561643942881983506066042818105864129178593001327423646717446545633525002218361750", 10)

	b := objects.Blob([]byte("blob found"))
	l := objects.NewList(b.Hash(), "blob", "Blobinlist")
	c := objects.NewCommit(l.Hash(), hkid)
	localfile.Instance.PostCommit(c)
	go func() {
		time.Sleep(1 * time.Millisecond)
		mcaddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:8000")
		Instance.receivemessage(fmt.Sprintf("{\"type\":\"commit\", \"hkid\": \"9bd1b3c9aeda7025068319c0a4af1d2b7b644066c9820d247b19f1b9bf40840c\", \"URL\": \"/c/9bd1b3c9aeda7025068319c0a4af1d2b7b644066c9820d247b19f1b9bf40840c/%d\"}", c.Version), mcaddr)
	}()

	output, err := Instance.GetCommit(c.Hkid())

	if err != nil {
		t.Errorf("Get Commit Failed \nError:%s", err)
	} else if !output.Verify() {
		//else if !bytes.Equal(output.Hash(), c.Hash()) {
		t.Errorf("Get Commit Failed \nExpected:%s \nGot: %s", c, output)
	}
	if output.Version() != c.Version() {
		log.Printf("Commit is stale %d", c.Version()-output.Version())
	}
}

func TestMulticastservice_GetTag(t *testing.T) {
	//t.Skipf("Come back to this test")
	//log.Printf("The key generated is, %d", KeyGen().D)
	//hkid := HKID{}
	hkid := objects.HkidFromDString("6450698573071574057685373503239926609554390924514830851922442833127942726436428023022500281659846836919706975681006884631876585143520956760217923400876937896", 10)
	b := objects.Blob([]byte("blob found"))
	tag_t := objects.NewTag(b.Hash(), "blob", "BlobinTag", hkid)
	localfile.Instance.PostTag(tag_t)
	go func() {
		time.Sleep(1 * time.Millisecond)
		mcaddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:8000")
		Instance.receivemessage(fmt.Sprintf("{\"type\":\"tag\", \"hkid\": \"%s\", \"namesegment\": \"%s\", \"URL\": \"/t/%s/%s/%d\"}", hkid, tag_t.NameSegment, hkid, tag_t.NameSegment, tag_t.Version), mcaddr)
	}()

	output, err := Instance.GetTag(tag_t.Hkid(), tag_t.NameSegment)

	if err != nil {
		t.Errorf("Get Tag Failed \nError:%s", err)
	} else if !output.Verify() {
		//!bytes.Equal(output.Hash(), tag_t.Hash()) {
		t.Errorf("Get Tag Failed \nExpected:%s \nGot: %s", tag_t, output)
	}
	if output.Version != tag_t.Version {
		log.Printf("Tag is stale %d", tag_t.Version-output.Version)
	}
}

//func TestMulticastservice_GetKey(t *testing.T) {
//	//t.Skipf("Come back to this test")
//	//log.Printf("The key generated is, %d", KeyGen().D)
//	//hkid := HKID{}
//	hkid := hkidFromDString("1404824859588041073678522358128083070212026722272843802627814234500338310965596733375766723741689655938709528394105700164973947017056140017071321292025240790", 10)
//	log.Printf("The HKID is, %s", hkid)

//	privkey, err := GetKey(hkid)
//	if err != nil {
//		log.Printf("Error for GetKey, %s", err)
//	}

//	go func() {
//		time.Sleep(1 * time.Millisecond)
//		mcaddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:8000")
//		multicastserviceInstance.receivemessage(fmt.Sprintf("{\"type\":\"key\", \"hkid\": \"%s\", \"URL\": \"/k/%s\"}", hkid, hkid), mcaddr)
//	}()

//	output, err := multicastserviceInstance.GetKey(hkid)

//	if err != nil {
//		t.Errorf("Get Key Failed \nError:%t", err)
//	}

//	if !bytes.Equal(output.Hash(), privkey.Hash()) || err != nil {
//		t.Errorf("Make URL Failed \nExpected:%s \nGot: %s", privkey, output)
//	}

//}
