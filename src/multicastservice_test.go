// multicastservice_test
package main

import (
	"bytes"
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
			t.Errorf("Make URL Failled \nExpected:%s \nGot: %s", answer.response, output)
		}

	}
}
