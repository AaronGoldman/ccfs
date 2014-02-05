// Modified from Go Authors
// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"log"
	"os"
	"os/exec"
	"os/signal"
)

func startFSintegration() {
	log.SetFlags(log.Lshortfile) //gives filename for every log statement
	mountpoint := "../mountpoint"
	err := os.MkdirAll(mountpoint, 0777)
	if err != nil {
		log.Printf("Unable to create directory in mountpoint: %s", err)
		return
	}
	c, err := fuse.Mount(mountpoint)
	if err != nil {
		log.Printf("unable to mount mountpoint: %s", err)
		return
	}

	go func() { //defining, calling and throwing to a different thread ...:O !
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, os.Kill)
		sig := <-ch //c is the name of the channel. usually there would be a target to receive the channel before the <-, but we don't need to use one
		log.Printf("Got signal: %s", sig)

		log.Printf("Exit unmount")
		cmd := exec.Command("fusermount", "-u", mountpoint)
		err = cmd.Run()
		if err != nil {
			log.Printf("Could not unmount: %s", err)
		}
		log.Printf("Exit-kill program")
		os.Exit(0)
	}() //end func
	fs.Serve(c, FS_from_HKID_string("c09b2765c6fd4b999d47c82f9cdf7f4b659bf7c29487cc0b357b8fc92ac8ad02"))

}

// FS implements the hello world file system.
type FS struct {
	hkid HKID
}

func FS_from_HKID_string(s string) FS {
	//get hkid from hex
	h, err := HkidFromHex(s)
	//check if err is not nil else return h = NULL
	if err != nil {
		log.Printf("Invalid initilizing filesystem FS: %s", err)
		return FS{}
	}
	//return filesystem
	return FS{h}
}

func (fs_obj FS) Root() (fs.Node, fuse.Error) { //returns a directory
	log.Println("Root func")
	_, err := GetKey(fs_obj.hkid)
	perm := os.FileMode(0555)
	if err != nil {
		perm = 0777
	}

	return Dir{path: "/",
		trunc:        fs_obj.hkid,
		branch:       fs_obj.hkid,
		permission:   perm,
		content_type: "commit",
		hash:         fs_obj.hkid,
	}, nil

}

// Dir implements both Node and Handle for the root directory.
type Dir struct {
	path         string
	trunc        HKID
	branch       HKID
	permission   os.FileMode
	content_type string
	hash         HID
	//add type field ?
}

func (d Dir) Attr() fuse.Attr {
	log.Println("Attr func")
	return fuse.Attr{Inode: 1, Mode: os.ModeDir | d.permission}
}

func (Dir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	log.Printf("string=%s\n", name)
	if name == "hello" {
		return File{}, nil
	}
	return nil, fuse.ENOENT
}

var dirDirs = []fuse.Dirent{
	{Inode: 2, Name: "hello", Type: fuse.DT_File},
}

func (d Dir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	log.Println("ReadDir func")
	var l list
	var err error
	if d.content_type == "tag" {
		return []fuse.Dirent{}, nil
	} else if d.content_type == "commit" {
		c, err := GetCommit(d.hash.(HKID))
		if err != nil {
			return nil, nil
		}
		l, err = GetList(c.listHash)
		if err == nil {
			return nil, nil
		}

	} else if d.content_type == "list" {
		l, err = GetList(d.hash.(HCID))
		if err == nil {
			return nil, nil
		}
	} else {
		return nil, nil
	}
	for name, entry := range l {
		if entry.TypeString == "blob" {
			append_to_list := fuse.Dirent{Inode: 2, Name: name, Type: fuse.DT_File}
			_ = append_to_list
		} // we need to append this to list + work on the next if(commit/list/tag? )
		// Type for the other one will be fuse.DT_DIR
	} // end if range
	return dirDirs, nil
}

// File implements both Node and Handle for the hello file.
type File struct{}

func (File) Attr() fuse.Attr {
	log.Println("Attr 0444")
	return fuse.Attr{Mode: 0444}
}

func (File) ReadAll(intr fs.Intr) ([]byte, fuse.Error) {
	log.Println("ReadAll func")
	return []byte("hello, world\n"), nil
}
