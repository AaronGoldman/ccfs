//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

//Package crawler implements a cralwer interface
//It is started by calling crawler.Start
package crawler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/AaronGoldman/ccfs/interfaces"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

var queuedTargets map[string]bool
var targetQueue chan target

// Start is the function that starts up the crawler for the CCFS
func Start() {
	fmt.Printf("Crawler Starting\n")
	queuedTargets = make(map[string]bool)
	targetQueue = make(chan target, 100)
	go processQueue()
	http.HandleFunc("/crawler/", webCrawlerHandler)
	http.HandleFunc("/index/", webIndexHandler)
	http.HandleFunc("/search/", webSearchHandler)
	seedQueue(interfaces.GetLocalSeed())
}

// This function handles web requests for the crawler
func webCrawlerHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(r.RequestURI[9:], "/", 2)
	hkidhex := parts[0]
	hexerr := seedQueue(hkidhex)

	t, err := template.New("WebIndex template").Parse(`Request Statistics:
	The HID received is: {{.HkidHex}}{{if .Err}}
	Error Parsing HID: {{.Err}}{{end}}
Queue Statistics:
	The current length of the queue is: {{.QueueLength}}
Index Statistics:
	{{range $keys, $values := .Queue}}
		{{$keys}}{{end}}`)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("HTTP Error 500 Internal Crawler server error\n%s\n", err),
			500,
		)
	} else {
		t.Execute(w, struct {
			HkidHex     string
			QueueLength int
			Queue       map[string]bool
			Err         error
		}{
			HkidHex:     hkidhex,
			QueueLength: len(targetQueue),
			Queue:       queuedTargets,
			Err:         hexerr,
		}) //merge template ‘t’ with content of ‘index’
	}
}

// This type, target, is used for the map and the queue
type target struct {
	typeString string
	hash       objects.HID
}

func (targ target) String() string {
	return fmt.Sprintf("%s %s", targ.typeString, targ.hash)
}

// This function seeds the queue from a web request
func seedQueue(hkidhex string) (err error) {
	h, hexerr := objects.HcidFromHex(hkidhex)
	if hexerr == nil {
		err = crawlList(objects.HCID(h.Bytes()))
		if err != nil {
			log.Println(err)
		}
		err = crawlBlob(objects.HCID(h.Bytes()))
		if err != nil {
			log.Println(err)
		}
		err = crawlCommit(objects.HKID(h.Bytes()))
		if err != nil {
			log.Println(err)
		}
		err = crawlhcidCommit(objects.HCID(h.Bytes()))
		if err != nil {
			log.Println(err)
		}
		err = crawlhcidTag(objects.HCID(h.Bytes()))
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Println(hexerr)
	}
	return hexerr
}

// This function scrapes a target for new targets to add to the queue
func crawlTarget(targ target) {
	err := fmt.Errorf("Attempted to crawl malformed typeString in target\n\t%v", targ)
	switch targ.typeString {
	case "commit":
		err = crawlCommit(targ.hash.(objects.HKID))
	case "list":
		err = crawlList(targ.hash.(objects.HCID))
	case "blob":
		err = crawlBlob(targ.hash.(objects.HCID))
	case "tag":
		err = crawlTag(targ.hash.(objects.HKID))
	case "hcid_commit":
		err = crawlhcidCommit(targ.hash.(objects.HCID))
	case "hcid_tag":
		err = crawlhcidTag(targ.hash.(objects.HCID))
	default:
	}
	if err != nil {
		log.Print(err)
	}
}

func crawlBlob(targHash objects.HCID) (err error) {
	nextBlob, blobErr := services.GetBlob(targHash)
	if blobErr != nil {
		return blobErr
	}
	indexBlob(nextBlob)
	return nil
}

func crawlList(targHash objects.HCID) (err error) {
	firstList, listErr := services.GetList(targHash)
	if listErr != nil {
		return listErr
	}

	for _, entry := range firstList {
		newlistHash := target{
			typeString: entry.TypeString,
			hash:       entry.Hash,
		}
		if !queuedTargets[newlistHash.String()] {
			targetQueue <- newlistHash
			queuedTargets[newlistHash.String()] = true
		}
	}
	indexList(firstList)
	return nil
}

func crawlCommit(targHash objects.HKID) (err error) {
	firstCommit, commitErr := services.GetCommit(targHash)
	if commitErr != nil {
		return commitErr
	}
	handleCommit(firstCommit)
	return nil
}

func crawlhcidCommit(targHash objects.HCID) (err error) {
	hcidCommit, commErr := services.GetCommitForHcid(targHash)
	if commErr != nil {
		return commErr
	}

	handleCommit(hcidCommit)
	return nil
}

// This function handles commits from HCID or HKID
func handleCommit(inCommit objects.Commit) {
	newHash := target{
		typeString: "list",
		hash:       inCommit.ListHash,
	}
	if !queuedTargets[newHash.String()] {
		targetQueue <- newHash
		queuedTargets[newHash.String()] = true
	}
	for _, cparent := range inCommit.Parents {
		newParent := target{
			typeString: "hcid_commit",
			hash:       cparent,
		}
		if !queuedTargets[newParent.String()] {
			targetQueue <- newParent
			queuedTargets[newParent.String()] = true
		}
	}
	indexCommit(inCommit)
}

func crawlTag(targHash objects.HKID) (err error) {
	tags, tagErr := services.GetTags(targHash)
	if tagErr != nil {
		return tagErr
	}
	for _, tag := range tags {
		handleTag(tag)
	}
	return nil
}

func crawlhcidTag(targHash objects.HCID) (err error) {
	hcidTag, tagErr := services.GetTagForHcid(targHash)
	if tagErr != nil {
		return tagErr
	}
	handleTag(hcidTag)
	return nil
}

// This function handles Tags from HCID or HKID
func handleTag(inTag objects.Tag) {
	newHash := target{
		typeString: inTag.TypeString,
		hash:       inTag.HashBytes,
	}

	if !queuedTargets[newHash.String()] {
		targetQueue <- newHash
		queuedTargets[newHash.String()] = true
	}
	for _, tparent := range inTag.Parents {
		newParent := target{
			typeString: "hcid_tag",
			hash:       tparent,
		}
		if !queuedTargets[newParent.String()] {
			targetQueue <- newParent
			queuedTargets[newParent.String()] = true
		}
	}
	indexTag(inTag)
}

func processQueue() {
	var targ target
	for {
		targ = <-targetQueue
		crawlTarget(targ)
	}
}
