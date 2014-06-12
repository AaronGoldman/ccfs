package crawler

import (
	"fmt"
	//"github.com/AaronGoldman/ccfs/services"
	"net/http"
)

func Start() {
	fmt.Printf("Crawler Starting")
	http.HandleFunc("/crawler/", WebCrawlerHandler)
}
func WebCrawlerHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "HTTP Error 500 Internal Crawler server error\n\n", 500)
}
