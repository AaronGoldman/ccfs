//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

//Package kademliadht is the kademliaservice
package kademliadht

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

//Instance is the instance of the kademliaservice
var Instance kademliaservice
var running bool

func init() {
	services.Registercommand(
		Instance,
		"kademliadht command", //This is the usage string
	)
	services.Registerrunner(Instance)
}

//Start registers kademliadhtservice instances
func Start() {
	Instance = kademliaservicefactory()
	services.Registerblobgeter(Instance)
	services.Registercommitgeter(Instance)
	services.Registertaggeter(Instance)
	//services.Registertagsgeter(Instance)
	services.Registerkeygeter(Instance)
	services.Registerblobposter(Instance)
	services.Registercommitposter(Instance)
	services.Registertagposter(Instance)
	services.Registerkeyposter(Instance)
	running = true

}

//Stop deregisters kademliadhtservice instances
func Stop() {
	services.DeRegisterblobgeter(Instance)
	services.DeRegistercommitgeter(Instance)
	services.DeRegistertaggeter(Instance)
	//services.DeRegistertagsgeter(Instance)
	services.DeRegisterkeygeter(Instance)
	services.DeRegisterblobposter(Instance)
	services.DeRegistercommitposter(Instance)
	services.DeRegistertagposter(Instance)
	services.DeRegisterkeyposter(Instance)
	running = false
}

type kademliaservice struct {
	url string
}

//Running returns a bool that indicates the registration status of the service
func (k kademliaservice) Running() bool {
	return running
}

//ID gets the ID string
func (k kademliaservice) ID() string {
	return "kademliadht"
}

func (k kademliaservice) Command(command string) {
	tokens := strings.FieldsFunc(command, unicode.IsSpace)
	if len(tokens) == 0 {
		return
	}
	switch tokens[0] {
	case "start":
		Start()

	case "stop":
		Stop()

	case "set":
		if len(tokens) > 1 {
			k.url = fmt.Sprintf("http://%s/?\n", tokens[1])
			fmt.Printf("URL set to: %s", k.url)
		}

	default:
		fmt.Printf("Kademlia Service Command Line\n")
		return
	}

}

func (k kademliaservice) GetBlob(h objects.HCID) (b objects.Blob, err error) {
	values := url.Values{}
	values.Add("type", "blob")
	values.Add("hcid", h.Hex())
	data, err := k.getobject(values)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data, err
}
func (k kademliaservice) GetCommit(h objects.HKID) (c objects.Commit, err error) {
	values := url.Values{}
	values.Add("type", "commit")
	values.Add("hkid", h.Hex())
	data, err := k.getobject(values)
	if err != nil {
		log.Println(err)
		return c, err
	}
	c, err = objects.CommitFromBytes(data)
	return c, err
}
func (k kademliaservice) GetTag(h objects.HKID, namesegment string) (t objects.Tag, err error) {
	values := url.Values{}
	values.Add("type", "tag")
	values.Add("hkid", h.Hex())
	values.Add("namesegment", namesegment)
	data, err := k.getobject(values)
	if err != nil {
		log.Println(err)
		return t, err
	}
	t, err = objects.TagFromBytes(data)
	return t, err
}
func (k kademliaservice) GetKey(h objects.HKID) (b objects.Blob, err error) {
	values := url.Values{}
	values.Add("type", "key")
	values.Add("hkid", h.Hex())
	data, err := k.getobject(values)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data, err
}
func (k kademliaservice) PostBlob(b objects.Blob) (err error) {
	values := url.Values{}
	values.Add("type", "blob")
	values.Add("hcid", b.Hash().Hex())
	_, err = k.postobject(values, b)
	if err != nil {
		log.Println(err)
		return err
	}
	//log.Printf("Responce: %s", data)
	return err
}
func (k kademliaservice) PostTag(t objects.Tag) (err error) {
	values := url.Values{}
	values.Add("type", "tag")
	values.Add("hkid", t.Hkid.Hex())
	values.Add("namesegment", t.NameSegment)
	_, err = k.postobject(values, t.Bytes())
	if err != nil {
		log.Println(err)
		return err
	}
	//log.Printf("Responce: %s", data)
	return err
}
func (k kademliaservice) PostCommit(c objects.Commit) (err error) {
	values := url.Values{}
	values.Add("type", "commit")
	values.Add("hkid", c.Hkid.Hex())
	data, err := k.postobject(values, c.Bytes())
	if err != nil {
		log.Println(err)
		log.Printf("%s", data)
		return err
	}
	return err
}
func (k kademliaservice) PostKey(p *objects.PrivateKey) (err error) {
	values := url.Values{}
	values.Add("type", "key")
	values.Add("hkid", p.Hkid().Hex())
	data, err := k.postobject(values, p.Bytes())
	if err != nil {
		log.Println(err)
		log.Printf("%s", data)
		return err
	}
	return err
}

func (k kademliaservice) getobject(values url.Values) (data []byte, err error) {
	resp, err := http.Get(k.url + values.Encode())
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	//log.Printf("[msg] %s", data)
	return data, nil

}

func (k kademliaservice) postobject(values url.Values, b []byte) (data []byte, err error) {
	//log.Printf("post:%s%s", k.url, values.Encode())
	resp, err := http.Post(k.url+values.Encode(), "application/octet-stream",
		bytes.NewReader(b))
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	return data, nil
}

func kademliaservicefactory() kademliaservice {
	return kademliaservice{url: "http://127.0.0.1:5000/?"}
}
