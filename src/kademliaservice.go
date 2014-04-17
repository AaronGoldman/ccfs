// kademliaservice
package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type kademliaservice struct {
	url string
}

func (k kademliaservice) GetBlob(h HCID) (b blob, err error) {
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
func (k kademliaservice) GetCommit(h HKID) (c commit, err error) {
	values := url.Values{}
	values.Add("type", "commit")
	values.Add("hkid", h.Hex())
	data, err := k.getobject(values)
	if err != nil {
		log.Println(err)
		return c, err
	}
	c, err = CommitFromBytes(data)
	return c, err
}
func (k kademliaservice) GetTag(h HKID, namesegment string) (t tag, err error) {
	values := url.Values{}
	values.Add("type", "tag")
	values.Add("hkid", h.Hex())
	values.Add("namesegment", namesegment)
	data, err := k.getobject(values)
	if err != nil {
		log.Println(err)
		return t, err
	}
	t, err = TagFromBytes(data)
	return t, err
}
func (k kademliaservice) GetKey(h HKID) (b blob, err error) {
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
func (k kademliaservice) PostBlob(b blob) (err error) {
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
func (k kademliaservice) PostTag(t tag) (err error) {
	values := url.Values{}
	values.Add("type", "tag")
	values.Add("hkid", t.Hkid().Hex())
	values.Add("namesegment", t.nameSegment)
	_, err = k.postobject(values, t.Bytes())
	if err != nil {
		log.Println(err)
		return err
	}
	//log.Printf("Responce: %s", data)
	return err
}
func (k kademliaservice) PostCommit(c commit) (err error) {
	values := url.Values{}
	values.Add("type", "commit")
	values.Add("hkid", c.Hkid().Hex())
	data, err := k.postobject(values, c.Bytes())
	if err != nil {
		log.Println(err)
		log.Printf("%s", data)
		return err
	}
	return err
}
func (k kademliaservice) PostKey(p *PrivateKey) (err error) {
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
	} else {
		//log.Printf("[msg] %s", data)
		return data, nil
	}
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
	} else {
		return data, nil
	}
}

func kademliaservicefactory() kademliaservice {
	return kademliaservice{url: "http://128.61.21.129:5000/?"}
}

func init() {
	kademliaserviceInstance = kademliaservicefactory()
	_ = kademliaserviceInstance
}

var kademliaserviceInstance kademliaservice
