package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

//BlobServerStart starts a server for the content services
func BlobServerStart() {
	http.Handle("/b/", http.StripPrefix("/b/",
		http.FileServer(http.Dir("bin/blobs"))))
	http.Handle("/t/", http.StripPrefix("/t/",
		http.FileServer(http.Dir("bin/tags"))))
	http.Handle("/c/", http.StripPrefix("/c/",
		http.FileServer(http.Dir("bin/commits"))))
	http.ListenAndServe(":8080", nil)
}

//CollectionServerStart starts a server for full CCFS queries HKID/path
func CollectionServerStart() {
	http.HandleFunc("/r/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.SplitN(r.RequestURI[3:], "/", 2)
		hkidhex := parts[0]
		path := ""
		if len(parts) > 1 {
			path = parts[1]
		}
		err := error(nil)
		if len(hkidhex) == 64 {
			h, err := objects.HkidFromHex(hkidhex)
			if err == nil {
				b, err := services.Get(h, path)
				//log.Printf("\n\t%s\n", path)
				if err == nil {
					w.Write(b.Bytes())
					return
				}
				http.Error(w, fmt.Sprint(
					"HTTP Error 500 Internal server error\n\n", err), 500)
			}
		}
		w.Write([]byte(fmt.Sprintf("Invalid HKID\nerr: %v", err)))
		return
	})
	http.HandleFunc("/d/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.SplitN(r.RequestURI[3:], "/", 2)
		hkidhex := parts[0]
		path := ""
		if len(parts) > 1 {
			path = parts[1]
		}
		err := error(nil)
		if len(hkidhex) == 64 {
			h, err := objects.HkidFromHex(hkidhex)
			if err == nil {
				b, err := services.GetD(h, path)
				//log.Printf("\n\t%s\n", path)
				if err == nil {
					w.Write(b.Bytes())
					return
				}
				http.Error(w, fmt.Sprint(
					"HTTP Error 500 Internal server error\n\n", err), 500)
			}
		}
		w.Write([]byte(fmt.Sprintf("Invalid HKID\nerr: %v", err)))
		return
	})
}

func composeQuery(typestring string, hash objects.HKID, namesegment string) (message string) {
	message = fmt.Sprintf("%s,%s", typestring, hash.String())
	if namesegment != "" {
		message = fmt.Sprintf("%s,%s", message, namesegment)
	}
	return message

}

func parseQuery(message string) (typestring string, hash objects.HKID, namesegment string) {
	arrMessage := strings.SplitN(message, ",", 3)
	if len(arrMessage) < 2 {
		panic("error: malformed parseQuery")
	}
	typestring = arrMessage[0]

	hash, err := objects.HkidFromHex(arrMessage[1])
	if err != nil {
		panic("error: malformed hexadecimal")
	}
	if len(arrMessage) > 2 {
		namesegment = arrMessage[2]
	} else {
		namesegment = ""
	}
	fmt.Println("Error")
	return typestring, hash, namesegment
}

//Start begins the listening for web reqwests on 8080
func Start() {
	go BlobServerStart()
	go CollectionServerStart()
}
