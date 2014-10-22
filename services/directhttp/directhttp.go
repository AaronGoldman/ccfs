//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package directhttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"unicode"

	"github.com/AaronGoldman/ccfs/objects"
)

var running bool
var hosts []string

//Instance is the instance of the directhttpservice
var Instance directhttpservice

//Start registers directhttpservice instances
func Start(remotes []string) {
	running = true
	hosts = remotes
}

type directhttpservice struct{}

func (d directhttpservice) GetBlob(h objects.HCID) (objects.Blob, error) {
	for host := range hosts {
		quarryurl := fmt.Sprintf(
			"https://%s/b/%s",
			host,
			h.Hex(),
		)
		resp, err := http.Get(quarryurl)
		if err != nil {
			return objects.Blob{}, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return objects.Blob{}, err
		}
		return objects.Blob(body), err
	}
	return objects.Blob{}, fmt.Errorf("No Hosts")
}
func (d directhttpservice) GetCommit(h objects.HKID) (objects.Commit, error) {
	for host := range hosts {
		urlVertions := fmt.Sprintf(
			"https://%s/c/%s/",
			host,
			h.Hex(),
		)
		respVertions, err := http.Get(urlVertions)
		if err != nil {
			return objects.Commit{}, err
		}
		defer respVertions.Body.Close()
		body, err := ioutil.ReadAll(respVertions.Body)
		if err != nil {
			return objects.Commit{}, err
		}
		vertionNumber := latestsVertion(body)
		urlCommit := fmt.Sprintf(
			"https://%s/c/%s/%s",
			host,
			h.Hex(),
			vertionNumber,
		)
		respCommit, err := http.Get(urlCommit)
		body, err = ioutil.ReadAll(respCommit.Body)
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
	for host := range hosts {
		quarryurl := fmt.Sprintf(
			"https://%s/t/%s/%s",
			host,
			h.Hex(),
			namesegment,
		)
		resp, err := http.Get(quarryurl)
		if err != nil {
			return objects.Tag{}, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return objects.Tag{}, err
		}
		//ToDo find and get latests vertion
		tag, err := objects.TagFromBytes(body)
		if err == nil {
			return tag, err
		}
	}
	return objects.Tag{}, fmt.Errorf("No Hosts")
}

func (d directhttpservice) GetTags(h objects.HKID) (
	tags []objects.Tag, err error) {
	for host := range hosts {
		quarryurl := fmt.Sprintf(
			"https://%s/t/%s/",
			host,
			h.Hex(),
		)
		resp, err := http.Get(quarryurl)
		if err != nil {
			return []objects.Tag{}, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []objects.Tag{}, err
		}
		//ToDo find and get latests vertion of all labels
		//vertionNumber := latestsVertion(body)

		tag, err := objects.TagFromBytes(body)
		if err == nil {
			return []objects.Tag{tag}, err
		}
	}
	return []objects.Tag{}, fmt.Errorf("No Hosts")
}

//ID gets the ID string
func (d directhttpservice) ID() string {
	return "directhttp"
}

//func latestsVertion(htmldoc bufio.Reader) (vertionNumber string, err error) {
func latestsVertion(htmldoc []byte) (vertionNumber []byte) {
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
