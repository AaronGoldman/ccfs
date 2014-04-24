package appsscript

import (
	"encoding/base64"
	"fmt"
	"github.com/AaronGoldman/ccfs/objects"
	"io/ioutil"
	"log"
	"net/http"
)

type appsscriptservice struct{}

func (a appsscriptservice) GetBlob(h objects.HCID) (b objects.Blob, err error) {
	quarryurl := fmt.Sprintf(
		"%s%s%s%s%s%s",
		"https://",
		"script.google.com",
		"/macros/s/AKfycbzl2R7UR2FGGVdgl_WbKabbIoku66ELRSnQ4pbkmBgDdWWvgh8b/exec?",
		//"/macros/s/AKfycbxyU7ABEmq4HS_8nb7E5ZbtJKRwuVlLBwnhUJ4VjSH0/dev?",
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
	data, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		log.Printf("%s\n", body)
		log.Println("error:", err)
		return nil, err
	}
	return data, err
}

func (a appsscriptservice) GetCommit(h objects.HKID) (c objects.Commit, err error) {
	quarryurl := fmt.Sprintf(
		"%s%s%s%s%s%s",
		"https://",
		"script.google.com",
		"/macros/s/AKfycbzl2R7UR2FGGVdgl_WbKabbIoku66ELRSnQ4pbkmBgDdWWvgh8b/exec?",
		//"/macros/s/AKfycbxyU7ABEmq4HS_8nb7E5ZbtJKRwuVlLBwnhUJ4VjSH0/dev?",
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
	data, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		log.Println("error:", err)
		return c, err
	}
	c, err = objects.CommitFromBytes(data)
	return c, err
}

func (a appsscriptservice) GetTag(h objects.HKID, namesegment string) (t objects.Tag, err error) {
	quarryurl := fmt.Sprintf(
		"%s%s%s%s%s%s%s",
		"https://",
		"script.google.com",
		"/macros/s/AKfycbzl2R7UR2FGGVdgl_WbKabbIoku66ELRSnQ4pbkmBgDdWWvgh8b/exec?",
		//"/macros/s/AKfycbxyU7ABEmq4HS_8nb7E5ZbtJKRwuVlLBwnhUJ4VjSH0/dev?",
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
	data, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		log.Println("error:", err)
		return t, err
	}
	t, err = objects.TagFromBytes(data)
	return t, err
}

func (a appsscriptservice) GetKey(h objects.HKID) (objects.Blob, error) {
	quarryurl := fmt.Sprintf(
		"%s%s%s%s%s",
		"https://",
		"script.google.com",
		"/macros/s/AKfycbzl2R7UR2FGGVdgl_WbKabbIoku66ELRSnQ4pbkmBgDdWWvgh8b/exec?",
		//"/macros/s/AKfycbxyU7ABEmq4HS_8nb7E5ZbtJKRwuVlLBwnhUJ4VjSH0/dev?",
		"type=key&hkid=",
		h.Hex(),
	)
	//log.Println(quarryurl)
	resp, err := http.Get(quarryurl)
	if err != nil {
		return objects.Blob{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return objects.Blob{}, err
	}
	data, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}
	return data, err
}

var Instance appsscriptservice = appsscriptservice{}

func init() {
	//Registerblobgeter(appsscriptserviceInstance)
	//Registercommitgeter(appsscriptserviceInstance)
	//Registertaggeter(appsscriptserviceInstance)
}
