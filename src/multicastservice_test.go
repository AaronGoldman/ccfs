// multicastservice_test
package main

import (
	"testing"
	"time"
)

func TestMulticastservice_GetBlob(t *testing.T) {
	t.Skipf("Come back to this test")
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
			multicastserviceInstance.receivemessage("{\"type\":\"blob\", \"hcid\": \"42cc3a4c4a9d9d3ee7de9322b45acb0e5a5c33550d9ad4791df6ae937a869e12\", \"URL\": \"http://localhost:8080/b/42cc3a4c4a9d9d3ee7de9322b45acb0e5a5c33550d9ad4791df6ae937a869e12\"}", multicastserviceInstance.mcaddr)
		}()

		output, err := multicastserviceInstance.GetBlob(answer.hcid)
		if err != nil {
			t.Errorf("Get Blob Failed \nError:%t", err)
		}
		if output.Hash().Hex() != answer.hcid.Hex() {
			t.Errorf("GetBlob failed \nExpected:%t \nGot: %t", answer.response, output)
		}

	}
}
