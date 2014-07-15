package crawler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
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
	http.HandleFunc("/index/", webIndexHandler)

}

// This function handles web requests for the crawler
func webCrawlerHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(r.RequestURI[9:], "/", 2)
	hkidhex := parts[0]
	h, hexerr := objects.HcidFromHex(hkidhex)
	if hexerr == nil {
		seedQueue(h)
	}
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

// This function handles web requests for the index of the crawler
func webIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Welcome to the web index!\n")))
	index := struct {
		Blobs   interface{}
		Commits interface{}
		Tags    interface{}
	}{
		Blobs:   blobIndex,
		Commits: commitIndex,
		Tags:    tagIndex,
	}

	t, err := template.New("WebIndex template").Parse(
		`{{define "sliceTemplate"}}{{range $lice:= .}}
				{{$lice}}{{end}}{{end}}{{define "mapTemplate"}}{{range $key, $value:= .}}
				#: {{$key}} HCID: {{$value}}{{end}}{{end}}{{define "blobIndexEntryTemplate"}}
		Type: {{.TypeString}}
		Size: {{.Size}}
		Name Segment: {{range $key, $value:= .NameSeg}}
			{{$key}}{{template "sliceTemplate" $value}}{{end}}
		Descendants: {{range $key, $value:= .Descendants}}
			#: {{$value}} HCID: {{$key}}{{end}}
{{end}}{{define "commitIndexEntryTemplate"}}
		Name Segment: {{range $key, $value:= .NameSeg}}
			{{$key}}{{template "sliceTemplate" $value}}{{end}}
		Version: {{range $key, $value:= .Version}}
			#: {{$key}} HCID: {{$value}}{{end}}
{{end}}{{define "tagIndexEntryTemplate"}}
		Name Segment: {{range $key, $value:= .NameSeg}}
			{{$key}}{{template "sliceTemplate" $value}}{{end}}
		Version: {{range $key, $value:= .Version}}
			Name Segment: {{$key}}{{template "mapTemplate" $value}}{{end}}
{{end}}Blobs	{{range $key, $value:= .Blobs}}
	{{$key}}{{template "blobIndexEntryTemplate" $value}}{{end}}
Commits	{{range $key, $value:= .Commits}}
	{{$key}}{{template "commitIndexEntryTemplate" $value}}{{end}}
Tags		{{range $key, $value:= .Tags}}
	{{$key}}{{template "tagIndexEntryTemplate" $value}}{{end}}
`)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("HTTP Error 500 Internal Indexer server error\n%s\n", err),
			500,
		)
	} else {
		t.Execute(w, index) //merge template ‘t’ with content of ‘index’
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

type blobIndexEntry struct {
	TypeString  string
	Size        int
	NameSeg     map[ /*nameSeg*/ string] /*referringHID*/ []string
	Descendants map[ /*referringHCID*/ string] /*versionNumber*/ int64
}

func (indexEntry blobIndexEntry) insertSize(size int) blobIndexEntry {
	indexEntry.Size = size
	return indexEntry
}

func (indexEntry blobIndexEntry) insertType(typeString string) blobIndexEntry {
	indexEntry.TypeString = typeString
	return indexEntry
}

func (indexEntry blobIndexEntry) insertNameSegment(
	nameSeg string,
	referHID string,
) blobIndexEntry {
	if indexEntry.NameSeg == nil {
		indexEntry.NameSeg = make(map[string][]string)
	}
	if _, present := indexEntry.NameSeg[nameSeg]; !present {
		indexEntry.NameSeg[nameSeg] = []string{referHID}
	} else {
		indexEntry.NameSeg[nameSeg] = append(indexEntry.NameSeg[nameSeg], referHID)
	}
	return indexEntry
}

func (indexEntry blobIndexEntry) insertDescendant(
	versionNumber int64,
	descendantHCID objects.HCID,
) blobIndexEntry {
	if indexEntry.Descendants == nil {
		indexEntry.Descendants = make(map[string]int64)
	}
	if _, present := indexEntry.Descendants[descendantHCID.Hex()]; !present {
		indexEntry.Descendants[descendantHCID.Hex()] = versionNumber
	}

	return indexEntry
}

func insertDescendantS(
	parents []objects.HCID,
	descendant objects.HCID,
	version int64,
) {
	if blobIndex == nil {
		blobIndex = make(map[string]blobIndexEntry)
	}
	for _, entryParent := range parents {
		if _, present := blobIndex[entryParent.Hex()]; !present {
			blobIndex[entryParent.Hex()] = blobIndexEntry{}
		}
		blobIndex[entryParent.Hex()] = blobIndex[entryParent.Hex()].insertDescendant(version, descendant)
	}
}

type commitIndexEntry struct {
	NameSeg map[string][]string
	Version map[ /*versionNumber*/ int64]objects.HCID
}

func (indexEntry commitIndexEntry) insertVersion(
	versionNumber int64,
	instanceHCID objects.HCID,
) commitIndexEntry {
	if indexEntry.Version == nil {
		indexEntry.Version = make(map[int64]objects.HCID)
	}
	indexEntry.Version[versionNumber] = instanceHCID
	return indexEntry
}

func (indexEntry commitIndexEntry) insertNameSegment(
	nameSeg string,
	referHID string,
) commitIndexEntry {
	if indexEntry.NameSeg == nil {
		indexEntry.NameSeg = make(map[string][]string)
	}
	if _, present := indexEntry.NameSeg[nameSeg]; !present {
		indexEntry.NameSeg[nameSeg] = []string{referHID}
	} else {
		indexEntry.NameSeg[nameSeg] = append(indexEntry.NameSeg[nameSeg], referHID)
	}
	return indexEntry
}

type tagIndexEntry struct {
	NameSeg map[string][]string
	Version map[ /*nameSeg*/ string]map[ /*versionNumber*/ int64]objects.HCID
}

func (indexEntry tagIndexEntry) insertNameSegment(
	nameSeg string,
	referHID string,
) tagIndexEntry {
	if indexEntry.NameSeg == nil {
		indexEntry.NameSeg = make(map[string][]string)
	}
	if _, present := indexEntry.NameSeg[nameSeg]; !present {
		indexEntry.NameSeg[nameSeg] = []string{referHID}
	} else {
		indexEntry.NameSeg[nameSeg] = append(indexEntry.NameSeg[nameSeg], referHID)
	}
	return indexEntry
}

func (indexEntry tagIndexEntry) insertVersion(
	versionNumber int64,
	nameSeg string,
	instanceHCID objects.HCID,
) tagIndexEntry {
	if indexEntry.Version == nil {
		indexEntry.Version = make(map[string]map[int64]objects.HCID)
	}
	if _, present := indexEntry.Version[nameSeg]; !present {
		indexEntry.Version[nameSeg] = make(map[int64]objects.HCID)
	}
	indexEntry.Version[nameSeg][versionNumber] = instanceHCID
	return indexEntry
}

func indexNameSegment(typeString, targetHex, referingHex, nameSeg string) {

	switch typeString {
	case "blob", "list":
		if blobIndex == nil {
			blobIndex = make(map[string]blobIndexEntry)
		}
		if _, present := blobIndex[targetHex]; !present {
			blobIndex[targetHex] = blobIndexEntry{}
		}
		blobIndex[targetHex] = blobIndex[targetHex].insertNameSegment(nameSeg, referingHex)

	case "commit":
		if commitIndex == nil {
			commitIndex = make(map[string]commitIndexEntry)
		}
		if _, present := commitIndex[targetHex]; !present {
			commitIndex[targetHex] = commitIndexEntry{}
		}
		commitIndex[targetHex] = commitIndex[targetHex].insertNameSegment(nameSeg, referingHex)
	case "tag":
		if tagIndex == nil {
			tagIndex = make(map[string]tagIndexEntry)
		}
		if _, present := tagIndex[targetHex]; !present {
			tagIndex[targetHex] = tagIndexEntry{}
		}
		tagIndex[targetHex] = tagIndex[targetHex].insertNameSegment(nameSeg, referingHex)
	default:
		log.Printf("Received invalid typestring: %s\n", typeString)
	}
}

var blobIndex map[string]blobIndexEntry

func indexBlob(inBlob objects.Blob) {
	if blobIndex == nil {
		blobIndex = make(map[string]blobIndexEntry)
	}
	hashHex := inBlob.Hash().Hex()
	if _, present := blobIndex[hashHex]; !present {
		blobIndex[hashHex] = blobIndexEntry{}
	}
	blobIndex[hashHex] = blobIndex[hashHex].insertSize(len(inBlob))
	blobIndex[hashHex] = blobIndex[hashHex].insertType("Blob")
}

func indexList(inList objects.List) {
	indexBlob(inList.Bytes()) //Indexing Lists as blobs because they are also blobs
	hashHex := inList.Hash().Hex()
	blobIndex[hashHex] = blobIndex[hashHex].insertType("List")
	for nameSeg, entry := range inList {
		indexNameSegment(
			entry.TypeString,
			entry.Hash.Hex(),
			inList.Hash().Hex(),
			nameSeg,
		)
	}
}

var commitIndex map[string]commitIndexEntry

func indexCommit(inCommit objects.Commit) {
	indexBlob(inCommit.Bytes())
	hashHex := inCommit.Hash().Hex()
	blobIndex[hashHex] = blobIndex[hashHex].insertType("Commit")
	if commitIndex == nil {
		commitIndex = make(map[string]commitIndexEntry)
	}
	if _, present := commitIndex[inCommit.Hkid.Hex()]; !present {
		commitIndex[inCommit.Hkid.Hex()] = commitIndexEntry{}
	}
	commitIndex[inCommit.Hkid.Hex()] = commitIndex[inCommit.Hkid.Hex()].insertVersion(
		inCommit.Version,
		inCommit.Hash(),
	)
	insertDescendantS(inCommit.Parents, inCommit.Hash(), inCommit.Version)
}

var tagIndex map[string]tagIndexEntry

func indexTag(inTag objects.Tag) {
	indexBlob(inTag.Bytes())
	hashHex := inTag.Hash().Hex()
	blobIndex[hashHex] = blobIndex[hashHex].insertType("Tag")
	indexNameSegment(
		inTag.TypeString,
		inTag.HashBytes.Hex(),
		inTag.Hash().Hex(),
		inTag.NameSegment,
	)

	if tagIndex == nil {
		tagIndex = make(map[string]tagIndexEntry)
	}
	if _, present := tagIndex[inTag.Hkid.Hex()]; !present {
		tagIndex[inTag.Hkid.Hex()] = tagIndexEntry{}
	}
	tagIndex[inTag.Hkid.Hex()] = tagIndex[inTag.Hkid.Hex()].insertVersion(
		inTag.Version,
		inTag.NameSegment,
		inTag.Hash(),
	)
	insertDescendantS(inTag.Parents, inTag.Hash(), inTag.Version)
}
