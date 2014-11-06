//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package appsscript

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

//Instance is the instance of the appsscriptservice
var Instance appsscriptservice
var running bool

func init() {
	services.Registercommand(
		Instance,
		"appsscript command", //This is the usage string
	)
	services.Registerrunner(Instance)
}

//Start registers appsscriptservice instances
func Start() {
	Instance = appsscriptservice{}
	services.Registerblobgeter(Instance)
	services.Registercommitgeter(Instance)
	services.Registertaggeter(Instance)
	services.Registerkeygeter(Instance)
	running = true
}

//Stop deregisters appsscriptservice instances
func Stop() {
	services.DeRegisterblobgeter(Instance)
	services.DeRegistercommitgeter(Instance)
	services.DeRegistertaggeter(Instance)
	services.DeRegisterkeygeter(Instance)
	running = false

}

type appsscriptservice struct{}

//Running returns a bool that indicates the registration status of the service
func (a appsscriptservice) Running() bool {
	return running
}

//ID gets the ID string
func (a appsscriptservice) ID() string {
	return "appsscript"
}

func (a appsscriptservice) Command(command string) {
	switch command {
	case "start":
		Start()

	case "stop":
		Stop()

	default:
		fmt.Printf("Appsscript Service Command Line\n")
		return
	}
}

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
