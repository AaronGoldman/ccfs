//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package timeout

import (
	"fmt"
	"time"

	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

//Start registers timeoutservice instances
func Start() {
	services.Registerblobgeter(Instance)
	services.Registercommitgeter(Instance)
	services.Registertaggeter(Instance)
	services.Registertagsgeter(Instance)
	services.Registerkeygeter(Instance)
}

//Stop deregisters timeoutservice instances
func Stop() {
	services.DeRegisterblobgeter(Instance)
	services.DeRegistercommitgeter(Instance)
	services.DeRegistertaggeter(Instance)
	services.DeRegistertagsgeter(Instance)
	services.DeRegisterkeygeter(Instance)
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
