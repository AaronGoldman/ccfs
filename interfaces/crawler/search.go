package crawler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

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
			case "list":
			case "commit":
				commit, err := services.GetCommitForHcid(hcid)
				if err == nil {
					ResponseStruct.CommitInfo,
						ResponseStruct.CommitInfoPresent =
						commitIndex[commit.Hkid.Hex()]
				}
			case "tag":
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
				log.Printf("Unrecognized Type",
					ResponseStruct.BlobInfo.TypeString,
				)
			}
		}
		log.Println(err)
	}

	nameSegFields := strings.Fields(ResponseStruct.Query)
	ResponseStruct.NameSegmentInfos = make(map[string]map[nameSegmentIndexEntry]int)
	for nameSegment, nameSegEntry := range nameSegmentIndex {
		for _, nameSegField := range nameSegFields {
			present := strings.Contains(nameSegment, nameSegField)
			if present {
				ResponseStruct.NameSegmentInfos[nameSegment] = nameSegEntry
				break
			}
		}
	}

	t, err := template.New("WebSearch template").Parse(`
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
				{{if .BlobInfoPresent}}
					</br>
					{{.BlobInfo.TypeString}}[{{.BlobInfo.Size}}]:
					{{with .BlobInfo.Collection}}
					<a href= "/search/?q={{.}}">{{.}}</a>:
					{{end}}
					{{template "NameSegTemp" .BlobInfo.NameSeg}}
					{{template "VersionTemp" .BlobInfo.Descendants}}
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
