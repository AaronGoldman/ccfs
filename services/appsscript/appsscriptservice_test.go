package appsscript

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
	objects.RegisterGeterPoster(
		services.GetPublicKeyForHkid,
		services.GetPrivateKeyForHkid,
		localfile.Instance.PostKey,
		localfile.Instance.PostBlob,
	)
	services.Registerblobgeter(Instance)
	services.Registerblobgeter(timeout.Instance)
}

func TestAppsscriptservice_getBlob(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	//t.Skip()
	h, err := objects.HcidFromHex(
		"ca4c4244cee2bd8b8a35feddcd0ba36d775d68637b7f0b4d2558728d0752a2a2",
	)
	b, err := Instance.GetBlob(h)
	if err != nil || !bytes.Equal(b.Hash(), h) {
		log.Println(string(b))
		t.FailNow()
	}
}

func TestAppsscriptservice_getCommit(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	//t.Skip("Cause")
	h, err := objects.HkidFromHex(
		"c09b2765c6fd4b999d47c82f9cdf7f4b659bf7c29487cc0b357b8fc92ac8ad02",
	)
	c, err := Instance.GetCommit(h)
	if err != nil || c.Verify() == false {
		log.Println(err, "\n", c)
		t.Fail()
	}
}

func TestAppsscriptservice_getTag(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	//t.Skip()
	h, err := objects.HkidFromHex(
		"f65b92b9ce15e167b98fc896f0a365c87c39565642a59ba0060db3b33be6d885",
	)
	tt, err := Instance.GetTag(h, "testBlob")
	if err != nil || tt.Verify() == false {
		log.Println(err, "\n", tt)
		t.Fail()
	}
}

func TestAppsscriptservice_getKey(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	//t.Skip()
	h, err := objects.HkidFromHex(
		"f65b92b9ce15e167b98fc896f0a365c87c39565642a59ba0060db3b33be6d885",
	)
	k, err := Instance.GetKey(h)
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	if k.Hash().Hex() != "478025f8e8d4f769986232ca120be2b9c51a184455f6de1925a62a6f46df15b1" {
		log.Println(k.Hash().Hex())
		t.Fail()
	}
}
