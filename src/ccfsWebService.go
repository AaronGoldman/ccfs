package main

import (
	"fmt"
	"log"
	"net/http"
)

type BlobServer struct{}

func (h BlobServer) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request) {
	fmt.Fprint(w, "Hello")
	log.Println(r)
}

func BlobServerStart() {
	http.Handle("/blob/", http.StripPrefix("/blob/",
		http.FileServer(http.Dir("../blobs"))))
	http.Handle("/tag/", http.StripPrefix("/tag/",
		http.FileServer(http.Dir("../tag"))))
	http.Handle("/commit/", http.StripPrefix("/commit/",
		http.FileServer(http.Dir("../commit"))))
	http.ListenAndServe(":8080", nil)
}
