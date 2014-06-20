package crawler

import (
	"fmt"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
	"log"
	"net/http"
	"strings"
)

var queuedTargets map[target]bool
var targetQueue chan target

func Start() {
	fmt.Printf("Crawler Starting")
	http.HandleFunc("/crawler/", WebCrawlerHandler)

	queuedTargets = make(map[target]bool)
	targetQueue = make(chan target, 100)
}

func WebCrawlerHandler(w http.ResponseWriter, r *http.Request) {

	parts := strings.SplitN(r.RequestURI[9:], "/", 2)
	hkidhex := parts[0]

	if len(hkidhex) == 64 {
		h, err := objects.HkidFromHex(hkidhex)
		if err == nil {
			response := fmt.Sprintf("The hkid gotten is %s", h)
			w.Write([]byte(response)) //converts the response to bytes from strings
			return
		} else {
			// User Interface (to ask user for valid HKID) to be added here
		}
	}
	http.Error(w, "HTTP Error 500 Internal Crawler server error\n\n", 500)
}

type target struct {
	typeString string
	hash       objects.HID
}

func ContentFiles(h objects.HID) {
	hlist := target{
		typeString: "list",
		hash:       objects.HCID(h.Bytes()),
	}
	targetQueue <- hlist

	hblob := target{
		typeString: "blob",
		hash:       objects.HCID(h.Bytes()),
	}
	targetQueue <- hblob

	hcommit := target{
		typeString: "commit",
		hash:       objects.HKID(h.Bytes()),
	}
	targetQueue <- hcommit

	h_hcid_commit := target{
		typeString: "hcid_commit",
		hash:       objects.HCID(h.Bytes()),
	}
	targetQueue <- h_hcid_commit

	h_hcid_tag := target{
		typeString: "hcid_tag",
		hash:       objects.HCID(h.Bytes()),
	}
	targetQueue <- h_hcid_tag
}

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
		if !queuedTargets[newlistHash] {
			targetQueue <- newlistHash
			queuedTargets[newlistHash] = true
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

func handleCommit(inCommit objects.Commit) {
	newHash := target{
		typeString: "list",
		hash:       inCommit.ListHash(),
	}
	newParent := target{
		typeString: "hcid_commit",
		hash:       inCommit.Parent(),
	}
	//commitHKID := firstCommit.Hkid()

	if !queuedTargets[newHash] {
		targetQueue <- newHash
		queuedTargets[newHash] = true
	}
	if !queuedTargets[newParent] {
		targetQueue <- newParent
		queuedTargets[newParent] = true
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

func handleTag(inTag objects.Tag) {
	newHash := target{
		typeString: inTag.TypeString,
		hash:       inTag.HashBytes,
	}
	//newParent := target{
	//	typeString: "hcid_tag",
	//	hash:       inTag.Parent,
	//}
	//commitHKID := inTag.Hkid()

	if !queuedTargets[newHash] {
		targetQueue <- newHash
		queuedTargets[newHash] = true
	}
	//if !queuedTargets[newParent] {
	//	targetQueue <- newParent
	//	queuedTargets[newParent] = true
	//}

}
