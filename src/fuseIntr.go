// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Hellofs implements a simple "hello world" file system.

package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
)

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s MOUNTPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

func startFSintegration() {
	log.SetFlags(log.Lshortfile) //gives filename for every log statement
	/*	log.Println("In main")
		flag.Usage = Usage
		flag.Parse()

		if flag.NArg() != 1 {
			Usage()
			os.Exit(2)
		}
		mountpoint := flag.Arg(0)
	*/
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
	fs.Serve(c,  FS_from_HKID_string("c09b2765c6fd4b999d47c82f9cdf7f4b659bf7c29487cc0b357b8fc92ac8ad02" ) )

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

func (fs_obj FS) Root() (fs.Node, fuse.Error) {
	log.Println("Root func")
	return Dir{path: "/",
		 trunc: fs_obj.hkid,
		 branch: fs_obj.hkid,
			}, nil
}

// Dir implements both Node and Handle for the root directory.
type Dir struct {
	path string
	trunc HKID
        branch HKID
}

func (Dir) Attr() fuse.Attr {
	log.Println("Attr func")
	return fuse.Attr{Inode: 1, Mode: os.ModeDir | 0555}
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

func (Dir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	log.Println("ReadDir func")
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
