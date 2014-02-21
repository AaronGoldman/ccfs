package main

import (
	"encoding/json"
	"fmt"
	"log"
	//"net"
)

//checks if I have the blob, it returns yes or no
func blobAvailable(hash HCID) bool {
	_, err := localfileserviceInstance.GetBlob(hash)
	if err != nil {
		return false
	}
	return true
}

////checks if I have the key, it returns yes or no
func keyAvailable(hash HKID) bool {
	_, err := localfileserviceInstance.GetKey(hash)
	if err != nil {
		return false
	}
	return true
}

//checks if I have the tag, it returns yes or no and the latest version
func tagAvailable(hash HKID, name string) (bool, int64) {
	t, err := localfileserviceInstance.GetTag(hash, name)
	if err != nil {
		return false, 0
	}
	return true, t.version
}

//checks if I have the commit, it returns yes or no and the latest version
func commitAvailable(hash HKID) (bool, int64) {
	c, err := localfileserviceInstance.GetCommit(hash)
	if err != nil {
		return false, 0
	}
	return true, c.version
}

func parseMessage(message string) (hkid HKID, hcid HCID, typeString string, nameSegment string, url string) {
	var Message map[string]interface{}

	err := json.Unmarshal([]byte(message), &Message)
	if err != nil {
		log.Printf("Error %s\n", err)
	}

	if Message["hcid"] != nil {
		hcid, err = HcidFromHex(Message["hcid"].(string))
	}
	if err != nil {
		log.Printf("Error with hex to string %s", err)
	}

	if Message["hkid"] != nil {
		hkid, err = HkidFromHex(Message["hkid"].(string))
	}
	if err != nil {
		log.Printf("Error with hex to string %s", err)
	}

	if Message["type"] != nil {
		typeString = Message["type"].(string)
	}

	if Message["nameSegment"] != nil {
		nameSegment = Message["nameSegment"].(string)
	}
	if Message["URL"] != nil {
		url = Message["URL"].(string)
	}
	return hkid, hcid, typeString, nameSegment, url

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
func buildResponse(hkid HKID, hcid HCID, typeString string, nameSegment string, version int64) (response string) {
	if typeString == "blob" {
		response = fmt.Sprintf("{\"type\": \"blob\", \"HCID\": \"%s\", \"URL\": \"%s\"}", hcid.Hex(),
			makeURL(hkid, hcid, typeString, nameSegment, version))
	} else if typeString == "commit" {
		response = fmt.Sprintf("{\"type\": \"commit\",\"HKID\": \"%s\", \"URL\": \"%s\"}", hkid.Hex(),
			makeURL(hkid, hcid, typeString, nameSegment, version))
	} else if typeString == "tag" {
		response = fmt.Sprintf("{\"type\": \"tag\", \"HKID\": \"%s\", \"namesegment\": \"%s\", \"URL\": \"%s\"}", hkid.Hex(), nameSegment,
			makeURL(hkid, hcid, typeString, nameSegment, version))
	} else if typeString == "key" {
		response = fmt.Sprintf("{\"type\": \"key\",\"HKID\": \"%s\", \"URL\": \"%s\"}", hkid.Hex(),
			makeURL(hkid, hcid, typeString, nameSegment, version))
	} else {
		return ""
	}
	return response

}

func makeURL(hkid HKID, hcid HCID, typeString string, nameSegment string, version int64) (response string) {
	//Path
	if typeString == "blob" {
		response = fmt.Sprintf("/b/%s" /*host,*/, hcid.Hex())
	} else if typeString == "commit" {
		response = fmt.Sprintf("/c/%s/%d" /*host,*/, hkid.Hex(), version)
	} else if typeString == "tag" {
		response = fmt.Sprintf("/t/%s/%s/%d" /*host,*/, hkid.Hex(), nameSegment, version)
	} else if typeString == "key" {
		response = fmt.Sprintf("/k/%s" /*host,*/, hkid.Hex())
	} else {
		response = ""
	}
	return response
}
func checkAndRespond(hkid HKID, hcid HCID, typeString string, nameSegment string) {
	response := ""
	if typeString == "blob" && blobAvailable(hcid) {
		response = buildResponse(hkid, hcid, typeString, nameSegment, 0)
	} else if typeString == "commit" {
		isAvailable, version := commitAvailable(hkid)
		if isAvailable {
			response = buildResponse(hkid, hcid, typeString, nameSegment, version)
		}
	} else if typeString == "tag" {
		isAvailable, version := tagAvailable(hkid, nameSegment)
		if isAvailable {
			response = buildResponse(hkid, hcid, typeString, nameSegment, version)
		}
	} else if typeString == "key" && keyAvailable(hkid) {
		response = buildResponse(hkid, hcid, typeString, nameSegment, 0)
	} else {
		return
	}
	err := multicastserviceInstance.sendmessage(response)
	if err != nil {
		log.Printf("check and responde failed to send message %s", err)
	}
	return
}
