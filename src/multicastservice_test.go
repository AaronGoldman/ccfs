// multicastservice_test
package main

import (
	"bytes"
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

		multicastserviceInstance.receivemessage("{\"type\":\"blob\", \"hcid\": \"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\", \"URL\": \"localhost:8080/b/e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\"}", multicastserviceInstance.mcaddr)
		if !bytes.Equal(output.Hash(), answer.hcid) || err != nil {
			t.Errorf("Make URL Failled \nExpected:%s \nGot: %s", answer.response, output)
		}

	}
}
