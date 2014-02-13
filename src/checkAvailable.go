package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
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

	hcid := HCID{}
	if Message["hcid"] != nil {
		hcid, err = HcidFromHex(Message["hcid"].(string))
	}
	if err != nil {
		log.Printf("Error with hex to string %s", err)
	}

	hkid := HKID{}
	if Message["hkid"] != nil {
		hkid, err = HkidFromHex(Message["hkid"].(string))
	}
	if err != nil {
		log.Printf("Error with hex to string %s", err)
	}
	typeString := ""
	if Message["type"] != nil {
		typeString = Message["type"].(string)
	}
	nameSegment := ""
	if Message["nameSegment"] != nil {
		nameSegment = Message["nameSegment"].(string)
	}
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
func getHostName() string {
	//ToDo
	return "localhost:8080"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Something meaningful... %s\n", err)
		return "localhost:8080"
	}
	for _, addr := range addrs {
		log.Printf("Network:%s  \nString:%s\n", addr.Network(), addr.String())
	}
	return "LAME"

}
func makeURL(hkid HKID, hcid HCID, typeString string, nameSegment string, version int64) (response string) {
	//Host Name
	host := getHostName()
	//Path
	if typeString == "blob" {
		response = fmt.Sprintf("%s/b/%s", host, hcid.Hex())
	} else if typeString == "commit" {
		response = fmt.Sprintf("%s/c/%s/%d", host, hkid.Hex(), version)
	} else if typeString == "tag" {
		response = fmt.Sprintf("%s/t/%s/%s/%d", host, hkid.Hex(), nameSegment, version)
	} else if typeString == "key" {
		response = fmt.Sprintf("%s/k/%s", host, hkid.Hex())
	} else {
		response = ""
	}
	return response
}
