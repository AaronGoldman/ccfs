// Modified from Go Authors
// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"fmt"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
	"log"
	"os"
	"os/signal"
)

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

	go func() { //defining, calling and throwing to a different thread!
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

// Dir implements both Node and Handle for the root directory.
type Dir struct {
	//path         string
	//trunc        HKID
	//branch       HKID
	leaf         objects.HID
	permission   os.FileMode
	content_type string
	parent       *Dir
	name         string
	openHandles  map[string]bool
	inode        uint64 //fuse.NodeID
}

func (d Dir) Attr() fuse.Attr {
	log.Printf("Directory attributes requested\n\tName:%s", d.name)
	return fuse.Attr{Inode: d.inode, Mode: os.ModeDir | d.permission}
}

func (d Dir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	log.Printf("Directory Lookup:\n\tName: %s\n\tHID: %s", name, d.leaf.Hex())
	new_nodeID := fs.GenerateDynamicInode(d.inode, name)
	if name == "hello" {
		return File{permission: os.FileMode(0444)}, nil
	}
	//in each case, call
	switch d.content_type {
	default:
		return nil, nil
	case "commit": // a commit has a list hash
		c, err := services.GetCommit(d.leaf.(objects.HKID))
		if err != nil {
			log.Printf("commit %s:", err)
			return nil, nil
		}
		//get list hash
		l, err := services.GetList(c.ListHash()) //l is the list object
		if err != nil {
			log.Printf("commit list retieval error %s:", err)
			return nil, nil
		}

		list_entry, present := l[name] //go through list entries and is it maps to the string you passed in present == 1
		if !present {
			return nil, fuse.ENOENT
		}
		//getKey to figure out permissions of the child
		_, err = services.GetKey(c.Hkid())
		//perm := fuse.Attr{Mode: 0555}//default read permissions
		perm := os.FileMode(0555)

		if err != nil {
			log.Printf("error not nil; change file Mode %s:", err)
			//perm =  fuse.Attr{Mode: 0755}
			perm = os.FileMode(0755)
		} else {
			//log.Printf("no private key %s:", err)
		}
		if list_entry.TypeString == "blob" {
			return File{
				contentHash: list_entry.Hash.(objects.HCID),
				permission:  perm,
				name:        name,
				parent:      &d,
				inode:       new_nodeID,
			}, nil
		}

		return Dir{
			//path:         d.path + "/" + name,
			//trunc:        d.trunc,
			//branch:       d.leaf.(HKID),
			leaf:         list_entry.Hash,
			permission:   perm,
			content_type: list_entry.TypeString,
			parent:       &d,
			name:         name,
			openHandles:  map[string]bool{},
			inode:        new_nodeID,
		}, nil

	case "list":
		l, err := services.GetList(d.leaf.(objects.HCID))
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
				contentHash: list_entry.Hash.(objects.HCID),
				permission:  d.permission,
				parent:      &d,
				name:        name,
			}, nil
		}
		return Dir{
			//path: d.path + "/" + name,
			//trunc:        d.trunc,
			//branch:       d.branch,
			leaf:         list_entry.Hash,
			permission:   d.permission,
			content_type: list_entry.TypeString,
			parent:       &d,
			openHandles:  map[string]bool{},
			inode:        new_nodeID,
			name:         name,
		}, nil
	case "tag":
		t, err := services.GetTag(d.leaf.(objects.HKID), name) //leaf is HID
		// no blobs because blobs are for file structure

		if err != nil {
			log.Printf("not a tag Err:%s", err)
			return nil, fuse.ENOENT
		}
		//getKey to figure out permissions of the child
		_, err = services.GetKey(t.Hkid())
		perm := os.FileMode(0555) //default read permissions
		if err == nil {
			perm = os.FileMode(0755)
		} else {
			log.Printf("no private key %s:", err)
		}
		if t.TypeString == "blob" {
			return File{
				contentHash: t.HashBytes.(objects.HCID),
				permission:  perm,
				name:        name,
				parent:      &d,
				inode:       new_nodeID,
			}, nil
		}
		return Dir{
			//path:         d.path + "/" + name,
			//trunc:        d.trunc,
			//branch:       t.hkid,
			leaf:         t.HashBytes,
			permission:   perm,
			content_type: t.TypeString,
			parent:       &d,
			openHandles:  map[string]bool{},
			inode:        new_nodeID,
			name:         name,
		}, nil

	}

}

func (d Dir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	log.Printf("ReadDir requested:\n\tName:%s", d.name)
	var l objects.List
	var err error
	var dirDirs = []fuse.Dirent{
		{Inode: fs.GenerateDynamicInode(d.inode, "hello"),
			Name: "hello",
			Type: fuse.DT_File,
		},
	}

	if d.content_type == "tag" {
		return dirDirs, nil
	} else if d.content_type == "commit" {
		c, err := services.GetCommit(d.leaf.(objects.HKID))
		//log.Printf("hash is: %s", d.leaf)
		if err != nil {
			log.Printf("commit %s:", err)
			return nil, nil
		}
		l, err = services.GetList(c.ListHash())
		if err != nil {
			log.Printf("commit list %s:", err)
			return nil, nil
		}

	} else if d.content_type == "list" {
		l, err = services.GetList(d.leaf.(objects.HCID))
		if err != nil {
			log.Printf("list %s:", err)
			return nil, nil
		}
	} else {
		return nil, nil
	}
	//log.Printf("list map: %s", l)

	for name, entry := range l {
		if entry.TypeString == "blob" {
			append_to_list := fuse.Dirent{
				Inode: fs.GenerateDynamicInode(d.inode, name),
				Name:  name,
				Type:  fuse.DT_File,
			}
			dirDirs = append(dirDirs, append_to_list)
			// we need to append this to list + work on the next if(commit/list/tag? )
			// Type for the other one will be fuse.DT_DIR
		} else {
			append_to_list := fuse.Dirent{
				Inode: fs.GenerateDynamicInode(d.inode, name),
				Name:  name,
				Type:  fuse.DT_Dir}
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
			dirDirs = append(
				dirDirs,
				fuse.Dirent{
					Inode: fs.GenerateDynamicInode(d.inode, openHandle),
					Name:  openHandle,
					Type:  fuse.DT_Dir,
				})
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
		contentHash: objects.Blob{}.Hash(),
		permission:  request.Mode,
		parent:      &d,
		name:        request.Name,
		inode:       fs.GenerateDynamicInode(d.inode, request.Name),
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
	contentHash objects.HCID
	permission  os.FileMode
	parent      *Dir
	name        string
	inode       uint64 //fuse.NodeID
}

func (f File) Attr() fuse.Attr {
	log.Println("File attributes requested")
	return fuse.Attr{Inode: 1, Mode: f.permission}
	//log.Println("Attr 0444")
	//return fuse.Attr{Mode: 0444}
}

func (f File) ReadAll(intr fs.Intr) ([]byte, fuse.Error) {
	log.Println("File ReadAll requested")
	if f.contentHash == nil { // default file
		return []byte("hello, world\n"), nil
	}
	b, err := services.GetBlob(f.contentHash) //
	if err != nil {
		return nil, fuse.ENOENT
	}
	return b, nil
}

//nodeopener interface contains open(). Node may be used for file or directory
func (f File) Open(request *fuse.OpenRequest, response *fuse.OpenResponse, intr fs.Intr) (fs.Handle, fuse.Error) {
	b, err := services.GetBlob(f.contentHash) //
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
	inode  fuse.NodeID //we're not using this field yet
}

//handleReader interface
func (o OpenFileHandle) Read(request *fuse.ReadRequest, response *fuse.ReadResponse, intr fs.Intr) fuse.Error {
	log.Println("FileHandle Read requested")
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

	//log.Printf("start:%d", start)
	//log.Printf("stop:%d", stop)
	//log.Printf("length of buffer:%d", len(bufptr))
	slice := bufptr[start:stop]
	response.Data = slice //address of buffer goes to response
	log.Printf("FileHandle data:%s", response.Data)
	return nil
}

func (o OpenFileHandle) Write(request *fuse.WriteRequest, response *fuse.WriteResponse, intr fs.Intr) fuse.Error {
	log.Printf("FileHandle Write requested:\n\t%s", request.Data)
	start := request.Offset
	writeData := request.Data
	//log.Printf("start:%d", start)
	//log.Printf("length of write data:%d", len(writeData))
	if writeData == nil {
		return fuse.ENOENT
	}
	lenData := int(start) + (len(writeData))
	if lenData > int(len(o.buffer)) {
		//set length and capacity of buffer
		newbfr := make([]byte, (lenData), (lenData))
		copy(newbfr, o.buffer)
		response.Size = copy(newbfr[start:lenData], writeData)
		//log.Printf("before copying to o.buffer: %s", newbfr)
		o.buffer = newbfr

	} else {
		num := copy(o.buffer[start:lenData], writeData)
		response.Size = num

	}
	//log.Printf("buffer: %s", o.buffer)
	//log.Printf("write request handle: %v", request.Handle)
	err := o.Publish() //get into loop on parent object
	if err != nil {
		return fuse.EIO
	}
	return nil
}

func (o OpenFileHandle) Release(request *fuse.ReleaseRequest, intr fs.Intr) fuse.Error {
	log.Println("FileHandle Release requested:\n\tName:", o.name)
	//log.Printf("buffer is released")
	//err := o.Publish()
	//log.Printf("%s has been released!", o.name)
	//if err != nil {
	//	return nil
	//}
	//delete(o.parent.openHandles, o.name)
	return nil //fuse.ENOENT
}

//func (o OpenfileHandle)
//////// flush() ////
func (o OpenFileHandle) Flush(request *fuse.FlushRequest, intr fs.Intr) fuse.Error {
	log.Println("FileHandle Flush requested:\n\tName:", o.name)
	//node := request.Header //header contains nodeid - how to access?????
	//FlushRequest asks for the current state of an open file to be flushed to storage, as when a file descriptor is being closed
	///recursion to traceback to parent tag or commit
	//post buffer
	//log.Printf("flush request handle: %s",request.Handle)
	//log.Printf("flush \n")
	//err := PostBlob(o.buffer)
	//if err != nil {
	//	return fuse.EIO
	//}
	//call parent recursively
	//err = o.Publish() //get into loop on parent object
	//if err != nil {
	//	return fuse.EIO
	//}
	return nil
}

//write out file using postblob
//func (o OpenFileHandle) Release(request *fuse.ReleaseRequest, intr fs.Intr) fuse.Error {
//	request.Handle
//}
func (o OpenFileHandle) Publish() error { //name=file name
	//log.Printf("buffer contains: %s", o.buffer)
	bfrblob := objects.Blob(o.buffer)
	log.Printf("Posting blob %s\n-----BEGIN BLOB-------\n%s\n-------END BLOB-------", bfrblob.Hash(), bfrblob)
	err := services.PostBlob(bfrblob)
	if err != nil {
		return err
	}
	o.parent.Publish(bfrblob.Hash(), o.name, "blob")
	return err
}
func (d Dir) Publish(h objects.HCID, name string, typeString string) (err error) { //name=file name

	switch d.content_type {
	//-----BEGIN LIST-------
	//-------END LIST-------
	//-----BEGIN COMMIT-----
	//-------END COMMIT-----
	//-----BEGIN TAG--------
	//-------END TAG--------
	//-----BEGIN BLOB-------
	//-------END BLOB-------
	default:
		log.Printf("unknown type: %s", d.content_type)
		return fmt.Errorf("unknown type: %s", d.content_type)
	case "commit":
		c, err := services.GetCommit(d.leaf.(objects.HKID))
		if err != nil {
			return err
		}
		l, err := services.GetList(c.ListHash())
		if err != nil {
			return err
		}
		newList := l.Add(name, h, typeString)
		//log.Printf()
		newCommit := c.Update(newList.Hash())
		log.Printf("Posting list %s\n-----BEGIN LIST-------\n%s\n-------END LIST-------", newList.Hash(), newList)
		el := services.PostList(newList)
		if el != nil {
			return err
		}
		log.Printf("Posting commit %s\n-----BEGIN COMMIT-----\n%s\n-------END COMMIT-----", newCommit.Hash(), newCommit)
		ec := services.PostCommit(newCommit)
		if ec != nil {
			return err
		}
		return nil
	case "tag":
		//log.Printf("entered tag block\n\thkid:%s\n\tnamesegment:%s", d.leaf, name)
		t, err := services.GetTag(d.leaf.(objects.HKID), name)
		var newTag objects.Tag
		if err == nil {
			newTag = t.Update(h, typeString)
		} else {
			log.Printf("Tag %s\\%s Not Found", d.leaf, name)
			newTag = objects.NewTag(h, typeString, name, d.leaf.(objects.HKID))
		}
		log.Printf("Posting tag %s\n-----BEGIN TAG--------\n%s\n-------END TAG--------", newTag.Hash(), newTag)
		et := services.PostTag(newTag)
		if et != nil {
			return err
		}
		return nil
	case "list":
		l, err := services.GetList(d.leaf.(objects.HCID))
		if err != nil {
			return err
		}
		newList := l.Add(name, h, typeString)
		log.Printf("Posting list %s\n-----BEGIN LIST-------\n%s\n-------END LIST-------", newList.Hash(), newList)
		el := services.PostList(newList)
		if el != nil {
			return err
		}
		d.parent.Publish(newList.Hash(), d.name, "list")
		return nil
	}
}

func Start() {
	go startFSintegration()
}
