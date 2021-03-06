//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

//Package web implements a web interface
//It is started by calling web.Start
package web

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

//BlobServerStart starts a server for the content services
func BlobServerStart() {
	http.Handle(
		"/b/",
		http.StripPrefix(
			"/b/",
			http.FileServer(http.Dir("bin/blobs")),
		),
	)
	http.Handle(
		"/t/",
		http.StripPrefix(
			"/t/",
			http.FileServer(http.Dir("bin/tags")),
		),
	)
	http.Handle(
		"/c/",
		http.StripPrefix("/c/",
			http.FileServer(http.Dir("bin/commits")),
		),
	)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Listening on :8080")
}

//CollectionServerStart starts a server for full CCFS queries HKID/path
func CollectionServerStart() {
	http.HandleFunc("/r/", webRepositoryHandler)
	http.HandleFunc("/d/", webDomainHandler)
}

func webRepositoryHandler(w http.ResponseWriter, r *http.Request) {
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
}
func webDomainHandler(w http.ResponseWriter, r *http.Request) {
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

var uploadHtml []byte

func UploadServerStart() {
	var fileReadErr error
	uploadHtml, fileReadErr = ioutil.ReadFile("./interfaces/web/upload.htm")
	if fileReadErr != nil {
		log.Printf("read the file %q with error %s\n", uploadHtml, fileReadErr)
	}
	http.HandleFunc("/upload/", UploadDomainHandler)
}

func UploadDomainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	//io.WriteString(w, fmt.Sprintf("%v\n",r))
	//io.WriteString(w, fmt.Sprintf("length of uploadHtml is %d\n",len(uploadHtml)))
	w.Write(uploadHtml)
	//http.Error(w, fmt.Sprint(
	//	"HTTP Error 500 Internal server error\n\n", err), 500)
	return
}

//Start begins the listening for web reqwests on 8080
func Start() {
	go BlobServerStart()
	go CollectionServerStart()
	go UploadServerStart()
}
