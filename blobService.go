package main

import (
	"fmt"
	"net/http"
)

type BlobServer struct{}

func (h BlobServer) ServeHTTP(
		w http.ResponseWriter,
		r *http.Request) {
	fmt.Fprint(w, "Hello %s", )
	fmt.Println(r)
	
}

func main() {
	//var h BlobServer
	http.Handle("/blob", http.HandlerFunc(blobHandler))
	http.Handle("/tag", http.HandlerFunc(tagHandler))
	http.Handle("/commit", http.HandlerFunc(commitHandler))
	http.ListenAndServe("localhost:8080",nil)
}

func blobHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "blob %s", r.URL)
}

func tagHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "tag %s", r.URL)
}

func commitHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "commit %s" ,r.URL)
}
