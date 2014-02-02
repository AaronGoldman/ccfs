package main

import (
	"encoding/json"
	"log"
)

//checks if I have the blob, it returns yes or no
func blobAvailable(hash HCID) bool {
	localfileserviceInstance.GetBlob(hash)
	return false
}

////checks if I have the key, it returns yes or no
func keyAvailable(hash HKID) bool {
	return false
}

//checks if I have the tag, it returns yes or no and the latest version
func tagAvailable(hash HKID, name string) (bool, int64) {
	return false, 0
}

//checks if I have the commit, it returns yes or no and the latest version
func commitAvailable(hash HKID) (bool, int64) {
	return false, 0
}
func parseMessage(message string) (HKID, HCID, string, string) {
	var Message map[string]interface{}

	err := json.Unmarshal([]byte(message), &Message)
	if err != nil {
		log.Printf("Error %s\n", err)
	}
	hcid, err := HcidFromHex(Message["hcid"].(string))
	if err != nil {
		log.Printf("Error with hex to string %s", err)
	}
	hkid, err := HkidFromHex(Message["hkid"].(string))
	if err != nil {
		log.Printf("Error with hex to string %s", err)
	}
	typeString := Message["type"].(string)
	nameSegment := Message["namesegment"].(string)
	return hkid, hcid, typeString, nameSegment
}

func responseAvaiable(hkid HKID, hcid HCID, typeString string, nameSegment string) (available bool, version int64) {

	if typeString == "blob" {
		if hcid == nil {
			log.Printf("Malformed json")
			return
		}
		available = blobAvailable(hcid)
		version = 0
		return

		//Might wanna validate laterrrr
	} else if typeString == "commit" {
		if hkid == nil {
			log.Printf("Malformed json")
			return
		}
		available, version = commitAvailable(hkid)
		return
		//localfileserviceInstance.getCommit(h)
	} else if typeString == "tag" {
		if hkid == nil || nameSegment == "" {
			log.Printf("Malformed json")
			return
		}
		available, version = tagAvailable(hkid, nameSegment)
		return
		//localfileserviceInstance.getTag(h, nameSegment.(string))
	} else if typeString == "key" {
		if hkid == nil {
			log.Printf("Malformed json")
			return
		}
		available = keyAvailable(hkid)
		version = 0
		return
		//localfileserviceInstance.getKey(h)
	} else {
		log.Printf("Malformed json")
		return
	}
}
