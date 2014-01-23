package main

import (
	"encoding/json"
	"log"
)

//checks if I have the blob, it returns yes or no
func blobAvailable(hash HCID) bool {
	return false
}

////checks if I have the key, it returns yes or no
func keyAvailable(hash HKID) bool {
	return false
}

//checks if I have the tag, it returns yes or no and the latest version
func tagAvailable(hash HKID, name string) (bool, int) {
	return false, 0
}

//checks if I have the commit, it returns yes or no and the latest version
func commitAvailable(hash HKID) (bool, int) {
	return false, 0
}
func parseMessage(message string) {
	var Message map[string]interface{}
	//dec :=json.NewDecoder(strings.NewReader(message))
	err := json.Unmarshal([]byte(message), &Message)
	if err != nil {
		log.Printf("Error %s\n", err)
	}
	log.Println(Message["type"], Message["hkid"], Message["hcid"], Message["namesegment"])
	log.Printf("Derp %s", Message)

	if Message["type"] == "blob" {
		if Message["hcid"] == nil {
			log.Printf("Malformed json")
			return
		}
		h, err := HcidFromHex(Message["hcid"].(string))
		if err != nil {
			log.Printf("Error with hex to string %s", err)
		}
		localfileserviceInstance.getBlob(h)
		//Might wanna validate date laterrrr
	} else if Message["type"] == "commit" {
		if Message["hkid"] == nil {
			log.Printf("Malformed json")
			return
		}
		h, err := HkidFromHex(Message["hkid"].(string))
		if err != nil {
			log.Printf("Error with hex to string %s", err)
		}
		localfileserviceInstance.getCommit(h)
	} else if Message["type"] == "tag" {
		if Message["hkid"] == nil || Message["namesegment"] == nil {
			log.Printf("Malformed json")
			return
		}
		h, err := HkidFromHex(Message["hkid"].(string))
		if err != nil {
			log.Printf("Error with hex to string %s", err)
		}
		localfileserviceInstance.getTag(h, Message["namesegment"].(string))
	} else if Message["type"] == "key" {
		if Message["hkid"] == nil {
			log.Printf("Malformed json")
			return
		}
		h, err := HkidFromHex(Message["hkid"].(string))
		if err != nil {
			log.Printf("Error with hex to string %s", err)
		}
		localfileserviceInstance.getKey(h)
	} else {
		log.Printf("Malformed json")
		return
	}
}
