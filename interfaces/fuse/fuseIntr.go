// Modified from Go Authors
// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	//"fmt"
	//"github.com/AaronGoldman/ccfs/objects"
	//"github.com/AaronGoldman/ccfs/services"
	"log"
	"os"
	//"os/signal"
)



func Start() {
	go startFSintegration()
}

func ccfsUnmount(mountpoint string) {
	err := fuse.Unmount(mountpoint)
	if err != nil {
		log.Printf("Could not unmount: %s", err)
	}
	log.Printf("Exit-kill program")
	os.Exit(0)
}


func GenerateInode(NodeID fuse.NodeID, name string) fuse.NodeID {
	return fuse.NodeID(fs.GenerateDynamicInode( uint64(NodeID), name))
}


