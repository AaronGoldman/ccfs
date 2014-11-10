//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package timeout

import (
	"fmt"
	"time"

	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

func init() {
	services.Registercommand(
		Instance,
		"timeout (Number)ms", //This is the usage string
	)
}

//Start registers timeoutservice instances
func Start() {
	services.Registerblobgeter(Instance)
	services.Registercommitgeter(Instance)
	services.Registertaggeter(Instance)
	services.Registertagsgeter(Instance)
	services.Registerkeygeter(Instance)
	running = true

}

//Stop deregisters timeoutservice instances
func Stop() {
	services.DeRegisterblobgeter(Instance)
	services.DeRegistercommitgeter(Instance)
	services.DeRegistertaggeter(Instance)
	services.DeRegistertagsgeter(Instance)
	services.DeRegisterkeygeter(Instance)
	running = false
}

type timeoutservice struct{}

//ID gets the ID string
func (timeoutservice) ID() string {
	return "timeout"
}

var waitTime = time.Millisecond * 6000

func (timeoutservice) GetBlob(objects.HCID) (objects.Blob, error) {
	time.Sleep(waitTime)
	return objects.Blob{}, fmt.Errorf("GetBlob Timeout")
}

//Running returns a bool that indicates the registration status of the service
func (timeoutservice) Running() bool {
	return running
}

func (timeoutservice) Command(command string) {
	setTimeOut, err := time.ParseDuration(command)
	fmt.Printf("Timeout Service Command Line\n")
	if err != nil || setTimeOut <= 0 {
		fmt.Printf("Please input a positive integer\n")
	} else {
		waitTime = setTimeOut
		fmt.Printf("The timeout is now %s\n", waitTime)
	}
}

func (timeoutservice) GetCommit(objects.HKID) (objects.Commit, error) {
	time.Sleep(waitTime)
	return objects.Commit{}, fmt.Errorf("GetCommit Timeout")
}
func (timeoutservice) GetTag(h objects.HKID, namesegment string) (objects.Tag, error) {
	time.Sleep(waitTime)
	return objects.Tag{}, fmt.Errorf("GetTag Timeout")
}

func (timeoutservice) GetTags(h objects.HKID) ([]objects.Tag, error) {
	time.Sleep(waitTime)
	return []objects.Tag{}, fmt.Errorf("GetTags Timeout")
}

func (timeoutservice) GetKey(objects.HKID) (objects.Blob, error) {
	time.Sleep(waitTime)
	return objects.Blob{}, fmt.Errorf("GetKey Timeout")
}

//Instance is the instance of the timeoutservice
var Instance timeoutservice
var running bool
