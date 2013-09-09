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
				if err == nil {
					w.Write(b.Bytes())
					return
				}
			}
		}
		w.Write([]byte(fmt.Sprintf("Invalid HKID\nerr: %v", err)))
		return
	})
}
