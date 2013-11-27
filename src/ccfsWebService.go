package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type BlobServer struct{}

func (h BlobServer) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request) {
	fmt.Fprint(w, "Hello")
	log.Println(r)
}

func BlobServerStart() {
	http.Handle("/b/", http.StripPrefix("/blob/",
		http.FileServer(http.Dir("../blobs"))))
	http.Handle("/t/", http.StripPrefix("/tag/",
		http.FileServer(http.Dir("../tag"))))
	http.Handle("/c/", http.StripPrefix("/commit/",
		http.FileServer(http.Dir("../commit"))))
	http.ListenAndServe(":8080", nil)
}

func RepoServerStart() {
	http.HandleFunc("/r/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.SplitN(r.RequestURI[3:], "/", 2)
		hkidhex := parts[0]
		path := ""
		if len(parts) > 1 {
			path = parts[1]
		}
		err := error(nil)
		if len(hkidhex) == 64 {
			h, err := HkidFromHex(hkidhex)
			if err == nil {
				b, err := Get(h, path)
				//log.Printf("\n\t%s\n", path)
				if err == nil {
					w.Write(b.Bytes())
					return
				} else {
					http.Error(w, fmt.Sprint(
						"HTTP Error 500 Internal server error\n\n", err), 500)
				}
			}
		}
		w.Write([]byte(fmt.Sprintf("Invalid HKID\nerr: %v", err)))
		return
	})
}

func composeQuery(typestring string, hash HKID, namesegment string) (message string) {
	message = fmt.Sprintf("%s,%s", typestring, hash.String())
	if namesegment != "" {
		message = fmt.Sprintf("%s,%s", message, namesegment)
	}
	return message

}

func parseQuery(message string) (typestring string, hash HKID, namesegment string) {
	arr_message := strings.SplitN(message, ",", 3)
	if len(arr_message) < 2 {
		panic("error: malformed parseQuery")
	}
	typestring = arr_message[0]

	hash, err := HkidFromHex(arr_message[1])
	if err != nil {
		panic("error: malformed hexadecimal")
	}
	if len(arr_message) > 2 {
		namesegment = arr_message[2]
	} else {
		namesegment = ""
	}
	fmt.Println("Error")
	return typestring, hash, namesegment
}
