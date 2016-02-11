//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package directhttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

var running bool
var hosts []string

//Instance is the instance of the directhttpservice
var Instance directhttpservice

func init() {
	services.Registercommand(
		Instance,
		"directhttp command", //This is the usage string
	)
	services.Registerrunner(Instance)
}

//Start registers directhttpservice instances
func Start( /*remotes []string*/ ) {
	services.Registerblobgeter(Instance)
	services.Registercommitgeter(Instance)
	services.Registertaggeter(Instance)
	running = true
	//hosts = remotes
}

//Stop deregisters directhttp instances
func Stop() {
	services.DeRegisterblobgeter(Instance)
	services.DeRegistercommitgeter(Instance)
	services.DeRegistertaggeter(Instance)
	running = false
}

func (d directhttpservice) Command(command string) {
	tokens := strings.FieldsFunc(command, unicode.IsSpace)
	if len(tokens) == 0 {
		return
	}
	switch tokens[0] {
	case "start":
		Start()

	case "stop":
		Stop()

	case "add":
		if len(tokens) > 1 {
			hosts = append(hosts, tokens[1])
		}
		d.List()

	case "remove":
		if len(tokens) > 1 {
			index, err := strconv.Atoi(tokens[1])
			if err == strconv.ErrSyntax {
				fmt.Printf("Please indicate the host by host number.\n")
				return
			} else if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			if index < 0 || index >= len(hosts) {
				fmt.Print("Index out of range\n")
				return
			}

			if index < len(hosts)-1 {
				hosts = append(hosts[:index], hosts[index+1:]...)
			} else {
				hosts = hosts[:index]
			}
		}

		d.List()

	case "list":
		d.List()

	default:
		fmt.Printf("Direct HTTP Service Command Line\n" +
			"start\n" +
			"stop\n" +
			"add [Domain]\n" +
			"remove [Domain]\n" +
			"list\n")
		return
	}

}

//Prints the list of directhttp addresses
func (d directhttpservice) List() {
	if len(hosts) < 1 {
		fmt.Println("No hosts to list")
	}
	for index, host := range hosts {
		fmt.Printf("Host %d: %s\n", index, host)
	}
}

//Running returns a bool that indicates the registration status of the service
func (d directhttpservice) Running() bool {
	return running
}

type directhttpservice struct{}

func (d directhttpservice) GetBlob(h objects.HCID) (objects.Blob, error) {
	for _, host := range hosts {
		quarryurl := fmt.Sprintf(
			"https://%s/b/%s",
			host,
			h.Hex(),
		)
		body, err := urlReadAll(quarryurl)
		if err != nil {
			return objects.Blob{}, err
		}
		return objects.Blob(body), err
	}
	return objects.Blob{}, fmt.Errorf("No Hosts")
}

func (d directhttpservice) GetCommit(h objects.HKID) (objects.Commit, error) {
	for _, host := range hosts {
		urlVersions := fmt.Sprintf(
			"https://%s/c/%s/",
			host,
			h.Hex(),
		)
		bodyVersions, errVersions := urlReadAll(urlVersions)
		if errVersions != nil {
			return objects.Commit{}, errVersions
		}
		versionNumber := latestsVersion(bodyVersions)
		urlCommit := fmt.Sprintf(
			"https://%s/c/%s/%s",
			host,
			h.Hex(),
			versionNumber,
		)
		body, err := urlReadAll(urlCommit)
		if err != nil {
			return objects.Commit{}, err
		}
		commit, err := objects.CommitFromBytes(body)
		if err == nil {
			return commit, err
		}
	}
	return objects.Commit{}, fmt.Errorf("No Hosts")
}

func (d directhttpservice) GetTag(h objects.HKID, namesegment string) (
	objects.Tag, error) {
	for _, host := range hosts {
		quarryurl := fmt.Sprintf(
			"https://%s/t/%s/%s",
			host,
			h.Hex(),
			namesegment,
		)
		bodyVersions, err := urlReadAll(quarryurl)
		if err != nil {
			return objects.Tag{}, err
		}
		versionNumber := latestsVersion(bodyVersions)
		tagurl := fmt.Sprintf(
			"https://%s/t/%s/%s/%s",
			host,
			h.Hex(),
			namesegment,
			versionNumber,
		)
		body, err := urlReadAll(tagurl)
		tag, err := objects.TagFromBytes(body)
		if err == nil {
			return tag, err
		}
	}
	return objects.Tag{}, fmt.Errorf("No Hosts")
}

func (d directhttpservice) GetTags(h objects.HKID) (
	tags []objects.Tag, err error) {
	if len(hosts) == 0 {
		return []objects.Tag{}, fmt.Errorf("No Hosts")
	}

	var maxTags map[string]objects.Tag
	for _, host := range hosts {
		quarryurl := fmt.Sprintf(
			"http://%s/t/%s/",
			host,
			h.Hex(),
		)
		bodyNameSegments, err := urlReadAll(quarryurl)
		if err != nil {
			continue
			//return []objects.Tag{}, err
		}
		//ToDo find and get latests version of all labels
		namesegments := allNameSegments(bodyNameSegments)
		for _, namesegment := range namesegments {
			versionurl := fmt.Sprintf(
				"https://%s/t/%s/%s",
				host,
				h.Hex(),
				namesegment,
			)
			versionBody, err := urlReadAll(versionurl)
			if err != nil {
				continue
			}
			versionNumber := latestsVersion(versionBody)

			tagurl := fmt.Sprintf(
				"https://%s/t/%s/%s/%s",
				host,
				h.Hex(),
				namesegment,
				versionNumber,
			)

			body, err := urlReadAll(tagurl)
			if err != nil {
				continue
			}

			tag, err := objects.TagFromBytes(body)
			if err != nil {
				continue
				//return []objects.Tag{tag}, err
			}
			maxTag, present := maxTags[string(namesegment)]
			if !present || maxTag.Version < tag.Version {
				maxTags[string(namesegment)] = tag
			}
		}
	}
	for _, maxTag := range maxTags {
		tags = append(tags, maxTag)
	}
	if len(tags) == 0 {
		return []objects.Tag{}, fmt.Errorf("No Tags found")
	}
	return tags, nil
}

//ID gets the ID string
func (d directhttpservice) ID() string {
	return "directhttp"
}

func urlReadAll(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

//func latestsVertion(htmldoc bufio.Reader) (vertionNumber string, err error) {
func latestsVersion(htmldoc []byte) (vertionNumber []byte) {
	maxvertionNumber := []byte{}
	for _, line := range bytes.Fields(htmldoc) {
		queryTokens := bytes.FieldsFunc(line, isNotDigit)
		if len(queryTokens) > 0 && bytes.Compare(queryTokens[0], maxvertionNumber) > 0 {
			maxvertionNumber = queryTokens[0]
		}
	}
	return maxvertionNumber
}

func isNotDigit(r rune) bool { return !unicode.IsDigit(r) }

func allNameSegments(html []byte) (nameSegments [][]byte) {
	lines := bytes.Split(html, []byte("\n"))
	for _, line := range lines {
		anchor := bytes.SplitAfterN(line, []byte("href=\""), 2)
		if len(anchor) > 1 {
			reference := bytes.SplitN(anchor[1], []byte("/\""), 2)
			nameSegments = append(nameSegments, reference[0])
		}
	}
	return
}
