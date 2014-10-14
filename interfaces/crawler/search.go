//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package crawler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

func webSearchHandler(w http.ResponseWriter, r *http.Request) {

	ResponseStruct := struct {
		Query             string
		NameSegmentInfos  map[string]map[nameSegmentIndexEntry]int
		BlobInfo          blobIndexEntry
		BlobInfoPresent   bool
		CommitInfo        commitIndexEntry
		CommitInfoPresent bool
		TagInfo           tagIndexEntry
		TagInfoPresent    bool
		Text              map[string][]string
	}{}

	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Println(err)
	}
	ResponseStruct.Query = parsedURL.Query().Get("q")
	hcid, err := objects.HcidFromHex(ResponseStruct.Query)
	if err == nil {
		ResponseStruct.BlobInfo, ResponseStruct.BlobInfoPresent =
			blobIndex[ResponseStruct.Query]
		if !ResponseStruct.BlobInfoPresent {
			log.Println("HID is not present")
		} else {
			switch ResponseStruct.BlobInfo.TypeString {
			case "Blob":
			case "List":
			case "Commit":
				commit, err := services.GetCommitForHcid(hcid)
				if err == nil {
					ResponseStruct.CommitInfo,
						ResponseStruct.CommitInfoPresent =
						commitIndex[commit.Hkid.Hex()]
				}
			case "Tag":
				tag, err := services.GetTagForHcid(hcid)
				if err == nil {
					ResponseStruct.TagInfo, ResponseStruct.TagInfoPresent =
						tagIndex[tag.Hkid.Hex()]
				}
			case "Repository":
				ResponseStruct.CommitInfo, ResponseStruct.CommitInfoPresent =
					commitIndex[ResponseStruct.Query]
			case "Domain":
				ResponseStruct.TagInfo, ResponseStruct.TagInfoPresent =
					tagIndex[ResponseStruct.Query]
			default:
				log.Printf("Unrecognized Type %s",
					ResponseStruct.BlobInfo.TypeString,
				)
			}
		}
		log.Println(err)
	}

	queryTokens := strings.FieldsFunc(ResponseStruct.Query, isSeperator)
	ResponseStruct.NameSegmentInfos = make(map[string]map[nameSegmentIndexEntry]int)
	for nameSegment, nameSegEntry := range nameSegmentIndex {
		for _, queryToken := range queryTokens {
			if present := strings.Contains(nameSegment, queryToken); present {
				ResponseStruct.NameSegmentInfos[nameSegment] = nameSegEntry
				break
			}
		}
	}

	tempMap := make(map[string][]objects.HCID)
	for _, queryToken := range queryTokens {
		if _, present := textIndex[queryToken]; present {
			tempMap[queryToken] = textIndex[queryToken]
		}
	}

	ResponseStruct.Text = make(map[string][]string)
	for token, tempHCIDs := range tempMap {
		for _, tempHCID := range tempHCIDs {
			if _, present := ResponseStruct.Text[tempHCID.Hex()]; !present {
				ResponseStruct.Text[tempHCID.Hex()] = []string{token}
			} else {
				ResponseStruct.Text[tempHCID.Hex()] = append(
					ResponseStruct.Text[tempHCID.Hex()],
					token,
				)
			}
		}
	}

	t, err := template.New("WebSearch template").Funcs(
		template.FuncMap{
			"GetCuratorsofBlob": GetCuratorsofBlob,
			"GetPathsForHCID":   GetPathsForHCID,
		},
	).Parse(`
	{{define "NameSegTemp"}}
		<dl>
			{{range $key, $value := .}}
				<dt> {{$key}}: </dt>
				{{range $key1:= $value}}
					<dd> <a href= "/b/{{$key1}}">{{$key1}}</a> </dd>
				{{end}}
			{{end}}
		</dl>
	{{end}}
	{{define "VersionTemp"}}
		<dl>
			{{range $key, $value := .}}
				<dt>
					<a href= "/search/?q={{$value}}"> {{$key}}</a>:
					 <a href= "/b/{{$value}}">{{$value}}</a>
				</dt>
			{{end}}
		</dl>
	{{end}}
	{{define "BlobInfoTemp"}}
		{{.TypeString}}[{{.Size}}]:
			{{with .SignedBy}}
				<a href= "/search/?q={{.}}">{{.}}</a>:
			{{end}}
		{{template "NameSegTemp" .NameSeg}}
		{{template "VersionTemp" .Descendants}}
		{{range $key, $value := GetCuratorsofBlob .HCID}}
			Curators: <a href= "/{{$value}}/{{$key}}">{{$key}}</a>
		{{end}}
		<dl>
			Paths:
			{{range $key, $value := GetPathsForHCID .HCID}}
				<dt><a href= "/{{$key}}">{{$key}}</a></dt>
			{{end}}
		</dl>
	{{end}}
	<html>
		<head>
			<title>
				Search - CCFS
			</title>
		</head>
		<body>
			</br>
			<form action = "./" method="get">
				<input type = "text" name = "q" value="{{.Query}}">
				<input type = "submit" value="Search">
			</form>
			</br>
			{{with .Query}}
				Search results for: {{.}}
			{{end}}
			</br>
			<dl>
				{{with .NameSegmentInfos}}
					{{range $key1, $value1:= .}}
						{{$key1}}
						<ul>
							{{range $key, $value:= $value1}}
								<li>
									<a href= "/search/?q={{$key.Hash}}">
									{{$key.TypeString}}
									</a>:
									<a href= "/b/{{$key.Hash}}">
									{{$key.Hash}}
									</a>
									{{$value}}
								</li>
							{{end}}
						</ul>
					{{end}}
				{{end}}

				{{with .Text}}
				<br>
				Full Text Search Results:
					<ul>
						{{range $key1, $value1:= .}}
						<a href= "/b/{{$key1}}">
							{{$key1}}
							</a>
							{{range $key:= $value1}}
								<li>
									{{$key}}
								</li>
							{{end}}
						{{end}}
					</ul>
				{{end}}
				{{if .BlobInfoPresent}}
					</br>
					{{template "BlobInfoTemp" .BlobInfo}}
				{{end}}
				{{if .CommitInfoPresent}}
					{{template "NameSegTemp" .CommitInfo.NameSeg}}
					{{template "VersionTemp" .CommitInfo.Version}}
				{{end}}
				{{if .TagInfoPresent}}
					<dl>
						<dt>Aliases: </dt>
						<dd>{{template "NameSegTemp" .TagInfo.NameSeg}}</dd>
						<dt>Sub Domains: </dt>
						<dd>
							<dl>
								{{range $key, $value := .TagInfo.Version}}
									<dt>{{$key}}: </dt>
									<dd>{{template "VersionTemp" $value}}</dd>
								{{end}}
							</dl>
						</dd>
					</dl>
				{{end}}
			</dl>
		</body>
	</html>
`)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("HTTP Error 500 Internal Search server error\n%s\n",
				err,
			),
			500,
		)
	} else {
		t.Execute(w, ResponseStruct)
	}
}

//GetCuratorsofBlob gets curators of the blob
func GetCuratorsofBlob(hcidString string) map[string]string {
	curators := make(map[string]string)
	info, present := blobIndex[hcidString]
	if !present {
		return make(map[string]string)
	}

	for _, RefCommitHCID := range info.RefCommits {
		entry, present := blobIndex[RefCommitHCID]
		if !present {
			continue
		}
		curators[entry.SignedBy] = "r"
	}

	for _, hidValues := range info.NameSeg {
		for _, hidValue := range hidValues {
			segInfo, present := blobIndex[hidValue]
			if !present {
				continue
			}
			switch segInfo.TypeString {
			case "List":
				for key, value := range GetCuratorsofBlob(hidValue) {
					curators[key] = value
				}
			case "Tag":
				curators[segInfo.SignedBy] = "d"
			default:
				log.Printf("cannot switch on typestring %s", segInfo.TypeString)
			}
		}
	}
	return curators
}

//GetPathsForHCID gets the paths for the HCID
func GetPathsForHCID(blobHCID string) map[string]bool {

	tempPaths := map[string]bool{}
	blobPath, present := blobIndex[blobHCID]
	log.Printf("%s", blobPath.TypeString)
	if !present {
		log.Printf("Blob is not present: %s", blobHCID)
		return map[string]bool{}
	}
	switch blobPath.TypeString {
	case "Commit":
		//commitVersion := commitIndex[blobPath.SignedBy].Version
		return map[string]bool{"r/" + blobPath.SignedBy: true}
	case "Tag":
		//TagVersion := tagIndex[blobPath.SignedBy].Version
		return map[string]bool{"d/" + blobPath.SignedBy: true}
	default:
	}
	for blobName, referingHCIDs := range blobPath.NameSeg {
		for _, referingHCID := range referingHCIDs {
			paths := GetPathsForHCID(referingHCID)
			for path, version := range paths {
				tempPaths[path+"/"+blobName] = version
			}
		}
	}

	for _, refCommits := range blobPath.RefCommits {
		paths := GetPathsForHCID(refCommits)
		for path, version := range paths {
			tempPaths[path] = version
		}
	}
	return tempPaths
}

func isSeperator(c rune) bool {
	return !unicode.IsLetter(c) && !unicode.IsNumber(c)
}
