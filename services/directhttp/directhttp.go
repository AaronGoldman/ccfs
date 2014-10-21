//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package directhttp

import (
	"fmt"
	"io/ioutil"
	"net/http"

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
		quarryurl := fmt.Sprintf(
			"https://%s/c/%s/",
			host,
			h.Hex(),
		)
		resp, err := http.Get(quarryurl)
		if err != nil {
			return objects.Commit{}, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return objects.Commit{}, err
		}
		//ToDo find and get latests vertion
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

//ID gets the ID string
func (d directhttpservice) ID() string {
	return "directhttp"
}
