package crawler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	//	"strings"

	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

func webSearchHandler(w http.ResponseWriter, r *http.Request) {

	ResponseStruct := struct {
		Query                  string
		NameSegmentInfo        map[nameSegmentIndexEntry]int
		NameSegmentInfoPresent bool
		BlobInfo               blobIndexEntry
		BlobInfoPresent        bool
		CommitInfo             commitIndexEntry
		CommitInfoPresent      bool
		TagInfo                tagIndexEntry
		TagInfoPresent         bool
	}{}

	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Println(err)
	}
	ResponseStruct.Query = parsedURL.Query().Get("q")
	log.Println(ResponseStruct.Query)
	log.Println(len(ResponseStruct.Query))
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
			case "Repository Key":
				ResponseStruct.CommitInfo, ResponseStruct.CommitInfoPresent =
					commitIndex[ResponseStruct.Query]
			case "Domain Key":
				ResponseStruct.TagInfo, ResponseStruct.TagInfoPresent =
					tagIndex[ResponseStruct.Query]
			default:
			}
		}
		log.Println(err)
	}

	ResponseStruct.NameSegmentInfo,
		ResponseStruct.NameSegmentInfoPresent =
		nameSegmentIndex[ResponseStruct.Query]
	if !ResponseStruct.NameSegmentInfoPresent {
		log.Println("Name Segment is not present")
	}

	t, err := template.New("WebSearch template").Parse(`
	<html>
		<head>
			<title>
			 Search - CCFS
			</title>
		</head>
		<body>
			<form action = "./" method="get">
				<input type = "text" name = "q">
				<input type = "submit" value = "search">
			</form>
			</br>
			{{if .Query}} Search results for: {{.Query}} {{end}}
			</br>
			</br>
			<dl>
			{{if .NameSegmentInfoPresent}}
			<ul>
			{{range $key, $value:= .NameSegmentInfo}}
			<li>
			{{$key.TypeString}}: <a href= "/b/{{$key.Hash}}">{{$key.Hash}}</a>
			</li>
			{{end}}
			</ul>
			{{end}}
			{{if .BlobInfoPresent}}
				</br>
				{{.BlobInfo.TypeString}}[{{.BlobInfo.Size}}]:
				<dl>
					{{range $key, $value := .BlobInfo.NameSeg}}
					<dt> {{$key}}: </dt>
						{{range $key1:= $value}}
						<dd> <a href= "/b/{{$key1}}">{{$key1}}</a> </dd>
						{{end}}
					{{end}}
				</dl>
				<dl>
					{{range $key, $value := .BlobInfo.Descendants}}
					<dt> <a href= "/search/?q={{$value}}"> {{$key}}</a>:
					 <a href= "/b/{{$value}}">{{$value}}</a></dt>
					{{end}}
				</dl>
			{{end}}
			{{if .CommitInfoPresent}} {{.CommitInfo}} {{end}}
			{{if .TagInfoPresent}} {{.TagInfo}} {{end}}
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
