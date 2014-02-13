package main

import (
	"testing"
)

func TestMakeURL(t *testing.T) {
	AnswerKey := []struct {
		hkid        HKID
		hcid        HCID
		typeString  string
		nameSegment string
		version     int64
		response    string
	}{
		{benchmarkRepo, nil, "commit", "", 100, "localhost:8080/c/549baa6497db3615332aae859680b511117e299879ee311fbac4d1a40f93b8d0/100"},
		{benchmarkRepo, nil, "tag", "Yolo", 100, "localhost:8080/t/549baa6497db3615332aae859680b511117e299879ee311fbac4d1a40f93b8d0/Yolo/100"},
		{nil, blob{}.Hash(), "blob", "", 0, "localhost:8080/b/e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{benchmarkRepo, nil, "key", "", 0, "localhost:8080/k/549baa6497db3615332aae859680b511117e299879ee311fbac4d1a40f93b8d0"},
	}

	for _, answer := range AnswerKey {
		output := makeURL(answer.hkid, answer.hcid, answer.typeString, answer.nameSegment, answer.version)
		if output != answer.response {
			t.Errorf("Make URL Failled \nExpected:%s \nGot: %s", answer.response, output)
		}
	}
}
func TestBuildResponse(t *testing.T) {
	AnswerKey := []struct {
		hkid        HKID
		hcid        HCID
		typeString  string
		nameSegment string
		version     int64
		response    string
	}{
		{benchmarkRepo, nil, "commit", "", 100, "{\"type\": \"commit\",\"HKID\": \"549baa6497db3615332aae859680b511117e299879ee311fbac4d1a40f93b8d0\", \"URL\": \"localhost:8080/c/549baa6497db3615332aae859680b511117e299879ee311fbac4d1a40f93b8d0/100\"}"},
		{benchmarkRepo, nil, "tag", "Yolo", 100, "{\"type\": \"tag\", \"HKID\": \"549baa6497db3615332aae859680b511117e299879ee311fbac4d1a40f93b8d0\", \"namesegment\": \"Yolo\", \"URL\": \"localhost:8080/t/549baa6497db3615332aae859680b511117e299879ee311fbac4d1a40f93b8d0/Yolo/100\"}"},
		{nil, blob{}.Hash(), "blob", "", 0, "{\"type\": \"blob\", \"HCID\": \"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\", \"URL\": \"localhost:8080/b/e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\"}"},
		{benchmarkRepo, nil, "key", "", 0, "{\"type\": \"key\",\"HKID\": \"549baa6497db3615332aae859680b511117e299879ee311fbac4d1a40f93b8d0\", \"URL\": \"localhost:8080/k/549baa6497db3615332aae859680b511117e299879ee311fbac4d1a40f93b8d0\"}"},
	}

	for _, answer := range AnswerKey {
		output := buildResponse(answer.hkid, answer.hcid, answer.typeString, answer.nameSegment, answer.version)
		if output != answer.response {
			t.Errorf("Build Response Failled \nExpected:%s \nGot:        %s", answer.response, output)
		}
	}
}