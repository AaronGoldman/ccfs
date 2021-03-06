//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package crawler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/AaronGoldman/ccfs/objects"
	//	"github.com/AaronGoldman/ccfs/services"
)

// This function handles web requests for the index of the crawler
func webIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Welcome to the web index!\n")))
	index := struct {
		NameSegments interface{}
		Blobs        interface{}
		Commits      interface{}
		Tags         interface{}
		Texts        interface{}
	}{
		NameSegments: nameSegmentIndex,
		Blobs:        blobIndex,
		Commits:      commitIndex,
		Tags:         tagIndex,
		Texts:        textIndex,
	}

	t, err := template.New("WebIndex template").Parse(
		`{{define "sliceTemplate"}}{{range $lice:= .}}
				{{$lice}}{{end}}{{end}}{{define "mapTemplate"}}{{range $key, $value:= .}}
				Version: {{$key}} HCID: {{$value}}{{end}}{{end}}{{define "nameSegmentIndexEntryTemplate"}}{{range $key, $value:= .}}
			Type: {{$key.TypeString}} HID: {{$key.Hash}} Count: {{$value}}{{end}}
{{end}}{{define "blobIndexEntryTemplate"}}
		Type: {{.TypeString}}
		Size: {{.Size}}
		HCID: {{.HCID}}
		Name Segment: {{range $key, $value:= .NameSeg}}
			{{$key}}{{template "sliceTemplate" $value}}{{end}}
		Referring Commits: {{range $key:= .RefCommits}}
			{{$key}}{{end}}
		Signed By: {{.SignedBy}}
		Descendants: {{range $key, $value:= .Descendants}}
			Version: {{$key}} HCID: {{$value}}{{end}}
{{end}}{{define "commitIndexEntryTemplate"}}
		Aliases: {{range $key, $value:= .NameSeg}}
			{{$key}}{{template "sliceTemplate" $value}}{{end}}
		Collection: {{range $key, $value:= .Version}}
			Version: {{$key}} HCID: {{$value}}{{end}}
{{end}}{{define "tagIndexEntryTemplate"}}
		Aliases: {{range $key, $value:= .NameSeg}}
			{{$key}}{{template "sliceTemplate" $value}}{{end}}
		Collection: {{range $key, $value:= .Version}}
			Name Segment: {{$key}}{{template "mapTemplate" $value}}{{end}}
{{end}}Name Segments	{{range $key, $value:= .NameSegments}}
	{{$key}}{{template "nameSegmentIndexEntryTemplate" $value}}{{end}}
Blobs	{{range $key, $value:= .Blobs}}
	{{$key}}{{template "blobIndexEntryTemplate" $value}}{{end}}
Commits	{{range $key, $value:= .Commits}}
	{{$key}}{{template "commitIndexEntryTemplate" $value}}{{end}}
Tags		{{range $key, $value:= .Tags}}
	{{$key}}{{template "tagIndexEntryTemplate" $value}}{{end}}
Texts {{range $key, $value := .Texts}}
	{{$key}}
		{{range $values:= $value}}
			{{$values}}{{end}}{{end}}
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

type blobIndexEntry struct {
	TypeString string
	Size       int
	HCID       string
	NameSeg    map[ /*nameSeg*/ string] /*referringHCID*/ []string
	RefCommits/*referringHCID*/ []string
	SignedBy    string
	Descendants map[ /*versionNumber*/ int64] /*referringHCID*/ string
}

func (indexEntry blobIndexEntry) insertSize(size int) blobIndexEntry {
	indexEntry.Size = size
	return indexEntry
}

func (indexEntry blobIndexEntry) insertType(typeString string) blobIndexEntry {
	if indexEntry.TypeString != "" && typeString == "blob" {
		return indexEntry
	}
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
		indexEntry.NameSeg[nameSeg] = append(
			indexEntry.NameSeg[nameSeg],
			referHID,
		)
	}
	return indexEntry
}

func (indexEntry blobIndexEntry) insertCollection(collectionKey string) blobIndexEntry {
	indexEntry.SignedBy = collectionKey
	return indexEntry
}

func (indexEntry blobIndexEntry) insertRefCommits(refCommit string) blobIndexEntry {
	indexEntry.RefCommits = append(
		indexEntry.RefCommits,
		refCommit,
	)
	return indexEntry
}

func (indexEntry blobIndexEntry) insertDescendant(
	versionNumber int64,
	descendantHCID objects.HCID,
) blobIndexEntry {
	if indexEntry.Descendants == nil {
		indexEntry.Descendants = make(map[int64]string)
	}
	if _, present := indexEntry.Descendants[versionNumber]; !present {
		indexEntry.Descendants[versionNumber] = descendantHCID.Hex()
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
			blobIndex[entryParent.Hex()] = blobIndexEntry{HCID: entryParent.Hex()}
		}
		blobIndex[entryParent.Hex()] =
			blobIndex[entryParent.Hex()].insertDescendant(version, descendant)
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
		indexEntry.NameSeg[nameSeg] = append(
			indexEntry.NameSeg[nameSeg],
			referHID,
		)
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
		indexEntry.NameSeg[nameSeg] = append(
			indexEntry.NameSeg[nameSeg],
			referHID,
		)
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
	if nameSegmentIndex == nil {
		nameSegmentIndex = make(map[string]map[nameSegmentIndexEntry]int)
	}

	if _, present := nameSegmentIndex[nameSeg]; !present {
		nameSegmentIndex[nameSeg] = make(map[nameSegmentIndexEntry]int)
	}
	nameSegmentIndex[nameSeg][nameSegmentIndexEntry{typeString, targetHex}]++

	switch typeString {
	case "blob", "list":
		if blobIndex == nil {
			blobIndex = make(map[string]blobIndexEntry)
		}
		if _, present := blobIndex[targetHex]; !present {
			blobIndex[targetHex] = blobIndexEntry{HCID: targetHex}
		}
		blobIndex[targetHex] = blobIndex[targetHex].insertNameSegment(
			nameSeg,
			referingHex,
		)

	case "commit":
		if commitIndex == nil {
			commitIndex = make(map[string]commitIndexEntry)
		}
		if _, present := commitIndex[targetHex]; !present {
			commitIndex[targetHex] = commitIndexEntry{}
		}
		commitIndex[targetHex] = commitIndex[targetHex].insertNameSegment(
			nameSeg,
			referingHex,
		)
	case "tag":
		if tagIndex == nil {
			tagIndex = make(map[string]tagIndexEntry)
		}
		if _, present := tagIndex[targetHex]; !present {
			tagIndex[targetHex] = tagIndexEntry{}
		}
		tagIndex[targetHex] = tagIndex[targetHex].insertNameSegment(
			nameSeg,
			referingHex,
		)
	default:
		log.Printf("Received invalid typestring: %s\n", typeString)
	}
}

var nameSegmentIndex map[string]map[nameSegmentIndexEntry]int

type nameSegmentIndexEntry struct {
	TypeString string
	Hash       string
}

var textIndex map[ /*Tokens in blob*/ string][]objects.HCID /*HCIDs containing token*/

//type textIndexEntry struct {
//	blobToken string
//	blobHCID objects.HCID
//}
func indexText(inBlob objects.Blob) {
	if textIndex == nil {
		textIndex = make(map[string][]objects.HCID)
	}
	tokens := strings.FieldsFunc(string(inBlob), isSeperator)
	for _, token := range tokens {
		if _, present := textIndex[token]; !present {
			textIndex[token] = []objects.HCID{inBlob.Hash()}
		} else {
			textIndex[token] = append(
				textIndex[token],
				inBlob.Hash(),
			)
		}
	}
}

var blobIndex map[ /*HCID*/ string]blobIndexEntry

func indexBlob(inBlob objects.Blob) {
	if blobIndex == nil {
		blobIndex = make(map[string]blobIndexEntry)
	}
	hashHex := inBlob.Hash().Hex()
	if _, present := blobIndex[hashHex]; !present {
		blobIndex[hashHex] = blobIndexEntry{HCID: hashHex}
	}
	blobIndex[hashHex] = blobIndex[hashHex].insertSize(len(inBlob))
	blobIndex[hashHex] = blobIndex[hashHex].insertType("Blob")

	indexText(inBlob)

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

var commitIndex map[ /*Hkid*/ string]commitIndexEntry

func indexCommit(inCommit objects.Commit) {
	indexBlob(inCommit.Bytes())
	hashHex := inCommit.Hash().Hex()
	blobIndex[hashHex] = blobIndex[hashHex].insertType("Commit")
	blobIndex[inCommit.Hkid.Hex()] =
		blobIndex[inCommit.Hkid.Hex()].insertType("Repository")
	blobIndex[hashHex] =
		blobIndex[hashHex].insertCollection(inCommit.Hkid.Hex())
	if _, present := blobIndex[inCommit.ListHash.Hex()]; !present {
		blobIndex[inCommit.ListHash.Hex()] =
			blobIndexEntry{HCID: inCommit.ListHash.Hex()}
	}
	blobIndex[inCommit.ListHash.Hex()] =
		blobIndex[inCommit.ListHash.Hex()].insertRefCommits(hashHex)

	if commitIndex == nil {
		commitIndex = make(map[string]commitIndexEntry)
	}
	if _, present := commitIndex[inCommit.Hkid.Hex()]; !present {
		commitIndex[inCommit.Hkid.Hex()] = commitIndexEntry{}
	}
	commitIndex[inCommit.Hkid.Hex()] =
		commitIndex[inCommit.Hkid.Hex()].insertVersion(
			inCommit.Version,
			inCommit.Hash(),
		)
	insertDescendantS(inCommit.Parents, inCommit.Hash(), inCommit.Version)
}

var tagIndex map[ /*Hkid*/ string]tagIndexEntry

func indexTag(inTag objects.Tag) {
	indexBlob(inTag.Bytes())
	hashHex := inTag.Hash().Hex()
	blobIndex[hashHex] = blobIndex[hashHex].insertType("Tag")
	blobIndex[inTag.Hkid.Hex()] =
		blobIndex[inTag.Hkid.Hex()].insertType("Domain")
	blobIndex[hashHex] =
		blobIndex[hashHex].insertCollection(inTag.Hkid.Hex())
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
