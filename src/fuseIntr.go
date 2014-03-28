// Modified from Go Authors
// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"fmt"
	"log"
	"os"
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
		//cmd := exec.Command("fusermount", "-u", mountpoint)
		ccfsUnmount(mountpoint)

	}() //end func
	fs.Serve(c, FS_from_HKID_string("c09b2765c6fd4b999d47c82f9cdf7f4b659bf7c29487cc0b357b8fc92ac8ad02", mountpoint))

}

func ccfsUnmount(mountpoint string) {
	err := fuse.Unmount(mountpoint)
	if err != nil {
		log.Printf("Could not unmount: %s", err)
	}
	log.Printf("Exit-kill program")
	os.Exit(0)
}

// FS implements the hello world file system.
type FS struct {
	hkid       HKID
	mountpoint string //fs object needs to know its mountpoint
}

func FS_from_HKID_string(HKIDstring string, mountpoint string) FS {
	//get hkid from hex
	h, err := HkidFromHex(HKIDstring)
	//check if err is not nil else return h = NULL
	if err != nil {
		log.Printf("Invalid initilizing filesystem FS: %s", err)
		return FS{}
	}
	//return filesystem
	return FS{h, mountpoint}
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
		leaf:         fs_obj.hkid,
		parent:       nil,
		name:         "",
		openHandles:  map[string]bool{},
	}, nil
}

// function to save writes before ejecting mountpoint
func (fs_obj FS) Destroy() {

	ccfsUnmount(fs_obj.mountpoint)

}

// Dir implements both Node and Handle for the root directory.
type Dir struct {
	path         string
	trunc        HKID
	branch       HKID
	leaf         HID
	permission   os.FileMode
	content_type string
	parent       *Dir
	name         string
	openHandles  map[string]bool
}

func (d Dir) Attr() fuse.Attr {
	log.Println("Attr func")
	return fuse.Attr{Inode: 1, Mode: os.ModeDir | d.permission}
}

func (d Dir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	log.Printf("string=%s\n", name)
	if name == "hello" {
		return File{permission: os.FileMode(0444)}, nil
	}

	log.Printf("d.leaf is %s", d.leaf.Hex())

	//in each case, call
	switch d.content_type {
	default:
		return nil, nil
	case "commit": // a commit has a list hash
		c, err := GetCommit(d.leaf.(HKID))
		if err != nil {
			log.Printf("commit %s:", err)
			return nil, nil
		}
		//get list hash
		l, err := GetList(c.listHash) //l is the list object
		if err != nil {
			log.Printf("commit list retieval error %s:", err)
			return nil, nil
		}

		list_entry, present := l[name] //go through list entries and is it maps to the string you passed in present == 1
		if !present {
			return nil, fuse.ENOENT
		}
		//getKey to figure out permissions of the child
		_, err = GetKey(c.hkid)
		//perm := fuse.Attr{Mode: 0555}//default read permissions
		perm := os.FileMode(0555)

		if err != nil {
			log.Printf("no private key %s:", err)
			//perm =  fuse.Attr{Mode: 0755}
			perm = os.FileMode(0755)
		} else {
			log.Printf("no private key %s:", err)
		}
		if list_entry.TypeString == "blob" {
			return File{
				contentHash: list_entry.Hash.(HCID),
				permission:  perm,
				name:        name,
				parent:      &d,
			}, nil
		}

		return Dir{
			path:         d.path + "/" + name,
			trunc:        d.trunc,
			branch:       d.leaf.(HKID),
			leaf:         list_entry.Hash,
			permission:   perm,
			content_type: list_entry.TypeString,
			parent:       &d,
			name:         name,
			openHandles:  map[string]bool{},
		}, nil

	case "list":
		l, err := GetList(d.leaf.(HCID))
		if err != nil {
			log.Printf("commit list %s:", err)
			return nil, nil
		}
		list_entry, present := l[name] //go through list entries and is it maps to the string you passed in present == 1
		if !present {
			return nil, fuse.ENOENT
		}
		if list_entry.TypeString == "blob" {
			return File{
				contentHash: list_entry.Hash.(HCID),
				permission:  d.permission,
				parent:      &d,
				name:        name,
			}, nil
		}
		return Dir{path: d.path + "/" + name,
			trunc:        d.trunc,
			branch:       d.branch,
			leaf:         list_entry.Hash,
			permission:   d.permission,
			content_type: list_entry.TypeString,
			parent:       &d,
			openHandles:  map[string]bool{},
		}, nil
	case "tag":
		t, err := GetTag(d.leaf.(HKID), name) //leaf is HID
		// no blobs because blobs are for file structure

		if err != nil {
			log.Printf("not a tag %s", err)
			return nil, fuse.ENOENT
		}
		//getKey to figure out permissions of the child
		_, err = GetKey(t.hkid)
		perm := os.FileMode(0555) //default read permissions
		if err == nil {
			perm = os.FileMode(0755)
		} else {
			log.Printf("no private key %s:", err)
		}
		if t.TypeString == "blob" {
			return File{
				contentHash: t.HashBytes.(HCID),
				permission:  perm,
				name:        name,
				parent:      &d,
			}, nil
		}
		return Dir{path: d.path + "/" + name,
			trunc:        d.trunc,
			branch:       t.hkid,
			leaf:         t.HashBytes,
			permission:   perm,
			content_type: t.TypeString,
			parent:       &d,
			openHandles:  map[string]bool{},
		}, nil

	}

}

func (d Dir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	log.Println("ReadDir func")
	var l list
	var err error
	var dirDirs = []fuse.Dirent{
		{Inode: 2, Name: "hello", Type: fuse.DT_File},
	}

	if d.content_type == "tag" {
		return dirDirs, nil
	} else if d.content_type == "commit" {
		c, err := GetCommit(d.leaf.(HKID))
		log.Printf("hash is: %s", d.leaf)
		if err != nil {
			log.Printf("commit %s:", err)
			return nil, nil
		}
		l, err = GetList(c.listHash)
		if err != nil {
			log.Printf("commit list %s:", err)
			return nil, nil
		}

	} else if d.content_type == "list" {
		l, err = GetList(d.leaf.(HCID))
		if err != nil {
			log.Printf("list %s:", err)
			return nil, nil
		}
	} else {
		return nil, nil
	}
	log.Printf("list map: %s", l)
	for name, entry := range l {
		if entry.TypeString == "blob" {
			append_to_list := fuse.Dirent{Inode: 2, Name: name, Type: fuse.DT_File}
			dirDirs = append(dirDirs, append_to_list)
			// we need to append this to list + work on the next if(commit/list/tag? )
			// Type for the other one will be fuse.DT_DIR
		} else {
			append_to_list := fuse.Dirent{Inode: 2, Name: name, Type: fuse.DT_Dir}
			dirDirs = append(dirDirs, append_to_list)
		}
	} // end if range
	//	log.Printf("return dirDirs: %s", dirDirs)

	//loop through openHandles
	for openHandle, _ := range d.openHandles {
		inList := false
		for _, dir_entry := range dirDirs {
			if openHandle == dir_entry.Name {
				inList = true
				break
			}
		}
		if !inList {
			dirDirs = append(dirDirs, fuse.Dirent{Inode: 2, Name: openHandle, Type: fuse.DT_Dir})
		}
	}
	return dirDirs, nil
}

//2 types of nodes for files and directories. So call rename twice?
//Create node (directory)

func (d Dir) Create(request *fuse.CreateRequest, response *fuse.CreateResponse, intr fs.Intr) (fs.Node, fs.Handle, fuse.Error) {
	log.Printf("create node")
	log.Printf("permission: %s", request.Mode)
	log.Printf("name: %s", request.Name)
	log.Printf("flags: %s", request.Flags)
	node := File{
		contentHash: blob{}.Hash(),
		permission:  request.Mode,
		parent:      &d,
		name:        request.Name,
	}
	handle := OpenFileHandle{
		buffer: []byte{},
		parent: &d,
		name:   request.Name,
	}
	d.openHandles[handle.name] = true
	return node, handle, nil
}

//For directory node

func (f File) Rename(r *fuse.RenameRequest, newDir fs.Node, intr fs.Intr) fuse.Error {
	log.Printf("print request: %s", r)
	return nil
}

// File implements both Node and Handle for the hello file.
type File struct {
	contentHash HCID
	permission  os.FileMode
	parent      *Dir
	name        string
}

func (f File) Attr() fuse.Attr {
	log.Println("File: Attr func")
	return fuse.Attr{Inode: 1, Mode: f.permission}
	//log.Println("Attr 0444")
	//return fuse.Attr{Mode: 0444}
}

func (f File) ReadAll(intr fs.Intr) ([]byte, fuse.Error) {
	log.Println("File: ReadAll func")
	if f.contentHash == nil { // default file
		return []byte("hello, world\n"), nil
	}
	b, err := GetBlob(f.contentHash) //
	if err != nil {
		return nil, fuse.ENOENT
	}
	return b, nil
}

//nodeopener interface contains open(). Node may be used for file or directory
func (f File) Open(request *fuse.OpenRequest, response *fuse.OpenResponse, intr fs.Intr) (fs.Handle, fuse.Error) {
	b, err := GetBlob(f.contentHash) //
	if err != nil {
		return nil, fuse.ENOENT
	}
	handle := OpenFileHandle{buffer: b, parent: f.parent, name: f.name}
	f.parent.openHandles[handle.name] = true
	return handle, nil
}

type OpenFileHandle struct {
	buffer []byte
	parent *Dir
	name   string
}

//handleReader interface
func (o OpenFileHandle) Read(request *fuse.ReadRequest, response *fuse.ReadResponse, intr fs.Intr) fuse.Error {
	//readresponse contains a string of bytes
	fmt.Println("file read works! ? ")
	start := request.Offset
	stop := start + int64(request.Size)
	bufptr := o.buffer

	if stop > int64(len(bufptr)-1) {
		stop = int64(len(bufptr))
	}
	if stop == start {
		response.Data = []byte{} //new gives you a pointer
		return nil
	}

	log.Printf("start:%d", start)
	log.Printf("stop:%d", stop)
	log.Printf("length of buffer:%d", len(bufptr))
	slice := bufptr[start:stop]
	response.Data = slice //address of buffer goes to response
	log.Printf("response data:%s", response.Data)
	return nil
}

func (o OpenFileHandle) Write(request *fuse.WriteRequest, response *fuse.WriteResponse, intr fs.Intr) fuse.Error {
	log.Println("file read works! Into write ")
	start := request.Offset
	writeData := request.Data
	log.Printf("start:%d", start)
	log.Printf("length of write data:%d", len(writeData))
	if writeData != nil {
		return fuse.ENOENT
	}
	num := copy(o.buffer[start:], writeData)
	response.Size = num
	return nil
}

func (o OpenFileHandle) Release(request *fuse.ReleaseRequest, intr fs.Intr) fuse.Error {
	err := o.Publish()
	log.Printf("%s has been released!", o.name)
	if err != nil {
		return nil
	}
	delete(o.parent.openHandles, o.name)
	return fuse.ENOENT
}

//func (o OpenfileHandle)

func (o OpenFileHandle) Flush(request *fuse.FlushRequest, intr fs.Intr) fuse.Error {

	//node := request.Header //header contains nodeid - how to access?????
	// FlushRequest asks for the current state of an open file to be flushed to storage, as when a file descriptor is being closed
	//recursion to traceback to parent tag or commit
	//post buffer
	err := PostBlob(o.buffer)
	if err != nil {
		return fuse.EIO
	}
	//call parent recursively
	err = o.Publish() //get into loop on parent object
	if err != nil {
		return fuse.EIO
	}
	return nil
}

//write out file using postblob
//func (o OpenFileHandle) Release(request *fuse.ReleaseRequest, intr fs.Intr) fuse.Error {
//	request.Handle
//}
func (o OpenFileHandle) Publish() error { //name=file name
	bfrblob := blob(o.buffer)
	err := PostBlob(bfrblob)
	if err != nil {
		return err
	}
	o.parent.Publish(bfrblob.Hash(), o.name, "blob")
	return err
}
func (d Dir) Publish(h HCID, name string, typeString string) (err error) { //name=file name

	switch d.content_type {
	default:
		log.Printf("unknown type: %s", d.content_type)
		return fmt.Errorf("unknown type: %s", d.content_type)
	case "commit":
		c, err := GetCommit(d.leaf.(HKID))
		if err != nil {
			return err
		}
		l, err := GetList(c.listHash)
		if err != nil {
			return err
		}
		newList := l.add(name, h, typeString)
		c.Update(l.Hash())
		el := PostList(newList)
		if el != nil {
			return err
		}
		ec := PostCommit(c)
		if ec != nil {
			return err
		}
		return nil
	case "tag":
		t, err := GetTag(d.leaf.(HKID), name)
		if err != nil {
			return err
		}
		t.Update(h, typeString)
		et := PostTag(t)
		if et != nil {
			return err
		}
		return nil
	case "list":
		l, err := GetList(d.leaf.(HCID))
		if err != nil {
			return err
		}
		newList := l.add(name, h, typeString)
		el := PostList(newList)
		if el != nil {
			return err
		}
		d.parent.Publish(newList.Hash(), d.name, "list")
		return nil
	}
}
