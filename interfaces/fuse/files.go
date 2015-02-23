//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package fuse

import (
	"fmt"
	"log"
	"os"
	//"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

type file struct {
	contentHash objects.HCID
	permission  os.FileMode
	parent      *dir
	name        string
	inode       fuse.NodeID
	//Mtime		time.
	flags fuse.OpenFlags
	size  uint64
}

func (f file) String() string {
	return fmt.Sprintf(
		"File Size:\t[%d]\nFile Hash:\t%s\nFile Mode:\t%s\nFile Flags:\t%s\nFile Id:\t%v\nFile Parent Information:\n%v\n",
		f.size,
		f.contentHash,
		f.permission,
		f.flags,
		f.inode,
		f.parent,
	)
}

func (f file) Attr() fuse.Attr {
	log.Printf("File Attributes Requested: %s\n%+v", f.name, f)
	att := fuse.Attr{
		Inode:  uint64(f.inode),
		Size:   f.size,
		Blocks: f.size / 4096,
		// 	Atime:0001-01-01 00:00:00 +0000 UTC
		// 	Mtime:0001-01-01 00:00:00 +0000 UTC
		// 	Ctime:0001-01-01 00:00:00 +0000 UTC
		// 	Crtime:0001-01-01 00:00:00 +0000 UTC
		Mode: f.permission,
		Uid:  uint32(os.Getuid()),
		Gid:  uint32(os.Getgid()),
		// 	Rdev:0
		// 	Nlink:0
		Flags: uint32(f.flags),
	}

	return att
	// files.go:31: file atributes:
	// {
	// 	Inode:10526737836144204806
	// 	Size:0 Blocks:0 				// size should be
	// 	Atime:0001-01-01 00:00:00 +0000 UTC
	// 	Mtime:0001-01-01 00:00:00 +0000 UTC
	// 	Ctime:0001-01-01 00:00:00 +0000 UTC
	// 	Crtime:0001-01-01 00:00:00 +0000 UTC
	// 	Mode:-rw-r--r--
	// 	Nlink:0
	// 	Uid:0
	// 	Gid:0
	// 	Rdev:0
	// 	Flags:0
	// }
}

func (f file) ReadAll(intr fs.Intr) ([]byte, fuse.Error) {
	log.Println("File ReadAll requested")
	select {
	case <-intr:
		return nil, fuse.EINTR
	default:
	}
	if f.contentHash == nil { // default file
		return []byte("hello, world\n"), nil
	}
	b, blobErr := services.GetBlob(f.contentHash) //
	if blobErr != nil {
		return nil, fuse.ENOENT
	}
	return b, nil
}

//nodeopener interface contains open(). Node may be used for file or directory
func (f *file) Open(request *fuse.OpenRequest, response *fuse.OpenResponse, intr fs.Intr) (fs.Handle, fuse.Error) {
	log.Printf("Node Open Request:\nHeader:\t%+v\nFlags:\t%+v\nObject:\n%+v", request.Header, request.Flags, f)
	//request.dir = 0
	//   O_RDONLY int = os.O_RDONLY // open the file read-only.
	//   O_WRONLY int = os.O_WRONLY // open the file write-only.
	//   O_RDWR   int = os.O_RDWR   // open the file read-write.
	//   O_APPEND int = os.O_APPEND // append data to the file when writing.
	//   O_CREATE int = os.O_CREAT  // create a new file if none exists.
	//   O_EXCL   int = os.O_EXCL   // used with O_CREATE, file must not exist
	//   O_SYNC   int = os.O_SYNC   // open for synchronous I/O.
	//   O_TRUNC  int = os.O_TRUNC  // if possible, truncate file when opened.
	select {
	case <-intr:
		return nil, fuse.EINTR
	default:
	}
	log.Printf("File Open Request: %q", f.name)
	b, blobErr := services.GetBlob(f.contentHash)
	if blobErr != nil {
		log.Printf("Get Blob error in opening handle: %s", blobErr)
		return nil, fuse.ENOENT
	}
	log.Printf("Request Information:\nReq Header:\t%v\nReq. Dir:\t%v\nReq Flags:\t%v\n", request.Header, request.Dir, request.Flags)

	//Change current file flags to match request flags
	f.flags = request.Flags & 7

	//Create handle
	handle := openFileHandle{
		buffer:    b,
		parent:    f.parent,
		name:      f.name,
		inode:     f.inode,
		publish:   true,
		handleNum: 1,
	}

	//Check to see if handle already exists
	if f.parent.openHandles[f.name] == true {
		log.Printf("File Handle Already Exists\n")
		handle.handleNum = f.parent.openHandlesList[f.name].handleNum + 1
		f.parent.openHandlesList[f.name].handleNum++
	} else {
		f.parent.openHandles[f.name] = true
		f.parent.openHandlesList[handle.name] = &handle
	}

	response.Handle = fuse.HandleID(handle.inode)
	response.Flags = fuse.OpenResponseFlags(request.Flags)
	log.Printf("Response Information:\nRes. ID:\t%v\nRes. Flags:\t%v\n", response.Handle, response.Flags)
	return &handle, nil
}
