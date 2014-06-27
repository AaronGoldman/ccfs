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
	//indexBlob(nextBlob)
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
	//indexList(firstList)
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
		//commitHKID := firstCommit.Hkid()

		if !queuedTargets[newParent.String()] {
			targetQueue <- newParent
			queuedTargets[newParent.String()] = true
		}
	}
	//indexCommit(inCommit)

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

	if !queuedTargets[newHash.String()] {
		targetQueue <- newHash
		queuedTargets[newHash.String()] = true
	}
	for _, tparent := range inTag.Parents {
		newParent := target{
			typeString: "hcid_tag",
			hash:       tparent,
		}
		//tagHKID := inTag.Hkid

		if !queuedTargets[newParent.String()] {
			targetQueue <- newParent
			queuedTargets[newParent.String()] = true
		}
	}
	//indexTag(inTag)
}

func processQueue() {
	var targ target
	for {
		targ = <-targetQueue
		crawlTarget(targ)
	}
}

var blobIndex struct {
	size    map[string]int
	nameSeg map[string]map[string]string
	//map[blobHcid string]map[nameSeg string]referringHID string
	descendants map[string]map[int64]objects.HCID
	//map[blobHcid string]map[versionNumber int64]referringHCID string
}

var commitIndex struct {
	nameSeg map[string]map[string]string
	version map[string]map[int64]objects.HCID
	//map[hkid string]map[versionNumber int64]HCIDversion string
}

var tagIndex struct {
	nameSeg map[string]map[string]string
	version map[string]map[string]map[int64]objects.HCID
	//version map[HKID string]map[nameSeg string]map[versionNumber int64]objects.HCID
}

func indexBlob(inBlob objects.Blob) {
	blobIndex.size[inBlob.Hash().Hex()] = len(inBlob)
}

func indexList(inList objects.List) {
	inListHex := inList.Hash().Hex()
	blobIndex.size[inListHex] = len(inList)
	for nameSeg, entry := range inList {
		if entry.TypeString == "list" || entry.TypeString == "blob" {
			blobIndex.nameSeg[entry.Hash.Hex()][nameSeg] = inListHex
		} else if entry.TypeString == "tag" {
			tagIndex.nameSeg[entry.Hash.Hex()][nameSeg] = inListHex
		} else if entry.TypeString == "commit" {
			commitIndex.nameSeg[entry.Hash.Hex()][nameSeg] = inListHex
		} else {
			log.Printf("Received invalid typestring: %s\n", entry.TypeString)
		}
	}
}

func indexCommit(inCommit objects.Commit) {
	commitIndex.version[inCommit.Hkid.Hex()][inCommit.Version] = inCommit.Hash()
	for _, inComParent := range inCommit.Parents {
		blobIndex.descendants[inComParent.Hex()][inCommit.Version] = inCommit.Hash()
	}
}

func indexTag(inTag objects.Tag) {
	inTagHex := inTag.Hash().Hex()
	if inTag.TypeString == "list" || inTag.TypeString == "blob" {
		blobIndex.nameSeg[inTag.HashBytes.Hex()][inTag.NameSegment] = inTagHex
	} else if inTag.TypeString == "tag" {
		tagIndex.nameSeg[inTag.HashBytes.Hex()][inTag.NameSegment] = inTagHex
	} else if inTag.TypeString == "commit" {
		commitIndex.nameSeg[inTag.HashBytes.Hex()][inTag.NameSegment] = inTagHex
	} else {
		log.Printf("Received invalid typestring: %s\n", inTag.TypeString)
	}

	tagIndex.version[inTag.Hkid.Hex()][inTag.NameSegment][inTag.Version] = inTag.Hash()
	for _, inTagParent := range inTag.Parents {
		blobIndex.descendants[inTagParent.Hex()][inTag.Version] = inTag.Hash()
	}
}
