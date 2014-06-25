package crawler

import (
	"fmt"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
	"log"
	"net/http"
	"sort"
	"strings"
)

var queuedTargets map[string]bool
var targetQueue chan target

// This function starts up the crawler for the CCFS
func Start() {
	fmt.Printf("Crawler Starting\n")
	queuedTargets = make(map[string]bool)
	targetQueue = make(chan target, 100)
	go processQueue()
	http.HandleFunc("/crawler/", webCrawlerHandler)

}

// This function handles web requests for the crawler
func webCrawlerHandler(w http.ResponseWriter, r *http.Request) {

	parts := strings.SplitN(r.RequestURI[9:], "/", 2)
	hkidhex := parts[0]
	requestStats := fmt.Sprintf("Request Statistics: \n")
	w.Write([]byte(requestStats))
	if len(hkidhex) == 64 {
		h, err := objects.HkidFromHex(hkidhex)
		if err == nil {
			response := fmt.Sprintf("\tThe hkid gotten is %s\n", h)
			w.Write([]byte(response)) //converts the response to bytes from strings
			seedQueue(h)
		} else {
			hc, err := objects.HcidFromHex(hkidhex)
			if err == nil {
				response := fmt.Sprintf("\tThe hcid gotten is %s\n", hc)
				w.Write([]byte(response))
				seedQueue(hc)
			}
		}

	}
	queueStats := fmt.Sprintf("Queue Statistics: \n")
	w.Write([]byte(queueStats))
	queuePrint := fmt.Sprintf("\tThe current length of the queue is %v\n",
		len(targetQueue))
	w.Write([]byte(queuePrint))
	indexStats := fmt.Sprintf("Index Statistics: \n")
	w.Write([]byte(indexStats))

	var keys []string
	for key := range queuedTargets {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, k := range keys {
		mapLine := fmt.Sprintf("\t%s %v\n", k, queuedTargets[k])
		w.Write([]byte(mapLine))
	}
	//http.Error(w, "HTTP Error 500 Internal Crawler server error\n\n", 500)
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
func seedQueue(h objects.HID) {
	crawlList(objects.HCID(h.Bytes()))
	crawlBlob(objects.HCID(h.Bytes()))
	crawlCommit(objects.HKID(h.Bytes()))
	crawlhcidCommit(objects.HCID(h.Bytes()))
	crawlhcidTag(objects.HCID(h.Bytes()))
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

	//firstKey, keyErr := services.GetKey(h.Bytes())
	//if keyErr != nil {
	//	// .....
	//	_ = firstKey
	//}
}

func crawlBlob(targHash objects.HCID) (err error) {
	nextBlob, blobErr := services.GetBlob(targHash)
	if blobErr != nil {
		return blobErr
	}
	// ......
	_ = nextBlob
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
	newParent := target{
		typeString: "hcid_commit",
		hash:       inCommit.Parent,
	}
	//commitHKID := firstCommit.Hkid()

	if !queuedTargets[newHash.String()] {
		targetQueue <- newHash
		queuedTargets[newHash.String()] = true
	}
	if !queuedTargets[newParent.String()] {
		targetQueue <- newParent
		queuedTargets[newParent.String()] = true
	}

}

func crawlTag(targHash objects.HKID) (err error) {
	// ......
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
	newParent := target{
		typeString: "hcid_tag",
		hash:       inTag.Parent,
	}
	//tagHKID := inTag.Hkid

	if !queuedTargets[newHash.String()] {
		targetQueue <- newHash
		queuedTargets[newHash.String()] = true
	}
	if !queuedTargets[newParent.String()] {
		targetQueue <- newParent
		queuedTargets[newParent.String()] = true
	}

}

func processQueue() {
	var targ target
	for {
		targ = <-targetQueue
		crawlTarget(targ)
	}
}
