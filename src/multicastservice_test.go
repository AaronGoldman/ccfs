// multicastservice_test
package main

import (
	"testing"
)

func TestMulticastservice_GetBlob(t *testing.T) {
	t.Skipf("Come back to this test")
	AnswerKey := []struct {
		hcid     HCID
		response blob
	}{
		{blob{}.Hash(), blob{}},
	}
	for _, answer := range AnswerKey {
		output, err := multicastserviceInstance.GetBlob(answer.hcid)
		multicastserviceInstance.receivemessage("{\"type\":\"blob\", \"hcid\": \"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\", \"URL\": \"localhost:8080/b/e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\"}", multicastserviceInstance.mcaddr)
		if output.Hash().Hex() != answer.hcid.Hex() || err != nil {
			t.Errorf("Make URL Failled \nExpected:%s \nGot: %s", answer.response, output)
		}
	}
}
