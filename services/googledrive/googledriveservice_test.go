//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file
// googledriveservice_test.go
package googledrive

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/AaronGoldman/ccfs/objects"
)

func TestGoogledriveservice_GetBlob(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	t.Skip("skipping google drive tests")
	//googledriveservice_setup()
	hb, err := objects.HcidFromHex("42cc3a4c4a9d9d3ee7de9322b45acb0e5a5c33550d9ad4791df6ae937a869e12")
	if err != nil {
		t.Fail()
	}
	b, err := Instance.GetBlob(hb)
	if err != nil || !bytes.Equal(b.Hash(), hb) {
		t.Fail()
	}
}

func TestGoogledriveservice_GetCommit(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	t.Skip("skipping google drive tests")
	//googledriveservice_setup()
	hc, err := objects.HkidFromHex("c09b2765c6fd4b999d47c82f9cdf7f4b659bf7c29487cc0b357b8fc92ac8ad02")
	c, err := Instance.GetCommit(hc)
	if err != nil || c.Verify() == false {
		log.Println(err, "\n", c)
		t.Fail()
	}
	//log.Printf("\n\tCommit Contents: %v \n\tError: %v", c.Verifiy(), err)
}

func TestGoogledriveservice_GetTag(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	t.Skip("skipping google drive tests")
	//googledriveservice_setup()
	ht, err := objects.HkidFromHex("f65b92b9ce15e167b98fc896f0a365c87c39565642a59ba0060db3b33be6d885")
	tt, err := Instance.GetTag(ht, "testBlob")
	if err != nil || tt.Verify() == false {
		log.Println(err, "\n", tt)
		t.Fail()
	}
	//log.Printf("\n\tTag Contents: %v \n\tError: %v", tt.Verifiy(), err)
}

func TestGoogledriveservice_GetKey(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	t.Skip("skipping google drive tests")
	//googledriveservice_setup()
	hk, err := objects.HkidFromHex("f65b92b9ce15e167b98fc896f0a365c87c39565642a59ba0060db3b33be6d885")
	k, err := Instance.GetKey(hk)
	if err != nil {
		t.Fail()
	}
	if k.Hash().Hex() != "478025f8e8d4f769986232ca120be2b9c51a184455f6de1925a62a6f46df15b1" {
		t.Fail()
	}
	//log.Printf("\n\tKey Contents: %s \n\tError: %v", k.Hash(), err)
}

func init() {
	os.Chdir("../../") //Changes the working directory to CcfsLink
	Start()
}
