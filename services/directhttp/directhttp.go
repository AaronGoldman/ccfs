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
		urlVertions := fmt.Sprintf(
			"https://%s/c/%s/",
			host,
			h.Hex(),
		)
		bodyVertions, errVertions := urlReadAll(urlVertions)
		if errVertions != nil {
			return objects.Commit{}, errVertions
		}
		vertionNumber := latestsVertion(bodyVertions)
		urlCommit := fmt.Sprintf(
			"https://%s/c/%s/%s",
			host,
			h.Hex(),
			vertionNumber,
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
		bodyVertions, err := urlReadAll(quarryurl)
		if err != nil {
			return objects.Tag{}, err
		}
		vertionNumber := latestsVertion(bodyVertions)
		tagurl := fmt.Sprintf(
			"https://%s/t/%s/%s/%s",
			host,
			h.Hex(),
			namesegment,
			vertionNumber,
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
	for _, host := range hosts {
		quarryurl := fmt.Sprintf(
			"https://%s/t/%s/",
			host,
			h.Hex(),
		)
		body, err := urlReadAll(quarryurl)
		if err != nil {
			return []objects.Tag{}, err
		}
		//ToDo find and get latests vertion of all labels
		vertionNumber := latestsVertion(bodyVertions)
		tagurl := fmt.Sprintf(
			"https://%s/t/%s/%s/%s",
			host,
			h.Hex(),
			namesegment,
			vertionNumber,
		)
		body, err := urlReadAll(tagurl)

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
