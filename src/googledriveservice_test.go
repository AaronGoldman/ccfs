// googledriveservice_test.go
package main

import (
	"log"
	"testing"
)

func TestGoogledriveservice_GetBlob(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	t.Skip("skipping google drive tests")
	//googledriveservice_setup()
	hb, err := HcidFromHex("42cc3a4c4a9d9d3ee7de9322b45acb0e5a5c33550d9ad4791df6ae937a869e12")
	if err != nil {
		t.Fail()
	}
	b, err := googledriveserviceInstance.getBlob(hb)
	if err != nil || b.Hash().String() != hb.Hex() {
		t.Fail()
	}
}

func TestGoogledriveservice_GetCommit(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	t.Skip("skipping google drive tests")
	//googledriveservice_setup()
	hc, err := HkidFromHex("c09b2765c6fd4b999d47c82f9cdf7f4b659bf7c29487cc0b357b8fc92ac8ad02")
	c, err := googledriveserviceInstance.getCommit(hc)
	log.Println(err, "\n", c)
	if err != nil || c.Verifiy() == false {
		t.Fail()
	}
	//log.Printf("\n\tCommit Contents: %v \n\tError: %v", c.Verifiy(), err)
}

func TestGoogledriveservice_GetTag(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	t.Skip("skipping google drive tests")
	//googledriveservice_setup()
	ht, err := HkidFromHex("f65b92b9ce15e167b98fc896f0a365c87c39565642a59ba0060db3b33be6d885")
	tt, err := googledriveserviceInstance.getTag(ht, "testBlob")
	log.Println(err, "\n", tt)
	if err != nil || tt.Verifiy() == false {
		t.Fail()
	}
	//log.Printf("\n\tTag Contents: %v \n\tError: %v", tt.Verifiy(), err)
}

func TestGoogledriveservice_GetKey(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	t.Skip("skipping google drive tests")
	//googledriveservice_setup()
	hk, err := HkidFromHex("f65b92b9ce15e167b98fc896f0a365c87c39565642a59ba0060db3b33be6d885")
	k, err := googledriveserviceInstance.getKey(hk)
	if err != nil {
		t.Fail()
	}
	if k.Hash().Hex() != "478025f8e8d4f769986232ca120be2b9c51a184455f6de1925a62a6f46df15b1" {
		t.Fail()
	}
	//log.Printf("\n\tKey Contents: %s \n\tError: %v", k.Hash(), err)
}
