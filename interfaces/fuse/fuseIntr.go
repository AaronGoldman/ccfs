//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package fuse

import (
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fs/fstestutil"
)

var running bool

//Start mounts the local seed on the local file system
func Start() {
	fstestutil.DebugByDefault()
	go startFSintegration()
	running = true
}

func ccfsUnmount(mountpoint string) {
	err := fuse.Unmount(mountpoint)
	if err != nil {
		log.Printf("Could not unmount: %s", err)
	}
	log.Printf("Exit-kill program")
	os.Exit(0)
}

func generateInode(NodeID fuse.NodeID, name string) fuse.NodeID {
	return fuse.NodeID(fs.GenerateDynamicInode(uint64(NodeID), name))
}
