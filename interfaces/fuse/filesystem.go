//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file
package fuse

import (
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/AaronGoldman/ccfs/interfaces"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

var instance FS

//testing push with new origin
func startFSintegration() {
	log.SetFlags(log.Lshortfile) //gives filename for every log statement
	mountpoint := "mountpoint"
	err := os.MkdirAll(mountpoint, 0777)
	if err != nil {
		log.Printf("Unable to create directory in mountpoint: %s", err)
		return
	}
	c, err := fuse.Mount(mountpoint)
	if err != nil {
		log.Printf("Unable to mount mountpoint: %s", err)
		return
	}

	//defer profile.Start(profile.CPUProfile).Stop()
	instance = FS_from_HKID_string(interfaces.GetLocalSeed(), mountpoint)
	fs.Serve(c, instance)
}

func Stop() {
	if running {
		ccfsUnmount(instance.mountpoint)
		running = false
	}
}

type FS struct {
	hkid       objects.HKID
	mountpoint string //fs object needs to know its mountpoint
}

func FS_from_HKID_string(HKIDstring string, mountpoint string) FS {
	//get hkid from hex
	h, err := objects.HkidFromHex(HKIDstring)
	//check if err is not nil else return h = NULL
	if err != nil {
		log.Printf("Invalid initilizing filesystem FS: %s", err)
		return FS{}
	}
	//return filesystem
	return FS{h, mountpoint}
}

func (fs_obj FS) Root() (fs.Node, fuse.Error) { //returns a directory
	log.Printf("Initilizing filesystem:\n\tHKID: %s", fs_obj.hkid)
	_, err := services.GetKey(fs_obj.hkid)
	perm := os.FileMode(0555)
	if err != nil {
		perm = 0777
	}

	//return Repository{
	//	contents:   fs_obj.hkid,
	//	inode:      1,
	//	name:       "",
	//	parent:     nil,
	//	permission: perm,
	//}

	return Dir{
		//path: "/",
		//trunc:        fs_obj.hkid,
		//branch:       fs_obj.hkid,
		permission:   perm,
		content_type: "commit",
		leaf:         fs_obj.hkid,
		parent:       nil,
		name:         "",
		openHandles:  map[string]bool{},
		inode:        1,
	}, nil
}

// function to save writes before ejecting mountpoint
func (fs_obj FS) Destroy() {
	ccfsUnmount(fs_obj.mountpoint)
}
