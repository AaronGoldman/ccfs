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
	http.HandleFunc("/index/", webIndexHandler)

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

// This function handles web requests for the index of the crawler
func webIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Welcome to the web index!\n")))

	w.Write([]byte(fmt.Sprintf("Blobs\n")))
	w.Write([]byte(blobMaptoString(blobIndex)))

	w.Write([]byte(fmt.Sprintf("Commits\n")))
	w.Write([]byte(commitMaptoString(commitIndex)))

	w.Write([]byte(fmt.Sprintf("Tags\n")))
	w.Write([]byte(tagMaptoString(tagIndex)))
}

func blobMaptoString(hashMap map[string]blobIndexEntry) string {
	document := ""
	var keys []string
	for key := range hashMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		mapLine := fmt.Sprintf("\t%s\n%v\n", key, hashMap[key])
		document += mapLine
		//w.Write([]byte(mapLine))
	}
	return document
}

func sliceStringMaptoString(hashMap map[string][]string) string {
	document := ""
	var keys []string
	for key := range hashMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		mapLine := fmt.Sprintf("\t\t\t%s\n%v\n", key, sliceStringtoString(hashMap[key]))
		document += mapLine
		//w.Write([]byte(mapLine))
	}
	return document
}

func sliceStringtoString(strings []string) string {
	document := ""
	for _, key := range strings {
		mapLine := fmt.Sprintf("\t\t\t\t%s\n", key)
		document += mapLine
	}
	return document
}

func commitMaptoString(hashMap map[string]commitIndexEntry) string {
	document := ""
	var keys []string
	for key := range hashMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		mapLine := fmt.Sprintf("\t%s\n%v\n", key, hashMap[key])
		document += mapLine
		//w.Write([]byte(mapLine))
	}
	return document
}

func tagMaptoString(hashMap map[string]tagIndexEntry) string {
	document := ""
	var keys []string
	for key := range hashMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		mapLine := fmt.Sprintf("\t%s\n%v\n", key, hashMap[key])
		document += mapLine
		//w.Write([]byte(mapLine))
	}
	return document
}

type int64arr []int64

func (data int64arr) Len() int { return len(data) }

func (data int64arr) Swap(i, j int) { data[i], data[j] = data[j], data[i] }

func (data int64arr) Less(i, j int) bool { return data[i] < data[j] }

func sortint64(sortData int64arr) { sort.Sort(sortData) }

func intMaptoString(hashMap map[int64]objects.HCID) string {
	document := ""
	var keys []int64
	for key := range hashMap {
		keys = append(keys, key)
	}
	sortint64(keys)
	for _, key := range keys {
		mapLine := fmt.Sprintf("\t\t\tVersion: %d HCID: %v\n", key, hashMap[key])
		document += mapLine
		//w.Write([]byte(mapLine))
	}
	return document
}

func stringIntMaptoString(hashMap map[string]map[int64]objects.HCID) string {
	document := ""
	var keys []string
	for key := range hashMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		mapLine := intMaptoString(hashMap[key])
		document += mapLine + "\n"
		//w.Write([]byte(mapLine))
	}
	return document
}

func stringMaptoInt64(hashMap map[string]int64) string {
	document := ""
	var keys []string
	for key := range hashMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		mapLine := fmt.Sprintf("\t%s\n%v\n", key, hashMap[key])
		document += mapLine
		//w.Write([]byte(mapLine))
	}
	return document
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
	//_ = nextBlob
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
		//commitHKID := firstCommit.Hkid()

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
		//tagHKID := inTag.Hkid

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

type blobIndexEntry struct { //map[blobHcid string]struct
	typeString string
	size       int
	nameSeg    map[string][]string
	//map[nameSeg string]referringHID string
	descendants map[string]int64
	//map[versionNumber int64]referringHCID string
}

func (indexEntry blobIndexEntry) String() string {
	return fmt.Sprintf(
		"\t\tType: %s\n\t\tSize: %d\n\t\tName Segments:\n%s\t\tDescendants:\n%s",
		indexEntry.typeString,
		indexEntry.size,
		sliceStringMaptoString(indexEntry.nameSeg),
		stringMaptoInt64(indexEntry.descendants),
	)
}

func (indexEntry blobIndexEntry) insertSize(size int) blobIndexEntry {
	indexEntry.size = size
	return indexEntry
}

func (indexEntry blobIndexEntry) insertType(typeString string) blobIndexEntry {
	indexEntry.typeString = typeString
	return indexEntry
}

func (indexEntry blobIndexEntry) insertNameSegment(
	nameSeg string,
	referHID string,
) blobIndexEntry {
	if indexEntry.nameSeg == nil {
		indexEntry.nameSeg = make(map[string][]string)
	}
	if _, present := indexEntry.nameSeg[nameSeg]; !present {
		indexEntry.nameSeg[nameSeg] = []string{referHID}
	} else {
		indexEntry.nameSeg[nameSeg] = append(indexEntry.nameSeg[nameSeg], referHID)
	}
	return indexEntry
}

func (indexEntry blobIndexEntry) insertDescendant(
	versionNumber int64,
	descendantHCID objects.HCID,
) blobIndexEntry {
	if indexEntry.descendants == nil {
		indexEntry.descendants = make(map[string]int64)
	}
	if _, present := indexEntry.descendants[descendantHCID.Hex()]; !present {
		indexEntry.descendants[descendantHCID.Hex()] = versionNumber
	}

	return indexEntry
}

func insertDescendantS(parents []objects.HCID, version int64) {
	if blobIndex == nil {
		blobIndex = make(map[string]blobIndexEntry)
	}
	for _, entryParent := range parents {

		if _, present := blobIndex[entryParent.Hex()]; !present {
			blobIndex[entryParent.Hex()] = blobIndexEntry{}
		}
		blobIndex[entryParent.Hex()].insertDescendant(version, entryParent)
	}
}

type commitIndexEntry struct { //map[hkid string] struct
	nameSeg map[string][]string
	version map[int64]objects.HCID
	//map[versionNumber int64]HCIDversion string
}

func (indexEntry commitIndexEntry) String() string {
	return fmt.Sprintf(
		"\t\tName Segments:\n%s\n\t\tVersion:\n%s\n",
		sliceStringMaptoString(indexEntry.nameSeg),
		intMaptoString(indexEntry.version),
	)
}

func (indexEntry commitIndexEntry) insertVersion(
	versionNumber int64,
	instanceHCID objects.HCID,
) commitIndexEntry {
	if indexEntry.version == nil {
		indexEntry.version = make(map[int64]objects.HCID)
	}
	indexEntry.version[versionNumber] = instanceHCID
	return indexEntry
}

func (indexEntry commitIndexEntry) insertNameSegment(
	nameSeg string,
	referHID string,
) commitIndexEntry {
	if indexEntry.nameSeg == nil {
		indexEntry.nameSeg = make(map[string][]string)
	}
	if _, present := indexEntry.nameSeg[nameSeg]; !present {
		indexEntry.nameSeg[nameSeg] = []string{referHID}
	} else {
		indexEntry.nameSeg[nameSeg] = append(indexEntry.nameSeg[nameSeg], referHID)
	}
	return indexEntry
}

type tagIndexEntry struct { //map[HKID string]struct
	nameSeg map[string][]string
	version map[string]map[int64]objects.HCID
	//version map[nameSeg string]map[versionNumber int64]objects.HCID
}

func (indexEntry tagIndexEntry) String() string {
	return fmt.Sprintf(
		"\t\tName Segments:\n%s\n\t\tVersion:\n%s\n",
		sliceStringMaptoString(indexEntry.nameSeg),
		stringIntMaptoString(indexEntry.version),
	)
}

func (indexEntry tagIndexEntry) insertNameSegment(
	nameSeg string,
	referHID string,
) tagIndexEntry {
	if indexEntry.nameSeg == nil {
		indexEntry.nameSeg = make(map[string][]string)
	}
	if _, present := indexEntry.nameSeg[nameSeg]; !present {
		indexEntry.nameSeg[nameSeg] = []string{referHID}
	} else {
		indexEntry.nameSeg[nameSeg] = append(indexEntry.nameSeg[nameSeg], referHID)
	}
	return indexEntry
}

func (indexEntry tagIndexEntry) insertVersion(
	versionNumber int64,
	nameSeg string,
	instanceHCID objects.HCID,
) tagIndexEntry {
	if indexEntry.version == nil {
		indexEntry.version = make(map[string]map[int64]objects.HCID)
	}
	if _, present := indexEntry.version[nameSeg]; !present {
		indexEntry.version[nameSeg] = make(map[int64]objects.HCID)
	}
	indexEntry.version[nameSeg][versionNumber] = instanceHCID
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
	insertDescendantS(inCommit.Parents, inCommit.Version)
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
		tagIndex[inTag.Hkid.Hex()] = tagIndexEntry{} //  make(map[string]map[int64]objects.HCID)
	}
	tagIndex[inTag.Hkid.Hex()] = tagIndex[inTag.Hkid.Hex()].insertVersion(
		inTag.Version,
		inTag.NameSegment,
		inTag.Hash(),
	)
	insertDescendantS(inTag.Parents, inTag.Version)
}
