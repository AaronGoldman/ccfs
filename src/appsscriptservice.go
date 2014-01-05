package main

import (
	"fmt"
	"io/ioutil"
	//"log"
	"net/http"
)

type appsscriptservice struct{}

func (a appsscriptservice) getBlob(h HCID) (b blob, err error) {
	quarryurl := fmt.Sprintf(
		"%s%s%s%s%s%s",
		"https://",
		"script.google.com",
		"/macros/s/AKfycbzl2R7UR2FGGVdgl_WbKabbIoku66ELRSnQ4pbkmBgDdWWvgh8b/exec?",
		"type=blob",
		"&hcid=",
		h.Hex(),
	)
	//log.Println(quarryurl)
	resp, err := http.Get(quarryurl)
	if err != nil {
		return b, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return b, err
	}
	return body, err
}

func (a appsscriptservice) getCommit(h HKID) (c commit, err error) {
	quarryurl := fmt.Sprintf(
		"%s%s%s%s%s%s",
		"https://",
		"script.google.com",
		"/macros/s/AKfycbzl2R7UR2FGGVdgl_WbKabbIoku66ELRSnQ4pbkmBgDdWWvgh8b/exec?",
		"type=commit",
		"&hkid=",
		h.Hex(),
	)
	//log.Println(quarryurl)
	resp, err := http.Get(quarryurl)
	if err != nil {
		return c, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c, err
	}
	c, err = CommitFromBytes(body)
	return c, err
}

func (a appsscriptservice) getTag(h HKID, namesegment string) (t tag, err error) {
	quarryurl := fmt.Sprintf(
		"%s%s%s%s%s%s%s",
		"https://",
		"script.google.com",
		"/macros/s/AKfycbzl2R7UR2FGGVdgl_WbKabbIoku66ELRSnQ4pbkmBgDdWWvgh8b/exec?",
		"type=tag&hkid=",
		h.Hex(),
		"&namesegment=",
		namesegment,
	)
	//log.Println(quarryurl)
	resp, err := http.Get(quarryurl)
	if err != nil {
		return t, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return t, err
	}
	t, err = TagFromBytes(body)
	return t, err
}

func (a appsscriptservice) getKey(h HKID) (blob, error) {
	quarryurl := fmt.Sprintf(
		"%s%s%s%s%s",
		"https://",
		"script.google.com",
		"/macros/s/AKfycbzl2R7UR2FGGVdgl_WbKabbIoku66ELRSnQ4pbkmBgDdWWvgh8b/exec?",
		"type=key&hkid=",
		h.Hex(),
	)
	//log.Println(quarryurl)
	resp, err := http.Get(quarryurl)
	if err != nil {
		return blob{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return blob{}, err
	}
	return body, err
}

var appsscriptserviceInstance appsscriptservice = appsscriptservice{}
